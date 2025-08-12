package countrycodes

import (
	"testing"
)

func TestCaseInsensitiveAlpha2ToAlpha3(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		ok       bool
	}{
		// Uppercase (original)
		{"US", "USA", true},
		{"GB", "GBR", true},
		// Lowercase
		{"us", "USA", true},
		{"gb", "GBR", true},
		// Mixed case
		{"Us", "USA", true},
		{"uS", "USA", true},
		{"Gb", "GBR", true},
		{"gB", "GBR", true},
		// Invalid
		{"xx", "", false},
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

func TestCaseInsensitiveAlpha3ToAlpha2(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		ok       bool
	}{
		// Uppercase (original)
		{"USA", "US", true},
		{"GBR", "GB", true},
		// Lowercase
		{"usa", "US", true},
		{"gbr", "GB", true},
		// Mixed case
		{"Usa", "US", true},
		{"UsA", "US", true},
		{"uSa", "US", true},
		{"Gbr", "GB", true},
		{"gBr", "GB", true},
		{"gbR", "GB", true},
		// Invalid
		{"xxx", "", false},
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

func TestCaseInsensitiveAlpha2ToNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		ok       bool
	}{
		{"US", "840", true},
		{"us", "840", true},
		{"Us", "840", true},
		{"uS", "840", true},
		{"DE", "276", true},
		{"de", "276", true},
		{"De", "276", true},
		{"xx", "", false},
		{"XX", "", false},
	}

	for _, tt := range tests {
		result, ok := Alpha2ToNumber(tt.input)
		if ok != tt.ok || result != tt.expected {
			t.Errorf("Alpha2ToNumber(%q) = %q, %v; want %q, %v",
				tt.input, result, ok, tt.expected, tt.ok)
		}
	}
}

func TestCaseInsensitiveAlpha3ToNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		ok       bool
	}{
		{"USA", "840", true},
		{"usa", "840", true},
		{"Usa", "840", true},
		{"UsA", "840", true},
		{"DEU", "276", true},
		{"deu", "276", true},
		{"Deu", "276", true},
		{"xxx", "", false},
		{"XXX", "", false},
	}

	for _, tt := range tests {
		result, ok := Alpha3ToNumber(tt.input)
		if ok != tt.ok || result != tt.expected {
			t.Errorf("Alpha3ToNumber(%q) = %q, %v; want %q, %v",
				tt.input, result, ok, tt.expected, tt.ok)
		}
	}
}

func TestCaseInsensitiveAlpha2ToNumberInt(t *testing.T) {
	tests := []struct {
		input    string
		expected int
		ok       bool
	}{
		{"US", 840, true},
		{"us", 840, true},
		{"Us", 840, true},
		{"GB", 826, true},
		{"gb", 826, true},
		{"Gb", 826, true},
		{"xx", 0, false},
		{"XX", 0, false},
	}

	for _, tt := range tests {
		result, ok := Alpha2ToNumberInt(tt.input)
		if ok != tt.ok || result != tt.expected {
			t.Errorf("Alpha2ToNumberInt(%q) = %d, %v; want %d, %v",
				tt.input, result, ok, tt.expected, tt.ok)
		}
	}
}

func TestCaseInsensitiveAlpha3ToNumberInt(t *testing.T) {
	tests := []struct {
		input    string
		expected int
		ok       bool
	}{
		{"USA", 840, true},
		{"usa", 840, true},
		{"Usa", 840, true},
		{"GBR", 826, true},
		{"gbr", 826, true},
		{"Gbr", 826, true},
		{"xxx", 0, false},
		{"XXX", 0, false},
	}

	for _, tt := range tests {
		result, ok := Alpha3ToNumberInt(tt.input)
		if ok != tt.ok || result != tt.expected {
			t.Errorf("Alpha3ToNumberInt(%q) = %d, %v; want %d, %v",
				tt.input, result, ok, tt.expected, tt.ok)
		}
	}
}

func TestCaseInsensitiveGetByAlpha2(t *testing.T) {
	tests := []struct {
		input          string
		expectedAlpha3 string
		ok             bool
	}{
		{"US", "USA", true},
		{"us", "USA", true},
		{"Us", "USA", true},
		{"GB", "GBR", true},
		{"gb", "GBR", true},
		{"Gb", "GBR", true},
		{"xx", "", false},
		{"XX", "", false},
	}

	for _, tt := range tests {
		country, ok := GetByAlpha2(tt.input)
		if ok != tt.ok {
			t.Errorf("GetByAlpha2(%q) ok = %v; want %v",
				tt.input, ok, tt.ok)
			continue
		}
		if ok && country.Alpha3 != tt.expectedAlpha3 {
			t.Errorf("GetByAlpha2(%q).Alpha3 = %q; want %q",
				tt.input, country.Alpha3, tt.expectedAlpha3)
		}
	}
}

func TestCaseInsensitiveGetByAlpha3(t *testing.T) {
	tests := []struct {
		input          string
		expectedAlpha2 string
		ok             bool
	}{
		{"USA", "US", true},
		{"usa", "US", true},
		{"Usa", "US", true},
		{"GBR", "GB", true},
		{"gbr", "GB", true},
		{"Gbr", "GB", true},
		{"xxx", "", false},
		{"XXX", "", false},
	}

	for _, tt := range tests {
		country, ok := GetByAlpha3(tt.input)
		if ok != tt.ok {
			t.Errorf("GetByAlpha3(%q) ok = %v; want %v",
				tt.input, ok, tt.ok)
			continue
		}
		if ok && country.Alpha2 != tt.expectedAlpha2 {
			t.Errorf("GetByAlpha3(%q).Alpha2 = %q; want %q",
				tt.input, country.Alpha2, tt.expectedAlpha2)
		}
	}
}

func TestCaseInsensitiveAlpha2ToName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		ok       bool
	}{
		{"US", "United States of America", true},
		{"us", "United States of America", true},
		{"Us", "United States of America", true},
		{"GB", "United Kingdom of Great Britain and Northern Ireland", true},
		{"gb", "United Kingdom of Great Britain and Northern Ireland", true},
		{"Gb", "United Kingdom of Great Britain and Northern Ireland", true},
		{"xx", "", false},
		{"XX", "", false},
	}

	for _, tt := range tests {
		result, ok := Alpha2ToName(tt.input)
		if ok != tt.ok || result != tt.expected {
			t.Errorf("Alpha2ToName(%q) = %q, %v; want %q, %v",
				tt.input, result, ok, tt.expected, tt.ok)
		}
	}
}

func TestCaseInsensitiveAlpha3ToName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		ok       bool
	}{
		{"USA", "United States of America", true},
		{"usa", "United States of America", true},
		{"Usa", "United States of America", true},
		{"GBR", "United Kingdom of Great Britain and Northern Ireland", true},
		{"gbr", "United Kingdom of Great Britain and Northern Ireland", true},
		{"Gbr", "United Kingdom of Great Britain and Northern Ireland", true},
		{"xxx", "", false},
		{"XXX", "", false},
	}

	for _, tt := range tests {
		result, ok := Alpha3ToName(tt.input)
		if ok != tt.ok || result != tt.expected {
			t.Errorf("Alpha3ToName(%q) = %q, %v; want %q, %v",
				tt.input, result, ok, tt.expected, tt.ok)
		}
	}
}

func TestCaseInsensitiveIsValidAlpha2(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"US", true},
		{"us", true},
		{"Us", true},
		{"uS", true},
		{"GB", true},
		{"gb", true},
		{"XX", false},
		{"xx", false},
		{"Xx", false},
	}

	for _, tt := range tests {
		result := IsValidAlpha2(tt.input)
		if result != tt.expected {
			t.Errorf("IsValidAlpha2(%q) = %v; want %v",
				tt.input, result, tt.expected)
		}
	}
}

func TestCaseInsensitiveIsValidAlpha3(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"USA", true},
		{"usa", true},
		{"Usa", true},
		{"UsA", true},
		{"uSa", true},
		{"GBR", true},
		{"gbr", true},
		{"XXX", false},
		{"xxx", false},
		{"Xxx", false},
	}

	for _, tt := range tests {
		result := IsValidAlpha3(tt.input)
		if result != tt.expected {
			t.Errorf("IsValidAlpha3(%q) = %v; want %v",
				tt.input, result, tt.expected)
		}
	}
}
