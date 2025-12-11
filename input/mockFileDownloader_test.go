package input

import (
	"github.com/stretchr/testify/mock"
)

// MockFileDownloader ...
type MockFileDownloader struct {
	mock.Mock
}

// Get ...
func (m *MockFileDownloader) Get(destination, source string) error {
	args := m.Called(destination, source)
	return args.Error(0)
}

// GetContents ...
func (m *MockFileDownloader) GetRemoteContents(source string) ([]byte, error) {
	args := m.Called(source)
	arg0, ok := args.Get(0).([]byte)
	if !ok {
		panic("unexpected type")
	}

	return arg0, args.Error(1)
}

// ReadLocalFile ...
func (m *MockFileDownloader) ReadLocalFile(path string) ([]byte, error) {
	args := m.Called(path)
	arg0, ok := args.Get(0).([]byte)
	if !ok {
		panic("unexpected type")
	}

	return arg0, args.Error(1)
}

// GivenGetFails ...
func (m *MockFileDownloader) GivenGetFails(reason error) *MockFileDownloader {
	m.On("Get", mock.Anything, mock.Anything).Return(reason)
	return m
}

// GivenGetSucceed ...
func (m *MockFileDownloader) GivenGetSucceed() *MockFileDownloader {
	m.On("Get", mock.Anything, mock.Anything).Return(nil)
	return m
}
