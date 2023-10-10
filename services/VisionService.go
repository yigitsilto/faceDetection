package services

import (
	vision "cloud.google.com/go/vision/apiv1"
	"cloud.google.com/go/vision/v2/apiv1/visionpb"
	"context"
	"errors"
	"fmt"
	"githup.com/makromusicCase/makromusic/config"
	"gorm.io/gorm"
	"io"
	"os"
)

type VisionService interface {
	CreateDetectedFaces(w io.Writer, file string) error
	UpdateDetectedFaces(w io.Writer, file string) error
}

type VisionServiceImpl struct {
	db          *gorm.DB
	redis       *config.RedisClient
	fileService FileService
}

func NewVisionService(
	db *gorm.DB, redis *config.RedisClient, fileService FileService,
) VisionService {
	return &VisionServiceImpl{db: db, redis: redis, fileService: fileService}
}

func (r *VisionServiceImpl) generateFaceDetections(w io.Writer, file string) ([]*visionpb.FaceAnnotation, error) {
	ctx := context.Background()
	credentialPath := "./credentials.json"

	err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", credentialPath)
	if err != nil {
		fmt.Println("Çevresel değişken ayarlanamadı:", err)
	}

	gcpCredentialPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if gcpCredentialPath != "" {
		fmt.Println("GOOGLE_APPLICATION_CREDENTIALS çevresel değişkeni ayarlandı:", gcpCredentialPath)
	} else {
		fmt.Println("GOOGLE_APPLICATION_CREDENTIALS çevresel değişkeni ayarlanmadı.")
	}

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	image, err := vision.NewImageFromReader(f)
	if err != nil {
		return nil, err
	}
	annotations, err := client.DetectFaces(ctx, image, nil, 10)
	if err != nil {
		return nil, err
	}

	return annotations, err
}

func (r *VisionServiceImpl) UpdateDetectedFaces(w io.Writer, file string) error {

	annotations, err := r.generateFaceDetections(w, file)
	if err != nil {
		return err
	}
	if len(annotations) == 0 {
		return errors.New("no faces found")
	}

	// file bul ve getir
	path, err := r.fileService.GetFileByPath(file)
	if err != nil || path.Path == "" {
		return err
	}

	// veritabanına feed leri ekle
	err = r.fileService.UpdateImageDetail(path.Id, annotations)
	if err != nil {
		return err
	}

	return nil
}

func (r *VisionServiceImpl) CreateDetectedFaces(w io.Writer, file string) error {

	annotations, err := r.generateFaceDetections(w, file)
	if err != nil {
		return err
	}
	if len(annotations) == 0 {
		return errors.New("no faces found")
	}

	// file bul ve getir
	path, err := r.fileService.GetFileByPath(file)
	if err != nil || path.Path == "" {
		return err
	}

	// veritabanına feed leri ekle
	err = r.fileService.CreateFeeds(path.Id, annotations)
	if err != nil {
		return err
	}

	return nil
}
