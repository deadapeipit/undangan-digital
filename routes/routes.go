package routes

import (
	"undangan-digital/controllers"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// SetupRoutes untuk mengatur semua routes termasuk untuk menambah acara dan undangan
func SetupRoutes(router *gin.Engine, dbClient *mongo.Client) {

	//API Route:

	// Route untuk mendapatkan acara berdasarkan ID
	router.GET("/acara", func(c *gin.Context) {
		controllers.GetAcaraByID(c, dbClient)
	})

	// Route untuk menambah acara
	router.POST("/acara", func(c *gin.Context) {
		controllers.CreateAcara(c, dbClient)
	})

	// Route untuk mendapatkan undangan berdasarkan Acara ID
	router.GET("/undangan", func(c *gin.Context) {
		controllers.GetUndanganByAcaraID(c, dbClient)
	})

	// Route untuk menambah undangan
	router.POST("/undangan", func(c *gin.Context) {
		controllers.CreateUndangan(c, dbClient)
	})

	// Route untuk melihat semua undangan
	router.GET("/invitations", func(c *gin.Context) {
		controllers.GetAllRSVPs(c, dbClient)
	})

	// Route konfirmasi kehadiran undangan
	router.POST("/rsvp/:id", func(c *gin.Context) {
		controllers.HandleRSVP(c, dbClient)
	})

	//HTML route:

	// Route HTML template undangan
	router.GET("/undangan/html/:id", func(c *gin.Context) {
		controllers.RenderInvitationHTML(c, dbClient)
	})

	// Route HTML template insert undangan
	router.GET("/undangan/html/insert/:id", func(c *gin.Context) {
		controllers.RenderInsertInvitationHTML(c, dbClient)
	})

	// Route HTML template insert undangan
	router.GET("/acara/html/insert", func(c *gin.Context) {
		controllers.RenderInsertAcaraHTML(c, dbClient)
	})
}
