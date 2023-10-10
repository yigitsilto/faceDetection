package services

import (
	"cloud.google.com/go/vision/v2/apiv1/visionpb"
	"errors"
	"github.com/stretchr/testify/assert"
	"githup.com/makromusicCase/makromusic/entities"
	services "githup.com/makromusicCase/makromusic/mocks"
	makromusic_proto "githup.com/makromusicCase/makromusic/proto-output/makromusic/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestFileServiceImpl_GetImageDetail(t *testing.T) {
	mock := &services.MockFileService{
		GetImageDetailFunc: func(
			ctx context.Context, request *makromusic_proto.GetImageDetailRequest,
		) (*makromusic_proto.ImageDetailResponse, error) {
			if request.GetId() == 1 {
				return &makromusic_proto.ImageDetailResponse{Face: "mock_data"}, nil
			}
			return nil, status.Error(codes.NotFound, "Image not found")
		},
	}

	ctx := context.TODO()

	// success
	request1 := &makromusic_proto.GetImageDetailRequest{Id: 1}
	response1, err1 := mock.GetImageDetail(ctx, request1)
	assert.NoError(t, err1)
	assert.NotNil(t, response1)
	assert.Equal(t, "mock_data", response1.Face)

	// error
	request2 := &makromusic_proto.GetImageDetailRequest{Id: 2}
	response2, err2 := mock.GetImageDetail(ctx, request2)
	assert.Error(t, err2)
	assert.Nil(t, response2)
	assert.Equal(t, codes.NotFound, status.Code(err2))
}

func TestFileServiceImpl_GetImageFeed(t *testing.T) {
	mock := &services.MockFileService{
		GetImageFeedFunc: func(
			ctx context.Context, request *makromusic_proto.GetImageFeedRequest,
		) (*makromusic_proto.GetImageFeedResponse, error) {
			if request.GetPage() == 1 {
				return &makromusic_proto.GetImageFeedResponse{}, nil
			}
			return nil, errors.New("Internal server error")
		},
	}

	ctx := context.TODO()

	// success
	request1 := &makromusic_proto.GetImageFeedRequest{Page: 1}
	response1, err1 := mock.GetImageFeed(ctx, request1)
	assert.NoError(t, err1)
	assert.NotNil(t, response1)

	// error
	request2 := &makromusic_proto.GetImageFeedRequest{Page: 2}
	response2, err2 := mock.GetImageFeed(ctx, request2)
	assert.Error(t, err2)
	assert.Nil(t, response2)
}

func TestFileServiceImpl_UpdateImageDetail(t *testing.T) {
	mock := &services.MockFileService{
		UpdateImageDetailFunc: func(
			fileId uint, annotations []*visionpb.FaceAnnotation,
		) error {
			if fileId == 1 {
				return nil
			}
			return errors.New("Failed to update image detail")
		},
	}

	// success
	fileId1 := uint(1)
	var annotations1 []*visionpb.FaceAnnotation
	err1 := mock.UpdateImageDetail(fileId1, annotations1)
	assert.NoError(t, err1)

	// errors
	fileId2 := uint(2)
	var annotations2 []*visionpb.FaceAnnotation
	err2 := mock.UpdateImageDetail(fileId2, annotations2)
	assert.Error(t, err2)
	assert.Equal(t, "Failed to update image detail", err2.Error())
}

func TestMockFileService_GetFileByPath(t *testing.T) {
	mock := &services.MockFileService{
		GetFileByPathFunc: func(path string) (entities.FileEntity, error) {
			if path == "existing_path" {
				return entities.FileEntity{}, nil
			}
			return entities.FileEntity{}, errors.New("File not found")
		},
	}

	// success
	path1 := "existing_path"
	fileEntity1, err1 := mock.GetFileByPath(path1)
	assert.NoError(t, err1)
	assert.NotNil(t, fileEntity1)

	// error
	path2 := "non_existing_path"
	fileEntity2, err2 := mock.GetFileByPath(path2)
	assert.Error(t, err2)
	assert.Nil(t, fileEntity2)
	assert.Equal(t, "File not found", err2.Error())
}
func TestMockFileService_CreateFeeds(t *testing.T) {
	mock := &services.MockFileService{
		CreateFeedsFunc: func(fileId uint, annotations []*visionpb.FaceAnnotation) error {
			if fileId == 1 {
				return nil
			}
			return errors.New("Failed to create feeds")
		},
	}

	// success
	fileId1 := uint(1)
	var annotations1 []*visionpb.FaceAnnotation
	err1 := mock.CreateFeeds(fileId1, annotations1)
	assert.NoError(t, err1)

	// error
	fileId2 := uint(2)
	var annotations2 []*visionpb.FaceAnnotation
	err2 := mock.CreateFeeds(fileId2, annotations2)
	assert.Error(t, err2)
	assert.Equal(t, "Failed to create feeds", err2.Error())
}

func TestMockFileService_UpdateProducer(t *testing.T) {
	mock := &services.MockFileService{
		UpdateProducerFunc: func(
			ctx context.Context, request *makromusic_proto.UpdateImageDetailRequest,
		) (*makromusic_proto.UpdateImageDetailResponse, error) {
			if request.GetId() == 1 {
				return &makromusic_proto.UpdateImageDetailResponse{Status: true}, nil
			}
			return nil, status.Error(codes.Internal, "Update failed")
		},
	}

	ctx := context.TODO()

	// success
	request1 := &makromusic_proto.UpdateImageDetailRequest{Id: 1}
	response1, err1 := mock.UpdateProducer(ctx, request1)
	assert.NoError(t, err1)
	assert.NotNil(t, response1)
	assert.True(t, response1.Status)

	// error
	request2 := &makromusic_proto.UpdateImageDetailRequest{Id: 2}
	response2, err2 := mock.UpdateProducer(ctx, request2)
	assert.Error(t, err2)
	assert.Nil(t, response2)
	assert.Equal(t, codes.Internal, status.Code(err2))
}
