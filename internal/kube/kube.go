
package kube

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	// "k8s.io/client-go/util/homedir"
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