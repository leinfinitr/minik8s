package cmd

import (
	"fmt"
	"os"
	"minik8s/pkg/config"
	"minik8s/internal/pkg/httpRequest"
	"minik8s/pkg/apiObject"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create pod [name] --image [image] --namespace [namespace]",
	Short: "Create a pod by command line",
	Long:  `Create a pod by specifying the name, image, and optional namespace directly on the command line.`,
	Run:   createHandler,
}

// CreatePod function creates a Pod resource with given parameters.
func CreatePod(name string, image string, namespace string) {
	fmt.Printf("Creating pod %s with image %s in namespace %s...\n", name, image, namespace)
	// Construct the Pod object
	pod := apiObject.Pod{
		TypeMeta: apiObject.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: apiObject.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: apiObject.PodSpec{
			Containers: []apiObject.Container{
				{
					Name:  name,
					Image: image,
				},
			},
			RestartPolicy: "Always",
		},
	}
	
	url := config.APIServerUrl() + "/api/v1/namespaces/" + namespace + "/pods"
	resp, err := httprequest.PostObjMsg(url, pod)
	if err != nil {
		fmt.Println("Error: Could not create the pod.")
		os.Exit(1)
	}
	ResultDisplay("Pod", resp)
}

// createHandler interprets the command line arguments and invokes the CreatePod function.
func createHandler(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("Error: You must specify a name for the pod.")
		os.Exit(1)
	}
	
	name := args[0]
	image, _ := cmd.Flags().GetString("image")
	if image == "" {
		fmt.Println("Error: You must specify an image for the pod to use.")
		os.Exit(1)
	}
	namespace, _ := cmd.Flags().GetString("namespace")
	if namespace == "" {
		namespace = "default" // Set a default namespace if not specified
	}
	
	CreatePod(name, image, namespace)
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().String("image", "", "Specify the image of the pod")
	createCmd.Flags().String("namespace", "", "Specify the namespace of the pod")
}
