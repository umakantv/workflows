package repo

import (
	"log"

	"github.com/google/uuid"
	"github.com/umakantv/workflows/common/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	Repo *ChangeRequestRepo
)

const (
	DB_PATH = "/Users/umakant.vashishtha/code/personal/tutorials/workflows/gorm.db"
)

func init() {

	db, err := gorm.Open(sqlite.Open(DB_PATH), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	} else {
		log.Println("Connected to database")
	}

	err = db.AutoMigrate(&model.ChangeRequest{})
	if err != nil {
		panic("failed to migrate database")
	}
	Repo = NewChangeRequestRepo(db)
}

func GetChangeRequestRepo() *ChangeRequestRepo {
	return Repo
}

type ChangeRequestRepo struct {
	db *gorm.DB
}

func NewChangeRequestRepo(db *gorm.DB) *ChangeRequestRepo {
	return &ChangeRequestRepo{
		db: db,
	}
}

func (repo *ChangeRequestRepo) InitiateChangeRequest() (string, error) {
	// initiate change request
	changeRequest := model.ChangeRequest{
		ID:     uuid.NewString(),
		Status: model.ChangeRequestStatusDraft,
	}

	result := repo.db.Create(&changeRequest)

	return changeRequest.ID, result.Error
}

func (repo *ChangeRequestRepo) GetChangeRequest(id string) (*model.ChangeRequest, error) {
	// get change request
	var changeRequest model.ChangeRequest

	result := repo.db.Debug().First(&changeRequest, "id = ?", id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &changeRequest, result.Error
}

func (repo *ChangeRequestRepo) UpdateChangeRequest(changeRequest *model.ChangeRequest) error {
	// get change request

	result := repo.db.Save(changeRequest)

	return result.Error
}
