package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/parthvinchhi/data-migration/pkg/db"
	"github.com/parthvinchhi/data-migration/pkg/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// User model for data migration
type User struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Email string
}

func main() {
	router := gin.Default()
	router.Static("/static", "./pkg/static")
	router.LoadHTMLGlob("pkg/templates/*")

	// Serve the form page
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// Handle form submission and data migration
	router.POST("/migrate", func(c *gin.Context) {
		// Parse form data for source and target DBs
		sourceDBType := c.PostForm("source_db_type")
		sourceHost := c.PostForm("source_host")
		sourcePort := c.PostForm("source_port")
		sourceUser := c.PostForm("source_user")
		sourcePassword := c.PostForm("source_password")
		sourceDBName := c.PostForm("source_dbname")

		targetDBType := c.PostForm("target_db_type")
		targetHost := c.PostForm("target_host")
		targetPort := c.PostForm("target_port")
		targetUser := c.PostForm("target_user")
		targetPassword := c.PostForm("target_password")
		targetDBName := c.PostForm("target_dbname")

		// Source DSN
		sourceDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", sourceUser, sourcePassword, sourceHost, sourcePort, sourceDBName)
		var sourceDB *gorm.DB
		var err error

		if sourceDBType == "mysql" {
			sourceDB, err = db.ConnectMySQL(sourceDSN)
		} else {
			sourceDSN = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", sourceHost, sourceUser, sourcePassword, sourceDBName, sourcePort)
			sourceDB, err = db.ConnectPostgres(sourceDSN)
		}

		if err != nil {
			log.Fatal("Failed to connect to source DB: ", err)
		}

		// Target DSN
		targetDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", targetUser, targetPassword, targetHost, targetPort, targetDBName)
		var targetDB *gorm.DB

		if targetDBType == "mysql" {
			targetDB, err = db.ConnectMySQL(targetDSN)
		} else {
			targetDSN = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", targetHost, targetUser, targetPassword, targetDBName, targetPort)
			targetDB, err = db.ConnectPostgres(targetDSN)
		}

		if err != nil {
			log.Fatal("Failed to connect to target DB: ", err)
		}

		// Auto-migrate User model in target database
		targetDB.AutoMigrate(&models.User{})

		// Fetch data from source database
		var users []models.User
		if err := sourceDB.Find(&users).Error; err != nil {
			log.Fatal("Failed to fetch data from source: ", err)
		}

		// Insert data into target database
		if err := targetDB.Create(&users).Error; err != nil {
			log.Fatal("Failed to migrate data: ", err)
		}

		c.String(http.StatusOK, "Data migration complete!")
	})

	// Run the server
	router.Run(":8080")
}
