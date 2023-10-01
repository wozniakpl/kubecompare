package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func getKubectlOutput(resourceType, resourceName, outputFormat string) (string, error) {
	cmd := exec.Command("kubectl", "get", fmt.Sprintf("%s/%s", resourceType, resourceName), "-o", outputFormat)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return out.String(), nil
}

func main() {
	outputFormat := "yaml" // default output format
	args := os.Args[1:]
	for i, arg := range args {
		if arg == "-o" {
			if i+1 < len(args) {
				outputFormat = args[i+1]
				args = append(args[:i], args[i+2:]...) // remove -o and its argument
				break
			}
		}
	}

	var resourceType, resourceName string
	if len(args) == 1 && strings.Contains(args[0], "/") {
		parts := strings.SplitN(args[0], "/", 2)
		resourceType, resourceName = parts[0], parts[1]
	} else if len(args) == 2 {
		resourceType, resourceName = args[0], args[1]
	} else {
		fmt.Println("Usage: kubecompare [-o output_format] <resource-type> <resource-name> or <resource-type>/<resource-name>")
		return
	}

	data, err := getKubectlOutput(resourceType, resourceName, outputFormat)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(data)
}
