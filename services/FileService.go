package services

import (
	"github.com/google/uuid"
	"githup.com/makromusicCase/makromusic/config"
	"githup.com/makromusicCase/makromusic/entities"
	"githup.com/makromusicCase/makromusic/producers"
	makromusic_proto "githup.com/makromusicCase/makromusic/proto-output/makromusic/proto"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"strconv"
	"time"
)

type FileService interface {
	UploadFile(server makromusic_proto.FileService_UploadFileServer) error
	GetImageDetail(
		ctx context.Context, request *makromusic_proto.GetImageDetailRequest,
	) (*makromusic_proto.ImageDetailResponse, error)
	GetImageFeed(
		ctx context.Context, request *makromusic_proto.GetImageFeedRequest,
	) (*makromusic_proto.GetImageFeedResponse, error)
	UpdateImageDetail(
		ctx context.Context, request *makromusic_proto.UpdateImageDetailRequest,
	) (*makromusic_proto.UpdateImageDetailResponse, error)
}

type FileServiceImpl struct {
	db    *gorm.DB
	redis *config.RedisClient
}

func NewFileService(
	db *gorm.DB, redis *config.RedisClient,
) FileService {
	return &FileServiceImpl{db: db, redis: redis}
}

func (f *FileServiceImpl) UpdateImageDetail(
	ctx context.Context, request *makromusic_proto.UpdateImageDetailRequest,
) (*makromusic_proto.UpdateImageDetailResponse, error) {

	var fileEntity entities.FileEntity

	err := f.db.Table("file_entities").Where("id =?", request.GetId()).Find(&fileEntity).Error

	if err != nil {
		return nil, err
	}

	f.sendToKafka(fileEntity.Path)

	return &makromusic_proto.UpdateImageDetailResponse{Status: true}, nil

}

func (f *FileServiceImpl) GetImageFeed(
	ctx context.Context, request *makromusic_proto.GetImageFeedRequest,
) (*makromusic_proto.GetImageFeedResponse, error) {

	var fileEntities []entities.FileEntity

	pageNumber := int(request.Page)
	pageSize := 10

	offset := (pageNumber - 1) * pageSize

	err := f.db.
		Table("file_entities").
		Preload("Values").
		Order(
			`
        created_at DESC
    `,
		).
		Offset(offset).
		Limit(pageSize).
		Find(&fileEntities).Error

	if err != nil {
		return nil, err
	}
	response := makromusic_proto.GetImageFeedResponse{}
	returnValue := f.buildResponseForImageFeed(&fileEntities, &response)

	return returnValue, nil

}

func (f *FileServiceImpl) buildResponseForImageFeed(
	fileEntities *[]entities.FileEntity, response *makromusic_proto.GetImageFeedResponse,
) *makromusic_proto.GetImageFeedResponse {
	for _, entity := range *fileEntities {
		imageItem := &makromusic_proto.ImageItem{
			CreatedAt: entity.CreatedAt.Format(time.RFC3339), // String olarak tarih ve saat
			UpdatedAt: entity.UpdatedAt.Format(time.RFC3339),
			Path:      entity.Path,
			Id:        int64(entity.Id),
		}

		// Values dizisini döngü içinde geziyoruz ve her bir öğesini ekliyoruz
		for _, valueEntity := range entity.Values {
			value := &makromusic_proto.ImageValues{
				Key:   valueEntity.Key,
				Value: valueEntity.Value,
			}
			imageItem.Values = append(imageItem.Values, value)
		}

		response.ImageItems = append(response.ImageItems, imageItem)
	}

	return response
}

func (f *FileServiceImpl) GetImageDetail(
	ctx context.Context, request *makromusic_proto.GetImageDetailRequest,
) (*makromusic_proto.ImageDetailResponse, error) {

	redisKey := "feeds-" + strconv.FormatInt(request.GetId(), 10)
	redisValue, err := f.redis.Get(redisKey)
	if err != nil {
		return nil, err
	}

	return &makromusic_proto.ImageDetailResponse{
		Face: redisValue,
	}, nil

}

func (f *FileServiceImpl) UploadFile(stream makromusic_proto.FileService_UploadFileServer) error {
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&makromusic_proto.FileUploadResponse{Status: true})
		}
		if err != nil {
			return err
		}

		err = f.saveToDatabase(chunk.GetByte())

		if err != nil {
			return err
		}

	}
	return nil
}

func (f *FileServiceImpl) saveToDatabase(byte []byte) error {
	filePrefix := uuid.New().String()
	filePath := "./public/" + filePrefix + ".jpg"

	fileEntity := entities.FileEntity{Path: filePath}

	result := f.db.Create(&fileEntity)

	if result.Error != nil {
		return result.Error
	}

	// diske kaydet
	err := f.saveToDisc(byte, filePath)

	return err
}

func (f *FileServiceImpl) saveToDisc(byte []byte, filePath string) error {

	// Veriyi diske kaydet
	err := ioutil.WriteFile(filePath, byte, 0644)
	if err != nil {
		return err
	}

	// send to producer
	f.sendToKafka(filePath)
	return nil
}

func (f *FileServiceImpl) sendToKafka(filePath string) {
	producers.FileUploadProducer(filePath)
}
