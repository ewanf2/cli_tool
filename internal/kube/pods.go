package kube

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Get list of pods based off selector and namespace
func GetPods(kubeClient kubernetes.Interface, namespace string, selector *metav1.LabelSelector) ([]corev1.Pod, error) {
	label, _ := metav1.LabelSelectorAsSelector(selector)
	pods, err := kubeClient.CoreV1().Pods(namespace).List(context.TODO(),
		metav1.ListOptions{LabelSelector: label.String()})

	if err != nil {
		return nil, err
	}
	return pods.Items, nil
}

// Get pod phases based off list of pods
func GetPodPhases(kubeClient kubernetes.Interface, namespace string, selector *metav1.LabelSelector) (r int, p int, f int, s int, u int, err error) {
	pods, err := GetPods(kubeClient, namespace, selector)
	if err != nil {
		err = fmt.Errorf("Error getting pods: %s", err)
		return
	}
	running, pending, failed, succeeded, unknown := 0, 0, 0, 0, 0
	for _, p := range pods {
		switch p.Status.Phase {
		case "Running":
			running++
		case "Pending":
			pending++
		case "Failed":
			failed++
		case "Succeeded":
			succeeded++
		case "Unknown":
			unknown++
		}
	}
	return running, pending, failed, succeeded, unknown, nil
}

type ContainerSummary struct {
	PodName      string
	ContainerName         string
	State        string
	Reason       string
	Message      string
	RestartCount int32
	LastTerminationReason string
}

func GetContainerStatuses(pods []corev1.Pod) ([]ContainerSummary, error) {
	summary := []ContainerSummary{}
	for _, pod := range pods {
		containerStatuses := pod.Status.ContainerStatuses
		for _, containerStat := range containerStatuses {
			switch {
			case containerStat.State.Running != nil:
				summary = append(summary, ContainerSummary{
					PodName:      pod.Name,
					ContainerName:         containerStat.Name,
					State:        fmt.Sprintf("%t", containerStat.Ready),
					Reason:       "",
					Message:      "",
					RestartCount: containerStat.RestartCount,
					LastTerminationReason: "",
				})
			case containerStat.State.Waiting != nil:
				summary = append(summary, ContainerSummary{
					PodName:      pod.Name,
					ContainerName:         containerStat.Name,
					State:        fmt.Sprintf("%t", containerStat.Ready),
					Reason:       fmt.Sprintf("[%s]", containerStat.State.Waiting.Reason),
					Message: 	  containerStat.State.Waiting.Message,
					RestartCount: containerStat.RestartCount,
					LastTerminationReason: "",
				})
			case containerStat.State.Terminated != nil:
				summary = append(summary, ContainerSummary{
					PodName:      pod.Name,
					ContainerName:         containerStat.Name,
					State:        fmt.Sprintf("%t", containerStat.Ready),
					Reason:       fmt.Sprintf("[%s]", containerStat.State.Terminated.Reason),
					Message: 	  containerStat.State.Waiting.Message,
					RestartCount: containerStat.RestartCount,
					LastTerminationReason: "",
				})
			}
			
		}
	}
	return summary, nil
}


