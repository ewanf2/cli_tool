/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (	
	"fmt"
	"github.com/spf13/cobra"
	// "k8s.io/apimachinery/pkg/api/errors"

	// "k8s.io/client-go/kubernetes"
	// "k8s.io/client-go/tools/clientcmd"
	// "k8s.io/client-go/util/homedir"
	// "github.com/ewanf2/cli_tool/internal/kube"
	kube "github.com/ewanf2/cli_tool/internal/kube"
)

var getpods = &cobra.Command{
	Use: "watch",
	Short: "Display workload status of Helm release",
	Long: "Shows all statefulsets,deployments and daemonsets under a helm release. Displays pod readiness, container statuses and pod events",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		releaseName := fmt.Sprintf(args[0])
		
		deployments,_ := kube.GetDeploy(kubeClient, kubeNamespace,releaseName)
		daemonsets,_ := kube.GetDaemonset(kubeClient, kubeNamespace, releaseName)
		statefulsets,_ := kube.GetStateful(kubeClient, kubeNamespace, releaseName)

		kube.DescribeWorkload(daemonsets, kubeClient, kubeNamespace)
		kube.DescribeWorkload(deployments, kubeClient, kubeNamespace)
		kube.DescribeWorkload(statefulsets, kubeClient, kubeNamespace)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getpods)
	getpods.Flags().String("n",kubeNamespace, "Namespace" )
}


