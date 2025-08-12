package countrycodes

import (
	"testing"
)

func TestAlpha2ToAlpha3(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		ok       bool
	}{
		{"US", "USA", true},
		{"GB", "GBR", true},
		{"DE", "DEU", true},
		{"JP", "JPN", true},
		{"XX", "", false},
	}

	for _, tt := range tests {
		result, ok := Alpha2ToAlpha3(tt.input)
		if ok != tt.ok || result != tt.expected {
			t.Errorf("Alpha2ToAlpha3(%q) = %q, %v; want %q, %v",
				tt.input, result, ok, tt.expected, tt.ok)
		}
	}
}

func TestAlpha3ToAlpha2(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		ok       bool
	}{
		{"USA", "US", true},
		{"GBR", "GB", true},
		{"DEU", "DE", true},
		{"JPN", "JP", true},
		{"XXX", "", false},
	}

	for _, tt := range tests {
		result, ok := Alpha3ToAlpha2(tt.input)
		if ok != tt.ok || result != tt.expected {
			t.Errorf("Alpha3ToAlpha2(%q) = %q, %v; want %q, %v",
				tt.input, result, ok, tt.expected, tt.ok)
		}
	}
}

func TestNumberConversions(t *testing.T) {
	// Test Number to Alpha2
	alpha2, ok := NumberToAlpha2("840")
	if !ok || alpha2 != "US" {
		t.Errorf("NumberToAlpha2(\"840\") = %q, %v; want \"US\", true", alpha2, ok)
	}

	// Test Number to Alpha3
	alpha3, ok := NumberToAlpha3("840")
	if !ok || alpha3 != "USA" {
		t.Errorf("NumberToAlpha3(\"840\") = %q, %v; want \"USA\", true", alpha3, ok)
	}

	// Test Alpha2 to Number
	num, ok := Alpha2ToNumber("US")
	if !ok || num != "840" {
		t.Errorf("Alpha2ToNumber(\"US\") = %q, %v; want \"840\", true", num, ok)
	}

	// Test Alpha3 to Number
	num, ok = Alpha3ToNumber("USA")
	if !ok || num != "840" {
		t.Errorf("Alpha3ToNumber(\"USA\") = %q, %v; want \"840\", true", num, ok)
	}
}

func TestGetByCode(t *testing.T) {
	// Test GetByAlpha2
	country, ok := GetByAlpha2("US")
	if !ok {
		t.Fatal("GetByAlpha2(\"US\") failed")
	}
	if country.Name != "United States of America" {
		t.Errorf("Expected country name \"United States of America\", got %q", country.Name)
	}
	if country.Alpha3 != "USA" {
		t.Errorf("Expected alpha3 \"USA\", got %q", country.Alpha3)
	}
	if country.ISO31662 != "ISO 3166-2:US" {
		t.Errorf("Expected ISO31662 \"ISO 3166-2:US\", got %q", country.ISO31662)
	}
	if country.Region != "Americas" {
		t.Errorf("Expected region \"Americas\", got %q", country.Region)
	}

	// Test GetByAlpha3
	country, ok = GetByAlpha3("GBR")
	if !ok {
		t.Fatal("GetByAlpha3(\"GBR\") failed")
	}
	if country.Alpha2 != "GB" {
		t.Errorf("Expected alpha2 \"GB\", got %q", country.Alpha2)
	}
	if country.ISO31662 != "ISO 3166-2:GB" {
		t.Errorf("Expected ISO31662 \"ISO 3166-2:GB\", got %q", country.ISO31662)
	}

	// Test GetByNumber
	country, ok = GetByNumber("276")
	if !ok {
		t.Fatal("GetByNumber(\"276\") failed")
	}
	if country.Alpha2 != "DE" {
		t.Errorf("Expected alpha2 \"DE\", got %q", country.Alpha2)
	}
	if country.ISO31662 != "ISO 3166-2:DE" {
		t.Errorf("Expected ISO31662 \"ISO 3166-2:DE\", got %q", country.ISO31662)
	}

	// Test country with intermediate region (e.g., Angola)
	country, ok = GetByAlpha2("AO")
	if !ok {
		t.Fatal("GetByAlpha2(\"AO\") failed")
	}
	if country.IntermediateRegion != "Middle Africa" {
		t.Errorf("Expected intermediate region \"Middle Africa\", got %q", country.IntermediateRegion)
	}
	if country.IntermediateRegionCode != "017" {
		t.Errorf("Expected intermediate region code \"017\", got %q", country.IntermediateRegionCode)
	}
}

func TestValidation(t *testing.T) {
	// Test valid codes
	if !IsValidAlpha2("US") {
		t.Error("IsValidAlpha2(\"US\") should be true")
	}
	if !IsValidAlpha3("USA") {
		t.Error("IsValidAlpha3(\"USA\") should be true")
	}
	if !IsValidNumber("840") {
		t.Error("IsValidNumber(\"840\") should be true")
	}

	// Test invalid codes
	if IsValidAlpha2("XX") {
		t.Error("IsValidAlpha2(\"XX\") should be false")
	}
	if IsValidAlpha3("XXX") {
		t.Error("IsValidAlpha3(\"XXX\") should be false")
	}
	if IsValidNumber("999") {
		t.Error("IsValidNumber(\"999\") should be false")
	}
}

func BenchmarkAlpha2ToAlpha3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Alpha2ToAlpha3("US")
	}
}

func BenchmarkAlpha3ToAlpha2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Alpha3ToAlpha2("USA")
	}
}

func BenchmarkNumberToAlpha2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NumberToAlpha2("840")
	}
}

func BenchmarkGetByAlpha2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetByAlpha2("US")
	}
}

func BenchmarkIsValidAlpha2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsValidAlpha2("US")
	}
}

// Benchmark to verify zero allocations
func BenchmarkAlpha2ToAlpha3_Allocs(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = Alpha2ToAlpha3("US")
	}
}

func BenchmarkGetByAlpha2_Allocs(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = GetByAlpha2("US")
	}
}
