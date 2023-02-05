package models

import "gorm.io/gorm"

type Advert struct {
	gorm.Model
	Title    string
	Price    string
	Location string
	New      bool
	Top      bool
}
