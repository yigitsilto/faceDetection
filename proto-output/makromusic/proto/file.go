package makromusic_proto

import "golang.org/x/net/context"

type FileService interface {
	UploadFile(server FileService_UploadFileServer) error
	GetImageDetail(ctx context.Context, request *GetImageDetailRequest) (*ImageDetailResponse, error)
	GetImageFeed(ctx context.Context, request *GetImageFeedRequest) (*GetImageFeedResponse, error)
	UpdateImageDetail(ctx context.Context, request *UpdateImageDetailRequest) (*UpdateImageDetailResponse, error)
}

type FileServer struct {
	fileService FileService
}

func NewFileServer(fileService FileService) FileServiceServer {
	return &FileServer{fileService: fileService}
}

func (f *FileServer) UpdateImageDetail(
	ctx context.Context, request *UpdateImageDetailRequest,
) (*UpdateImageDetailResponse, error) {
	return f.fileService.UpdateImageDetail(ctx, request)
}

func (f *FileServer) UploadFile(server FileService_UploadFileServer) error {
	return f.fileService.UploadFile(server)
}

func (f *FileServer) GetImageDetail(ctx context.Context, request *GetImageDetailRequest) (*ImageDetailResponse, error) {
	return f.fileService.GetImageDetail(ctx, request)
}

func (f *FileServer) GetImageFeed(ctx context.Context, request *GetImageFeedRequest) (*GetImageFeedResponse, error) {
	return f.fileService.GetImageFeed(ctx, request)
}

func (f *FileServer) mustEmbedUnimplementedFileServiceServer() {
	//TODO implement me
	panic("implement me")
}
