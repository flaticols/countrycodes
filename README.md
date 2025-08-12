# countrycodes

Ultra-fast, zero-allocation, zero-dependency ISO 3166-1 country code library for Go. All lookups are compile-time optimized with O(1) complexity and absolutely no memory allocations.

## Features

- **Zero dependencies** - Absolutely no external packages in go.mod
- **Zero allocations** - All operations are allocation-free
- **Blazing fast** - ~0.3-8 nanoseconds per lookup
- **Compile-time optimized** - Uses switch statements and lookup arrays
- **Type safe** - All conversions verified at compile time
- **Case-insensitive** - All lookups accept any case and return proper ISO format
- **Name-based lookups** - Get country codes from country names
- **Complete ISO 3166-1 coverage** - All 249 countries included
- **Regional data** - UN M.49 regions and sub-regions included

## Installation

```bash
go get github.com/flaticols/countrycodes
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/flaticols/countrycodes"
)

func main() {
    // Convert between alpha-2 and alpha-3 codes (case-insensitive)
    alpha3, ok := countrycodes.Alpha2ToAlpha3("fr")  // lowercase works!
    fmt.Println(alpha3, ok) // FRA true

    alpha2, ok := countrycodes.Alpha3ToAlpha2("deu")  // lowercase works!
    fmt.Println(alpha2, ok) // DE true

    // Convert to/from UN M.49 numeric codes (string)
    number, ok := countrycodes.Alpha2ToNumber("IT")
    fmt.Println(number, ok) // 380 true

    alpha2, ok = countrycodes.NumberToAlpha2("528")
    fmt.Println(alpha2, ok) // NL true

    // Convert using integer country codes (faster!)
    numInt, ok := countrycodes.Alpha2ToNumberInt("ES")
    fmt.Println(numInt, ok) // 724 true

    alpha2, ok = countrycodes.NumberIntToAlpha2(56)
    fmt.Println(alpha2, ok) // BE true

    // Get country names (case-insensitive)
    name, ok := countrycodes.Alpha2ToName("pl")  // lowercase works!
    fmt.Println(name, ok) // Poland true

    // Get country codes from names (case-insensitive)
    alpha2, ok = countrycodes.NameToAlpha2("netherlands, kingdom of the")
    fmt.Println(alpha2, ok) // NL true

    alpha3, ok = countrycodes.NameToAlpha3("ITALY")  // uppercase works!
    fmt.Println(alpha3, ok) // ITA true

    num, ok := countrycodes.NameToNumber("Germany")
    fmt.Println(num, ok) // 276 true

    // Get complete country information (case-insensitive)
    country, ok := countrycodes.GetByAlpha2("se")  // lowercase works!
    fmt.Println(country.Name)          // Sweden
    fmt.Println(country.Alpha3)        // SWE
    fmt.Println(country.CountryCode)   // 752
    fmt.Println(country.ISO31662)      // ISO 3166-2:SE
    fmt.Println(country.Region)        // Europe
    fmt.Println(country.SubRegion)     // Northern Europe

    // Get country by name (case-insensitive)
    country, ok = countrycodes.GetByName("switzerland")
    fmt.Println(country.Alpha2)        // CH
    fmt.Println(country.Alpha3)        // CHE

    // Country with intermediate region
    country, ok = countrycodes.GetByAlpha2("MX")
    fmt.Println(country.IntermediateRegion) // Central America

    // Validate codes (case-insensitive)
    fmt.Println(countrycodes.IsValidAlpha2("dk"))  // true (lowercase works!)
    fmt.Println(countrycodes.IsValidAlpha3("xxx")) // false
    fmt.Println(countrycodes.IsValidName("spain")) // true (case-insensitive!)
}
```

## Performance

Benchmark results on Apple M1 Pro:

```
BenchmarkAlpha2ToNumberInt-10           131171708                9.192 ns/op           0 B/op          0 allocs/op
BenchmarkNumberIntToAlpha2-10           1000000000               0.3152 ns/op          0 B/op          0 allocs/op
BenchmarkGetByNumberInt-10              158944388                7.537 ns/op           0 B/op          0 allocs/op
BenchmarkIsValidNumberInt-10            1000000000               0.3146 ns/op          0 B/op          0 allocs/op
BenchmarkNumberToAlpha2_String-10       100000000               15.17 ns/op            0 B/op          0 allocs/op
BenchmarkNumberToAlpha2_Int-10          1000000000               0.3128 ns/op          0 B/op          0 allocs/op
BenchmarkGetByNumber_String-10          153581832                7.538 ns/op           0 B/op          0 allocs/op
BenchmarkGetByNumber_Int-10             161488239                7.470 ns/op           0 B/op          0 allocs/op
BenchmarkAlpha2ToAlpha3-10              174179803                6.867 ns/op           0 B/op          0 allocs/op
BenchmarkAlpha3ToAlpha2-10              189879138                6.270 ns/op           0 B/op          0 allocs/op
BenchmarkNumberToAlpha2-10              87254749                13.79 ns/op            0 B/op          0 allocs/op
BenchmarkGetByAlpha2-10                 167205633                7.179 ns/op           0 B/op          0 allocs/op
BenchmarkIsValidAlpha2-10               174320918                6.855 ns/op           0 B/op          0 allocs/op
BenchmarkAlpha2ToAlpha3_Allocs-10       174757227                6.857 ns/op           0 B/op          0 allocs/op
BenchmarkGetByAlpha2_Allocs-10          84646881                14.07 ns/op            0 B/op          0 allocs/op
BenchmarkNameToAlpha2-10                23234931                51.57 ns/op            0 B/op          0 allocs/op
BenchmarkNameToAlpha3-10                23277316                52.61 ns/op            0 B/op          0 allocs/op
BenchmarkNameToNumberInt-10             23410104                51.28 ns/op            0 B/op          0 allocs/op
BenchmarkGetByName-10                   14574456                81.51 ns/op            0 B/op          0 allocs/op
BenchmarkIsValidName-10                 23539681                50.89 ns/op            0 B/op          0 allocs/op
BenchmarkNameToAlpha2_Allocs-10         23208399                51.67 ns/op            0 B/op          0 allocs/op
BenchmarkGetByName_Allocs-10            14151721                84.46 ns/op            0 B/op          0 allocs/op
```

## Architecture

This library maintains **zero dependencies** through a clever architecture:

- **Main module** (`github.com/flaticols/countrycodes`) - Zero dependencies, contains only the library code
- **Data files** (`data/*.json`) - Pre-committed JSON files with country data
- **Code generator** (`cmd/generator`) - Generates optimized Go code from JSON files
- **Data scrubber** (`cmd/scrubber`) - Separate module for updating data (with its own dependencies)

## Development

```bash
# Show all available commands
make help

# Generate optimized code from data/*.json files
make generate

# Run tests
make test

# Run benchmarks
make bench

# Update country data from Wikipedia/UN (requires network)
make scrub

# Update data and regenerate code
make update

# Clean all generated files
make clean
```

### Updating Country Data

The country data is stored in JSON files under the `data/` directory. To update it:

```bash
# This fetches latest data from Wikipedia and UN sources
make scrub

# Then regenerate the Go code
make generate
```

Note: The scrubber tool is in a separate module (`cmd/scrubber`) to keep the main library dependency-free.

## Data Sources

Country data is sourced from:

- **ISO 3166-1** - Country codes from Wikipedia
- **UN M.49** - Numeric codes and regional classifications from the UN Statistics Division

The data is pre-committed to ensure reproducible builds and zero runtime dependencies.

## License

MIT
