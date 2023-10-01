package main

import (
	"bytes"
	"flag"
	"fmt"
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
	var outputFormat string
	flag.StringVar(&outputFormat, "o", "yaml", "Output format. Allowed values: json, yaml")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		fmt.Println("Usage: kubecompare <resource-type> <resource-name>")
		return
	}

	resourceType, resourceName := args[0], args[1]
	if strings.Contains(resourceName, "/") {
		parts := strings.SplitN(resourceName, "/", 2)
		resourceType, resourceName = parts[0], parts[1]
	}

	data, err := getKubectlOutput(resourceType, resourceName, outputFormat)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(data)
}
