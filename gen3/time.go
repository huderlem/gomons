package gen3

import (
	"encoding/binary"
)

// GetLocalTimeOffset gets the save's local time offset, which is used
// with the RTC to determine the game's local time.
func (s *SaveData) GetLocalTimeOffset() (int16, int8, int8, int8, error) {
	section, err := s.getGameSaveSection(0)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	days := int16(binary.LittleEndian.Uint16(section.data[0x98:0x9A]))
	hours := int8(section.data[0x9A])
	minutes := int8(section.data[0x9B])
	seconds := int8(section.data[0x9C])
	return days, hours, minutes, seconds, nil
}

// SetLocalTimeOffset sets the save's local time offset, which is used
// with the RTC to determine the game's local time.
func (s *SaveData) SetLocalTimeOffset(days int16, hours, minutes, seconds int8) error {
	section, err := s.getGameSaveSection(0)
	if err != nil {
		return err
	}
	binary.LittleEndian.PutUint16(section.data[0x98:0x9A], uint16(days))
	section.data[0x9A] = uint8(hours)
	section.data[0x9B] = uint8(minutes)
	section.data[0x9C] = uint8(seconds)
	return nil
}

// GetLastBerryTreeUpdate gets the save's last berry tree update time offset,
// which is compared to the game's local time to determine if the berry trees
// need to be updated.
func (s *SaveData) GetLastBerryTreeUpdate() (int16, int8, int8, int8, error) {
	section, err := s.getGameSaveSection(0)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	days := int16(binary.LittleEndian.Uint16(section.data[0xA0:0xA2]))
	hours := int8(section.data[0xA2])
	minutes := int8(section.data[0xA3])
	seconds := int8(section.data[0xA4])
	return days, hours, minutes, seconds, nil
}

// SetLastBerryTreeUpdate sets the save's last berry tree update time offset,
// which is compared to the game's local time to determine if the berry trees
// need to be updated.
func (s *SaveData) SetLastBerryTreeUpdate(days int16, hours, minutes, seconds int8) error {
	section, err := s.getGameSaveSection(0)
	if err != nil {
		return err
	}
	binary.LittleEndian.PutUint16(section.data[0xA0:0xA2], uint16(days))
	section.data[0xA2] = uint8(hours)
	section.data[0xA3] = uint8(minutes)
	section.data[0xA4] = uint8(seconds)
	return nil
}
