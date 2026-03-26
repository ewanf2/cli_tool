
package kube

import (
	"os"
	"context"
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// "k8s.io/client-go/util/homedir"
	"github.com/olekukonko/tablewriter"
	
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

func GetDeploy(kubeClient kubernetes.Interface, namespace string, release string) ( error) {
	ctx := context.TODO()
	filter := metav1.ListOptions{ 
		LabelSelector : fmt.Sprintf("release=%s",release),
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Deployment name","Type","Status","Reason"})
	d,_ := kubeClient.AppsV1().Deployments(namespace).List(ctx,filter)
	for _, d:= range d.Items {
			for _, c := range d.Status.Conditions {
				table.Append([]string{d.Name,string(c.Type),string(c.Status),string(c.Reason)})
			}
		}
	table.Render()
	return nil
	}

func GetStateful(kubeClient kubernetes.Interface, namespace string, release string) ( error) {
	ctx := context.TODO()
	filter := metav1.ListOptions{ 
		LabelSelector : fmt.Sprintf("release=%s",release),
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"resource name","Type","Status","Reason"})
	d,_ := kubeClient.AppsV1().StatefulSets(namespace).List(ctx,filter)
	for _, d:= range d.Items {
			for _, c := range d.Status.Conditions {
				table.Append([]string{d.Name,string(c.Type),string(c.Status),string(c.Reason)})
			}
		}
	table.Render()
	return nil
	}