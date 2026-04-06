package kube

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// GetDeploy retrieves all deployments associated with a helm release. Returns a list of items - appv1.Deployment
func GetDeploy(kubeClient kubernetes.Interface, namespace string, release string) ([]appsv1.Deployment, error) {
	ctx := context.TODO()
	filter := metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/managed-by=Helm",
	}

	d, _ := kubeClient.AppsV1().Deployments(namespace).List(ctx, filter)
	var deployments []appsv1.Deployment
	for _, d := range d.Items {
		if d.Annotations["meta.helm.sh/release-name"] == release {
			deployments = append(deployments, d)
		}
	}

	if deployments == nil {
		return nil, nil
	}
	return deployments, nil
}

// Get Stateful retrieves all statefulsets associated with a helm release
func GetStateful(kubeClient kubernetes.Interface, namespace string, release string) ([]appsv1.StatefulSet, error) {
	ctx := context.TODO()
	filter := metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/managed-by=Helm",
	}

	d, _ := kubeClient.AppsV1().StatefulSets(namespace).List(ctx, filter)
	var statefuls []appsv1.StatefulSet

	for _, d := range d.Items {
		if d.Annotations["meta.helm.sh/release-name"] == release {
			statefuls = append(statefuls, d)
		}
	}
	if statefuls == nil {
		return nil,nil
	}

	

	return statefuls,nil
}

// GetDaemonset retrieves all daemonsets associated with a helm release
func GetDaemonset(kubeClient kubernetes.Interface, namespace string, release string) ([]appsv1.DaemonSet, error) {
	ctx := context.TODO()
	filter := metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/managed-by=Helm",
	}

	d, _ := kubeClient.AppsV1().DaemonSets(namespace).List(ctx, filter)
	var daemons []appsv1.DaemonSet
	for _, d := range d.Items {
		if d.Annotations["meta.helm.sh/release-name"] == release {
			daemons = append(daemons, d)
		}
	}
	if daemons == nil {
		return nil, nil
	}

	

	return daemons,nil
}
