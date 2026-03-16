/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
	kube "github.com/ewanf2/cli_tool/internal/kube"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/kubernetes"
)

var (
	kubeConfig clientcmd.ClientConfig
	kubeClient kubernetes.Interface
	kubeNamespace string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "greet",
	Short: "CLI that does stuff",
	Long: `CLI that does a lot of stuff For example:`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		kubeConfig = kube.GetKubeConfig()
		kubeClient,_ = kube.GetClientset(kubeConfig)
		kubeNamespace, _ = kube.GetNamespace(kubeConfig)
	},
	Run: func(cmd *cobra.Command, args []string) { 
		fmt.Println("Welcome to my CLI. hey lol")
	},
}
var greetCMD = &cobra.Command{
	Use: "greet",
	Short: "Greet the user innit",
	Long: "Function that greets the user i guess",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello,", args[0])
	},
}

var quickmaths = &cobra.Command{
	Use: "add",
	Short: "maths innit",
	Long: "maths innit bruv",
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		x, err := strconv.Atoi(args[0])
		if err != nil { return fmt.Errorf("invalid x: %w", err)}
		y, err2 := strconv.Atoi(args[1])
		if err2 != nil { return fmt.Errorf("invalid y: %w", err2)}
		fmt.Println(x+y)
		return nil
	},
}
// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.AddCommand(greetCMD)
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cli_tool.yaml)")
	rootCmd.AddCommand(quickmaths)
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


