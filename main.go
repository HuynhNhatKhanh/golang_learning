package main

import (
	"Fooddelivery/common"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Restaurant struct {
	common.SQLModel
	Name string `json:"name" gorm:"column:name;"`
	Addr string `json:"addr" gorm:"column:addr;"`
}

func (Restaurant) TableName() string { return "restaurants" }

type RestaurantCreate struct {
	common.SQLModel
	Name string `json:"name" gorm:"column:name;"`
	Addr string `json:"addr" gorm:"column:addr;"`
}

func (RestaurantCreate) TableName() string { return Restaurant{}.TableName() }

type RestaurantUpdate struct {
	Name *string `json:"name" gorm:"column:name;"`
	Addr *string `json:"addr" gorm:"column:addr;"`
}

func (RestaurantUpdate) TableName() string { return Restaurant{}.TableName() }

func main() {
	dsn := os.Getenv("MYSQL_CONN_STRING")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}
	log.Println(db)

	db = db.Debug()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Version, Group
	v1 := r.Group("/v1")
	{
		restaurants := v1.Group("/restaurants")
		{
			// API Post
			restaurants.POST("", func(c *gin.Context) {
				var newData RestaurantCreate

				if err := c.ShouldBind(&newData); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				if err := db.Create(&newData).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "err.Error()"})
					return
				}

				c.JSON(http.StatusOK, gin.H{"data": newData.Id})
			})

			// API get by id
			restaurants.GET("/:id", func(c *gin.Context) {
				var data Restaurant

				id, err := strconv.Atoi(c.Param("id"))

				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				if err := db.Where("id=?", id).First(&data).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "err.Error()"})
					return
				}

				c.JSON(http.StatusOK, gin.H{"data": data})
			})

			// API get all
			restaurants.GET("", func(c *gin.Context) {
				var data []Restaurant

				type Paging struct {
					Page  int `json:"page" form:"page"`
					Limit int `json:"limit" form:"limit"`
				}

				var pagingData Paging

				if err := c.ShouldBind(&pagingData); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				if pagingData.Page <= 0 {
					pagingData.Page = 1
				}

				if pagingData.Limit <= 0 {
					pagingData.Limit = 5
				}

				offset := pagingData.Page - 1

				if err := db.Offset(offset).Limit(pagingData.Limit).Order("id desc").Find(&data).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "err.Error()"})
					return
				}

				c.JSON(http.StatusOK, gin.H{"data": data})
			})

			// API update
			restaurants.PATCH("/:id", func(c *gin.Context) {
				id, err := strconv.Atoi(c.Param("id"))

				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				var data RestaurantUpdate

				if err := c.ShouldBind(&data); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				if err := db.Where("id=?", id).Updates(&data).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "err.Error()"})
					return
				}

				c.JSON(http.StatusOK, gin.H{"data": true})
			})

			// API delete
			restaurants.DELETE("/:id", func(c *gin.Context) {
				id, err := strconv.Atoi(c.Param("id"))

				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				if err := db.Table(Restaurant{}.TableName()).Where("id=?", id).Delete(nil).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "err.Error()"})
					return
				}

				c.JSON(http.StatusOK, gin.H{"data": true})
			})
		}
	}

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
