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
//go:embed pages/static/js/*
//go:embed docs/*
var files embed.FS

//go:embed pages/index.html
var Indexhtml []byte

func ReadFileByName(name string) []byte {
	return readFileByPath("pages/" + name)
}

func readFileByPath(name string) []byte {
	var fileData, err = files.ReadFile(name)

	if err != nil {
		log.TLogln("Failed to read file", name, err)
		return nil
	}
	return fileData
}

func handleSwaggerfunc(fileName string, c *gin.Context) {
	fileData := readFileByPath("docs" + fileName)
	if fileData == nil {
		c.Status(http.StatusNotFound)
		return
	}
	contentType := http.DetectContentType(fileData)
	c.Data(http.StatusOK, contentType, fileData)
}

func RouteWebPages(route gin.IRouter, engine *gin.Engine) {
	route.GET("/", func(c *gin.Context) {
		etag := fmt.Sprintf("%x", md5.Sum(Indexhtml))
		c.Header("Cache-Control", "public, max-age=31536000")
		c.Header("ETag", etag)
		c.Data(200, "text/html; charset=utf-8", Indexhtml)
	})

	route.GET("/swagger-ui/*any", func(c *gin.Context) {
		fileName := c.Param("any")
		if fileName == "/" || strings.EqualFold(fileName, "/index.html") {
			fileName = "/redoc.html"
		}
		handleSwaggerfunc(fileName, c)
	})

	//route.GET("/swagger-ui/swagger.yaml", func(c *gin.Context) {
	//	handleSwaggerfunc("swagger.yaml", c)
	//})
	//
	//route.GET("/swagger-ui/redoc.standalone.js", func(c *gin.Context) {
	//	handleSwaggerfunc("redoc.standalone.js", c)
	//})

	engine.NoRoute(func(c *gin.Context) {
		filePath := c.Request.URL.Path
		fileData, err := files.ReadFile("pages" + filePath)
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
