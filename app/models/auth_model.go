package models

import (
	"uidealist/pkg/repository"

	"github.com/google/uuid"
)

// Credentials Auth model
type AuthCredentials struct {
	ID       uuid.UUID `gorm:"column:auth_id;primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Username string    `gorm:"column:auth_username" json:"username"`
	Password string    `gorm:"column:auth_password" json:"password"`
	UserID   uuid.UUID `gorm:"column:usr_id;type:uuid" json:"userId"`
}

func (credentials *AuthCredentials) TableName() string {
	return "authcredentials"
}

// Provider Auth model
type AuthProvider struct {
	ID             uuid.UUID           `gorm:"column:auth_id;primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Email          string              `gorm:"column:auth_prov_email" json:"email"`
	ProviderUserID string              `gorm:"column:auth_password" json:"password"`
	ProviderType   repository.Provider `gorm:"column:auth_prov_type" json:"type"`
	UserID         uuid.UUID           `gorm:"column:usr_id;type:uuid" json:"userId"`
}

func (provider *AuthProvider) TableName() string {
	return "authprovider"
}
