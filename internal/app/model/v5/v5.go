package v5

import (
	"time"
)

type Credential struct {
	RoleId        string `gorm:"primaryKey;not null"`
	ApiKey        []byte
	EncryptedHash []byte
}

type Secret struct {
	Version    int `gorm:"primaryKey;not null"`
	Value      []byte
	ResourceId string `gorm:"primaryKey;not null"`
	ExpiresAt  time.Time
}

type SlosiloKeystore struct {
	Id          string `gorm:"primaryKey;not null"`
	Key         []byte `gorm:"not null;type:bytea"`
	Fingerprint string `gorm:"not null"`
}

// TableName v5 slosilo_keystore table name
func (SlosiloKeystore) TableName() string {
	return "slosilo_keystore"
}
