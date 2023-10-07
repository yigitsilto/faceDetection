package services

import (
	vision "cloud.google.com/go/vision/apiv1"
	"cloud.google.com/go/vision/v2/apiv1/visionpb"
	"context"
	"encoding/json"
	"fmt"
	"githup.com/makromusicCase/makromusic/config"
	"githup.com/makromusicCase/makromusic/entities"
	"gorm.io/gorm"
	"io"
	"os"
	"strconv"
)

type VisionService interface {
	DetectFaces(w io.Writer, file string) error
}

type VisionServiceImpl struct {
	db    *gorm.DB
	redis *config.RedisClient
}

func NewVisionService(
	db *gorm.DB, redis *config.RedisClient,
) VisionService {
	return &VisionServiceImpl{db: db, redis: redis}
}

func (r *VisionServiceImpl) DetectFaces(w io.Writer, file string) error {

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
		return err
	}
	defer client.Close()

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	image, err := vision.NewImageFromReader(f)
	if err != nil {
		return err
	}
	annotations, err := client.DetectFaces(ctx, image, nil, 10)
	if err != nil {
		return err
	}
	if len(annotations) == 0 {
		fmt.Fprintln(w, "No faces found.")
	} else {

		// file bul ve getir
		path, err := r.getFileByPath(file)
		if err != nil || path.Path == "" {
			return err
		}

		// veritabanına feed leri ekle
		err = r.createFeeds(path.Id, annotations, w)
		if err != nil {
			return err
		}

	}
	return nil
}

func (r *VisionServiceImpl) getFileByPath(path string) (entities.FileEntity, error) {
	var fileEntity entities.FileEntity
	err := r.db.Table("file_entities").Where("path = ?", path).Find(&fileEntity).Error

	if err != nil {
		return fileEntity, err
	}

	return fileEntity, nil

}

func (r *VisionServiceImpl) createFeeds(fileId uint, annotations []*visionpb.FaceAnnotation, w io.Writer) error {

	tx := r.db.Begin()
	defer tx.Rollback()

	if err := r.db.Where("file_id = ?", fileId).Delete(&entities.FileFeedValueEntity{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, annotation := range annotations {

		keys := []string{
			"Anger", "Joy", "Surprise", "UnderExposed", "Blurred", "Headwear", "Sorrow",
		}

		var err error

		// Iterate through the keys and create valueEntity objects
		for _, key := range keys {
			valueEntity := entities.FileFeedValueEntity{
				Key:    key,
				Value:  r.getFieldLikelihood(annotation, key),
				FileId: fileId,
			}

			err = r.db.Create(&valueEntity).Error
		}

		if err != nil {
			tx.Rollback()
		}

		tx.Commit()
	}

	// başarıyla sql e kaydettikten sonra redise kaydet
	err := r.addFeedJsonToRedis(fileId, annotations)
	if err != nil {
		return err
	}
	return nil
}

func (r *VisionServiceImpl) getFieldLikelihood(annotation *visionpb.FaceAnnotation, key string) string {
	switch key {
	case "Anger":
		return annotation.AngerLikelihood.String()
	case "Joy":
		return annotation.JoyLikelihood.String()
	case "Surprise":
		return annotation.SurpriseLikelihood.String()
	case "UnderExposed":
		return annotation.UnderExposedLikelihood.String()
	case "Blurred":
		return annotation.BlurredLikelihood.String()
	case "Headwear":
		return annotation.HeadwearLikelihood.String()
	case "Sorrow":
		return annotation.SorrowLikelihood.String()
	default:
		return ""
	}
}

func (r *VisionServiceImpl) addFeedJsonToRedis(fileId uint, annotations []*visionpb.FaceAnnotation) error {
	annotationsJSON, err := json.Marshal(annotations)
	if err != nil {
		fmt.Println("JSON dönüşüm hatası:", err)
		return err
	}

	err = r.redis.Set("feeds-"+strconv.Itoa(int(fileId)), string(annotationsJSON))
	if err != nil {
		return err
	}

	return nil
}
