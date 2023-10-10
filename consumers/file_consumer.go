package consumers

import (
	"github.com/IBM/sarama"
	"githup.com/makromusicCase/makromusic/services"
	"log"
	"os"
	"os/signal"
	"sync"
)

type FileUploadConsumerReceiver struct {
	visionService services.VisionService
}

func NewFileUploadConsumer(service services.VisionService) FileUploadConsumerReceiver {
	return FileUploadConsumerReceiver{visionService: service}
}

func (r *FileUploadConsumerReceiver) FileStoreConsumer() {
	broker := os.Getenv("KAFKA_BROKER")
	topic := os.Getenv("KAFKA_FILE_TOPIC")

	processMessage := func(msg string) error {
		return r.visionService.CreateDetectedFaces(os.Stdout, msg)
	}

	r.consumeKafkaMessages(broker, topic, processMessage)
}

func (r *FileUploadConsumerReceiver) FileUpdateConsumer() {
	broker := os.Getenv("KAFKA_BROKER")
	topic := os.Getenv("UPDATE_IMAGE_DETAIL_TOPIC")

	processMessage := func(msg string) error {
		return r.visionService.UpdateDetectedFaces(os.Stdout, msg)
	}

	r.consumeKafkaMessages(broker, topic, processMessage)
}

func (r *FileUploadConsumerReceiver) consumeKafkaMessages(broker, topic string, processMessage func(string) error) {
	// Kafka broker adresleri
	brokers := []string{broker}

	// Kafka consumer yapılandırması
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Kafka consumer'ı oluştur
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %v", err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalf("Error closing Kafka consumer: %v", err)
		}
	}()

	// Kafka topic üzerindeki tüm partition'ları al
	partitions, err := consumer.Partitions(topic)
	if err != nil {
		log.Fatalf("Error getting partitions: %v", err)
	}

	var wg sync.WaitGroup

	for _, partition := range partitions {
		wg.Add(1)
		go func(partition int32) {
			defer wg.Done()

			partitionConsumer, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
			if err != nil {
				log.Printf("Error creating partition consumer: %v", err)
				return
			}
			defer func() {
				if err := partitionConsumer.Close(); err != nil {
					log.Printf("Error closing partition consumer: %v", err)
				}
			}()

			// Sinyal yakala ve işlemi sonlandır
			signals := make(chan os.Signal, 1)
			signal.Notify(signals, os.Interrupt)

			for {
				select {
				case msg := <-partitionConsumer.Messages():
					log.Printf(
						"Received message from partition %d at offset %d: %s\n", partition, "Veri alındı",
						string(msg.Value),
					)

					// işlem için fonksiyon çağır
					err := processMessage(string(msg.Value))
					if err != nil {
						log.Printf("Error processing message: %v", err)
					}

				case err := <-partitionConsumer.Errors():
					log.Printf("Error: %v\n", err)
				case <-signals:
					log.Printf("Error processing message: %v", err)
				}
			}
		}(partition)
	}

	wg.Wait()
}
