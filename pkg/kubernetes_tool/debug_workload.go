package kubernetes_tool

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	mcp_golang "github.com/metoro-io/mcp-golang"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DebugWorkloadArgs struct {
	Namespace string `json:"namespace" jsonschema:"required,description=Namespace where the workload resides"`
	Name      string `json:"name" jsonschema:"required,description=Name of the workload"`
}

func (k *KubernetesTool) DebugWorkload(ctx context.Context, args DebugWorkloadArgs) (*mcp_golang.ToolResponse, error) {
	fmt.Println("debug_workload")

	var (
	// opts corev1.PodLogOptions

	// k8sNS = &corev1.Namespace{
	// 	ObjectMeta: metav1.ObjectMeta{
	// 		Name: args.Name,
	// 	},
	// }

	// _, err = k.client.
	// 	CoreV1().
	// 	Namespaces().
	// 	Create(ctx, k8sNS, opts)
	)

	// if err != nil {
	// 	return nil, err
	// }

	// TODO:
	// - current pod state
	// - describe pod
	// - pod logs

	// k.client.CoreV1().Pods(args.Namespace).GetLogs(args.Name, nil)
	// k.client.CoreV1().Events(args.Namespace).List()

	var podlist, err = k.client.CoreV1().Pods(args.Namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var events []Event
	var pods []string

	// TODO: need to handle multiple pods
	for _, v := range podlist.Items {
		if !strings.Contains(v.Name, args.Name) {
			continue
		}

		fmt.Println("v.Name:", v.Name)
		fmt.Println("v.Status:", v.Status.String())

		fmt.Println("uid:", v.UID)

		eventList, err := k.client.CoreV1().Events(args.Namespace).List(ctx, metav1.ListOptions{
			FieldSelector: fmt.Sprintf("involvedObject.uid=%s", v.UID),
		})

		if err != nil {
			return nil, err
		}

		pods = append(pods, v.Name)

		for _, event := range eventList.Items {
			events = append(events, Event{
				CreationTimestamp: event.CreationTimestamp,
				Message:           event.Message,
				Type:              event.Type,
				Reason:            event.Reason,
			})
		}

		break
	}

	if len(pods) > 0 {
		logs, err := k.getLogs(ctx, pods[0], args.Namespace)
		if err != nil {
			if err != ErrWaitingToStart {
				return nil, err
			}

			// TODO: do something here to explain to the LLM that we don't have any logs
		}

		if logs != nil {
			fmt.Println("logs:", logs)
		}
	}

	// 	// Copy logs to stdout (or process however you want)
	// 	_, err = io.Copy(os.Stdout, logStream)
	// if err != nil {
	// 	panic(err.Error())
	// }

	b, err := json.Marshal(events)
	if err != nil {
		return nil, err
	}

	fmt.Println("events b:", string(b))

	return &mcp_golang.ToolResponse{
		Content: []*mcp_golang.Content{
			{
				Type: mcp_golang.ContentTypeText,
				TextContent: &mcp_golang.TextContent{
					Text: fmt.Sprintf("Here are the events from the pod in JSON form. They are similar to `kubectl describe pod ... :` %s", string(b)),
				},
			},
		},
	}, nil
}

type Event struct {
	CreationTimestamp metav1.Time `json:"creation_timestamp,omitempty"`
	Message           string      `json:"message,omitempty"`
	Type              string      `json:"type,omitempty"`
	Reason            string      `json:"reason,omitempty"`
}

var ErrWaitingToStart = errors.New("waiting to start")

func (k *KubernetesTool) getLogs(ctx context.Context, name, namespace string) ([]string, error) {
	// TODO: need to handle multiple pods
	// TODO: need to handle multiple containers
	var (
		result = k.client.CoreV1().Pods(namespace).GetLogs(name, &corev1.PodLogOptions{}).Do(ctx)
		err    = result.Error()
	)

	if err != nil {
		if strings.Contains(err.Error(), "is waiting to start") {
			return nil, ErrWaitingToStart
		}

		return nil, err
	}

	b, err := result.Raw()
	if err != nil {
		return nil, err
	}

	fmt.Println("b:", string(b))

	return []string{""}, nil
}
