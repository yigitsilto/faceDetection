package producers

import (
	"github.com/IBM/sarama"
	"log"
	"os"
)

func FileUploadProducer(kafkaMessage string) {
	broker := os.Getenv("KAFKA_BROKER")
	topic := os.Getenv("KAFKA_FILE_TOPIC")

	// Kafka broker adresleri
	brokers := []string{broker}

	// Kafka producer yapılandırması
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	// Kafka producer'ı oluştur
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)
	}

	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatalf("Error closing Kafka producer: %v", err)
		}
	}()

	// Kafka'ya gönderilecek mesaj
	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(kafkaMessage),
	}

	// Mesajı Kafka'ya gönder
	partition, offset, err := producer.SendMessage(message)
	if err != nil {
		log.Fatalf("Error sending message: %v", err)
	}

	log.Printf("Sent message to partition %d at offset %d\n", partition, offset)
}
