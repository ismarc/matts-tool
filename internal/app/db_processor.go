package app

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"

	v5 "github.com/ismarc/matts-tool/internal/app/model/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type dbProcessor struct {
	sourceFilename     string
	sourceDataKey      string
	destination        *gorm.DB
	destinationDataKey string
	destinationAccount string
	noAct              bool
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (db *dbProcessor) init(config DBConfig) {
	db.sourceFilename = config.SourceFilename
	db.sourceDataKey = config.SourceDataKey
	db.destinationDataKey = config.DestinationDataKey
	db.destinationAccount = config.DestinationAccount
	db.noAct = config.NoAct

	if !db.noAct {
		dbConnection, err := gorm.Open(postgres.Open(config.DestinationDSN), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		db.destination = dbConnection
	}
}

type V4User struct {
	login          string
	api_key        string
	encrypted_hash string
}

func (user *V4User) toV5(dataKey string, account string) (result v5.Credential) {
	if strings.HasPrefix(user.login, "host/") {
		base := strings.TrimPrefix(user.login, "host/")
		result.RoleId = fmt.Sprintf("%s:host:%s", account, base)
	} else {
		result.RoleId = fmt.Sprintf("%s:user:%s", account, user.login)
	}

	if len(user.encrypted_hash) > 0 {
		hash, err := AES256GCMEncrypt(dataKey, user.encrypted_hash, result.RoleId)
		if err != nil {
			panic(err)
		}
		raw, err := hex.DecodeString(hash)
		if err != nil {
			panic(err)
		}
		result.EncryptedHash = raw
	}

	if len(user.api_key) > 0 {
		apiKey, err := AES256GCMEncrypt(dataKey, user.api_key, result.RoleId)
		if err != nil {
			panic(err)
		}
		raw, err := hex.DecodeString(apiKey)
		if err != nil {
			panic(err)
		}
		result.ApiKey = raw
	}
	return
}

func (db *dbProcessor) readData() (result []v5.Credential, err error) {
	inFile, err := os.Open(db.sourceFilename)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer inFile.Close()

	startProcessing := false
	stopProcessing := false
	inScanner := bufio.NewScanner(inFile)
	for inScanner.Scan() {
		line := inScanner.Text()
		if line == "COPY authn.users (login, api_key, encrypted_hash, cidr) FROM stdin;" {
			startProcessing = true
			continue
		}

		if startProcessing && line == "\\." {
			stopProcessing = true
		}

		if startProcessing && !stopProcessing {
			values := strings.Split(line, "\t")
			user := V4User{}
			user.login = values[0]
			if values[1] != "\\N" {
				user.api_key = values[1]
			}
			if values[2] != "\\N" {
				user.encrypted_hash = values[2]
			}

			if !strings.HasPrefix(user.login, "host/i-") && !strings.HasPrefix(user.login, "host/azure-linux-agent-v2") && user.login != "admin" {
				decrypted_api_key, err := AES256GCMDecrypt(db.sourceDataKey, user.api_key[3:], user.login)
				if err != nil {
					panic(err)
				}
				user.api_key = string(decrypted_api_key)
				if len(user.encrypted_hash) > 3 {
					decrypted_hash, err := AES256GCMDecrypt(db.sourceDataKey, user.encrypted_hash[3:], user.login)
					if err != nil {
						panic(err)
					}
					user.encrypted_hash = string(decrypted_hash)
				}
				result = append(result, user.toV5(db.destinationDataKey, db.destinationAccount))
			}
		}
	}

	if err = inScanner.Err(); err != nil {
		log.Fatal(err)
		return
	}

	return
}

func (db *dbProcessor) updateData(data []v5.Credential) {
	if db.destination == nil {
		fmt.Printf("NoAct set, would have written:\n")
	}
	for _, credential := range data {
		if db.destination != nil {
			result := db.destination.Save(credential)
			if result.Error != nil {
				panic(result.Error)
			}
		} else {
			fmt.Printf("%s\n", credential.RoleId)
		}
	}
}
