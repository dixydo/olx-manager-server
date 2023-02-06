package db

import (
	"github.com/dixydo/olxmanager-server/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDatabase() *gorm.DB {

	// opts := basicauth.Options{
	// 	Allow: basicauth.AllowUsers(map[string]string{
	// 		"username": "password",
	// 	}),
	// 	Realm:        "Authorization Required",
	// 	ErrorHandler: basicauth.DefaultErrorHandler,
	// 	// [...more options]
	// }

	// auth := basicauth.New(opts)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "user=yevhen password=password dbname=olxmanager port=5432",
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		panic("DB Error")
	}

	dberr := db.AutoMigrate(&models.Advert{})
	if dberr != nil {
		panic("Migration failed")
	}

	return db
}
