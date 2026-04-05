package kube

import (
	"fmt"

	"github.com/pterm/pterm"
	corev1 "k8s.io/api/core/v1"
)

func CreateEventTreeNode(events []corev1.Event) ([]pterm.TreeNode, error) {
	var eventNode []pterm.TreeNode
	if len(events) == 0 { 
		eventNode = append(eventNode, pterm.TreeNode{Text: pterm.FgLightBlue.Sprint("No events found")})
		return eventNode,nil }
	for _,e := range events[:(min(3,len(events)))] {
		
		if e.Type == "Warning" { 
			text := fmt.Sprintf("[%s] %s %s (x%d) %s",e.Type,e.Reason, e.Message, e.Count, e.LastTimestamp)
			text = pterm.FgRed.Sprint(text)
			eventNode = append(eventNode, pterm.TreeNode{Text: text})
		} else { 
			text := fmt.Sprintf("[%s] %s (x%d) %s",e.Type,e.Reason, e.Count, e.LastTimestamp)
			text = pterm.FgGreen.Sprint(text)
			eventNode = append(eventNode, pterm.TreeNode{Text: text})
		}
	}
	return eventNode, nil
}