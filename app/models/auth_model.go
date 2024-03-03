package models

// Credentials Auth model
type AuthCredentials struct {
	ID       string `gorm:"column:auth_id;primaryKey" json:"id"`
	Username string `gorm:"column:auth_username" json:"username"`
	Password string `gorm:"column:auth_password" json:"password"`
	UserID   string `gorm:"column:usr_id" json:"userId"`
}

func (credentials *AuthCredentials) TableName() string {
	return "authcredentials"
}

// Provider Auth model
type AuthProvider struct {
	ID             string `gorm:"column:auth_id;primaryKey" json:"id"`
	Email          string `gorm:"column:auth_prov_email" json:"email"`
	ProviderUserID string `gorm:"column:auth_password" json:"password"`
	ProviderType   string `gorm:"column:auth_prov_type" json:"type"`
	UserID         string `gorm:"column:usr_id" json:"userId"`
}

func (provider *AuthProvider) TableName() string {
	return "authprovider"
}
