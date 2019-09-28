package gen3

// DexSort is how the pokedex entries are ordered.
type DexSort byte

// DexSort values
const (
	DexSortNumeric          DexSort = 0
	DexSortAlphabetic       DexSort = 1
	DexSortWeightDescending DexSort = 2
	DexSortWeightAscending  DexSort = 3
	DexSortHeightDescending DexSort = 4
	DexSortHeightAscending  DexSort = 5
)

func (mode DexSort) String() string {
	switch mode {
	case DexSortNumeric:
		return "Numerical"
	case DexSortAlphabetic:
		return "A to Z"
	case DexSortWeightDescending:
		return "Heaviest"
	case DexSortWeightAscending:
		return "Lightest"
	case DexSortHeightDescending:
		return "Tallest"
	case DexSortHeightAscending:
		return "Smallest"
	default:
		return "Invalid"
	}
}

// DexRegion is the region display mode of the pokedex.
type DexRegion byte

// DexRegion values
const (
	DexRegionHoenn    DexRegion = 0
	DexRegionNational DexRegion = 1
)

func (mode DexRegion) String() string {
	switch mode {
	case DexRegionHoenn:
		return "Hoenn"
	case DexRegionNational:
		return "National"
	default:
		return "Invalid"
	}
}
