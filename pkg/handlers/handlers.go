package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/parthvinchhi/data-migration/pkg/db"
	"github.com/parthvinchhi/data-migration/pkg/models"
	"github.com/parthvinchhi/data-migration/pkg/services"
	"gorm.io/gorm"
)

type SnT models.SourceAndTarget

func LoadHtml(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

type Handler struct {
	StorageService services.StorageService
}

func NewHandler(storageService services.StorageService) *Handler {
	return &Handler{StorageService: storageService}
}

func (h *Handler) HandleMigrate(c *gin.Context) {
	source := &models.SourceAndTarget{
		SourceDBType:   c.PostForm("source_db_type"),
		SourceHost:     c.PostForm("source_host"),
		SourcePort:     c.PostForm("source_port"),
		SourceUser:     c.PostForm("source_user"),
		SourcePassword: c.PostForm("source_password"),
		SourceDBName:   c.PostForm("source_dbname"),
	}

	target := &models.SourceAndTarget{
		TargetDBType:   c.PostForm("target_db_type"),
		TargetHost:     c.PostForm("target_host"),
		TargetPort:     c.PostForm("target_port"),
		TargetUser:     c.PostForm("target_user"),
		TargetPassword: c.PostForm("target_password"),
		TargetDBName:   c.PostForm("target_dbname"),
	}

	// Save the data using the injected storage service
	if err := h.StorageService.SourceDetails(source); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save source data"})
		return
	}

	if err := h.StorageService.TargetDetails(target); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save target data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data migration request received"})

}

func (s *SnT) MigrateHandler(c *gin.Context) {
	s.getSourceDetails(c)
	s.getTargetDetails(c)
	sourceDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", s.SourceUser, s.SourcePassword, s.SourceHost, s.SourcePort, s.SourceDBName)
	var sourceDB *gorm.DB
	var err error

	if s.SourceDBType == "mysql" {
		sourceDB, err = db.ConnectMySQL(sourceDSN)
	} else {
		sourceDSN = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", s.SourceHost, s.SourceUser, s.SourcePassword, s.SourceDBName, s.SourcePort)
		sourceDB, err = db.ConnectPostgres(sourceDSN)
	}

	if err != nil {
		log.Fatal("Failed to connect to source DB: ", err)
	}

	// Target DSN
	targetDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", s.TargetUser, s.TargetPassword, s.TargetHost, s.TargetPort, s.TargetDBName)
	var targetDB *gorm.DB

	if s.TargetDBType == "mysql" {
		targetDB, err = db.ConnectMySQL(targetDSN)
	} else {
		targetDSN = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", s.TargetHost, s.TargetUser, s.TargetPassword, s.TargetDBName, s.TargetPort)
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

func (s *SnT) getSourceDetails(c *gin.Context) {
	s.SourceDBType = c.PostForm("source_db_type")
	s.SourceHost = c.PostForm("source_host")
	s.SourcePort = c.PostForm("source_port")
	s.SourceUser = c.PostForm("source_user")
	s.SourcePassword = c.PostForm("source_password")
	s.SourceDBName = c.PostForm("source_dbname")
}

func (s *SnT) getTargetDetails(c *gin.Context) {
	s.TargetDBType = c.PostForm("target_db_type")
	s.TargetHost = c.PostForm("target_host")
	s.TargetPort = c.PostForm("target_port")
	s.TargetUser = c.PostForm("target_user")
	s.TargetPassword = c.PostForm("target_password")
	s.TargetDBName = c.PostForm("target_dbname")
}
