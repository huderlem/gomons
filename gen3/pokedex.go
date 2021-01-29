package gen3

import (
	"encoding/binary"
	"fmt"

	"github.com/huderlem/gomons/util"
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
func (s *SaveData) GetPokedexSortMode() (DexSort, error) {
	section, err := s.getGameSaveSection(0)
	if err != nil {
		return DexSortAlphabetic, err
	}
	return DexSort(section.data[0x18]), nil
}

// SetPokedexSortMode sets the sort mode for the Pokedex.
func (s *SaveData) SetPokedexSortMode(sortMode DexSort) error {
	section, err := s.getGameSaveSection(0)
	if err != nil {
		return err
	}
	if sortMode > DexSortHeightAscending {
		return fmt.Errorf("Invalid Pokedex sort mode %d", sortMode)
	}
	section.data[0x18] = byte(sortMode)
	return nil
}

// GetPokedexRegionMode gets the region mode for the Pokedex.
func (s *SaveData) GetPokedexRegionMode() (DexRegion, error) {
	section, err := s.getGameSaveSection(0)
	if err != nil {
		return DexRegionNational, err
	}
	return DexRegion(section.data[0x19]), nil
}

// SetPokedexRegionMode sets the region mode for the Pokedex.
func (s *SaveData) SetPokedexRegionMode(mode DexRegion) error {
	section, err := s.getGameSaveSection(0)
	if err != nil {
		return err
	}
	if mode > DexRegionNational {
		return fmt.Errorf("Invalid Pokedex region mode %d", mode)
	}
	section.data[0x19] = byte(mode)
	if mode == DexRegionNational {
		section.data[0x1A] = 0xDA
	}
	return nil
}

// GetPokedexUnownPersonality gets the personality value used when viewing Unown in the pokedex.
func (s *SaveData) GetPokedexUnownPersonality() (uint32, error) {
	section, err := s.getGameSaveSection(0)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(section.data[0x1C:0x20]), nil
}

// SetPokedexUnownPersonality sets the personality value used when viewing Unown in the pokedex.
func (s *SaveData) SetPokedexUnownPersonality(personality uint32) error {
	section, err := s.getGameSaveSection(0)
	if err != nil {
		return err
	}
	binary.LittleEndian.PutUint32(section.data[0x1C:0x20], personality)
	return nil
}

// GetPokedexSpindaPersonality gets the personality value used when viewing Spinda in the pokedex.
func (s *SaveData) GetPokedexSpindaPersonality() (uint32, error) {
	section, err := s.getGameSaveSection(0)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(section.data[0x20:0x24]), nil
}

// SetPokedexSpindaPersonality sets the personality value used when viewing Spinda in the pokedex.
func (s *SaveData) SetPokedexSpindaPersonality(personality uint32) error {
	section, err := s.getGameSaveSection(0)
	if err != nil {
		return err
	}
	binary.LittleEndian.PutUint32(section.data[0x20:0x24], personality)
	return nil
}

// GetNumOwnedSpecies gets the total number of owned pokemon in the national Pokedex.
func (s *SaveData) GetNumOwnedSpecies() (int, error) {
	section, err := s.getGameSaveSection(0)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, nationalDexIndex := range NationalDexOrder {
		nationalDexIndex--
		if util.CheckBitSetInArray(section.data[0x28:], nationalDexIndex) {
			count++
		}
	}
	return count, nil
}

// GetNumSeenSpecies gets the total number of seen pokemon in the national Pokedex.
func (s *SaveData) GetNumSeenSpecies() (int, error) {
	section, err := s.getGameSaveSection(0)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, nationalDexIndex := range NationalDexOrder {
		nationalDexIndex--
		if util.CheckBitSetInArray(section.data[0x5C:], nationalDexIndex) {
			count++
		}
	}
	return count, nil
}

// GetOwnedSpecies gets whether or not the given species is owned in the Pokedex.
func (s *SaveData) GetOwnedSpecies(species int) (bool, error) {
	section, err := s.getGameSaveSection(0)
	if err != nil {
		return false, err
	}
	nationalDexIndex, ok := NationalDexOrder[species]
	if !ok {
		return false, fmt.Errorf("Invalid species '%d'. It isn't part of the national Pokedex", species)
	}
	nationalDexIndex--
	isOwned := util.CheckBitSetInArray(section.data[0x28:], nationalDexIndex)
	return isOwned, nil
}

// SetOwnedSpecies sets whether or not the given species is owned in the Pokedex.
func (s *SaveData) SetOwnedSpecies(species int, isOwned bool) error {
	section, err := s.getGameSaveSection(0)
	if err != nil {
		return err
	}
	nationalDexIndex, ok := NationalDexOrder[species]
	if !ok {
		return fmt.Errorf("Invalid species '%d'. It isn't part of the national Pokedex", species)
	}
	nationalDexIndex--
	if isOwned {
		util.SetBitInArray(section.data[0x28:], nationalDexIndex)
	} else {
		util.ClearBitInArray(section.data[0x28:], nationalDexIndex)
	}
	return nil
}

// GetSeenSpecies gets whether or not the given species is seen in the Pokedex.
func (s *SaveData) GetSeenSpecies(species int) (bool, error) {
	section, err := s.getGameSaveSection(0)
	if err != nil {
		return false, err
	}
	nationalDexIndex, ok := NationalDexOrder[species]
	if !ok {
		return false, fmt.Errorf("Invalid species '%d'. It isn't part of the national Pokedex", species)
	}
	nationalDexIndex--
	isSeen := util.CheckBitSetInArray(section.data[0x5C:], nationalDexIndex)
	return isSeen, nil
}

// SetSeenSpecies sets whether or not the given species is seen in the Pokedex.
// If parameter isSeen is false, then also mark the species as not owned.
func (s *SaveData) SetSeenSpecies(species int, isSeen bool) error {
	section0, err := s.getGameSaveSection(0)
	if err != nil {
		return err
	}
	section1, err := s.getGameSaveSection(1)
	if err != nil {
		return err
	}
	section4, err := s.getGameSaveSection(4)
	if err != nil {
		return err
	}
	nationalDexIndex, ok := NationalDexOrder[species]
	if !ok {
		return fmt.Errorf("Invalid species '%d'. It isn't part of the national Pokedex", species)
	}
	nationalDexIndex--
	// The "seen" data is duplicated three times in the save file, presumably
	// for anti-cheat reasons.
	if isSeen {
		util.SetBitInArray(section0.data[0x5C:], nationalDexIndex)
		util.SetBitInArray(section1.data[0x988:], nationalDexIndex)
		util.SetBitInArray(section4.data[0xCA4:], nationalDexIndex)
	} else {
		if err := s.SetOwnedSpecies(species, false); err != nil {
			return err
		}
		util.ClearBitInArray(section0.data[0x5C:], nationalDexIndex)
		util.ClearBitInArray(section1.data[0x988:], nationalDexIndex)
		util.ClearBitInArray(section4.data[0xCA4:], nationalDexIndex)
	}
	return nil
}
