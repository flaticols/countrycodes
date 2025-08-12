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
		{"templates/constants_alpha2.tmpl", "constants_alpha2_generated.go"},
		{"templates/constants_alpha3.tmpl", "alpha3/constants_generated.go"},
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
		"upper":          strings.ToUpper,
		"lower":          strings.ToLower,
		"normalizeKey":   normalizeKey,
		"toInt":          countryCodeToInt,
		"packAlpha2":     packAlpha2,
		"packAlpha3":     packAlpha3,
		"toGoIdentifier": toGoIdentifier,
	}).Parse(string(tmplContent)))

	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(config.OutputPath)
	if outputDir != "." && outputDir != "" {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory %s: %w", outputDir, err)
		}
	}

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

// toGoIdentifier converts a country name to a valid Go identifier
func toGoIdentifier(name string) string {
	// Special cases for common country names
	replacements := map[string]string{
		"United States of America":                             "UnitedStates",
		"United Kingdom of Great Britain and Northern Ireland": "UnitedKingdom",
		"Russian Federation":                                   "Russia",
		"Venezuela, Bolivarian Republic of":                    "Venezuela",
		"Bolivia, Plurinational State of":                      "Bolivia",
		"Czechia":                                              "CzechRepublic",
		"Korea, Republic of":                                   "SouthKorea",
		"Korea, Democratic People's Republic of":               "NorthKorea",
		"Lao People's Democratic Republic":                     "Laos",
		"Viet Nam":                                             "Vietnam",
		"Iran, Islamic Republic of":                            "Iran",
		"Syrian Arab Republic":                                 "Syria",
		"Tanzania, United Republic of":                         "Tanzania",
		"Moldova, Republic of":                                 "Moldova",
		"Netherlands, Kingdom of the":                          "Netherlands",
		"Palestine, State of":                                  "Palestine",
		"Holy See":                                             "VaticanCity",
		"Micronesia, Federated States of":                      "Micronesia",
		"Congo, Democratic Republic of the":                    "CongoDRC",
		"Congo":                                                "CongoBrazzaville",
		"Guinea-Bissau":                                        "GuineaBissau",
		"Timor-Leste":                                          "TimorLeste",
		"Côte d'Ivoire":                                        "CoteDIvoire",
		"Saint Martin (French part)":                           "SaintMartinFrench",
		"Sint Maarten (Dutch part)":                            "SintMaartenDutch",
		"Virgin Islands (British)":                             "BritishVirginIslands",
		"Virgin Islands (U.S.)":                                "USVirginIslands",
		"Bonaire, Sint Eustatius and Saba":                     "CaribbeanNetherlands",
		"Brunei Darussalam":                                    "Brunei",
		"Cabo Verde":                                           "CapeVerde",
		"Taiwan, Province of China":                            "Taiwan",
		"Hong Kong":                                            "HongKong",
		"Macao":                                                "Macau",
		"Réunion":                                              "Reunion",
		"Curaçao":                                              "Curacao",
		"Saint Barthélemy":                                     "SaintBarthelemy",
		"Cocos (Keeling) Islands":                              "CocosIslands",
		"Falkland Islands (Malvinas)":                          "FalklandIslands",
		"Faroe Islands":                                        "FaroeIslands",
		"French Southern Territories":                          "FrenchSouthernTerritories",
		"Heard Island and McDonald Islands":                    "HeardAndMcDonaldIslands",
		"Saint Helena, Ascension and Tristan da Cunha":         "SaintHelena",
		"Saint Pierre and Miquelon":                            "SaintPierreAndMiquelon",
		"South Georgia and the South Sandwich Islands":         "SouthGeorgiaAndSouthSandwichIslands",
		"Svalbard and Jan Mayen":                               "SvalbardAndJanMayen",
		"Turks and Caicos Islands":                             "TurksAndCaicosIslands",
		"United States Minor Outlying Islands":                 "USMinorOutlyingIslands",
		"Wallis and Futuna":                                    "WallisAndFutuna",
		"Åland Islands":                                        "AlandIslands",
	}

	// Check for special cases first
	if replacement, ok := replacements[name]; ok {
		return replacement
	}

	// Remove parenthetical content
	if idx := strings.Index(name, "("); idx > 0 {
		name = strings.TrimSpace(name[:idx])
	}

	// Split by common separators and process each word
	words := strings.FieldsFunc(name, func(r rune) bool {
		return r == ' ' || r == '-' || r == ',' || r == '\'' || r == '.'
	})

	result := ""
	for _, word := range words {
		// Skip articles and prepositions
		lower := strings.ToLower(word)
		if lower == "the" || lower == "of" || lower == "and" || lower == "da" || lower == "de" || lower == "del" || lower == "la" || lower == "las" || lower == "los" {
			continue
		}

		// Capitalize first letter
		if len(word) > 0 {
			result += strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}

	// Handle special characters
	result = strings.ReplaceAll(result, "ç", "c")
	result = strings.ReplaceAll(result, "é", "e")
	result = strings.ReplaceAll(result, "è", "e")
	result = strings.ReplaceAll(result, "ê", "e")
	result = strings.ReplaceAll(result, "ë", "e")
	result = strings.ReplaceAll(result, "á", "a")
	result = strings.ReplaceAll(result, "à", "a")
	result = strings.ReplaceAll(result, "â", "a")
	result = strings.ReplaceAll(result, "ä", "a")
	result = strings.ReplaceAll(result, "å", "a")
	result = strings.ReplaceAll(result, "ã", "a")
	result = strings.ReplaceAll(result, "ñ", "n")
	result = strings.ReplaceAll(result, "ó", "o")
	result = strings.ReplaceAll(result, "ò", "o")
	result = strings.ReplaceAll(result, "ô", "o")
	result = strings.ReplaceAll(result, "ö", "o")
	result = strings.ReplaceAll(result, "õ", "o")
	result = strings.ReplaceAll(result, "ú", "u")
	result = strings.ReplaceAll(result, "ù", "u")
	result = strings.ReplaceAll(result, "û", "u")
	result = strings.ReplaceAll(result, "ü", "u")
	result = strings.ReplaceAll(result, "í", "i")
	result = strings.ReplaceAll(result, "ì", "i")
	result = strings.ReplaceAll(result, "î", "i")
	result = strings.ReplaceAll(result, "ï", "i")
	result = strings.ReplaceAll(result, "ý", "y")

	// Ensure it starts with a letter
	if len(result) > 0 && result[0] >= '0' && result[0] <= '9' {
		result = "Country" + result
	}

	// If empty, return the alpha2 code
	if result == "" {
		return "UnknownCountry"
	}

	return result
}
