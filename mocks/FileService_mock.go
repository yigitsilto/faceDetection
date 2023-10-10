package services

import (
	"cloud.google.com/go/vision/v2/apiv1/visionpb"
	"context"
	entities "githup.com/makromusicCase/makromusic/entities"
	makromusic_proto "githup.com/makromusicCase/makromusic/proto-output/makromusic/proto"
)

type MockFileService struct {
	GetImageDetailFunc func(
		context.Context, *makromusic_proto.GetImageDetailRequest,
	) (*makromusic_proto.ImageDetailResponse, error)
	GetImageFeedFunc func(
		context.Context, *makromusic_proto.GetImageFeedRequest,
	) (*makromusic_proto.GetImageFeedResponse, error)
	UpdateImageDetailFunc func(uint, []*visionpb.FaceAnnotation) error
	GetFileByPathFunc     func(string) (entities.FileEntity, error)
	CreateFeedsFunc       func(uint, []*visionpb.FaceAnnotation) error
	UpdateProducerFunc    func(
		context.Context, *makromusic_proto.UpdateImageDetailRequest,
	) (*makromusic_proto.UpdateImageDetailResponse, error)
	UploadFileFunc func(makromusic_proto.FileService_UploadFileServer) error
}

func (m *MockFileService) GetImageDetail(
	ctx context.Context, request *makromusic_proto.GetImageDetailRequest,
) (*makromusic_proto.ImageDetailResponse, error) {
	return m.GetImageDetailFunc(ctx, request)
}

func (m *MockFileService) GetImageFeed(
	ctx context.Context, request *makromusic_proto.GetImageFeedRequest,
) (*makromusic_proto.GetImageFeedResponse, error) {
	return m.GetImageFeedFunc(ctx, request)
}

func (m *MockFileService) UpdateImageDetail(fileId uint, annotations []*visionpb.FaceAnnotation) error {
	return m.UpdateImageDetailFunc(fileId, annotations)
}

func (m *MockFileService) GetFileByPath(path string) (entities.FileEntity, error) {
	return m.GetFileByPathFunc(path)
}

func (m *MockFileService) CreateFeeds(fileId uint, annotations []*visionpb.FaceAnnotation) error {
	return m.CreateFeedsFunc(fileId, annotations)
}

func (m *MockFileService) UpdateProducer(
	ctx context.Context, request *makromusic_proto.UpdateImageDetailRequest,
) (*makromusic_proto.UpdateImageDetailResponse, error) {
	return m.UpdateProducerFunc(ctx, request)
}

func (m *MockFileService) UploadFile(server makromusic_proto.FileService_UploadFileServer) error {
	return m.UploadFileFunc(server)
}
