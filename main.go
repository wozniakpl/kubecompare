package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type KubectlInterface interface {
	getRolloutHistory(resourceType, resourceName string) (string, error)
	getRolloutHistoryWithRevision(resourceType, resourceName string, revision int) (string, error)
}

type RealKubectl struct{}

type OutputWriter interface {
	Write(string)
}

type StdoutWriter struct{}

func (sw StdoutWriter) Write(output string) {
	fmt.Println(output)
}

func (r RealKubectl) getRolloutHistory(resourceType, resourceName string) (string, error) {
	return r.runKubectlCommand("rollout", "history", fmt.Sprintf("%s/%s", resourceType, resourceName))
}

func (r RealKubectl) getRolloutHistoryWithRevision(resourceType, resourceName string, revision int) (string, error) {
	return r.runKubectlCommand("rollout", "history", fmt.Sprintf("%s/%s", resourceType, resourceName), fmt.Sprintf("--revision=%d", revision))
}

func (r RealKubectl) runKubectlCommand(args ...string) (string, error) {
	cmd := exec.Command("kubectl", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

func writeTempFile(data string) (string, error) {
	tempFile, err := os.CreateTemp("", "revision")
	if err != nil {
		return "", err
	}
	err = os.WriteFile(tempFile.Name(), []byte(data), 0644)
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

func usage() string {
	return "Usage: kubecompare [ <resource-type> <resource-name> | <resource-type>/<resource-name> ] [ <previous-revision> <next-revision> ]"
}

func mainLogic(k KubectlInterface, writer OutputWriter, args []string) (int, error) {
	if len(args) == 0 {
		writer.Write(usage())
		return 0, nil
	}

	var resourceType, resourceName, previousRevisionArg, nextRevisionArg string

	if strings.Contains(args[0], "/") {
		parts := strings.SplitN(args[0], "/", 2)
		resourceType, resourceName = parts[0], parts[1]
		args = args[1:]
	} else {
		resourceType, resourceName = args[0], args[1]
		args = args[2:]
	}

	if len(args) == 0 {
		history, err := k.getRolloutHistory(resourceType, resourceName)
		if err != nil {
			return 1, err
		}
		writer.Write(history)
		return 0, nil
	} else if len(args) != 2 {
		return 1, fmt.Errorf("wrong invocation")
	}

	if len(args) == 0 {
		history, err := k.getRolloutHistory(resourceType, resourceName)
		if err != nil {
			return 1, err
		}
		writer.Write(history)
		return 0, nil
	} else if len(args) == 2 {
		previousRevisionArg, nextRevisionArg = args[0], args[1]
	} else {
		writer.Write(usage())
		return 1, fmt.Errorf("invalid number of arguments")
	}

	previousRevision, err := strconv.Atoi(previousRevisionArg)
	if err != nil {
		return 1, err
	}

	nextRevision, err := strconv.Atoi(nextRevisionArg)
	if err != nil {
		return 1, err
	}

	previousData, err := k.getRolloutHistoryWithRevision(resourceType, resourceName, previousRevision)
	if err != nil {
		return 1, err
	}

	nextData, err := k.getRolloutHistoryWithRevision(resourceType, resourceName, nextRevision)
	if err != nil {
		return 1, err
	}

	previousFile, err := writeTempFile(previousData)
	if err != nil {
		return 1, err
	}
	defer os.Remove(previousFile)

	nextFile, err := writeTempFile(nextData)
	if err != nil {
		return 1, err
	}
	defer os.Remove(nextFile)

	diff, err := getDiff(previousFile, nextFile)
	if err != nil {
		return 1, err
	}

	writer.Write(diff)

	return 0, nil
}

func parseFlags() (bool, []string) {
	helpFlag := flag.Bool("h", false, "Show usage information")
	flag.Parse()
	return *helpFlag, flag.Args()
}

func main() {
	helpFlag, args := parseFlags()
	if helpFlag {
		fmt.Println(usage())
		os.Exit(0)
	}

	kubectl := RealKubectl{}
	writer := StdoutWriter{}
	status, err := mainLogic(kubectl, writer, args)

	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println(usage())
	}
	os.Exit(status)
}
