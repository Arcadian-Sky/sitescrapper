package scrapper

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Sitemap struct {
	XMLName xml.Name `xml:"urlset"`
	URLs    []URL    `xml:"url"`
}

type SitemapIndex struct {
	XMLName  xml.Name     `xml:"sitemapindex"`
	Sitemaps []SitemapRef `xml:"sitemap"`
}

type SitemapRef struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod"`
}

type URL struct {
	Loc string `xml:"loc"`
}

type RobotsTxt struct {
	Sitemaps []string
}

func ParseRobotsTxt(domain string) (*RobotsTxt, error) {
	resp, err := http.Get(fmt.Sprintf("https://%s/robots.txt", domain))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var robots RobotsTxt
	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(strings.ToLower(line), "sitemap:") {
			sitemapURL := strings.TrimSpace(strings.TrimPrefix(line, "Sitemap:"))
			robots.Sitemaps = append(robots.Sitemaps, sitemapURL)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &robots, nil
}

func ParseSitemap(sitemapURL string) ([]string, error) {
	resp, err := http.Get(sitemapURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Сначала пробуем распарсить как sitemapindex
	var sitemapIndex SitemapIndex
	if err := xml.NewDecoder(resp.Body).Decode(&sitemapIndex); err == nil && len(sitemapIndex.Sitemaps) > 0 {
		// Если это sitemapindex, рекурсивно парсим каждый sitemap
		var allURLs []string
		for _, sitemapRef := range sitemapIndex.Sitemaps {
			urls, err := ParseSitemap(sitemapRef.Loc)
			if err != nil {
				log.Printf("Ошибка при парсинге вложенного sitemap %s: %v", sitemapRef.Loc, err)
				continue
			}
			allURLs = append(allURLs, urls...)
		}
		return allURLs, nil
	}

	// Если это не sitemapindex, пробуем распарсить как обычный sitemap
	resp, err = http.Get(sitemapURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var sitemap Sitemap
	if err := xml.NewDecoder(resp.Body).Decode(&sitemap); err != nil {
		return nil, fmt.Errorf("не удалось распарсить sitemap: %v", err)
	}

	var urls []string
	for _, url := range sitemap.URLs {
		urls = append(urls, url.Loc)
	}
	return urls, nil
}
