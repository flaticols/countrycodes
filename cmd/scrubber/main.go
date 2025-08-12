package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

const (
	wikipediaURI = "https://en.wikipedia.org/wiki/ISO_3166-1"
	unURI        = "https://unstats.un.org/unsd/methodology/m49/overview"
)

// Country represents the complete country data
type Country struct {
	Name                   string `json:"name"`
	Alpha2                 string `json:"alpha_2"`
	Alpha3                 string `json:"alpha_3"`
	CountryCode            string `json:"country_code"`
	ISO31662               string `json:"iso_3166_2"`
	Region                 string `json:"region,omitempty"`
	SubRegion              string `json:"sub_region,omitempty"`
	IntermediateRegion     string `json:"intermediate_region,omitempty"`
	RegionCode             string `json:"region_code,omitempty"`
	SubRegionCode          string `json:"sub_region_code,omitempty"`
	IntermediateRegionCode string `json:"intermediate_region_code,omitempty"`
}

// SlimCountry2 represents the slim-2 format
type SlimCountry2 struct {
	Name        string `json:"name"`
	Alpha2      string `json:"alpha_2"`
	CountryCode string `json:"country_code"`
}

// SlimCountry3 represents the slim-3 format
type SlimCountry3 struct {
	Name        string `json:"name"`
	Alpha3      string `json:"alpha_3"`
	CountryCode string `json:"country_code"`
}

func main() {
	fmt.Println("Fetching data from Wikipedia", wikipediaURI, "...")

	data, err := fetchWikipediaData()
	if err != nil {
		log.Fatal("Error fetching Wikipedia data:", err)
	}

	fmt.Printf("  Data for %d countries found\n", len(data))
	fmt.Println("Fetching data from UN", unURI, "...")

	if err := enrichWithUNData(data); err != nil {
		log.Fatal("Error fetching UN data:", err)
	}

	// Check for missing data
	var blanks []Country
	for _, country := range data {
		if country.Alpha3 == "" || country.Region == "" || country.SubRegion == "" ||
			country.RegionCode == "" || country.SubRegionCode == "" {
			blanks = append(blanks, country)
		}
	}

	if len(blanks) > 0 {
		fmt.Println()
		fmt.Printf("There is missing data for %d countries\n", len(blanks))
		fmt.Printf("You may want to manually check %s\n", unURI)
		fmt.Println()
		for _, b := range blanks {
			fmt.Printf("%+v\n", b)
		}
	}

	fmt.Println()
	fmt.Println("Writing files...")

	// Create data directory if it doesn't exist (relative to project root)
	dataDir := "../../data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatal("Error creating data directory:", err)
	}

	// Write all data
	if err := writeAllData(data); err != nil {
		log.Fatal("Error writing all data:", err)
	}

	// Write slim-2 data
	if err := writeSlim2Data(data); err != nil {
		log.Fatal("Error writing slim-2 data:", err)
	}

	// Write slim-3 data
	if err := writeSlim3Data(data); err != nil {
		log.Fatal("Error writing slim-3 data:", err)
	}

	// Write last updated timestamp
	if err := os.WriteFile("../../data/LAST_UPDATED.txt", []byte(time.Now().String()), 0644); err != nil {
		log.Fatal("Error writing LAST_UPDATED.txt:", err)
	}

	fmt.Println("Done")
}

func fetchWikipediaData() ([]Country, error) {
	resp, err := http.Get(wikipediaURI)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println("Extracting data from Wikipedia")

	var data []Country

	doc.Find("table.sortable tr").Each(func(i int, row *goquery.Selection) {
		cells := row.Find("td")
		if cells.Length() < 5 {
			return
		}

		var country Country

		// Extract country name (remove flag icons)
		nameCell := cells.Eq(0)
		nameCell.Find("span.flagicon").Remove()
		if link := nameCell.Find("a").First(); link.Length() > 0 {
			country.Name = strings.TrimSpace(decodeHTML(link.Text()))
		}

		// Extract ISO codes
		if span := cells.Eq(1).Find("span").First(); span.Length() > 0 {
			country.Alpha2 = strings.TrimSpace(decodeHTML(span.Text()))
		}

		if span := cells.Eq(2).Find("span").First(); span.Length() > 0 {
			country.Alpha3 = strings.TrimSpace(decodeHTML(span.Text()))
		}

		if span := cells.Eq(3).Find("span").First(); span.Length() > 0 {
			country.CountryCode = strings.TrimSpace(decodeHTML(span.Text()))
		}

		if link := cells.Eq(4).Find("a").First(); link.Length() > 0 {
			country.ISO31662 = strings.TrimSpace(decodeHTML(link.Text()))
		}

		// Only add if we have all required fields
		if country.Name != "" && country.Alpha2 != "" && country.Alpha3 != "" &&
			country.CountryCode != "" && country.ISO31662 != "" {
			data = append(data, country)
		}
	})

	return data, nil
}

func enrichWithUNData(data []Country) error {
	resp, err := http.Get(unURI)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("Extracting data from UN")

	// Create a map for faster lookups
	countryMap := make(map[string]*Country)
	for i := range data {
		countryMap[data[i].Alpha2] = &data[i]
	}

	doc.Find("table#downloadTableEN tbody tr").Each(func(i int, row *goquery.Selection) {
		cells := row.Find("td")
		if cells.Length() < 11 {
			return
		}

		// Extract fields from UN data
		regionCode := strings.TrimSpace(cells.Eq(2).Text())
		regionName := strings.TrimSpace(cells.Eq(3).Text())
		subRegionCode := strings.TrimSpace(cells.Eq(4).Text())
		subRegionName := strings.TrimSpace(cells.Eq(5).Text())
		intermediateRegionCode := strings.TrimSpace(cells.Eq(6).Text())
		intermediateRegionName := strings.TrimSpace(cells.Eq(7).Text())
		countryName := strings.TrimSpace(cells.Eq(8).Text())
		alpha2 := strings.TrimSpace(cells.Eq(10).Text())

		// Find matching country and update it
		if country, exists := countryMap[alpha2]; exists {
			country.Region = decodeHTML(regionName)
			country.SubRegion = decodeHTML(subRegionName)
			country.IntermediateRegion = decodeHTML(intermediateRegionName)
			country.RegionCode = regionCode
			country.SubRegionCode = subRegionCode
			country.IntermediateRegionCode = intermediateRegionCode
		} else if countryName != "" && alpha2 != "" {
			fmt.Printf("  %s found in UN source but not in Wikipedia source\n", decodeHTML(countryName))
		}
	})

	return nil
}

func decodeHTML(s string) string {
	return html.UnescapeString(s)
}

func writeAllData(data []Country) error {
	// Write JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("../../data/aggregated.json", jsonData, 0644)
}

func writeSlim2Data(data []Country) error {
	slim2 := make([]SlimCountry2, len(data))
	for i, c := range data {
		slim2[i] = SlimCountry2{
			Name:        c.Name,
			Alpha2:      c.Alpha2,
			CountryCode: c.CountryCode,
		}
	}

	// Write JSON
	jsonData, err := json.MarshalIndent(slim2, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("../../data/iso3166-alpha2.json", jsonData, 0644)
}

func writeSlim3Data(data []Country) error {
	slim3 := make([]SlimCountry3, len(data))
	for i, c := range data {
		slim3[i] = SlimCountry3{
			Name:        c.Name,
			Alpha3:      c.Alpha3,
			CountryCode: c.CountryCode,
		}
	}

	// Write JSON
	jsonData, err := json.MarshalIndent(slim3, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("../../data/iso3166-alpha3.json", jsonData, 0644)
}
