package app

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"gorm.io/gorm"
)

type dbProcessor struct {
	sourceFilename     string
	sourceDataKey      string
	destination        *gorm.DB
	destinationVersion string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (db *dbProcessor) init(config DBConfig) {
	db.sourceFilename = config.SourceFilename
	db.sourceDataKey = config.SourceDataKey
	// destinationDSN, err := dburl.Parse(config.DestinationConnectionString)
	// if err != nil {
	// 	panic(err)
	// }

	// db.destination, err = gorm.Open(postgres.Open(destinationDSN.DSN), &gorm.Config{})
	// if err != nil {
	// 	panic(err)
	// }
}

type V4User struct {
	login          string
	api_key        string
	encrypted_hash string
}

func (db *dbProcessor) readData() (result []V4User, err error) {
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

			if !strings.HasPrefix(user.login, "host/i-") && !strings.HasPrefix(user.login, "host/azure-linux-agent-v2") {
				result = append(result, user)
			}

			decrypted, err := AES256GCMDecrypt(db.sourceDataKey, user.api_key[3:], "conjur:user:"+user.login)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Decrypted: %+v\n", decrypted)
		}
	}

	if err = inScanner.Err(); err != nil {
		log.Fatal(err)
		return
	}

	return
}
