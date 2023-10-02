package main

import (
    "github.com/stretchr/testify/mock"
)

type MockKubectl struct {
    mock.Mock
}

func (m *MockKubectl) getRolloutHistory(resourceType, resourceName string, revision int) (string, error) {
	args := m.Called(resourceType, resourceName, revision)
	return args.String(0), args.Error(1)
}


type MockWriter struct {
	mock.Mock
	output string
}

func (mw *MockWriter) Write(output string) error {
	args := mw.Called(output)
	mw.output = output
	return args.Error(0)
}

func (mw *MockWriter) GetOutput() string {
	return mw.output
}
