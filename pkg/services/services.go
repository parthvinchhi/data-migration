package services

import (
	"fmt"

	"github.com/parthvinchhi/data-migration/pkg/models"
)

type StorageService interface {
	SourceDetails(s *models.SourceAndTarget) error
	TargetDetails(t *models.SourceAndTarget) error
}

type MockStorageService struct{}

func (s *MockStorageService) SaveSource(source *models.SourceAndTarget) error {
	// Mock saving source
	fmt.Printf("Saving source: %+v\n", source)
	return nil
}

// SaveTarget saves the target data
func (s *MockStorageService) SaveTarget(target *models.SourceAndTarget) error {
	// Mock saving target
	fmt.Printf("Saving target: %+v\n", target)
	return nil
}
