package gen3

import (
	"fmt"
)

// GetNumPartyMons gets the number of pokemon in the player's party.
func (s *SaveData) GetNumPartyMons() (int, error) {
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
	count, err := s.GetNumPartyMons()
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
	count, err := s.GetNumPartyMons()
	if err != nil {
		return err
	}
	if index >= count {
		return fmt.Errorf("Invalid party index '%d' because there are only %d Pokemon in the player's party", index, count)
	}
	offset := 0x238 + (index * 100)
	return writePokemonData(section.data[offset:offset+100], mon)
}

// AddPartyMon add the given pokemon to the end of the player's party.
func (s *SaveData) AddPartyMon(mon Pokemon) error {
	section, err := s.getGameSaveSection(1)
	if err != nil {
		return err
	}
	count, err := s.GetNumPartyMons()
	if err != nil {
		return err
	}
	if count == 6 {
		return fmt.Errorf("Cannot add a Pokemon to the player's party because it is already full")
	}
	offset := 0x238 + (count * 100)
	section.data[0x234] = uint8(count + 1)
	return writePokemonData(section.data[offset:offset+100], mon)
}

// RemovePartyMon removes the pokemon at the given index in the player's party, and
// then shifts any following pokemon up in the party.
func (s *SaveData) RemovePartyMon(index int) error {
	section, err := s.getGameSaveSection(1)
	if err != nil {
		return err
	}
	count, err := s.GetNumPartyMons()
	if err != nil {
		return err
	}
	if index >= count {
		return nil
	}
	for index < count-1 {
		offset := 0x238 + (index * 100)
		nextMon, err := readPokemonData(section.data[offset+100 : offset+200])
		if err != nil {
			return err
		}
		writePokemonData(section.data[offset:offset+100], nextMon)
		index++
	}
	offset := 0x238 + (index * 100)
	section.data[0x234] = uint8(count - 1)
	return writePokemonData(section.data[offset:offset+100], Pokemon{})
}
