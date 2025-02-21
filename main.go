package main

import (
	"log"
	"undangan-digital/routes"
	"undangan-digital/utils"

	"github.com/gin-gonic/gin"
)

func main() {

	// Menghubungkan ke MongoDB menggunakan fungsi ConnectDB dari utils
	client, err := utils.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDB(client) // Pastikan untuk menutup koneksi saat aplikasi selesai

	// Inisialisasi Gin
	r := gin.Default()

	// Menetapkan folder templates untuk Gin
	r.LoadHTMLGlob("templates/*")

	// Setup Routes
	routes.SetupRoutes(r, client)

	// Menjalankan server
	r.Run("0.0.0.0:8080")
}
