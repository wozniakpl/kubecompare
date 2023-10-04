package main

import (
	"errors"
	"testing"
)

func TestShowingTheErrorWhenKubectlFails(t *testing.T) {
	mockKubectl := new(MockKubectl)
	mockKubectl.On("getRolloutHistory", "deployment", "nginx-that-does-not-exist", "").Return("", errors.New("error"))

	mockWriter := new(MockWriter)
	namespace := ""

	_, err := mainLogic(mockKubectl, mockWriter, namespace, []string{"deployment/nginx-that-does-not-exist"})

	if err == nil {
		t.Errorf("Expected error")
	}
}

func TestFailingWhenResourceNameIsNotSpecified(t *testing.T) {
	mockKubectl := new(MockKubectl)
	mockWriter := new(MockWriter)
	namespace := ""

	_, err := mainLogic(mockKubectl, mockWriter, namespace, []string{"deployment"})

	if err == nil {
		t.Errorf("Expected error")
	}
}

func TestFailingIfRevisionDoesNotExist(t *testing.T) {
	mockKubectl := new(MockKubectl)
	mockKubectl.On("getRolloutHistoryWithRevision", "deployment", "nginx", 1, "").Return("", errors.New("error"))

	mockWriter := new(MockWriter)
	namespace := ""

	_, err := mainLogic(mockKubectl, mockWriter, namespace, []string{"deployment", "nginx", "1", "2"})

	if err == nil {
		t.Errorf("Expected error")
	}
}
