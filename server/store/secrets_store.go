package store

import (
	"errors"
	"openhands-go/server/models"
)

type SecretsStore struct{}

func NewSecretsStore() *SecretsStore {
	return &SecretsStore{}
}

func (s *SecretsStore) Save(secret *models.SecretInfo) error {
	if DB == nil {
		return errors.New("database not initialized")
	}
	return DB.Save(secret).Error
}

func (s *SecretsStore) Get(name string) (*models.SecretInfo, error) {
	if DB == nil {
		return nil, errors.New("database not initialized")
	}
	var secret models.SecretInfo
	if err := DB.First(&secret, "name = ?", name).Error; err != nil {
		return nil, err
	}
	return &secret, nil
}

func (s *SecretsStore) Delete(name string) error {
	if DB == nil {
		return errors.New("database not initialized")
	}
	return DB.Delete(&models.SecretInfo{}, "name = ?", name).Error
}

func (s *SecretsStore) GetAll() ([]models.SecretInfo, error) {
	if DB == nil {
		return nil, errors.New("database not initialized")
	}
	var secrets []models.SecretInfo
	if err := DB.Find(&secrets).Error; err != nil {
		return nil, err
	}
	return secrets, nil
}
