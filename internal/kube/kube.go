package kube

import (
	"context"
	"fmt"
	"sort"
	"github.com/pterm/pterm"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	// "k8s.io/client-go/tools/events"
)

// Creating kubeconfig file, contains things like namespaces and stuff
func GetKubeConfig() (clientcmd.ClientConfig) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverr := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverr)
	
	return kubeConfig
}

// Creating clientset based off kubeconfig. Object with methods of
func GetClientset(kubeconfig clientcmd.ClientConfig) (kubernetes.Interface, error) {

	restConfig, err := kubeconfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("Unable to create rest config: %v", err)
	}

	client, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("Unable to create a client: %v", err)
	}

	return client, nil
}

func GetNamespace(kubeconfig clientcmd.ClientConfig) (string, error) {
	
	ns, _, err := kubeconfig.Namespace()
	if err != nil {
		return "Unable to get namespace", err
	}
	return ns, nil
}

func GetEvents(kubeClient kubernetes.Interface, namespace string, name string, kind string) ([]corev1.Event, error) {
	ctx := context.TODO()
	filter := metav1.ListOptions{ 
		FieldSelector : fmt.Sprintf("involvedObject.name=%s,involvedObject.kind=%s", name, kind),
	}
	events, err := kubeClient.CoreV1().Events(namespace).List(ctx,filter)
	if err != nil {
		return nil, err
	}

	sort.Slice(events.Items, func(i, j int) bool {
		return events.Items[i].LastTimestamp.After(events.Items[j].LastTimestamp.Time)
	})
	return events.Items, nil
}

// GetPods retrieves all pods matching a selector
// and returns a slice of pods
func GetPods(kubeClient kubernetes.Interface,namespace string, selector string) ([]corev1.Pod,error) {
	pods, err := kubeClient.CoreV1().Pods(namespace).List(context.TODO(), 
	metav1.ListOptions{ LabelSelector: selector}, )
	if err != nil {
		return nil, err
	}
	return pods.Items,nil
}

func FindPodEvents(kubeClient kubernetes.Interface ,selector *metav1.LabelSelector, namespace string ) ([]corev1.Event, error) {
	label := labels.Set(selector.MatchLabels).AsSelector()
	pods, err := GetPods(kubeClient, namespace, label.String())
	if err != nil { return nil,err}
	var eventList []corev1.Event
	for _,p := range pods {
		events,_ := GetEvents(kubeClient,namespace,p.Name, "Pod")
		eventList = append(eventList, events...)
	}
	sort.Slice(eventList, func(i, j int) bool {
		return eventList[i].LastTimestamp.After(eventList[j].LastTimestamp.Time)
	})
	return eventList, nil
}

func CreateEventTreeNode(events []corev1.Event) ([]pterm.TreeNode, error) {
	var eventNode []pterm.TreeNode
	if len(events) == 0 { return nil,nil }
	for _,e := range events[:(min(3,len(events)))] {
		
		if e.Type == "Warning" { 
			text := fmt.Sprintf("[%s] %s %s (x%d) %s",e.Type,e.Reason, e.Message, e.Count, e.LastTimestamp)
			text = pterm.FgRed.Sprint(text)
			eventNode = append(eventNode, pterm.TreeNode{Text: text})
		} else { 
			text := fmt.Sprintf("[%s] %s (x%d) %s",e.Type,e.Reason, e.Count, e.LastTimestamp)
			text = pterm.FgGreen.Sprint(text)
			eventNode = append(eventNode, pterm.TreeNode{Text: text})
		}
	}
	return eventNode, nil
}

// GetDeploy retrieves all statefulsets associated with a helm release
// and reports the health status of each StatefulSet
func GetDeploy(kubeClient kubernetes.Interface, namespace string, release string) ( error) {
	ctx := context.TODO()
	filter := metav1.ListOptions{ 
		LabelSelector : "app.kubernetes.io/managed-by=Helm",
	}

	d,_ := kubeClient.AppsV1().Deployments(namespace).List(ctx,filter)
	var deployments []appsv1.Deployment
	for _, d:= range d.Items {
		if d.Annotations["meta.helm.sh/release-name"] == release {
			deployments = append(deployments,d )
		}
	}

	if deployments == nil {
		return nil
	}
	for _, d:= range deployments {
		selector := d.Spec.Selector
		events, _ := FindPodEvents(kubeClient, selector, namespace)

		var status pterm.TreeNode
		for _,c := range ss.Status.Conditions {
			fmt.Println("yep")
			if c.Type == "Ready" {
				status = pterm.TreeNode{ Text: fmt.Sprintf("%s - %s",c.Status,c.Reason)}
			}
		}

		ready := fmt.Sprintf("%d/%d",d.Status.Replicas,*d.Spec.Replicas)
		eventNode,_ := CreateEventTreeNode(events)
		tree := pterm.TreeNode{ Text: d.Name,
			Children: []pterm.TreeNode{
				{Text: "Kind: Deployment"},
				{Text: fmt.Sprintf("Image: %s",d.Spec.Template.Spec.Containers[0].Image) },
				{Text: fmt.Sprintf("Replicas: %s",ready)},
				{Text: "Events", Children: eventNode,
			},
		},}
		pterm.DefaultTree.WithRoot(tree).Render()
		}

	return nil
	}

// Get Stateful retrieves all statefulsets associated with a helm release
// and reports the health status of each StatefulSet
func GetStateful(kubeClient kubernetes.Interface, namespace string, release string) ( error) {
	ctx := context.TODO()
	filter := metav1.ListOptions{ 
		LabelSelector : "app.kubernetes.io/managed-by=Helm",
	}

	d,_ := kubeClient.AppsV1().StatefulSets(namespace).List(ctx,filter)
	var statefuls []appsv1.StatefulSet
	
	for _, d:= range d.Items {
		if d.Annotations["meta.helm.sh/release-name"] == release {
			statefuls = append(statefuls,d )
		}
	}
	if statefuls == nil {return nil}

	for _, ss:= range statefuls {
		selector := ss.Spec.Selector
	
		events, _ := FindPodEvents(kubeClient, selector, namespace)
		ready := fmt.Sprintf("%d/%d",ss.Status.Replicas,*ss.Spec.Replicas)

		eventNode,_ := CreateEventTreeNode(events)
		
		tree := pterm.TreeNode{ Text: ss.Name,
			Children: []pterm.TreeNode{
				{Text: "Kind: StatefulSet"},
				{Text: fmt.Sprintf("Image: %s",ss.Spec.Template.Spec.Containers[0].Image) },
				{Text: fmt.Sprintf("Ready Replicas: %s",ready)},
				{Text: "Events", Children: eventNode,
			},
		},}
		pterm.DefaultTree.WithRoot(tree).Render()
		}
		
	return nil
	}

// GetDaemonset retrieves all daemonsets associated with a helm release
// and reports the health status of each StatefulSet
func GetDaemonset(kubeClient kubernetes.Interface, namespace string, release string) ( error) {
	ctx := context.TODO()
	filter := metav1.ListOptions{ 
		LabelSelector : "app.kubernetes.io/managed-by=Helm",
	}

	d,_ := kubeClient.AppsV1().DaemonSets(namespace).List(ctx,filter)
	var daemons []appsv1.DaemonSet
	for _, d:= range d.Items {
		if d.Annotations["meta.helm.sh/release-name"] == release {
			daemons = append(daemons,d )
		}}
	if daemons == nil {
		return nil
	}

	for _, ds:= range daemons {
		selector := ds.Spec.Selector
	
		events, _ := FindPodEvents(kubeClient, selector, namespace)
	
		ready := fmt.Sprintf("%d/%d",ds.Status.NumberReady,ds.Status.DesiredNumberScheduled)
		eventNode,_ := CreateEventTreeNode(events)
		if eventNode == nil { eventNode = append(eventNode, pterm.TreeNode{Text: "No events found"})}
		tree := pterm.TreeNode{ Text: ds.Name,
			Children: []pterm.TreeNode{
				{Text: "Kind: Deployment"},
				{Text: fmt.Sprintf("Image: %s",ds.Spec.Template.Spec.Containers[0].Image) },
				{Text: fmt.Sprintf("Replicas: %s",ready)},
				{Text: "Events: ", Children: eventNode,
			},
		},}
		pterm.DefaultTree.WithRoot(tree).Render()
		}

	return nil
}