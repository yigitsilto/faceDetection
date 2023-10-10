package services

import (
	"githup.com/makromusicCase/makromusic/services"
	"io"
)

type MockVisionService struct {
	CreateDetectedFacesFunc func(w io.Writer, file string) error
	UpdateDetectedFacesFunc func(w io.Writer, file string) error
}

func NewMockVisionService() services.VisionService {
	return &MockVisionService{}
}

func (m *MockVisionService) CreateDetectedFaces(w io.Writer, file string) error {
	return m.CreateDetectedFacesFunc(w, file)
}

func (m *MockVisionService) UpdateDetectedFaces(w io.Writer, file string) error {
	return m.UpdateDetectedFacesFunc(w, file)
}
