package gen3

import (
	"fmt"
)

const maxScriptSize = 995
const magicRAMScriptValue = 51

// SetRAMScript sets the RAM script in the save file. This allows attaching arbitrary scripts to
// arbitrary map event objects.
func (s *SaveData) SetRAMScript(mapGroup byte, mapNum byte, objectID byte, script []byte) error {
	if len(script) > maxScriptSize {
		return fmt.Errorf("Maximum size of RAM script is %d bytes, but the provided script is %d bytes", maxScriptSize, len(script))
	}
	s.writeU8(magicRAMScriptValue, 0x46AC)
	s.writeU8(mapGroup, 0x46AD)
	s.writeU8(mapNum, 0x46AE)
	s.writeU8(objectID, 0x46AF)
	s.writeU8Slice(script, 0x46B0)
	buffer := make([]uint8, 0x3E8)
	s.readU8Slice(buffer, 0x46AC)
	checksum := calculateCRC16WithTable(buffer)
	s.writeU32(uint32(checksum), 0x46A8)
	return nil
}
