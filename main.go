package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"githup.com/makromusicCase/makromusic/config"
	"githup.com/makromusicCase/makromusic/consumers"
	"githup.com/makromusicCase/makromusic/database"
	"githup.com/makromusicCase/makromusic/entities"
	makromusic_proto "githup.com/makromusicCase/makromusic/proto-output/makromusic/proto"
	"githup.com/makromusicCase/makromusic/services"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	loadEnv()
	loadDatabase()
	db := database.Database
	redisClient := config.NewRedisClient()

	visionServiceImpl := services.NewVisionService(db, redisClient)
	fileConsumer := consumers.NewFileUploadConsumer(visionServiceImpl)
	go fileConsumer.FileUploadConsumer()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9000))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fileServiceImpl := services.NewFileService(db, redisClient)
	file := makromusic_proto.NewFileServer(fileServiceImpl)
	s := grpc.NewServer()

	makromusic_proto.RegisterFileServiceServer(s, file)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

}

func loadDatabase() {
	database.Connect()
	database.Database.AutoMigrate(&entities.FileEntity{}, &entities.FileFeedValueEntity{})
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
