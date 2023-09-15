package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"strconv"
)

var db *gorm.DB

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Book{})

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/books", NewBook)
	r.GET("/books", ListBook)
	r.GET("/books/:id", GetBook)
	r.Run()
}

type Book struct {
	gorm.Model
	Name   string `json:"name"`
	Author string `json:"author"`
}

func NewBook(c *gin.Context) {
	var book Book
	if err := c.Bind(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	result := db.Create(&book)
	if err := result.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	c.Status(http.StatusCreated)
}

func ListBook(c *gin.Context) {
	var books []Book
	result := db.Find(&books)
	if err := result.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, books)
}

func GetBook(c *gin.Context) {
	var book Book
	id := c.Param("id")
	n, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	result := db.First(&book, n)
	if err := result.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, book)
}
