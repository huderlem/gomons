package gen3

import (
	"fmt"
)

// PlayerGender represents the player's in-game gender.
type PlayerGender byte

// PlayerGender values
const (
	Male    PlayerGender = 0
	Female  PlayerGender = 1
	Invalid PlayerGender = 0xFF
)

func (p PlayerGender) String() string {
	switch p {
	case Male:
		return "Male"
	case Female:
		return "Female"
	default:
		return "Invalid Gender"
	}
}

// GetPlayerGender gets the player's gender.
func (s *SaveData) GetPlayerGender() PlayerGender {
	return PlayerGender(s.readU8(0x8))
}

// SetPlayerGender sets the player's gender.
func (s *SaveData) SetPlayerGender(gender PlayerGender) {
	s.writeU8(uint8(gender), 0x8)
}

// GetPlayerName gets the player's OT name.
func (s *SaveData) GetPlayerName() string {
	nameBuffer := make([]uint8, 7)
	s.readU8Slice(nameBuffer, 0x0)
	return readGameString(nameBuffer)
}

// SetPlayerName gets the player's OT name.
func (s *SaveData) SetPlayerName(name string) error {
	if len(name) < 1 || len(name) > 7 {
		return fmt.Errorf("Player name must be between 1 and 7 characters long")
	}
	nameBuffer := make([]uint8, 7)
	if err := writeGameString(nameBuffer, name); err != nil {
		return err
	}
	s.writeU8Slice(nameBuffer, 0x0)
	return nil
}

// GetPlayerTrainerID gets the player's raw trainer id. The trainer id is composed of two parts--public and secret.
// The public id is viewable in-game, while the secret id is not. The trainer id is always used by the game
// engine as its full 4-byte value.
func (s *SaveData) GetPlayerTrainerID() uint32 {
	return s.readU32(0xA)
}

// GetPlayerPublicID gets the player's public trainer id part. This is the id viewable in-game.
func (s *SaveData) GetPlayerPublicID() uint16 {
	return uint16(s.GetPlayerTrainerID())
}

// GetPlayerSecretID gets the player's secret trainer id part. This id is not viewable in-game, but is
// still a part of the player's trainer identity. It is used, for example, when determining a mon's shininess.
func (s *SaveData) GetPlayerSecretID() uint16 {
	otid := s.GetPlayerTrainerID()
	return uint16(otid >> 16)
}

// SetPlayerTrainerID sets the player's raw trainer id. The trainer id is composed of two parts--public and secret.
// The public id is viewable in-game, while the secret id is not. The trainer id is always used by the game
// engine as its full 4-byte value.
func (s *SaveData) SetPlayerTrainerID(id uint32) {
	s.writeU32(id, 0xA)
}

// SetPlayerPublicID sets the player's public trainer id part. This is the id viewable in-game.
func (s *SaveData) SetPlayerPublicID(id uint16) {
	s.writeU16(id, 0xA)
}

// SetPlayerSecretID gets the player's secret trainer id part. This id is not viewable in-game, but is
// still a part of the player's trainer identity. It is used, for example, when determining a mon's shininess.
func (s *SaveData) SetPlayerSecretID(id uint16) {
	s.writeU16(id, 0xC)
}

// GetPlayTime gets the play time hours, seconds, minutes, and vblanks. One VBlank is roughly 1/60 of a second.
func (s *SaveData) GetPlayTime() (uint16, uint8, uint8, uint8, error) {
	hours := s.readU16(0xE)
	minutes := s.readU8(0x10)
	seconds := s.readU8(0x11)
	vblanks := s.readU8(0x12)
	return hours, minutes, seconds, vblanks, nil
}

// SetPlayTime gets the play time hours, seconds, minutes, and vblanks. One VBlank is roughly 1/60 of a second.
func (s *SaveData) SetPlayTime(hours uint16, minutes uint8, seconds uint8, vblanks uint8) error {
	if hours > 999 {
		return fmt.Errorf("Hours must be in range 0-999, but got %d", minutes)
	}
	if minutes > 59 {
		return fmt.Errorf("Minutes must be in range 0-59, but got %d", minutes)
	}
	if seconds > 59 {
		return fmt.Errorf("Seconds must be in range 0-59, but got %d", seconds)
	}
	if vblanks > 59 {
		return fmt.Errorf("VBlanks must be in range 0-59, but got %d", vblanks)
	}
	s.writeU16(hours, 0xE)
	s.writeU8(minutes, 0x10)
	s.writeU8(seconds, 0x11)
	s.writeU8(vblanks, 0x12)
	return nil
}

// GetRegionMapZoomedIn gets whether or not the region map is zoomed in.
func (s *SaveData) GetRegionMapZoomedIn() bool {
	return (s.readU8(0x15) & 0x8) != 0
}

// SetRegionMapZoomedIn sets whether or not the region map is zoomed in.
func (s *SaveData) SetRegionMapZoomedIn(isZoomedIn bool) {
	s.writeBit(isZoomedIn, 0x15, 3)
}

// GetEncryptionKey gets the encryption key used to decrypt various save data.
func (s *SaveData) GetEncryptionKey() uint32 {
	return s.readU32(0xAC)
}

// GetMoney gets the player's current money amount.
func (s *SaveData) GetMoney() uint32 {
	encryptedMoney := s.readU32(0x1410)
	key := s.GetEncryptionKey()
	return encryptedMoney ^ key
}

// SetMoney sets the player's current money amount.
func (s *SaveData) SetMoney(money uint32) {
	key := s.GetEncryptionKey()
	if money > 999999 {
		money = 999999
	}
	encryptedMoney := money ^ key
	s.writeU32(encryptedMoney, 0x1410)
}

// GetCoins gets the player's current coins amount.
func (s *SaveData) GetCoins() uint16 {
	encryptedCoins := s.readU16(0x1414)
	key := s.GetEncryptionKey()
	return encryptedCoins ^ uint16(key)
}

// SetCoins sets the player's current coins amount.
func (s *SaveData) SetCoins(coins uint16) {
	key := s.GetEncryptionKey()
	if coins > 9999 {
		coins = 9999
	}
	encryptedCoins := coins ^ uint16(key)
	s.writeU16(encryptedCoins, 0x1414)
}

// GetPlayerCoordinates gets the player's current x/y coordinates.
func (s *SaveData) GetPlayerCoordinates() (int16, int16) {
	x := s.readS16(0xF80)
	y := s.readS16(0xF82)
	return x, y
}
