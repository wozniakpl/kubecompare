package main

import (
	"github.com/stretchr/testify/mock"
	"strings"
	"testing"
)


func TestMainLogicNoArgs(t *testing.T) {
	mockKubectl := new(MockKubectl)
	mockWriter := new(MockWriter)
	mockWriter.On("Write", mock.Anything).Return(nil)
	_, err := mainLogic(mockKubectl, mockWriter, []string{})

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
	mockKubectl.On("getRolloutHistory", "deployment", "nginx", 1).Return("some output 1", nil)
	mockKubectl.On("getRolloutHistory", "deployment", "nginx", 2).Return("some output 2", nil)

	mockWriter := new(MockWriter)
	mockWriter.On("Write", mock.Anything).Return(nil)

	_, err := mainLogic(mockKubectl, mockWriter, []string{"deployment", "nginx", "1", "2"})

	if err != nil {
		t.Errorf("Expected no error")
	}

	assertThereIsSomeDiff(t, mockWriter)
}
