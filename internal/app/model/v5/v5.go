package v5

type Credential struct {
	RoleId        string `gorm:"primaryKey;not null"`
	ApiKey        []byte
	EncryptedHash []byte
}
