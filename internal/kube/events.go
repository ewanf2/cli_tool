package kube

import (
	"context"
	"fmt"
	"sort"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

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

func FindPodEvents(kubeClient kubernetes.Interface ,selector *metav1.LabelSelector, namespace string ) ([]corev1.Event, error) {

	pods, err := GetPods(kubeClient, namespace, selector)
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
