package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/parthvinchhi/data-migration/pkg/db"
	"github.com/parthvinchhi/data-migration/pkg/models"
	"gorm.io/gorm"
)

func LoadHtml(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func MigrateHandler(c *gin.Context) {
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
}
