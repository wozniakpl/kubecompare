package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func getKubectlRolloutHistory(resourceType, resourceName string, revision int) (string, error) {
	cmd := exec.Command("kubectl", "rollout", "history", fmt.Sprintf("%s/%s", resourceType, resourceName), fmt.Sprintf("--revision=%d", revision))
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return out.String(), nil
}

func writeTempFile(data string) (string, error) {
	tempFile, err := ioutil.TempFile("", "revision")
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(tempFile.Name(), []byte(data), 0644)
	if err != nil {
		return "", err
	}
	return tempFile.Name(), nil
}

func getDiff(file1, file2 string) (string, error) {
	cmd := exec.Command("diff", "-u", file1, file2)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil && err.Error() != "exit status 1" {
		return "", err
	}
	return out.String(), nil
}

func main() {
	args := os.Args[1:]

	if len(args) < 4 {
		fmt.Println("Usage: kubecompare <resource-type> <resource-name> or <resource-type>/<resource-name> <previous-revision> <next-revision>")
		return
	}

	var resourceType, resourceName string
	if strings.Contains(args[0], "/") {
		parts := strings.SplitN(args[0], "/", 2)
		resourceType, resourceName = parts[0], parts[1]
	} else {
		resourceType, resourceName = args[0], args[1]
	}

	previousRevisionArg, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println("Error: Invalid previous revision number")
		return
	}

	nextRevisionArg, err := strconv.Atoi(args[3])
	if err != nil {
		fmt.Println("Error: Invalid next revision number")
		return
	}

	previousData, err := getKubectlRolloutHistory(resourceType, resourceName, previousRevisionArg)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	nextData, err := getKubectlRolloutHistory(resourceType, resourceName, nextRevisionArg)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	previousFile, err := writeTempFile(previousData)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer os.Remove(previousFile)

	nextFile, err := writeTempFile(nextData)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer os.Remove(nextFile)

	diff, err := getDiff(previousFile, nextFile)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(diff)
}
