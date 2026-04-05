package kube

import (
	"context"
	"fmt"
	"sort"
	"github.com/pterm/pterm"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// "k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	// "k8s.io/client-go/tools/events"
)

// GetDeploy retrieves all deployments associated with a helm release. Returns a list of items - appv1.Deployment
func GetDeploy(kubeClient kubernetes.Interface, namespace string, release string) ([]appsv1.Deployment, error) {
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

	if deployments == nil {return nil,nil}
	return deployments,nil 
}

// Get Stateful retrieves all statefulsets associated with a helm release
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

		running, pending, failed,_,_, _:= GetPodPhases(kubeClient, namespace,selector)
		podPhases := fmt.Sprintf("%d Running | %d Pending | %d Failed",running,pending,failed)

		eventNode,_ := CreateEventTreeNode(events)
		
		tree := pterm.TreeNode{ Text: ss.Name,
			Children: []pterm.TreeNode{
				{Text: "Kind: StatefulSet"},
				{Text: fmt.Sprintf("Image: %s",ss.Spec.Template.Spec.Containers[0].Image) },
				{Text: fmt.Sprintf("Ready Replicas: %s",ready)},
				{Text: fmt.Sprintf("Pod phases: %s",podPhases)},
				{Text: "Events", Children: eventNode,
			},
		},}
		pterm.DefaultTree.WithRoot(tree).Render()
		}
		
	return nil
	}

// GetDaemonset retrieves all daemonsets associated with a helm release
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

		running, pending, failed,_,_, _:= GetPodPhases(kubeClient, namespace,selector)
		podPhases := fmt.Sprintf("%d Running | %d Pending | %d Failed",running,pending,failed)
	
		ready := fmt.Sprintf("%d/%d",ds.Status.NumberReady,ds.Status.DesiredNumberScheduled)
		eventNode,_ := CreateEventTreeNode(events)
		if eventNode == nil { eventNode = append(eventNode, pterm.TreeNode{Text: "No events found"})}
		tree := pterm.TreeNode{ Text: ds.Name,
			Children: []pterm.TreeNode{
				{Text: "Kind: Deployment"},
				{Text: fmt.Sprintf("Image: %s",ds.Spec.Template.Spec.Containers[0].Image) },
				{Text: fmt.Sprintf("Replicas: %s",ready)},
				{Text: fmt.Sprintf("Pod phases: %s",podPhases)},
				{Text: "Events: ", Children: eventNode,
			},
		},}
		pterm.DefaultTree.WithRoot(tree).Render()
		}

	return nil
}