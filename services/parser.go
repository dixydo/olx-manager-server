package services

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/dixydo/olxmanager-server/db"
	"github.com/dixydo/olxmanager-server/models"
	"github.com/dixydo/olxmanager-server/structs"
	"log"
	"net/http"
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

	for item := range attributeResults {
		advert := models.Advert{}
		advert.Title = item.Find("h6").Text()
		item.Find("p").Each(func(i int, s *goquery.Selection) {
			a := structs.Attribute{Key: "data-testid", Value: "ad-price"}

			attr, ok := s.Attr(a.Key)
			if ok && attr == a.Value {
				advert.Price = s.Text()
			}

			a = structs.Attribute{Key: "data-testid", Value: "location-date"}

			attr, ok = s.Attr(a.Key)
			if ok && attr == a.Value {
				advert.Location = s.Text()
			}
		})

		advert.Top = false
		advert.New = false

		item.Find("div").Each(func(i int, s *goquery.Selection) {
			a := structs.Attribute{Key: "title", Value: "Нові"}

			attr, ok := s.Attr(a.Key)
			if ok && attr == a.Value {
				advert.New = true
			}

			a = structs.Attribute{Key: "title", Value: "Б/в"}

			attr, ok = s.Attr(a.Key)
			if ok && attr == a.Value {
				advert.New = false
			}

			a = structs.Attribute{Key: "data-testid", Value: "adCard-featured"}

			attr, ok = s.Attr(a.Key)
			if ok && attr == a.Value {
				advert.Top = true
			}
		})

		orm := db.GetDatabase()
		orm.Create(&advert)
	}

}

func FindByAttribute(a structs.Attribute, s *goquery.Selection, attributeResults chan *goquery.Selection) {
	s.Each(func(i int, s *goquery.Selection) {

		attr, ok := s.Attr(a.Key)
		if ok && attr == a.Value {
			attributeResults <- s
		}
	})
}
