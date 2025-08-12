package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

type CountryFull struct {
	Name                   string `json:"name"`
	Alpha2                 string `json:"alpha_2"`
	Alpha3                 string `json:"alpha_3"`
	CountryCode            string `json:"country_code"`
	ISO3166_2              string `json:"iso_3166_2"`
	Region                 string `json:"region"`
	SubRegion              string `json:"sub_region"`
	IntermediateRegion     string `json:"intermediate_region"`
	RegionCode             string `json:"region_code"`
	SubRegionCode          string `json:"sub_region_code"`
	IntermediateRegionCode string `json:"intermediate_region_code"`
}

type CountrySlim2 struct {
	Name        string `json:"name"`
	Alpha2      string `json:"alpha_2"`
	CountryCode string `json:"country_code"`
}

type CountrySlim3 struct {
	Name        string `json:"name"`
	Alpha3      string `json:"alpha_3"`
	CountryCode string `json:"country_code"`
}

type GeneratorConfig struct {
	TemplatePath string
	OutputPath   string
}

func main() {
	// Read all JSON files from data directory
	fullData := readJSON[CountryFull]("data/aggregated.json")

	// Sort countries for consistent output
	sort.Slice(fullData, func(i, j int) bool {
		return fullData[i].Alpha2 < fullData[j].Alpha2
	})

	// Define which templates generate which files
	configs := []GeneratorConfig{
		{"templates/types.tmpl", "types_generated.go"},
		{"templates/conversions.tmpl", "conversions_generated.go"},
		{"templates/conversions_int.tmpl", "conversions_int_generated.go"},
		{"templates/names.tmpl", "names_generated.go"},
		{"templates/getters.tmpl", "getters_generated.go"},
		{"templates/validators.tmpl", "validators_generated.go"},
		{"templates/name_lookups_hash.tmpl", "name_lookups_generated.go"},
	}

	// Generate each file
	for _, config := range configs {
		if err := generateFile(config, fullData); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating %s: %v\n", config.OutputPath, err)
			os.Exit(1)
		}
		fmt.Printf("Generated %s\n", config.OutputPath)
	}
}

func readJSON[T any](filename string) []T {
	file, err := os.Open(filename)
	if err != nil {
		panic(fmt.Sprintf("Failed to open %s: %v", filename, err))
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		panic(fmt.Sprintf("Failed to read %s: %v", filename, err))
	}

	var result []T
	if err := json.Unmarshal(data, &result); err != nil {
		panic(fmt.Sprintf("Failed to parse %s: %v", filename, err))
	}

	return result
}

func generateFile(config GeneratorConfig, countries []CountryFull) error {
	// Read template file
	tmplContent, err := os.ReadFile(config.TemplatePath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", config.TemplatePath, err)
	}

	// Parse template
	tmpl := template.Must(template.New(filepath.Base(config.TemplatePath)).Funcs(template.FuncMap{
		"upper":        strings.ToUpper,
		"lower":        strings.ToLower,
		"normalizeKey": normalizeKey,
		"toInt":        countryCodeToInt,
		"packAlpha2":   packAlpha2,
		"packAlpha3":   packAlpha3,
	}).Parse(string(tmplContent)))

	// Create output file
	file, err := os.Create(config.OutputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %w", config.OutputPath, err)
	}
	defer file.Close()

	// Check if this is the hash-based name lookups template
	if strings.Contains(config.TemplatePath, "name_lookups_hash") {
		// Generate hash table data
		data := generateHashTableData(countries)
		return tmpl.Execute(file, data)
	}

	// Execute template with countries directly
	return tmpl.Execute(file, countries)
}

func normalizeKey(s string) string {
	// Normalize country names for lookup (remove special chars, lowercase)
	s = strings.ToUpper(s)
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "-", "")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, "'", "")
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, "(", "")
	s = strings.ReplaceAll(s, ")", "")
	return s
}

// countryCodeToInt converts country code string to integer
func countryCodeToInt(code string) int {
	// Remove leading zeros and convert to int
	code = strings.TrimLeft(code, "0")
	if code == "" {
		return 0
	}
	var result int
	fmt.Sscanf(code, "%d", &result)
	return result
}

// packAlpha2 generates the packed uint16 key for a 2-character code
func packAlpha2(code string) string {
	if len(code) != 2 {
		return "0"
	}
	// Convert to lowercase for the key
	a := strings.ToLower(code)[0]
	b := strings.ToLower(code)[1]
	key := uint16(a) | uint16(b)<<8
	return fmt.Sprintf("0x%04x", key)
}

// packAlpha3 generates the packed uint32 key for a 3-character code
func packAlpha3(code string) string {
	if len(code) != 3 {
		return "0"
	}
	// Convert to lowercase for the key
	a := strings.ToLower(code)[0]
	b := strings.ToLower(code)[1]
	c := strings.ToLower(code)[2]
	key := uint32(a) | uint32(b)<<8 | uint32(c)<<16
	return fmt.Sprintf("0x%06x", key)
}

// FNV-1a hash with case folding for country names
func hashCI(s string) uint64 {
	var h uint64 = 14695981039346656037 // FNV-1a offset basis
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		h ^= uint64(c)
		h *= 1099511628211 // FNV-1a prime
	}
	return h
}

// HashSlot represents a slot in the hash table
type HashSlot struct {
	Hash  string // Hash value as string (for template)
	Index string // Index as string (for template)
}

// HashTableData contains the data for hash table generation
type HashTableData struct {
	Countries []CountryFull
	HashSlots []HashSlot
}

// generateHashTableData generates hash table data for name lookups
func generateHashTableData(countries []CountryFull) HashTableData {
	const tableSize = 512
	const tableMask = tableSize - 1

	// Initialize hash table
	hashSlots := make([]HashSlot, tableSize)
	for i := range hashSlots {
		hashSlots[i] = HashSlot{
			Hash:  "0",
			Index: "0xFFFF",
		}
	}

	// Insert countries into hash table using linear probing
	for idx, country := range countries {
		h := hashCI(country.Name)
		i := int(h & tableMask)

		// Linear probing to find empty slot
		for steps := 0; steps < tableSize; steps++ {
			if hashSlots[i].Index == "0xFFFF" {
				// Found empty slot
				hashSlots[i].Hash = fmt.Sprintf("%d", h)
				hashSlots[i].Index = fmt.Sprintf("%d", idx)
				break
			}
			i = (i + 1) & tableMask
		}
	}

	return HashTableData{
		Countries: countries,
		HashSlots: hashSlots,
	}
}
