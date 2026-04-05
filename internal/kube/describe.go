package kube

import (
	"fmt"

	"github.com/pterm/pterm"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

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

func DescribeWorkload(workloads any, kubeClient kubernetes.Interface, namespace string, selector *metav1.LabelSelector) error {
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
	}

	return nil
}
