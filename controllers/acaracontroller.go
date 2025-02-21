package controllers

import (
	"context"
	"net/http"
	"undangan-digital/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Render insert acara HTML
func RenderInsertAcaraHTML(c *gin.Context, dbClient *mongo.Client) {

	// Render template undangan dan tampilkan data acara
	c.HTML(http.StatusOK, "insertacara.html", gin.H{})
}

// Menampilkan acara berdasarkan ID
func GetAcaraByID(c *gin.Context, dbClient *mongo.Client) {
	id := c.DefaultQuery("id", "")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID acara tidak ditemukan"})
		return
	}

	acaraCollection := dbClient.Database("undangan_db").Collection("acara")

	var acara models.Acara
	err := acaraCollection.FindOne(context.Background(), bson.M{"id": id}).Decode(&acara)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Acara tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, acara)
}

// Menambah acara baru ke MongoDB
func CreateAcara(c *gin.Context, dbClient *mongo.Client) {
	var acara models.Acara
	if err := c.ShouldBindJSON(&acara); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	acaraCollection := dbClient.Database("undangan_db").Collection("acara")
	_, err := acaraCollection.InsertOne(context.Background(), acara)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add acara"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Acara created successfully"})
}
