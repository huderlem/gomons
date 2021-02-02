package gen3

import (
	"fmt"
)

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

// GetPokedexSortMode gets the sort mode for the Pokedex.
func (s *SaveData) GetPokedexSortMode() DexSort {
	return DexSort(s.readU8(0x18))
}

// SetPokedexSortMode sets the sort mode for the Pokedex.
func (s *SaveData) SetPokedexSortMode(sortMode DexSort) {
	s.writeU8(uint8(sortMode), 0x18)
}

// GetPokedexRegionMode gets the region mode for the Pokedex.
func (s *SaveData) GetPokedexRegionMode() DexRegion {
	return DexRegion(s.readU8(0x19))
}

// SetPokedexRegionMode sets the region mode for the Pokedex.
func (s *SaveData) SetPokedexRegionMode(mode DexRegion) {
	s.writeU8(uint8(mode), 0x19)
	if mode == DexRegionNational {
		s.writeU8(0xDA, 0x1A)
	}
}

// GetPokedexUnownPersonality gets the personality value used when viewing Unown in the pokedex.
func (s *SaveData) GetPokedexUnownPersonality() uint32 {
	return s.readU32(0x1C)
}

// SetPokedexUnownPersonality sets the personality value used when viewing Unown in the pokedex.
func (s *SaveData) SetPokedexUnownPersonality(personality uint32) {
	s.writeU32(personality, 0x1C)
}

// GetPokedexSpindaPersonality gets the personality value used when viewing Spinda in the pokedex.
func (s *SaveData) GetPokedexSpindaPersonality() uint32 {
	return s.readU32(0x20)
}

// SetPokedexSpindaPersonality sets the personality value used when viewing Spinda in the pokedex.
func (s *SaveData) SetPokedexSpindaPersonality(personality uint32) {
	s.writeU32(personality, 0x20)
}

// GetNumOwnedSpecies gets the total number of owned pokemon in the national Pokedex.
func (s *SaveData) GetNumOwnedSpecies() int {
	count := 0
	for _, nationalDexIndex := range NationalDexOrder {
		nationalDexIndex--
		if s.readBit(0x28, uint(nationalDexIndex)) {
			count++
		}
	}
	return count
}

// GetNumSeenSpecies gets the total number of seen pokemon in the national Pokedex.
func (s *SaveData) GetNumSeenSpecies() int {
	count := 0
	for _, nationalDexIndex := range NationalDexOrder {
		nationalDexIndex--
		if s.readBit(0x5C, uint(nationalDexIndex)) {
			count++
		}
	}
	return count
}

// GetOwnedSpecies gets whether or not the given species is owned in the Pokedex.
func (s *SaveData) GetOwnedSpecies(species int) (bool, error) {
	nationalDexIndex, ok := NationalDexOrder[species]
	if !ok {
		return false, fmt.Errorf("Invalid species '%d'. It isn't part of the national Pokedex", species)
	}
	nationalDexIndex--
	isOwned := s.readBit(0x28, uint(nationalDexIndex))
	return isOwned, nil
}

// SetOwnedSpecies sets whether or not the given species is owned in the Pokedex.
func (s *SaveData) SetOwnedSpecies(species int, isOwned bool) error {
	nationalDexIndex, ok := NationalDexOrder[species]
	if !ok {
		return fmt.Errorf("Invalid species '%d'. It isn't part of the national Pokedex", species)
	}
	nationalDexIndex--
	s.writeBit(isOwned, 0x28, uint(nationalDexIndex))
	return nil
}

// GetSeenSpecies gets whether or not the given species is seen in the Pokedex.
func (s *SaveData) GetSeenSpecies(species int) (bool, error) {
	nationalDexIndex, ok := NationalDexOrder[species]
	if !ok {
		return false, fmt.Errorf("Invalid species '%d'. It isn't part of the national Pokedex", species)
	}
	nationalDexIndex--
	isSeen := s.readBit(0x5C, uint(nationalDexIndex))
	return isSeen, nil
}

// SetSeenSpecies sets whether or not the given species is seen in the Pokedex.
// If parameter isSeen is false, then also mark the species as not owned.
func (s *SaveData) SetSeenSpecies(species int, isSeen bool) error {
	nationalDexIndex, ok := NationalDexOrder[species]
	if !ok {
		return fmt.Errorf("Invalid species '%d'. It isn't part of the national Pokedex", species)
	}
	nationalDexIndex--
	// The "seen" data is duplicated three times in the save file, presumably
	// for anti-cheat reasons.
	if isSeen {
		s.writeBit(true, 0x5C, uint(nationalDexIndex))
		s.writeBit(true, 0x1908, uint(nationalDexIndex))
		s.writeBit(true, 0x3B24, uint(nationalDexIndex))
	} else {
		if err := s.SetOwnedSpecies(species, false); err != nil {
			return err
		}
		s.writeBit(false, 0x5C, uint(nationalDexIndex))
		s.writeBit(false, 0x1908, uint(nationalDexIndex))
		s.writeBit(false, 0x3B24, uint(nationalDexIndex))
	}
	return nil
}
