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

type RealKubectl struct{}
type StdoutWriter struct{}

func (sw StdoutWriter) Write(output string) {
	fmt.Println(output)
}

func (r RealKubectl) getRolloutHistory(resourceType, resourceName, namespace string) (string, error) {
	args := []string{"rollout", "history", fmt.Sprintf("%s/%s", resourceType, resourceName)}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	return r.runKubectlCommand(args...)
}

func (r RealKubectl) getRolloutHistoryWithRevision(resourceType, resourceName string, revision int, namespace string) (string, error) {
	args := []string{"rollout", "history", fmt.Sprintf("%s/%s", resourceType, resourceName), fmt.Sprintf("--revision=%d", revision)}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	return r.runKubectlCommand(args...)
}

func (r RealKubectl) runKubectlCommand(args ...string) (string, error) {
	cmd := exec.Command("kubectl", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

func getDiff(file1, file2 string) (string, error) {
	cmd := exec.Command("diff", "-u", file1, file2)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
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

func usage() string {
	return "Usage: kubecompare [ <resource-type> <resource-name> | <resource-type>/<resource-name> ] [ <previous-revision> <next-revision> ]"
}

func mainLogic(k KubectlInterface, writer OutputWriter, namespace string, args []string) (int, error) {
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
		if len(args) < 2 {
			return 2, fmt.Errorf("invalid number of arguments")
		}
		resourceType, resourceName = args[0], args[1]
		args = args[2:]
	}

	if len(args) == 0 {
		history, err := k.getRolloutHistory(resourceType, resourceName, namespace)
		if err != nil {
			return 1, err
		}
		writer.Write(history)
		return 0, nil
	} else if len(args) == 2 {
		previousRevisionArg, nextRevisionArg = args[0], args[1]
	} else {
		return 2, fmt.Errorf("invalid number of arguments")
	}

	previousRevision, err := strconv.Atoi(previousRevisionArg)
	if err != nil {
		return 2, err
	}

	nextRevision, err := strconv.Atoi(nextRevisionArg)
	if err != nil {
		return 2, err
	}

	previousData, err := k.getRolloutHistoryWithRevision(resourceType, resourceName, previousRevision, namespace)
	if err != nil {
		return 1, err
	}

	nextData, err := k.getRolloutHistoryWithRevision(resourceType, resourceName, nextRevision, namespace)
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

func parseFlags() (bool, string, []string) {
	helpFlag := flag.Bool("h", false, "Show usage information")

	var namespace string
	flag.StringVar(&namespace, "n", "", "Specify namespace")
	flag.StringVar(&namespace, "namespace", "", "Specify namespace")

	flag.Parse()

	args := flag.Args()

	for _, arg := range []string{"-n", "--namespace"} {
		for i, a := range args {
			if a == arg {
				args = append(args[:i], args[i+2:]...)
				break
			}
		}
	}

	return *helpFlag, namespace, args
}

func main() {
	helpFlag, namespace, args := parseFlags()
	if helpFlag {
		fmt.Println(usage())
		os.Exit(0)
	}

	kubectl := RealKubectl{}
	writer := StdoutWriter{}
	status, err := mainLogic(kubectl, writer, namespace, args)

	if err != nil {
		fmt.Println("Error:", err)
		if status == 2 {
			fmt.Println(usage())
		}
	}
	os.Exit(status)
}
