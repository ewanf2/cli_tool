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
	Use: "get",
	Short: "Get pods",
	Long: "Get pods init",
	// Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		kube.GetDeploy(kubeClient, kubeNamespace,"elastic-ef170")
		kube.GetStateful(kubeClient, kubeNamespace, "elastic-ef170")

		fmt.Println("worked")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getpods)
	getpods.Flags().String("n",kubeNamespace, "Namespace" )
}


