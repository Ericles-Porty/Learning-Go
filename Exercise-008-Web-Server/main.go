package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums = []Album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func main() {

	router := gin.Default()
	// router.SetTrustedProxies(nil)
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Hello World!",
		})
	})

	router.GET("/album", getAlbumsHandler)

	router.GET("/album/:id", getAlbumByIdHandler)

	router.POST("/album", postAlbumHandler)

	router.PUT("/album/:id", putAlbumByIdHandler)

	router.DELETE("/album/:id", deleteAlbumByIdHandler)

	router.Run(":8080")
}

func getAlbumsHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, albums)
}

func getAlbumByIdHandler(ctx *gin.Context) {
	id := ctx.Param("id")

	for _, album := range albums {
		if album.ID == id {
			ctx.JSON(http.StatusOK, album)
			return
		}
	}

	ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
}

func postAlbumHandler(ctx *gin.Context) {
	var newAlbum Album

	if err := ctx.BindJSON(&newAlbum); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	albums = append(albums, newAlbum)
	ctx.JSON(http.StatusOK, newAlbum)
}

func putAlbumByIdHandler(ctx *gin.Context) {
	var album Album
	id := ctx.Param("id")

	if err := ctx.BindJSON(&album); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, alb := range albums {
		if alb.ID == id {
			albums[i] = album
			ctx.JSON(http.StatusOK, album)
			return
		}
	}

	ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
}

func deleteAlbumByIdHandler(ctx *gin.Context) {
	id := ctx.Param("id")

	for i, album := range albums {
		if album.ID == id {
			albums = append(albums[:i], albums[i+1:]...)
			ctx.JSON(http.StatusOK, gin.H{"message": "deleted", "album": album})
			return
		}
	}

	ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
}
