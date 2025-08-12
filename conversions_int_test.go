package countrycodes

import (
	"testing"
)

func TestIntegerConversions(t *testing.T) {
	// Test Alpha2ToNumberInt
	num, ok := Alpha2ToNumberInt("US")
	if !ok || num != 840 {
		t.Errorf("Alpha2ToNumberInt(\"US\") = %d, %v; want 840, true", num, ok)
	}

	// Test Alpha3ToNumberInt
	num, ok = Alpha3ToNumberInt("USA")
	if !ok || num != 840 {
		t.Errorf("Alpha3ToNumberInt(\"USA\") = %d, %v; want 840, true", num, ok)
	}

	// Test NumberIntToAlpha2
	alpha2, ok := NumberIntToAlpha2(840)
	if !ok || alpha2 != "US" {
		t.Errorf("NumberIntToAlpha2(840) = %q, %v; want \"US\", true", alpha2, ok)
	}

	// Test NumberIntToAlpha3
	alpha3, ok := NumberIntToAlpha3(840)
	if !ok || alpha3 != "USA" {
		t.Errorf("NumberIntToAlpha3(840) = %q, %v; want \"USA\", true", alpha3, ok)
	}

	// Test NumberIntToName
	name, ok := NumberIntToName(840)
	if !ok || name != "United States of America" {
		t.Errorf("NumberIntToName(840) = %q, %v; want \"United States of America\", true", name, ok)
	}

	// Test GetByNumberInt
	country, ok := GetByNumberInt(276)
	if !ok {
		t.Fatal("GetByNumberInt(276) failed")
	}
	if country.Alpha2 != "DE" {
		t.Errorf("Expected alpha2 \"DE\", got %q", country.Alpha2)
	}
	if country.Name != "Germany" {
		t.Errorf("Expected name \"Germany\", got %q", country.Name)
	}

	// Test IsValidNumberInt
	if !IsValidNumberInt(840) {
		t.Error("IsValidNumberInt(840) should be true")
	}
	if IsValidNumberInt(999) {
		t.Error("IsValidNumberInt(999) should be false")
	}

	// Test edge cases
	_, ok = NumberIntToAlpha2(0)
	if ok {
		t.Error("NumberIntToAlpha2(0) should return false")
	}

	_, ok = NumberIntToAlpha2(-1)
	if ok {
		t.Error("NumberIntToAlpha2(-1) should return false")
	}

	_, ok = NumberIntToAlpha2(99999)
	if ok {
		t.Error("NumberIntToAlpha2(99999) should return false")
	}
}

func BenchmarkAlpha2ToNumberInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Alpha2ToNumberInt("US")
	}
}

func BenchmarkNumberIntToAlpha2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = NumberIntToAlpha2(840)
	}
}

func BenchmarkGetByNumberInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GetByNumberInt(840)
	}
}

func BenchmarkIsValidNumberInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = IsValidNumberInt(840)
	}
}

// Compare string vs int performance
func BenchmarkNumberToAlpha2_String(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = NumberToAlpha2("840")
	}
}

func BenchmarkNumberToAlpha2_Int(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = NumberIntToAlpha2(840)
	}
}

func BenchmarkGetByNumber_String(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GetByNumber("840")
	}
}

func BenchmarkGetByNumber_Int(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GetByNumberInt(840)
	}
}
