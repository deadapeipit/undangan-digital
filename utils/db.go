package utils

import (
	"context"
	"fmt"
	"log"
	"time"
	"undangan-digital/config" // Import config untuk mengambil Mongo URI

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Fungsi untuk menghubungkan ke MongoDB
func ConnectDB() (*mongo.Client, error) {
	// Memuat konfigurasi MongoDB dari file config.yaml
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error memuat konfigurasi MongoDB: ", err)
		return nil, err
	}

	// Menggunakan URI yang diambil dari config.yaml
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.Mongo.URI))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Test the connection with Ping
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("MongoDB connection failed: %v", err)
	}

	log.Println("Connected to MongoDB successfully!")
	return client, nil
}

// Fungsi untuk menutup koneksi ke MongoDB
func CloseDB(client *mongo.Client) {
	err := client.Disconnect(context.Background())
	if err != nil {
		log.Fatal("Error saat menutup koneksi MongoDB: ", err)
	}
}
