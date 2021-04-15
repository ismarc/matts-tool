package app

import (
	v4 "github.com/ismarc/matts-tool/internal/app/model/v4"
	"github.com/xo/dburl"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type dbProcessor struct {
	source             *gorm.DB
	sourceVersion      string
	destination        *gorm.DB
	destinationVersion string
}

func (db *dbProcessor) init(config DBConfig) {
	sourceDSN, err := dburl.Parse(config.SourceConnectionString)
	if err != nil {
		panic(err)
	}

	// destinationDSN, err := dburl.Parse(config.DestinationConnectionString)
	// if err != nil {
	// 	panic(err)
	// }

	db.source, err = gorm.Open(postgres.Open(sourceDSN.DSN), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	// db.destination, err = gorm.Open(postgres.Open(destinationDSN.DSN), &gorm.Config{})
	// if err != nil {
	// 	panic(err)
	// }
}

type dataResult struct {
	V4Users            []v4.User
	V4SlosiloKeystores []v4.SlosiloKeystore
}

func (db *dbProcessor) fetchData() (result dataResult, err error) {
	switch db.sourceVersion {
	case "4":
		users := []v4.User{}
		err = db.source.Find(&users).Error
		if err != nil {
			return
		}
		slosilo := []v4.SlosiloKeystore{}
		err = db.source.Find(&slosilo).Error
		if err != nil {
			return
		}
		result.V4Users = users
		result.V4SlosiloKeystores = slosilo
	}

	return
}
