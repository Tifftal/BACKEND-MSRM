package repository

import (
	"MSRM/internal/app/ds"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
	}, nil
}

func (repository *Repository) GetSampleByID(id int) (*ds.Samples, error) {
	sample := &ds.Samples{}

	// Assuming "Sample_status" is the correct field name for the status
	err := repository.db.First(sample, "Id_sample = ?", id).Error
	if err != nil {
		return nil, err
	}

	return sample, nil
}

func (repository *Repository) GetAllSamples() ([]ds.Samples, error) {
	samples := []ds.Samples{}

	err := repository.db.Where("Sample_status = ?", "Active").Order("Id_sample ASC").Find(&samples).Error
	if err != nil {
		return nil, err
	}

	return samples, nil
}

func (repository *Repository) GetSampleByName(name string) ([]ds.Samples, error) {
	var samples []ds.Samples
	err := repository.db.Where("Name LIKE ?", "%"+name+"%").Order("Sample_status ASC").Order("Id_sample ASC").Find(&samples).Error
	if err != nil {
		return nil, err
	}
	return samples, nil
}

func (r *Repository) DeleteSampleByID(id int) error {
	if err := r.db.Exec("UPDATE samples SET sample_status='Deleted' WHERE Id_sample= ?", id).Error; err != nil {
		return err
	}
	return nil
}
