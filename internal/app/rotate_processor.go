package app

import (
	"encoding/hex"
	"fmt"

	v5 "github.com/ismarc/matts-tool/internal/app/model/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type rotateProcessor struct {
	Connection         *gorm.DB
	SourceDataKey      string
	DestinationDataKey string
	NoAct              bool
}

func (rp *rotateProcessor) process() {
	fmt.Printf("--- Matt's Tool::Conjur Data Key Rotator ---\n")
	rp.checkDataKeys()
	rp.reencryptSloSiloKey()
	rp.reencryptUsers()
	rp.reencryptSecrets()
	fmt.Printf("--- Conjur Data Key Rotation Complete ---\n")
}

func (rp *rotateProcessor) checkDataKeys() {
	if len(rp.SourceDataKey) < 32 {
		panic(fmt.Sprintf("❌ Source data key has invalid length: %d should be 44. Check IN_CONJUR_DATA_KEY env var.", len(rp.SourceDataKey)))
	}
	if len(rp.DestinationDataKey) < 32 {
		panic(fmt.Sprintf("❌ Destination data key has invalid length: %d should be 44. Check OUT_CONJUR_DATA_KEY env var.", len(rp.DestinationDataKey)))
	}
}

func (rp *rotateProcessor) reencyrptBytes(input []byte, additionalData string) []byte {
	inputString := hex.EncodeToString(input)
	decrypted, err := AES256GCMDecrypt(rp.SourceDataKey, inputString, additionalData)
	check(err)
	recrryptedString, err := AES256GCMEncrypt(rp.DestinationDataKey, string(decrypted), additionalData)
	check(err)
	decryptCheck, err := AES256GCMDecrypt(rp.DestinationDataKey, recrryptedString, additionalData)
	check(err)
	if string(decryptCheck) != string(decrypted) {
		fmt.Printf("❌Decrypt check failed, %s != %s\n", decryptCheck, decrypted)
	}
	recryptedBytes, error := hex.DecodeString(recrryptedString)
	check(error)
	return recryptedBytes
}

func (rp *rotateProcessor) reencryptSloSiloKey() {
	var sloSiloKeys []v5.SlosiloKeystore

	fmt.Printf("♻️🔑 Re-Encrypting Slo Silo Keys\n")
	// Find and re-encrypt slosilo key
	result := rp.Connection.Find(&sloSiloKeys)
	check(result.Error)
	for i, key := range sloSiloKeys {
		if len(key.Key) > 0 {
			fmt.Printf("🔑 SloSilo Key %d %s size: %d\n", i, key.Id, len(key.Key))
			newKey := rp.reencyrptBytes(key.Key, key.Id)
			key.Key = newKey
			if !rp.NoAct {
				result = rp.Connection.Save(&key)
				check(result.Error)
			}
		}
	}
}

func (rp *rotateProcessor) reencryptUsers() {
	var creds []v5.Credential

	fmt.Printf("♻️👤 Re-Encrypting User API-Keys and Passwords\n")
	// Find and re-encrypt Login Credentials
	result := rp.Connection.Find(&creds)
	check(result.Error)
	for i, cred := range creds {
		fmt.Printf("👤 Login Credential %d %s\n", i, cred.RoleId)
		if len(cred.ApiKey) > 0 {
			newApiKey := rp.reencyrptBytes(cred.ApiKey, cred.RoleId)
			cred.ApiKey = newApiKey
		}
		if len(cred.EncryptedHash) > 0 {
			newHash := rp.reencyrptBytes(cred.EncryptedHash, cred.RoleId)
			cred.EncryptedHash = newHash

		}
		if !rp.NoAct {
			result = rp.Connection.Save(&cred)
			check(result.Error)
		}
	}
}

func (rp *rotateProcessor) reencryptSecrets() {
	var secrets []v5.Secret
	fmt.Printf("♻️🔒 Re-Encrypting Secrets\n")
	// Find and re-encrypt secrets
	result := rp.Connection.Find(&secrets)
	check(result.Error)
	for i, secret := range secrets {
		fmt.Printf("🔒 Secret %d %s\n", i, secret.ResourceId)
		if len(secret.Value) > 0 {
			newValue := rp.reencyrptBytes(secret.Value, secret.ResourceId)
			secret.Value = newValue
			if !rp.NoAct {
				result = rp.Connection.Save(&secret)
				check(result.Error)
			}
		}
	}
}

func (rp *rotateProcessor) init(config RotateConfig) {
	rp.SourceDataKey = config.SourceDataKey
	rp.DestinationDataKey = config.DestinationDataKey
	rp.NoAct = config.NoAct

	dbConnection, err := gorm.Open(postgres.Open(config.DSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	rp.Connection = dbConnection
}
