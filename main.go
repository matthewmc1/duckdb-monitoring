package main

import (
	"context"
	"net/http"
	"validator/monitoring/internal"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	internal.CreateDB()
	r.POST("/v1/create", createNoteHandler)
	r.GET("/v1/notes", getAllNotes)
	r.Run(":8080")
}

func createNoteHandler(c *gin.Context) {
	ctx := context.Background()
	value := c.Request.Header.Get("X-Auth-Token")
	if value != "home" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized request"})
	}
	internal.CreateNote(ctx, "First Note", "Create")
	c.JSON(http.StatusAccepted, gin.H{"status": "OK"})
}

func getAllNotes(c *gin.Context) {
	ctx := context.Background()
	value := c.Request.Header.Get("X-Auth-Token")
	if value != "home" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized request"})
	}

	c.JSON(http.StatusOK, internal.GetAllNotes(ctx))
}
