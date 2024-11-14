package template

import (
	"crypto/md5"
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"server/log"
	"strings"
)

//go:embed pages/*
var PageFiles embed.FS

//go:embed pages/static/js/*
var StaticJsFiles embed.FS

//go:embed pages/index.html
var Indexhtml []byte

func ReadFile(name string) []byte {
	var fileData []byte
	var err error
	if strings.HasPrefix(name, "static/js") {
		fileData, err = StaticJsFiles.ReadFile("pages/" + name)
	} else {
		// Read file from PageFiles
		fileData, err = PageFiles.ReadFile("pages/" + name)
	}
	if err != nil {
		log.TLogln("Failed to read Dlnaicon48.png: %v", err)
		return nil
	}
	return fileData
}

func RouteWebPages(route gin.IRouter, engine *gin.Engine) {
	route.GET("/", func(c *gin.Context) {
		etag := fmt.Sprintf("%x", md5.Sum(Indexhtml))
		c.Header("Cache-Control", "public, max-age=31536000")
		c.Header("ETag", etag)
		c.Data(200, "text/html; charset=utf-8", Indexhtml)
	})

	engine.NoRoute(func(c *gin.Context) {
		filePath := c.Request.URL.Path
		fileData, err := PageFiles.ReadFile("pages" + filePath)
		if err != nil {
			// If the file doesn't exist, return 404
			c.Status(http.StatusNotFound)
			return
		}
		contentType := http.DetectContentType(fileData)

		// Set headers for caching and ETag
		etag := fmt.Sprintf("%x", md5.Sum(fileData))
		c.Header("Cache-Control", "public, max-age=31536000")
		c.Header("ETag", etag)

		// Send the file with the correct MIME type
		c.Data(http.StatusOK, contentType, fileData)
	})

}
