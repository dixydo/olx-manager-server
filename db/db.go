package db

import (
	"github.com/dixydo/olxmanager-server/models"
	"gorm.io/driver/mysql"
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

	dsn := "yevhen:password@tcp(127.0.0.1:3306)/olxmanager?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("DB Error")
	}

	dberr := db.AutoMigrate(&models.Advert{})
	if dberr != nil {
		panic("Migration failed")
	}

	return db
}
