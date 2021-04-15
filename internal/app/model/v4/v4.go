package v4

import (
	"github.com/jackc/pgtype"
)

// User v4 user type
type User struct {
	Login         string `gorm:"primaryKey;not null"`
	ApiKey        []byte
	EncryptedHash []byte
	Cidr          []pgtype.Inet `gorm:"type:cidr[]"`
}

// TableName v4 user table name generation since it's in a different schema
func (User) TableName() string {
	return "authn.users"
}

// SlosiloKeystore v4 slosilo_keystore type
type SlosiloKeystore struct {
	ID          string `gorm:"not null"`
	Key         []byte `gorm:"not null;type:bytea"`
	Fingerprint string `gorm:"not null"`
}

// TableName v4 slosilo_keystore table name
func (SlosiloKeystore) TableName() string {
	return "authn.slosilo_keystore"
}

// Data a full database of data
type Data struct {
	Users           []User
	SlosiloKeystore []SlosiloKeystore
}
