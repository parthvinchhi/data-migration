package models

type User struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Email string
}

type SourceAndTarget struct {
	SourceDBType   string
	SourceHost     string
	SourcePort     string
	SourceUser     string
	SourcePassword string
	SourceDBName   string
	TargetDBType   string
	TargetHost     string
	TargetPort     string
	TargetUser     string
	TargetPassword string
	TargetDBName   string
}
