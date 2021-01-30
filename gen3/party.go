package gen3

import (
	"fmt"
)

// GetPartyCount gets the number of pokemon in the player's party.
func (s *SaveData) GetPartyCount() (int, error) {
	section, err := s.getGameSaveSection(1)
	if err != nil {
		return 0, err
	}
	return int(section.data[0x234]), nil
}

// GetPartyMon gets the pokemon at the given index in the player's party.
func (s *SaveData) GetPartyMon(index int) (Pokemon, error) {
	section, err := s.getGameSaveSection(1)
	if err != nil {
		return Pokemon{}, err
	}
	count, err := s.GetPartyCount()
	if err != nil {
		return Pokemon{}, err
	}
	if index < 0 {
		return Pokemon{}, fmt.Errorf("Party index must be a positive integer. Got %d instead", index)
	}
	if index >= count {
		return Pokemon{}, fmt.Errorf("Invalid party index '%d' because there are only %d Pokemon in the player's party", index, count)
	}
	offset := 0x238 + (index * 100)
	return readPokemonData(section.data[offset : offset+100])
}

// SetPartyMon sets the pokemon at the given index in the player's party.
func (s *SaveData) SetPartyMon(mon Pokemon, index int) error {
	section, err := s.getGameSaveSection(1)
	if err != nil {
		return err
	}
	offset := 0x238 + (index * 100)
	return writePokemonData(section.data[offset:offset+100], mon)
}
