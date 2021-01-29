package gen3

import (
	"encoding/binary"
	"fmt"
)

const maxScriptSize = 995
const magicRAMScriptValue = 51

// SetRAMScript sets the RAM script in the save file. This allows attaching arbitrary scripts to
// arbitrary map event objects.
func (s *SaveData) SetRAMScript(mapGroup byte, mapNum byte, objectID byte, script []byte) error {
	section, err := s.getGameSaveSection(4)
	if err != nil {
		return err
	}
	if len(script) > maxScriptSize {
		return fmt.Errorf("Maximum size of RAM script is %d bytes, but the provided script is %d bytes", maxScriptSize, len(script))
	}
	section.data[0x8AC] = magicRAMScriptValue
	section.data[0x8AD] = mapGroup
	section.data[0x8AE] = mapNum
	section.data[0x8AF] = objectID
	copy(section.data[0x8B0:0x8B0+len(script)], script)
	checksum := calculateCRC16WithTable(section.data[0x8AC:0xC94])
	binary.LittleEndian.PutUint32(section.data[0x8A8:0x8AC], uint32(checksum))
	return nil
}
