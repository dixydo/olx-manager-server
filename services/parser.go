package services

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/dixydo/olxmanager-server/db"
	"github.com/dixydo/olxmanager-server/models"
	"github.com/dixydo/olxmanager-server/structs"
)

func Parse() {
	attributeResults := make(chan *goquery.Selection)
	res, err := http.Get("https://www.olx.ua/d/uk/list/q-macbook-air-m1-16Gb/")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".listing-grid-container").Children().Each(func(i int, s *goquery.Selection) {
		attr, ok := s.Attr("data-testid")

		if ok && attr == "listing-grid" {
			a := structs.Attribute{Key: "data-cy", Value: "l-card"}
			go FindByAttribute(a, s.Children(), attributeResults)
		}
	})

	orm := db.GetDatabase()

	var adverts []models.Advert

	for item := range attributeResults {
		advert := models.Advert{}
		advert.Title = item.Find("h6").Text()

		item.Children().Find("p").Each(func(i int, s *goquery.Selection) {
			result, ok := s.Attr("data-testid")

			if ok && result == "ad-price" {
				advert.Price = s.Text()
			}
		})

		item.Children().Find("div").Each(func(i int, s *goquery.Selection) {
			result, ok := s.Attr("data-testid")

			if ok && result == "adCard-featured" {
				advert.Top = true
			}
		})

		item.Children().Find("span").Each(func(i int, s *goquery.Selection) {
			result, ok := s.Attr("title")

			if ok && result == "Вживане" {
				advert.New = false
			}

			if ok && result == "Нове" {
				advert.New = true
			}
		})

		adverts = append(adverts, advert)
	}

	orm.Create(&adverts)
}

func FindByAttribute(a structs.Attribute, s *goquery.Selection, attributeResults chan *goquery.Selection) {
	s.Each(func(i int, s *goquery.Selection) {

		value, ok := s.Attr(a.Key)
		if ok && value == a.Value {
			attributeResults <- s
		}
	})

	close(attributeResults)
}
