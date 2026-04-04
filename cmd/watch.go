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
	Short: "Get pods",
	Long: "Get pods init",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		releaseName := fmt.Sprintf(args[0])
		kube.GetDeploy(kubeClient, kubeNamespace,releaseName)
		kube.GetStateful(kubeClient, kubeNamespace, releaseName)
		kube.GetDaemonset(kubeClient, kubeNamespace, releaseName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getpods)
	getpods.Flags().String("n",kubeNamespace, "Namespace" )
}


