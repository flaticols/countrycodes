package countrycodes

import (
	"testing"
)

func TestNameToAlpha2(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		ok       bool
	}{
		// Normal case
		{"United States of America", "US", true},
		{"United Kingdom of Great Britain and Northern Ireland", "GB", true},
		{"Germany", "DE", true},
		{"Japan", "JP", true},
		// Lowercase
		{"united states of america", "US", true},
		{"united kingdom of great britain and northern ireland", "GB", true},
		{"germany", "DE", true},
		{"japan", "JP", true},
		// Uppercase
		{"UNITED STATES OF AMERICA", "US", true},
		{"GERMANY", "DE", true},
		{"JAPAN", "JP", true},
		// Mixed case
		{"UnItEd StAtEs Of AmErIcA", "US", true},
		{"UNITED kingdom OF great BRITAIN and NORTHERN ireland", "GB", true},
		// Invalid
		{"Invalid Country", "", false},
		{"", "", false},
	}

	for _, tt := range tests {
		result, ok := NameToAlpha2(tt.input)
		if ok != tt.ok || result != tt.expected {
			t.Errorf("NameToAlpha2(%q) = %q, %v; want %q, %v",
				tt.input, result, ok, tt.expected, tt.ok)
		}
	}
}

func TestNameToAlpha3(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		ok       bool
	}{
		// Normal case
		{"United States of America", "USA", true},
		{"United Kingdom of Great Britain and Northern Ireland", "GBR", true},
		{"Germany", "DEU", true},
		{"Japan", "JPN", true},
		// Lowercase
		{"united states of america", "USA", true},
		{"germany", "DEU", true},
		// Uppercase
		{"UNITED STATES OF AMERICA", "USA", true},
		{"GERMANY", "DEU", true},
		// Mixed case
		{"UnItEd StAtEs Of AmErIcA", "USA", true},
		// Invalid
		{"Invalid Country", "", false},
		{"", "", false},
	}

	for _, tt := range tests {
		result, ok := NameToAlpha3(tt.input)
		if ok != tt.ok || result != tt.expected {
			t.Errorf("NameToAlpha3(%q) = %q, %v; want %q, %v",
				tt.input, result, ok, tt.expected, tt.ok)
		}
	}
}

func TestNameToNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		ok       bool
	}{
		// Normal case
		{"United States of America", "840", true},
		{"United Kingdom of Great Britain and Northern Ireland", "826", true},
		{"Germany", "276", true},
		{"Japan", "392", true},
		// Lowercase
		{"united states of america", "840", true},
		{"germany", "276", true},
		// Uppercase
		{"UNITED STATES OF AMERICA", "840", true},
		{"GERMANY", "276", true},
		// Mixed case
		{"UnItEd StAtEs Of AmErIcA", "840", true},
		// Invalid
		{"Invalid Country", "", false},
		{"", "", false},
	}

	for _, tt := range tests {
		result, ok := NameToNumber(tt.input)
		if ok != tt.ok || result != tt.expected {
			t.Errorf("NameToNumber(%q) = %q, %v; want %q, %v",
				tt.input, result, ok, tt.expected, tt.ok)
		}
	}
}

func TestNameToNumberInt(t *testing.T) {
	tests := []struct {
		input    string
		expected int
		ok       bool
	}{
		// Normal case
		{"United States of America", 840, true},
		{"United Kingdom of Great Britain and Northern Ireland", 826, true},
		{"Germany", 276, true},
		{"Japan", 392, true},
		// Lowercase
		{"united states of america", 840, true},
		{"germany", 276, true},
		// Uppercase
		{"UNITED STATES OF AMERICA", 840, true},
		{"GERMANY", 276, true},
		// Mixed case
		{"UnItEd StAtEs Of AmErIcA", 840, true},
		// Invalid
		{"Invalid Country", 0, false},
		{"", 0, false},
	}

	for _, tt := range tests {
		result, ok := NameToNumberInt(tt.input)
		if ok != tt.ok || result != tt.expected {
			t.Errorf("NameToNumberInt(%q) = %d, %v; want %d, %v",
				tt.input, result, ok, tt.expected, tt.ok)
		}
	}
}

func TestGetByName(t *testing.T) {
	tests := []struct {
		input          string
		expectedAlpha2 string
		expectedAlpha3 string
		ok             bool
	}{
		// Normal case
		{"United States of America", "US", "USA", true},
		{"United Kingdom of Great Britain and Northern Ireland", "GB", "GBR", true},
		{"Germany", "DE", "DEU", true},
		// Lowercase
		{"united states of america", "US", "USA", true},
		{"germany", "DE", "DEU", true},
		// Uppercase
		{"UNITED STATES OF AMERICA", "US", "USA", true},
		{"GERMANY", "DE", "DEU", true},
		// Mixed case
		{"UnItEd StAtEs Of AmErIcA", "US", "USA", true},
		// Invalid
		{"Invalid Country", "", "", false},
		{"", "", "", false},
	}

	for _, tt := range tests {
		country, ok := GetByName(tt.input)
		if ok != tt.ok {
			t.Errorf("GetByName(%q) ok = %v; want %v",
				tt.input, ok, tt.ok)
			continue
		}
		if ok {
			if country.Alpha2 != tt.expectedAlpha2 {
				t.Errorf("GetByName(%q).Alpha2 = %q; want %q",
					tt.input, country.Alpha2, tt.expectedAlpha2)
			}
			if country.Alpha3 != tt.expectedAlpha3 {
				t.Errorf("GetByName(%q).Alpha3 = %q; want %q",
					tt.input, country.Alpha3, tt.expectedAlpha3)
			}
		}
	}
}

func TestIsValidName(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		// Normal case
		{"United States of America", true},
		{"United Kingdom of Great Britain and Northern Ireland", true},
		{"Germany", true},
		{"Japan", true},
		// Lowercase
		{"united states of america", true},
		{"germany", true},
		{"japan", true},
		// Uppercase
		{"UNITED STATES OF AMERICA", true},
		{"GERMANY", true},
		{"JAPAN", true},
		// Mixed case
		{"UnItEd StAtEs Of AmErIcA", true},
		{"gErMaNy", true},
		// Invalid
		{"Invalid Country", false},
		{"", false},
		{"USA", false}, // This is an alpha-3 code, not a name
		{"US", false},  // This is an alpha-2 code, not a name
	}

	for _, tt := range tests {
		result := IsValidName(tt.input)
		if result != tt.expected {
			t.Errorf("IsValidName(%q) = %v; want %v",
				tt.input, result, tt.expected)
		}
	}
}

// Benchmark tests for name lookups
func BenchmarkNameToAlpha2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NameToAlpha2("United States of America")
	}
}

func BenchmarkNameToAlpha3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NameToAlpha3("United States of America")
	}
}

func BenchmarkNameToNumberInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NameToNumberInt("United States of America")
	}
}

func BenchmarkGetByName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetByName("United States of America")
	}
}

func BenchmarkIsValidName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsValidName("United States of America")
	}
}

// Benchmark to verify zero allocations
func BenchmarkNameToAlpha2_Allocs(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = NameToAlpha2("United States of America")
	}
}

func BenchmarkGetByName_Allocs(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = GetByName("United States of America")
	}
}
