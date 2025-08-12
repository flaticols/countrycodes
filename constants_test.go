package countrycodes

import (
	"testing"

	"github.com/flaticols/countrycodes/alpha3"
)

func TestAlpha2Constants(t *testing.T) {
	tests := []struct {
		constant string
		expected string
	}{
		{Netherlands, "NL"},
		{France, "FR"},
		{Germany, "DE"},
		{Italy, "IT"},
		{Spain, "ES"},
		{UnitedKingdom, "GB"},
		{UnitedStates, "US"},
	}

	for _, tt := range tests {
		if tt.constant != tt.expected {
			t.Errorf("constant %s = %s, want %s", tt.expected, tt.constant, tt.expected)
		}

		// Test that constants work with conversion functions
		name, ok := Alpha2ToName(tt.constant)
		if !ok {
			t.Errorf("Alpha2ToName(%s) failed", tt.constant)
		}
		if name == "" {
			t.Errorf("Alpha2ToName(%s) returned empty name", tt.constant)
		}

		alpha3Code, ok := Alpha2ToAlpha3(tt.constant)
		if !ok {
			t.Errorf("Alpha2ToAlpha3(%s) failed", tt.constant)
		}
		if alpha3Code == "" {
			t.Errorf("Alpha2ToAlpha3(%s) returned empty code", tt.constant)
		}
	}
}

func TestAlpha3Constants(t *testing.T) {
	tests := []struct {
		constant string
		expected string
	}{
		{alpha3.Netherlands, "NLD"},
		{alpha3.France, "FRA"},
		{alpha3.Germany, "DEU"},
		{alpha3.Italy, "ITA"},
		{alpha3.Spain, "ESP"},
		{alpha3.UnitedKingdom, "GBR"},
		{alpha3.UnitedStates, "USA"},
	}

	for _, tt := range tests {
		if tt.constant != tt.expected {
			t.Errorf("constant %s = %s, want %s", tt.expected, tt.constant, tt.expected)
		}

		// Test that constants work with conversion functions
		name, ok := Alpha3ToName(tt.constant)
		if !ok {
			t.Errorf("Alpha3ToName(%s) failed", tt.constant)
		}
		if name == "" {
			t.Errorf("Alpha3ToName(%s) returned empty name", tt.constant)
		}

		alpha2Code, ok := Alpha3ToAlpha2(tt.constant)
		if !ok {
			t.Errorf("Alpha3ToAlpha2(%s) failed", tt.constant)
		}
		if alpha2Code == "" {
			t.Errorf("Alpha3ToAlpha2(%s) returned empty code", tt.constant)
		}
	}
}

func TestConstantsCompileTimeSafety(t *testing.T) {
	// These tests demonstrate compile-time type safety
	// The constants can be used directly without quotes

	// Alpha2 constants
	name, ok := Alpha2ToName(Netherlands)
	if !ok || name != "Netherlands, Kingdom of the" {
		t.Errorf("Expected Netherlands, got %s", name)
	}

	alpha3Code, ok := Alpha2ToAlpha3(France)
	if !ok || alpha3Code != "FRA" {
		t.Errorf("Expected FRA, got %s", alpha3Code)
	}

	// Alpha3 constants
	name, ok = Alpha3ToName(alpha3.Germany)
	if !ok || name != "Germany" {
		t.Errorf("Expected Germany, got %s", name)
	}

	alpha2Code, ok := Alpha3ToAlpha2(alpha3.Italy)
	if !ok || alpha2Code != "IT" {
		t.Errorf("Expected IT, got %s", alpha2Code)
	}
}
