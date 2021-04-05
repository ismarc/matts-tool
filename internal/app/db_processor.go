package app

import (
	v4 "github.com/ismarc/policy-handler/internal/app/model/v4"
	"github.com/xo/dburl"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type dbProcessor struct {
	v4     *gorm.DB
	v4Data v4.Data
	v5     *gorm.DB
}

func (db *dbProcessor) init(v4Url string, v5Url string) {
	v4DSN, err := dburl.Parse(v4Url)
	if err != nil {
		panic(err.Error())
	}

	db.v4, err = gorm.Open(postgres.Open(v4DSN.DSN), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	// db.v4.AutoMigrate(&v4.Resource{})
	// var data v4.SchemaMigration
	// var data struct{ Filename string }
	// result := db.v4.Table("schema_migrations").Find(&data)
	// if result.Error != nil {
	// 	panic(result.Error)
	// }
	// fmt.Printf("Migration: %+v\n", data)
}

func (db *dbProcessor) loadData() {
	result := db.v4.Find(&db.v4Data.Resources)
	handleError(result.Error)

	result = db.v4.Find(&db.v4Data.Roles)
	handleError(result.Error)

	result = db.v4.Find(&db.v4Data.Annotations)
	handleError(result.Error)

	result = db.v4.Find(&db.v4Data.Credentials)
	handleError(result.Error)

	result = db.v4.Find(&db.v4Data.HostFactoryTokens)
	handleError(result.Error)

	result = db.v4.Find(&db.v4Data.Permissions)
	handleError(result.Error)

	result = db.v4.Find(&db.v4Data.PolicyVersions)
	handleError(result.Error)

	result = db.v4.Find(&db.v4Data.ResourcesTextsearch)
	handleError(result.Error)

	result = db.v4.Find(&db.v4Data.RoleMemberships)
	handleError(result.Error)

	result = db.v4.Find(&db.v4Data.SchemaMigrations)
	handleError(result.Error)

	result = db.v4.Find(&db.v4Data.Secrets)
	handleError(result.Error)

	result = db.v4.Find(&db.v4Data.SlosiloKeystore)
	handleError(result.Error)
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
