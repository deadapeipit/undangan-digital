package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"undangan-digital/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Render invitation HTML
func RenderInvitationHTML(c *gin.Context, dbClient *mongo.Client) {
	undanganID := c.Param("id")

	// Cari undangan berdasarkan ID
	undanganCollection := dbClient.Database("undangan_db").Collection("undangan")
	var undangan models.Undangan
	err := undanganCollection.FindOne(context.Background(), bson.M{"id": undanganID}).Decode(&undangan)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Undangan tidak ditemukan"})
		return
	}

	// Cari acara yang terkait dengan acara_id di undangan
	acaraCollection := dbClient.Database("undangan_db").Collection("acara")
	var acara models.Acara
	err = acaraCollection.FindOne(context.Background(), bson.M{"id": undangan.AcaraID}).Decode(&acara)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Acara terkait tidak ditemukan"})
		return
	}

	// Acara.TanggalAcara sudah berupa tipe waktu (time.Time), jadi langsung ambil informasi yang dibutuhkan
	tanggal := acara.TanggalAcara.Day()
	bulan := acara.TanggalAcara.Month().String()
	tahun := acara.TanggalAcara.Year()
	jam := acara.TanggalAcara.Hour()
	menit := acara.TanggalAcara.Minute()

	// Format jam dan menit agar selalu dua digit
	jamStr := fmt.Sprintf("%02d", jam)
	menitStr := fmt.Sprintf("%02d", menit)

	// Render template undangan dan tampilkan data acara
	c.HTML(http.StatusOK, acara.Template, gin.H{
		"NamaAcara":          acara.NamaAcara,
		"TempatAcara":        acara.TempatAcara,
		"LokasiPinGoogleMap": acara.LokasiPinGoogleMap,
		"GoogleMapsEmbed":    acara.GoogleMapsEmbed,
		"Tanggal":            fmt.Sprintf("%02d", tanggal),
		"Bulan":              bulan,
		"Tahun":              tahun,
		"Jam":                jamStr,
		"Menit":              menitStr,
		"Foto":               acara.Foto,
		"ID":                 undangan.ID,
		"NamaUndangan":       undangan.Nama,
		"JumlahOrang":        undangan.JumlahOrang,
		"Hadir":              undangan.Hadir,
	})
}

// Render insert invitation HTML
func RenderInsertInvitationHTML(c *gin.Context, dbClient *mongo.Client) {
	acaraID := c.Param("id")

	// Cari acara yang terkait dengan acara_id di undangan
	acaraCollection := dbClient.Database("undangan_db").Collection("acara")
	var acara models.Acara
	err := acaraCollection.FindOne(context.Background(), bson.M{"id": acaraID}).Decode(&acara)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Acara terkait tidak ditemukan"})
		return
	}

	// Acara.TanggalAcara sudah berupa tipe waktu (time.Time), jadi langsung ambil informasi yang dibutuhkan
	tanggal := acara.TanggalAcara.Day()
	bulan := acara.TanggalAcara.Month().String()
	tahun := acara.TanggalAcara.Year()
	jam := acara.TanggalAcara.Hour()
	menit := acara.TanggalAcara.Minute()

	// Format jam dan menit agar selalu dua digit
	jamStr := fmt.Sprintf("%02d", jam)
	menitStr := fmt.Sprintf("%02d", menit)

	// Render template undangan dan tampilkan data acara
	c.HTML(http.StatusOK, "insertinvitation.html", gin.H{
		"IDAcara":            acaraID,
		"NamaAcara":          acara.NamaAcara,
		"TempatAcara":        acara.TempatAcara,
		"Tanggal":            fmt.Sprintf("%02d", tanggal),
		"Bulan":              bulan,
		"Tahun":              tahun,
		"Jam":                jamStr,
		"Menit":              menitStr,
		"Foto":               acara.Foto,
		"LokasiPinGoogleMap": acara.LokasiPinGoogleMap,
		"GoogleMapsEmbed":    acara.GoogleMapsEmbed,
	})
}

// Endpoint to get all RSVP details including wishes
func GetAllRSVPs(c *gin.Context, dbClient *mongo.Client) {
	// Create a context with timeout or use one from the request (this is an example of best practice)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := dbClient.Database("undangan_db").Collection("undangan")
	// Fetch all undangan (invitations) from the collection
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch invitations"})
		return
	}
	defer cursor.Close(context.Background())

	var filteredInvitations []models.Undangan
	for cursor.Next(ctx) {
		var undangan models.Undangan
		if err := cursor.Decode(&undangan); err != nil {
			log.Println("Error decoding invitation:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode invitation"})
			return
		}

		// Filter out invitations with empty wishes while looping
		if strings.TrimSpace(undangan.Wishes) != "" {
			filteredInvitations = append(filteredInvitations, undangan)
		}
	}

	// Check for any errors after the loop
	if err := cursor.Err(); err != nil {
		log.Println("Error iterating cursor:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to iterate through invitations"})
		return
	}

	// Respond with the filtered invitations
	c.JSON(http.StatusOK, filteredInvitations)
}

// Handle RSVP response
func HandleRSVP(c *gin.Context, dbClient *mongo.Client) {
	// Get the invitation ID from URL
	id := c.Param("id")

	type rsvp struct {
		Response string `json:"response" bson:"response"`
		Wishes   string `json:"wishes" bson:"wishes"`
	}

	var undangan rsvp
	if err := c.ShouldBindJSON(&undangan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	undanganResponse := undangan.Response
	wishes := undangan.Wishes

	// Get the RSVP response from the form (yes or no)
	response := false

	// Update RSVP status based on the response
	if undanganResponse == "yes" {
		response = true
	} else if undanganResponse == "no" {
		response = false
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid RSVP response"})
		return
	}

	collection := dbClient.Database("undangan_db").Collection("undangan")
	// Update the database
	_, err := collection.UpdateOne(
		context.Background(),
		bson.M{"id": id},
		bson.M{"$set": bson.M{"hadir": response, "wishes": wishes}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update RSVP"})
		return
	}

	// Respond with success message
	c.JSON(http.StatusOK, gin.H{"message": "RSVP updated successfully"})
}

// Menampilkan daftar undangan berdasarkan Acara ID
func GetUndanganByAcaraID(c *gin.Context, dbClient *mongo.Client) {
	acaraID := c.DefaultQuery("acara_id", "")
	if acaraID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID acara tidak ditemukan"})
		return
	}

	undanganCollection := dbClient.Database("undangan_db").Collection("undangan")

	cursor, err := undanganCollection.Find(context.Background(), bson.M{"acara_id": acaraID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data undangan"})
		return
	}
	defer cursor.Close(context.Background())

	var undangans []models.Undangan
	for cursor.Next(context.Background()) {
		var undangan models.Undangan
		if err := cursor.Decode(&undangan); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengdecode data undangan"})
			return
		}
		undangans = append(undangans, undangan)
	}

	c.JSON(http.StatusOK, undangans)
}

// Menambah undangan baru ke MongoDB
func CreateUndangan(c *gin.Context, dbClient *mongo.Client) {
	var undangan models.Undangan
	if err := c.ShouldBindJSON(&undangan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// set default kehadiran false dan wishes empty
	undangan.Hadir = false
	undangan.Wishes = ""

	undanganCollection := dbClient.Database("undangan_db").Collection("undangan")
	_, err := undanganCollection.InsertOne(context.Background(), undangan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add undangan"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Undangan created successfully"})
}
