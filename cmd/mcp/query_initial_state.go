package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	mcp_golang "github.com/metoro-io/mcp-golang"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/restmapper"
)

type QueryInitialStateArgs struct{}

/*
	- List Namespaces
	- List Deployments in those namespaces
	- List the pods in those namespaces
	- DON'T GRAG THE LOGS JUST YET. LET THE LLM ASK IF IT NEEDS IT. SOMETHING WILL HAVE TO BE DONE WITH THE CONTEXT. Grab the events/logs for each pod (maybe just the entire namespace)
	- List ALL entities
	- Get cluster-info
*/

func (k *KubernetesTool) QueryInitialState(ctx context.Context, _ QueryInitialStateArgs) (*mcp_golang.ToolResponse, error) {
	fmt.Println("query_initial_state")

	// Get all API resources
	apiGroupResources, err := restmapper.GetAPIGroupResources(k.discoveryClient)
	if err != nil {
		return nil, err
	}

	// Flatten into REST mappings
	// var mapper = restmapper.NewDiscoveryRESTMapper(apiGroupResources)

	// var wg sync.WaitGroup
	// var printChan = make(chan string)
	// var workerChan = make(chan string, 10)

	// go func() {
	// 	for _, line := range <-printChan {
	// 		fmt.Println(line)
	// 	}
	// }()

	var entities = map[string][]string{}

	// Iterate over all mappable resources
	for _, list := range apiGroupResources {
		for _, resources := range list.VersionedResources {
			for _, resource := range resources {
				gvr := schema.GroupVersionResource{
					Group:    list.Group.Name,
					Version:  list.Group.PreferredVersion.Version,
					Resource: resource.Name,
				}

				// Skip subresources like "pods/status"
				if strings.Contains(resource.Name, "/") {
					continue
				}

				// wg.Add(1)
				// go func() {
				// 	defer func() {
				// 		wg.Done()
				// 		<-workerChan
				// 	}()

				// 	workerChan <- ""

				// List the resources
				res, err := k.dynClient.Resource(gvr).Namespace("").List(ctx, metav1.ListOptions{})
				if err != nil {
					// TODO: handle this

					continue // often RBAC or unsupported resource
				}

				var subEntities []string
				if len(res.Items) == 0 {
					continue
				}

				fmt.Printf("%d %s found:\n", len(res.Items), gvr.Resource)
				// printChan <- fmt.Sprintf("%d %s found:\n", len(res.Items), gvr.Resource)
				for _, item := range res.Items {
					fmt.Printf("  - %s/%s\n", item.GetNamespace(), item.GetName())
					subEntities = append(subEntities, fmt.Sprintf("%s/%s", item.GetNamespace(), item.GetName()))
				}

				entities[gvr.Resource] = subEntities
				// }()
			}
		}
	}

	// wg.Wait()
	// close(printChan)

	b, err := json.Marshal(entities)
	if err != nil {
		return nil, err
	}

	return &mcp_golang.ToolResponse{
		Content: []*mcp_golang.Content{
			{
				Type: mcp_golang.ContentTypeText,
				TextContent: &mcp_golang.TextContent{
					Text: string(b),
				},
			},
		},
	}, nil
}
