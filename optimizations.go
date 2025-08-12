package countrycodes

// lowerASCII converts an ASCII uppercase letter to lowercase without allocations.
// Non-letter bytes are returned unchanged.
//
//go:inline
func lowerASCII(b byte) byte {
	if b >= 'A' && b <= 'Z' {
		return b + ('a' - 'A')
	}
	return b
}

// lowerASCIIBranchless converts ASCII uppercase to lowercase using branchless operations.
// This can be faster on modern CPUs by avoiding branch mispredictions.
//
//go:inline
func lowerASCIIBranchless(b byte) byte {
	// If b is between 'A' and 'Z', we want to add 0x20
	// mask will be 0xFF if uppercase, 0x00 otherwise
	isUpper := ((b - 'A') | ('Z' - b)) >> 31 // Sign bit propagation
	return b | (isUpper & 0x20)
}

// upperASCII converts an ASCII lowercase letter to uppercase without allocations.
// Non-letter bytes are returned unchanged.
//
//go:inline
func upperASCII(b byte) byte {
	if b >= 'a' && b <= 'z' {
		return b - ('a' - 'A')
	}
	return b
}

// key2 packs a 2-character string into a uint16 for fast map lookups.
// Returns the packed key and true if valid, or 0 and false if invalid.
func key2(s string) (uint16, bool) {
	if len(s) != 2 {
		return 0, false
	}
	a, b := lowerASCII(s[0]), lowerASCII(s[1])
	if a < 'a' || a > 'z' || b < 'a' || b > 'z' {
		return 0, false
	}
	return uint16(a) | uint16(b)<<8, true
}

// key2Inline packs a 2-character string into a uint16 inline for maximum performance.
// This version skips validation for use in hot paths where input is known to be valid.
//
//go:inline
func key2Inline(s string) uint16 {
	return uint16(lowerASCII(s[0])) | uint16(lowerASCII(s[1]))<<8
}

// key3 packs a 3-character string into a uint32 for fast map lookups.
// Returns the packed key and true if valid, or 0 and false if invalid.
func key3(s string) (uint32, bool) {
	if len(s) != 3 {
		return 0, false
	}
	a, b, c := lowerASCII(s[0]), lowerASCII(s[1]), lowerASCII(s[2])
	if a < 'a' || a > 'z' || b < 'a' || b > 'z' || c < 'a' || c > 'z' {
		return 0, false
	}
	return uint32(a) | uint32(b)<<8 | uint32(c)<<16, true
}

// key3Inline packs a 3-character string into a uint32 inline for maximum performance.
// This version skips validation for use in hot paths where input is known to be valid.
//
//go:inline
func key3Inline(s string) uint32 {
	return uint32(lowerASCII(s[0])) | uint32(lowerASCII(s[1]))<<8 | uint32(lowerASCII(s[2]))<<16
}

// stringsEqualFoldASCII compares two ASCII strings case-insensitively without allocations.
// This is optimized for country names which only contain ASCII characters.
func stringsEqualFoldASCII(s, t string) bool {
	if len(s) != len(t) {
		return false
	}
	for i := 0; i < len(s); i++ {
		if lowerASCII(s[i]) != lowerASCII(t[i]) {
			return false
		}
	}
	return true
}

// hashCI performs FNV-1a hash on ASCII string with case folding.
// Compact and inlineable for maximum performance.
//
//go:inline
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

// eqCI performs ASCII case-insensitive equality check without allocations.
//
//go:inline
func eqCI(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		aa, bb := a[i], b[i]
		if aa >= 'A' && aa <= 'Z' {
			aa += 'a' - 'A'
		}
		if bb >= 'A' && bb <= 'Z' {
			bb += 'a' - 'A'
		}
		if aa != bb {
			return false
		}
	}
	return true
}

// Pre-allocated strings for numbers 0-999 with leading zeros
var numberStrings [1000]string

func init() {
	// Pre-generate all 3-digit number strings
	for i := 0; i < 1000; i++ {
		h := i / 100
		t := (i / 10) % 10
		o := i % 10
		numberStrings[i] = string([]byte{'0' + byte(h), '0' + byte(t), '0' + byte(o)})
	}
}

// formatNumber formats an integer as a 3-digit string with leading zeros
//
//go:inline
func formatNumber(n int) string {
	if n >= 0 && n < 1000 {
		return numberStrings[n]
	}
	// Fallback for numbers outside the expected range
	return ""
}
