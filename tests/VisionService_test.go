package services

import (
	"github.com/stretchr/testify/assert"
	services "githup.com/makromusicCase/makromusic/mocks"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"os"
	"testing"
)

func TestVisionService_CreateDetectedFaces_Error(t *testing.T) {
	mock := &services.MockVisionService{
		CreateDetectedFacesFunc: func(
			w io.Writer, file string,
		) error {
			return status.Error(codes.NotFound, "Image not found")
		},
	}

	// error
	err := mock.CreateDetectedFaces(os.Stdout, "test")
	assert.Equal(t, codes.NotFound, status.Code(err))
}

func TestVisionService_CreateDetectedFaces_Success(t *testing.T) {
	mock := &services.MockVisionService{
		CreateDetectedFacesFunc: func(
			w io.Writer, file string,
		) error {
			return nil
		},
	}

	// success
	err := mock.CreateDetectedFaces(os.Stdout, "test")
	assert.NoError(t, err)
}

func TestVisionService_UpdateDetectedFaces_Error(t *testing.T) {
	mock := &services.MockVisionService{
		UpdateDetectedFacesFunc: func(
			w io.Writer, file string,
		) error {
			return status.Error(codes.NotFound, "Image not found")
		},
	}

	// error
	err := mock.UpdateDetectedFaces(os.Stdout, "test")
	assert.Equal(t, codes.NotFound, status.Code(err))
}

func TestVisionService_UpdateDetectedFaces_Success(t *testing.T) {
	mock := &services.MockVisionService{
		UpdateDetectedFacesFunc: func(
			w io.Writer, file string,
		) error {
			return nil
		},
	}

	// success
	err := mock.UpdateDetectedFaces(os.Stdout, "test")
	assert.NoError(t, err)
}
