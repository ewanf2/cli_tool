package kube

import (
	"fmt"

	"github.com/pterm/pterm"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

// Helper function that creates a pterm.TreeNode item out of kubernetes events
func CreateEventTreeNode(events []corev1.Event) ([]pterm.TreeNode, error) {
	var eventNode []pterm.TreeNode
	if len(events) == 0 {
		eventNode = append(eventNode, pterm.TreeNode{Text: pterm.FgLightBlue.Sprint("No events found")})
		return eventNode, nil
	}
	for _, e := range events[:(min(3, len(events)))] {

		if e.Type == "Warning" {
			text := fmt.Sprintf("[%s] %s %s (x%d) %s", e.Type, e.Reason, e.Message, e.Count, e.LastTimestamp)
			text = pterm.FgRed.Sprint(text)
			eventNode = append(eventNode, pterm.TreeNode{Text: text})
		} else {
			text := fmt.Sprintf("[%s] %s (x%d) %s", e.Type, e.Reason, e.Count, e.LastTimestamp)
			text = pterm.FgGreen.Sprint(text)
			eventNode = append(eventNode, pterm.TreeNode{Text: text})
		}
	}
	return eventNode, nil
}

func DescribeDeployment(deployments []appsv1.Deployment, kubeClient kubernetes.Interface, namespace string) ([]pterm.TreeNode, error) {
	var treeList []pterm.TreeNode
	for _, d := range deployments {
		selector := d.Spec.Selector
		events, _ := FindPodEvents(kubeClient, selector, namespace)

		running, pending, failed, _, _, _ := GetPodPhases(kubeClient, namespace, selector)
		podPhases := fmt.Sprintf("%d Running | %d Pending | %d Failed", running, pending, failed)

		ready := fmt.Sprintf("%d/%d", d.Status.Replicas, *d.Spec.Replicas)
		eventNode, _ := CreateEventTreeNode(events)

		tree := pterm.TreeNode{Text: d.Name,
			Children: []pterm.TreeNode{
				{Text: "Kind: Deployment"},
				{Text: fmt.Sprintf("Image: %s", d.Spec.Template.Spec.Containers[0].Image)},
				{Text: fmt.Sprintf("Replicas: %s", ready)},
				{Text: fmt.Sprintf("Pod phases: %s", podPhases)},
				{Text: "Events", Children: eventNode},
			}}
		pterm.DefaultTree.WithRoot(tree).Render()
		treeList = append(treeList, tree)
	}
	return treeList, nil
}

func DescribeStatefulset(statefuls []appsv1.StatefulSet, kubeClient kubernetes.Interface, namespace string) ([]pterm.TreeNode, error) {
	var treeList []pterm.TreeNode
	for _, ss := range statefuls {
		selector := ss.Spec.Selector

		events, _ := FindPodEvents(kubeClient, selector, namespace)
		ready := fmt.Sprintf("%d/%d", ss.Status.Replicas, *ss.Spec.Replicas)

		running, pending, failed, _, _, _ := GetPodPhases(kubeClient, namespace, selector)
		podPhases := fmt.Sprintf("%d Running | %d Pending | %d Failed", running, pending, failed)

		eventNode, _ := CreateEventTreeNode(events)

		tree := pterm.TreeNode{Text: ss.Name,
			Children: []pterm.TreeNode{
				{Text: "Kind: StatefulSet"},
				{Text: fmt.Sprintf("Image: %s", ss.Spec.Template.Spec.Containers[0].Image)},
				{Text: fmt.Sprintf("Ready Replicas: %s", ready)},
				{Text: fmt.Sprintf("Pod phases: %s", podPhases)},
				{Text: "Events", Children: eventNode},
			}}
		pterm.DefaultTree.WithRoot(tree).Render()
	}
	return treeList, nil
}

func DescribeDaemonset(daemons []appsv1.DaemonSet, kubeClient kubernetes.Interface, namespace string) ([]pterm.TreeNode, error) {
	var treeList []pterm.TreeNode
	for _, ds := range daemons {
		selector := ds.Spec.Selector

		events, _ := FindPodEvents(kubeClient, selector, namespace)

		running, pending, failed, _, _, _ := GetPodPhases(kubeClient, namespace, selector)
		podPhases := fmt.Sprintf("%d Running | %d Pending | %d Failed", running, pending, failed)

		ready := fmt.Sprintf("%d/%d", ds.Status.NumberReady, ds.Status.DesiredNumberScheduled)
		eventNode, _ := CreateEventTreeNode(events)
		if eventNode == nil {
			eventNode = append(eventNode, pterm.TreeNode{Text: "No events found"})
		}
		tree := pterm.TreeNode{Text: ds.Name,
			Children: []pterm.TreeNode{
				{Text: "Kind: Deployment"},
				{Text: fmt.Sprintf("Image: %s", ds.Spec.Template.Spec.Containers[0].Image)},
				{Text: fmt.Sprintf("Replicas: %s", ready)},
				{Text: fmt.Sprintf("Pod phases: %s", podPhases)},
				{Text: "Events: ", Children: eventNode},
			}}
		pterm.DefaultTree.WithRoot(tree).Render()
	}
	return treeList, nil
}

// Redner pterm.TreeNode to visualize Kubernetes workloads in the terminal
func DescribeWorkload(workloads any, kubeClient kubernetes.Interface, namespace string) error {
	switch w := workloads.(type) {
	case []appsv1.Deployment:
		treeList, err := DescribeDeployment(w, kubeClient, namespace)
		if err != nil {
			return err
		}
		for _, t := range treeList {
			pterm.DefaultTree.WithRoot(t).Render()
		}
	case []appsv1.StatefulSet:
		treeList, err := DescribeStatefulset(w, kubeClient, namespace)
		if err != nil {
			return err
		}
		for _, t := range treeList {
			pterm.DefaultTree.WithRoot(t).Render()
		}
	case []appsv1.DaemonSet:
		treeList, err := DescribeDaemonset(w, kubeClient, namespace)
		if err != nil {
			return err
		}
		for _, t := range treeList {
			pterm.DefaultTree.WithRoot(t).Render()
		}
	default:
		return fmt.Errorf("Unsupported type for %T", workloads)
	}

	return nil
}
