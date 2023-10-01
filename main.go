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

func showUsage() {
	fmt.Println("Usage: kubecompare [<resource-type> <resource-name> | <resource-type>/<resource-name>] <previous-revision> <next-revision>")
}

func handleError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		showUsage()
		os.Exit(1)
	}
}

func main() {
	args := os.Args[1:]

	if len(args) < 3 {
		showUsage()
		os.Exit(1)
	}

	var resourceType, resourceName string
	var previousRevisionArg, nextRevisionArg string

	if strings.Contains(args[0], "/") {
		parts := strings.SplitN(args[0], "/", 2)
		resourceType, resourceName = parts[0], parts[1]
		previousRevisionArg, nextRevisionArg = args[1], args[2]
	} else if len(args) >= 4 {
		resourceType, resourceName = args[0], args[1]
		previousRevisionArg, nextRevisionArg = args[2], args[3]
	} else {
		showUsage()
		os.Exit(1)
	}

	previousRevision, err := strconv.Atoi(previousRevisionArg)
	handleError(err)

	nextRevision, err := strconv.Atoi(nextRevisionArg)
	handleError(err)

	previousData, err := getKubectlRolloutHistory(resourceType, resourceName, previousRevision)
	handleError(err)

	nextData, err := getKubectlRolloutHistory(resourceType, resourceName, nextRevision)
	handleError(err)

	previousFile, err := writeTempFile(previousData)
	handleError(err)
	defer os.Remove(previousFile)

	nextFile, err := writeTempFile(nextData)
	handleError(err)
	defer os.Remove(nextFile)

	diff, err := getDiff(previousFile, nextFile)
	handleError(err)
	
	fmt.Println(diff)
}
