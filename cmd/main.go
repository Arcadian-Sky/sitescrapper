package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Arcadian-Sky/scrapper/internal/models"
	"github.com/Arcadian-Sky/scrapper/internal/scrapper"
	"github.com/Arcadian-Sky/scrapper/internal/storage"
	"github.com/gocolly/colly/v2"
)

func main() {
	domain := "mekka.spb.ru"

	// Инициализируем подключение к ClickHouse
	ch, err := storage.NewClickHouseStorage("localhost", "9000", "scrapper", "admin", "admin123")
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	product := models.Product{
		Name:        "1111",
		Description: "2222",
		Price:       "3333",
		URL:         "4444",
	}
	fmt.Printf("Найден товар: %+v\n", product)

	// Сохраняем товар в ClickHouse
	if err := ch.SaveProduct(context.Background(), product); err != nil {
		log.Printf("Ошибка сохранения товара в ClickHouse: %v", err)
	}
	os.Exit(0)

	// Парсим robots.txt
	robots, err := scrapper.ParseRobotsTxt(domain)
	if err != nil {
		log.Fatal(err)
	}

	// Создаем коллектор для парсинга страниц каталога
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.MaxDepth(2),
	)

	// Обработчик для страниц каталога
	c.OnHTML("div[itemscope]", func(e *colly.HTMLElement) {
		itemType := e.Attr("itemtype")

		if itemType == "http://schema.org/Product" {
			product := models.Product{
				Name:        e.ChildText("div[itemprop=\"name\"]"),
				Description: e.ChildText("div[itemprop=\"description\"]"),
				Price:       e.ChildAttr("[itemprop=\"price\"]", "content"),
				URL:         e.ChildAttr("[itemprop=\"url\"]", "href"),
			}
			fmt.Printf("Найден товар: %+v\n", product)

			// Сохраняем товар в ClickHouse
			if err := ch.SaveProduct(context.Background(), product); err != nil {
				log.Printf("Ошибка сохранения товара в ClickHouse: %v", err)
			}
		}
	})

	// Обработка ошибок
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Ошибка при запросе %s: %v", r.Request.URL, err)
	})

	// Парсим каждый sitemap
	for _, sitemapURL := range robots.Sitemaps {
		urls, err := scrapper.ParseSitemap(sitemapURL)
		if err != nil {
			log.Printf("Ошибка при парсинге sitemap %s: %v", sitemapURL, err)
			continue
		}

		// Посещаем каждую страницу каталога
		for _, url := range urls {
			if strings.Contains(url, "/catalog/") { // Проверяем, что это страница каталога
				err := c.Visit(url)
				if err != nil {
					log.Printf("Ошибка при посещении %s: %v", url, err)
				}
			}
		}
	}
}
