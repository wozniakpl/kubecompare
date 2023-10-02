package main

import (
	"github.com/stretchr/testify/mock"
)

type MockKubectl struct {
	mock.Mock
}

func (m *MockKubectl) getRolloutHistory(resourceType, resourceName string) (string, error) {
	args := m.Called(resourceType, resourceName)
	return args.String(0), args.Error(1)
}

func (m *MockKubectl) getRolloutHistoryWithRevision(resourceType, resourceName string, revision int) (string, error) {
	args := m.Called(resourceType, resourceName, revision)
	return args.String(0), args.Error(1)
}

type MockWriter struct {
	mock.Mock
	output string
}

func (mw *MockWriter) Write(output string) {
	mw.output = output
}

func (mw *MockWriter) GetOutput() string {
	return mw.output
}
