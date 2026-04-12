package kube

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Creating kubeconfig file, contains things like namespaces and stuff
func GetKubeConfig() clientcmd.ClientConfig {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverr := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverr)

	return kubeConfig
}

// Creating clientset based off kubeconfig.
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

// Get Current namespace
func GetNamespace(kubeconfig clientcmd.ClientConfig) (string, error) {

	ns, _, err := kubeconfig.Namespace()
	if err != nil {
		return "Unable to get namespace", err
	}
	return ns, nil
}

// Helper function for defaults
func DefaultString(val string, def string) string {
	if val == "" { 
		return def
	}
	return val
}