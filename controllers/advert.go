package controllers

import (
	"github.com/dixydo/olxmanager-server/db"
	"github.com/dixydo/olxmanager-server/models"
	"github.com/dixydo/olxmanager-server/services"
	"github.com/kataras/iris/v12"
)

type AdvertController struct {
	ctx iris.Context
}

func (c *AdvertController) Get() []models.Advert {
	var adverts []models.Advert
	database := db.GetDatabase()
	database.Find(&adverts)

	return adverts
}

func (c *AdvertController) GetParse(ctx iris.Context) {
	services.Parse()
	ctx.JSON("Parsed")
}
