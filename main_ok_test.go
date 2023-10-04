package main

import (
	"strings"
	"testing"
)

func TestMainLogicNoArgs(t *testing.T) {
	mockKubectl := new(MockKubectl)
	mockWriter := new(MockWriter)
	namespace := ""
	_, err := mainLogic(mockKubectl, mockWriter, namespace, []string{})

	if err != nil {
		t.Errorf("Expected no error")
	}
}

func assertThereIsSomeDiff(t *testing.T, writer *MockWriter) {
	output := writer.GetOutput()
	lines := strings.Split(output, "\n")
	plusFound := false
	minusFound := false
	for _, line := range lines {
		if strings.HasPrefix(line, "-") {
			minusFound = true
		} else if strings.HasPrefix(line, "+") {
			plusFound = true
		}
	}
	if !(plusFound && minusFound) {
		t.Errorf("Output does not contain both + and - line prefixes")
	}
}

func TestShowingDiffBetweenTwoRevisions(t *testing.T) {
	mockKubectl := new(MockKubectl)
	mockKubectl.On("getRolloutHistoryWithRevision", "deployment", "nginx", 1, "").Return("some output 1", nil)
	mockKubectl.On("getRolloutHistoryWithRevision", "deployment", "nginx", 2, "").Return("some output 2", nil)

	mockWriter := new(MockWriter)
	namespace := ""

	_, err := mainLogic(mockKubectl, mockWriter, namespace, []string{"deployment", "nginx", "1", "2"})

	if err != nil {
		t.Errorf("Expected no error")
	}

	assertThereIsSomeDiff(t, mockWriter)
}

func TestShowingRollbackHistoryWhenNoRevisionIsSpecified(t *testing.T) {
	mockKubectl := new(MockKubectl)
	mockKubectl.On("getRolloutHistory", "deployment", "nginx", "").Return("history", nil)

	mockWriter := new(MockWriter)
	namespace := ""

	_, err := mainLogic(mockKubectl, mockWriter, namespace, []string{"deployment/nginx"})

	if err != nil {
		t.Errorf("Expected no error")
	}

	output := mockWriter.GetOutput()
	if !strings.Contains(output, "history") {
		t.Errorf("Expected output to contain history")
	}
}

func TestShowingRollbackHistoryWhenNoRevisionIsSpecifiedAndNamespaceIs(t *testing.T) {
	mockKubectl := new(MockKubectl)
	mockKubectl.On("getRolloutHistory", "daemonset", "busybox", "default").Return("history", nil)

	mockWriter := new(MockWriter)
	namespace := "default"

	_, err := mainLogic(mockKubectl, mockWriter, namespace, []string{"daemonset/busybox"})

	if err != nil {
		t.Errorf("Expected no error")
	}

	output := mockWriter.GetOutput()
	if !strings.Contains(output, "history") {
		t.Errorf("Expected output to contain history")
	}
}
