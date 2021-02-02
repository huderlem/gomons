package gen3

import (
	"fmt"
)

// GetNumPartyMons gets the number of pokemon in the player's party.
func (s *SaveData) GetNumPartyMons() int {
	return int(s.readU32(0x11B4))
}

// GetPartyMon gets the pokemon at the given index in the player's party.
func (s *SaveData) GetPartyMon(index int) (Pokemon, error) {
	count := s.GetNumPartyMons()
	if index < 0 {
		return Pokemon{}, fmt.Errorf("Party index must be a positive integer. Got %d instead", index)
	}
	if index >= count {
		return Pokemon{}, fmt.Errorf("Invalid party index '%d' because there are only %d Pokemon in the player's party", index, count)
	}
	offset := 0x11B8 + uint(index)*100
	buffer := make([]uint8, 100)
	s.readU8Slice(buffer, offset)
	return readPokemonData(buffer)
}

// SetPartyMon sets the pokemon at the given index in the player's party.
func (s *SaveData) SetPartyMon(mon Pokemon, index int) error {
	count := s.GetNumPartyMons()
	if index >= count {
		return fmt.Errorf("Invalid party index '%d' because there are only %d Pokemon in the player's party", index, count)
	}
	buffer := make([]uint8, 100)
	if err := writePokemonData(buffer, mon); err != nil {
		return err
	}
	offset := 0x11B8 + uint(index)*100
	s.writeU8Slice(buffer, offset)
	return nil
}

// AddPartyMon add the given pokemon to the end of the player's party.
func (s *SaveData) AddPartyMon(mon Pokemon) error {
	count := s.GetNumPartyMons()
	if count >= 6 {
		return fmt.Errorf("Cannot add a Pokemon to the player's party because it is already full")
	}
	buffer := make([]uint8, 100)
	if err := writePokemonData(buffer, mon); err != nil {
		return err
	}
	offset := 0x11B8 + uint(count)*100
	s.writeU8Slice(buffer, offset)
	s.writeU8(uint8(count+1), 0x11B4)
	return nil
}

// RemovePartyMon removes the pokemon at the given index in the player's party, and
// then shifts any following pokemon up in the party.
func (s *SaveData) RemovePartyMon(index int) error {
	count := s.GetNumPartyMons()
	if index >= count {
		return nil
	}
	for index < count-1 {
		offset := 0x11B8 + uint(index)*100
		nextMonBuffer := make([]uint8, 100)
		s.readU8Slice(nextMonBuffer, offset+100)
		nextMon, err := readPokemonData(nextMonBuffer)
		if err != nil {
			return err
		}
		buffer := make([]uint8, 100)
		if err := writePokemonData(buffer, nextMon); err != nil {
			return err
		}
		s.writeU8Slice(buffer, offset)
		index++
	}
	buffer := make([]uint8, 100)
	if err := writePokemonData(buffer, Pokemon{}); err != nil {
		return err
	}
	offset := 0x11B8 + uint(index)*100
	s.writeU8Slice(buffer, offset)
	s.writeU8(uint8(count-1), 0x11B4)
	return nil
}
