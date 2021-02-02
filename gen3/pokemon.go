package gen3

import (
	"encoding/binary"

	"github.com/huderlem/gomons/util"
)

// Pokemon represents a fully-formed pokemon instance.
type Pokemon struct {
	Box                 BoxPokemon
	SleepTurnsRemaining uint8
	IsPoisoned          bool
	IsBurned            bool
	IsFrozen            bool
	IsParalyzed         bool
	IsToxicPoisoned     bool
	ToxicTurnCounter    uint8
	Level               uint8
	Mail                uint8
	Hp                  uint16
	MaxHp               uint16
	Attack              uint16
	Defense             uint16
	Speed               uint16
	SpecialAttack       uint16
	SpecialDefense      uint16
}

// BoxPokemon represents the box data of a pokemon instance.
type BoxPokemon struct {
	Personality          uint32
	OtID                 uint32
	Nickname             string
	Language             GameLanguage
	IsBadEgg             bool
	HasSpecies           bool
	IsEgg                bool
	OtName               string
	Markings             uint8
	Checksum             uint16
	Species              uint16
	HeldItem             uint16
	Experience           uint32
	PpBonusesMove1       uint8
	PpBonusesMove2       uint8
	PpBonusesMove3       uint8
	PpBonusesMove4       uint8
	Friendship           uint8
	Moves                []uint16
	Pp                   []uint8
	HpEV                 uint8
	AttackEV             uint8
	DefenseEV            uint8
	SpeedEV              uint8
	SpecialAttackEV      uint8
	SpecialDefenseEV     uint8
	Cool                 uint8
	Beauty               uint8
	Cute                 uint8
	Smart                uint8
	Tough                uint8
	Sheen                uint8
	Pokerus              uint8
	MetLocation          uint8
	MetLevel             uint8
	MetGame              GameCode
	PokeballType         uint8
	OtGender             PlayerGender
	HpIV                 uint8
	AttackIV             uint8
	DefenseIV            uint8
	SpeedIV              uint8
	SpecialAttackIV      uint8
	SpecialDefenseIV     uint8
	WhichAbility         uint8
	CoolRibbonLevel      uint8
	BeautyRibbonLevel    uint8
	CuteRibbonLevel      uint8
	SmartRibbonLevel     uint8
	ToughRibbonLevel     uint8
	HasChampionRibbon    bool
	HasWinningRibbon     bool
	HasVictoryRibbon     bool
	HasArtistRibbon      bool
	HasEffortRibbon      bool
	HasGiftRibbon1       bool
	HasGiftRibbon2       bool
	HasGiftRibbon3       bool
	HasGiftRibbon4       bool
	HasGiftRibbon5       bool
	HasGiftRibbon6       bool
	HasGiftRibbon7       bool
	FatefulEncounter     uint8
	FollowObedienceRules bool // When this is false, Deoxys and Mew will always disobey. Anti-cheat measure, presumably.
}

func readPokemonData(data []byte) (Pokemon, error) {
	mon := Pokemon{}
	mon.Box.Personality = binary.LittleEndian.Uint32(data[0x0:0x4])
	mon.Box.OtID = binary.LittleEndian.Uint32(data[0x4:0x8])
	mon.Box.Nickname = readGameString(data[0x8:0x12])
	mon.Box.Language = GameLanguage(data[0x12])
	mon.Box.IsBadEgg = util.CheckBit(data[0x13], 0)
	mon.Box.HasSpecies = util.CheckBit(data[0x13], 1)
	mon.Box.IsEgg = util.CheckBit(data[0x13], 2)
	mon.Box.OtName = readGameString(data[0x14:0x1B])
	mon.Box.Markings = data[0x1B]
	mon.Box.Checksum = binary.LittleEndian.Uint16(data[0x1C:0x1E])
	substructData := decryptSubstructs(data[0x20:0x50], mon.Box.Personality, mon.Box.OtID)
	mon.Box.Species = binary.LittleEndian.Uint16(substructData[0x0:0x2])
	mon.Box.HeldItem = binary.LittleEndian.Uint16(substructData[0x2:0x4])
	mon.Box.Experience = binary.LittleEndian.Uint32(substructData[0x4:0x8])
	mon.Box.PpBonusesMove1 = substructData[0x8] & 0x3
	mon.Box.PpBonusesMove2 = (substructData[0x8] & 0xC) >> 2
	mon.Box.PpBonusesMove3 = (substructData[0x8] & 0x30) >> 4
	mon.Box.PpBonusesMove4 = (substructData[0x8] & 0xC0) >> 6
	mon.Box.Friendship = substructData[0x9]
	mon.Box.Moves = make([]uint16, 4)
	mon.Box.Moves[0] = binary.LittleEndian.Uint16(substructData[0xC:0xE])
	mon.Box.Moves[1] = binary.LittleEndian.Uint16(substructData[0xE:0x10])
	mon.Box.Moves[2] = binary.LittleEndian.Uint16(substructData[0x10:0x12])
	mon.Box.Moves[3] = binary.LittleEndian.Uint16(substructData[0x12:0x14])
	mon.Box.Pp = make([]uint8, 4)
	mon.Box.Pp[0] = substructData[0x14]
	mon.Box.Pp[1] = substructData[0x15]
	mon.Box.Pp[2] = substructData[0x16]
	mon.Box.Pp[3] = substructData[0x17]
	mon.Box.HpEV = substructData[0x18]
	mon.Box.AttackEV = substructData[0x19]
	mon.Box.DefenseEV = substructData[0x1A]
	mon.Box.SpeedEV = substructData[0x1B]
	mon.Box.SpecialAttackEV = substructData[0x1C]
	mon.Box.SpecialDefenseEV = substructData[0x1D]
	mon.Box.Cool = substructData[0x1E]
	mon.Box.Beauty = substructData[0x1F]
	mon.Box.Cute = substructData[0x20]
	mon.Box.Smart = substructData[0x21]
	mon.Box.Tough = substructData[0x22]
	mon.Box.Sheen = substructData[0x23]
	mon.Box.Pokerus = substructData[0x24]
	mon.Box.MetLocation = substructData[0x25]
	mon.Box.MetLevel = substructData[0x26] & 0x7F
	mon.Box.MetGame = GameCode((binary.LittleEndian.Uint16(substructData[0x26:0x28]) & 0x780) >> 7)
	mon.Box.PokeballType = uint8((binary.LittleEndian.Uint16(substructData[0x26:0x28]) & 0x7800) >> 11)
	mon.Box.OtGender = PlayerGender((binary.LittleEndian.Uint16(substructData[0x26:0x28]) & 0x8000) >> 15)
	mon.Box.HpIV = substructData[0x28] & 0x1F
	mon.Box.AttackIV = uint8((binary.LittleEndian.Uint32(substructData[0x28:0x2C]) & 0x3E0) >> 5)
	mon.Box.DefenseIV = uint8((binary.LittleEndian.Uint32(substructData[0x28:0x2C]) & 0x7C00) >> 10)
	mon.Box.SpeedIV = uint8((binary.LittleEndian.Uint32(substructData[0x28:0x2C]) & 0xF8000) >> 15)
	mon.Box.SpecialAttackIV = uint8((binary.LittleEndian.Uint32(substructData[0x28:0x2C]) & 0x1f00000) >> 20)
	mon.Box.SpecialDefenseIV = uint8((binary.LittleEndian.Uint32(substructData[0x28:0x2C]) & 0x3E000000) >> 25)
	mon.Box.WhichAbility = uint8((binary.LittleEndian.Uint32(substructData[0x28:0x2C]) & 0x80000000) >> 31)
	mon.Box.CoolRibbonLevel = substructData[0x2C] & 0x7
	mon.Box.BeautyRibbonLevel = uint8((binary.LittleEndian.Uint32(substructData[0x2C:0x30]) & 0x38) >> 3)
	mon.Box.CuteRibbonLevel = uint8((binary.LittleEndian.Uint32(substructData[0x2C:0x30]) & 0x1C0) >> 6)
	mon.Box.SmartRibbonLevel = uint8((binary.LittleEndian.Uint32(substructData[0x2C:0x30]) & 0xE00) >> 9)
	mon.Box.ToughRibbonLevel = uint8((binary.LittleEndian.Uint32(substructData[0x2C:0x30]) & 0x7000) >> 12)
	mon.Box.HasChampionRibbon = util.CheckBitU32(binary.LittleEndian.Uint32(substructData[0x2C:0x30]), 15)
	mon.Box.HasWinningRibbon = util.CheckBitU32(binary.LittleEndian.Uint32(substructData[0x2C:0x30]), 16)
	mon.Box.HasVictoryRibbon = util.CheckBitU32(binary.LittleEndian.Uint32(substructData[0x2C:0x30]), 17)
	mon.Box.HasArtistRibbon = util.CheckBitU32(binary.LittleEndian.Uint32(substructData[0x2C:0x30]), 18)
	mon.Box.HasEffortRibbon = util.CheckBitU32(binary.LittleEndian.Uint32(substructData[0x2C:0x30]), 19)
	mon.Box.HasGiftRibbon1 = util.CheckBitU32(binary.LittleEndian.Uint32(substructData[0x2C:0x30]), 20)
	mon.Box.HasGiftRibbon2 = util.CheckBitU32(binary.LittleEndian.Uint32(substructData[0x2C:0x30]), 21)
	mon.Box.HasGiftRibbon3 = util.CheckBitU32(binary.LittleEndian.Uint32(substructData[0x2C:0x30]), 22)
	mon.Box.HasGiftRibbon4 = util.CheckBitU32(binary.LittleEndian.Uint32(substructData[0x2C:0x30]), 23)
	mon.Box.HasGiftRibbon5 = util.CheckBitU32(binary.LittleEndian.Uint32(substructData[0x2C:0x30]), 24)
	mon.Box.HasGiftRibbon6 = util.CheckBitU32(binary.LittleEndian.Uint32(substructData[0x2C:0x30]), 25)
	mon.Box.HasGiftRibbon7 = util.CheckBitU32(binary.LittleEndian.Uint32(substructData[0x2C:0x30]), 26)
	mon.Box.FatefulEncounter = uint8((binary.LittleEndian.Uint32(substructData[0x2C:0x30]) & 0x78000000) >> 12)
	mon.Box.FollowObedienceRules = util.CheckBitU32(binary.LittleEndian.Uint32(substructData[0x2C:0x30]), 31)
	mon.SleepTurnsRemaining = uint8(binary.LittleEndian.Uint32(data[0x50:0x54]) & 0x7)
	mon.IsPoisoned = util.CheckBitU32(binary.LittleEndian.Uint32(data[0x50:0x54]), 3)
	mon.IsBurned = util.CheckBitU32(binary.LittleEndian.Uint32(data[0x50:0x54]), 4)
	mon.IsFrozen = util.CheckBitU32(binary.LittleEndian.Uint32(data[0x50:0x54]), 5)
	mon.IsParalyzed = util.CheckBitU32(binary.LittleEndian.Uint32(data[0x50:0x54]), 6)
	mon.IsToxicPoisoned = util.CheckBitU32(binary.LittleEndian.Uint32(data[0x50:0x54]), 7)
	mon.ToxicTurnCounter = uint8((binary.LittleEndian.Uint32(data[0x50:0x54]) & 0xF00) >> 8)
	mon.Level = data[0x54]
	mon.Mail = data[0x55]
	mon.Hp = binary.LittleEndian.Uint16(data[0x56:0x58])
	mon.MaxHp = binary.LittleEndian.Uint16(data[0x58:0x5A])
	mon.Attack = binary.LittleEndian.Uint16(data[0x5A:0x5C])
	mon.Defense = binary.LittleEndian.Uint16(data[0x5C:0x5E])
	mon.Speed = binary.LittleEndian.Uint16(data[0x5E:0x60])
	mon.SpecialAttack = binary.LittleEndian.Uint16(data[0x60:0x62])
	mon.SpecialDefense = binary.LittleEndian.Uint16(data[0x62:0x64])
	return mon, nil
}

func writePokemonData(data []byte, mon Pokemon) error {
	binary.LittleEndian.PutUint32(data[0x0:0x4], mon.Box.Personality)
	binary.LittleEndian.PutUint32(data[0x4:0x8], mon.Box.OtID)
	if err := writeGameString(data[0x8:0x12], mon.Box.Nickname); err != nil {
		return err
	}
	data[0x12] = byte(mon.Box.Language)
	data[0x13] = util.WriteBit(data[0x13], 0, mon.Box.IsBadEgg)
	data[0x13] = util.WriteBit(data[0x13], 1, mon.Box.HasSpecies)
	data[0x13] = util.WriteBit(data[0x13], 2, mon.Box.IsEgg)
	if err := writeGameString(data[0x14:0x1B], mon.Box.OtName); err != nil {
		return err
	}
	data[0x1B] = mon.Box.Markings
	substructData := make([]byte, 0x30)
	binary.LittleEndian.PutUint16(substructData[0x0:0x2], mon.Box.Species)
	binary.LittleEndian.PutUint16(substructData[0x2:0x4], mon.Box.HeldItem)
	binary.LittleEndian.PutUint32(substructData[0x4:0x8], mon.Box.Experience)
	substructData[0x8] = (mon.Box.PpBonusesMove1 & 0x3) |
		((mon.Box.PpBonusesMove2 & 0x3) << 2) |
		((mon.Box.PpBonusesMove3 & 0x3) << 4) |
		((mon.Box.PpBonusesMove4 & 0x3) << 6)
	substructData[0x9] = mon.Box.Friendship
	for i := 0; i < 4; i++ {
		offset := 0xC + i*2
		if i < len(mon.Box.Moves) {
			binary.LittleEndian.PutUint16(substructData[offset:offset+2], mon.Box.Moves[i])
		} else {
			binary.LittleEndian.PutUint16(substructData[offset:offset+2], 0)
		}
	}
	for i := 0; i < 4; i++ {
		offset := 0x14 + i
		if i < len(mon.Box.Moves) {
			substructData[offset] = mon.Box.Pp[i]
		} else {
			substructData[offset] = 0
		}
	}
	substructData[0x18] = mon.Box.HpEV
	substructData[0x19] = mon.Box.AttackEV
	substructData[0x1A] = mon.Box.DefenseEV
	substructData[0x1B] = mon.Box.SpeedEV
	substructData[0x1C] = mon.Box.SpecialAttackEV
	substructData[0x1D] = mon.Box.SpecialDefenseEV
	substructData[0x1E] = mon.Box.Cool
	substructData[0x1F] = mon.Box.Beauty
	substructData[0x20] = mon.Box.Cute
	substructData[0x21] = mon.Box.Smart
	substructData[0x22] = mon.Box.Tough
	substructData[0x23] = mon.Box.Sheen
	substructData[0x24] = mon.Box.Pokerus
	substructData[0x25] = mon.Box.MetLocation
	substructData[0x26] = util.WriteBits(substructData[0x26], mon.Box.MetLevel, 0, 7)
	binary.LittleEndian.PutUint16(substructData[0x26:0x28], util.WriteBitsU16(binary.LittleEndian.Uint16(substructData[0x26:0x28]), uint16(mon.Box.MetGame), 7, 4))
	binary.LittleEndian.PutUint16(substructData[0x26:0x28], util.WriteBitsU16(binary.LittleEndian.Uint16(substructData[0x26:0x28]), uint16(mon.Box.PokeballType), 11, 4))
	binary.LittleEndian.PutUint16(substructData[0x26:0x28], util.WriteBitsU16(binary.LittleEndian.Uint16(substructData[0x26:0x28]), uint16(mon.Box.OtGender), 15, 1))
	binary.LittleEndian.PutUint32(substructData[0x28:0x2C], util.WriteBitsU32(binary.LittleEndian.Uint32(substructData[0x28:0x2C]), uint32(mon.Box.HpIV), 0, 5))
	binary.LittleEndian.PutUint32(substructData[0x28:0x2C], util.WriteBitsU32(binary.LittleEndian.Uint32(substructData[0x28:0x2C]), uint32(mon.Box.AttackIV), 5, 5))
	binary.LittleEndian.PutUint32(substructData[0x28:0x2C], util.WriteBitsU32(binary.LittleEndian.Uint32(substructData[0x28:0x2C]), uint32(mon.Box.DefenseIV), 10, 5))
	binary.LittleEndian.PutUint32(substructData[0x28:0x2C], util.WriteBitsU32(binary.LittleEndian.Uint32(substructData[0x28:0x2C]), uint32(mon.Box.SpeedIV), 15, 5))
	binary.LittleEndian.PutUint32(substructData[0x28:0x2C], util.WriteBitsU32(binary.LittleEndian.Uint32(substructData[0x28:0x2C]), uint32(mon.Box.SpecialAttackIV), 20, 5))
	binary.LittleEndian.PutUint32(substructData[0x28:0x2C], util.WriteBitsU32(binary.LittleEndian.Uint32(substructData[0x28:0x2C]), uint32(mon.Box.SpecialDefenseIV), 25, 5))
	binary.LittleEndian.PutUint32(substructData[0x28:0x2C], util.WriteBitsU32(binary.LittleEndian.Uint32(substructData[0x28:0x2C]), util.BoolToU32(mon.Box.IsEgg), 30, 1))
	binary.LittleEndian.PutUint32(substructData[0x28:0x2C], util.WriteBitsU32(binary.LittleEndian.Uint32(substructData[0x28:0x2C]), uint32(mon.Box.WhichAbility), 31, 1))
	binary.LittleEndian.PutUint32(substructData[0x2C:0x30], util.WriteBitsU32(binary.LittleEndian.Uint32(substructData[0x2C:0x30]), uint32(mon.Box.CoolRibbonLevel), 0, 3))
	binary.LittleEndian.PutUint32(substructData[0x2C:0x30], util.WriteBitsU32(binary.LittleEndian.Uint32(substructData[0x2C:0x30]), uint32(mon.Box.BeautyRibbonLevel), 3, 3))
	binary.LittleEndian.PutUint32(substructData[0x2C:0x30], util.WriteBitsU32(binary.LittleEndian.Uint32(substructData[0x2C:0x30]), uint32(mon.Box.CuteRibbonLevel), 6, 3))
	binary.LittleEndian.PutUint32(substructData[0x2C:0x30], util.WriteBitsU32(binary.LittleEndian.Uint32(substructData[0x2C:0x30]), uint32(mon.Box.SmartRibbonLevel), 9, 3))
	binary.LittleEndian.PutUint32(substructData[0x2C:0x30], util.WriteBitsU32(binary.LittleEndian.Uint32(substructData[0x2C:0x30]), uint32(mon.Box.ToughRibbonLevel), 12, 3))
	binary.LittleEndian.PutUint32(substructData[0x2C:0x30], util.WriteBitsU32(binary.LittleEndian.Uint32(substructData[0x2C:0x30]), util.BoolToU32(mon.Box.HasChampionRibbon), 15, 1))
	binary.LittleEndian.PutUint32(substructData[0x2C:0x30], util.WriteBitsU32(binary.LittleEndian.Uint32(substructData[0x2C:0x30]), util.BoolToU32(mon.Box.HasWinningRibbon), 16, 1))
	binary.LittleEndian.PutUint32(substructData[0x2C:0x30], util.WriteBitsU32(binary.LittleEndian.Uint32(substructData[0x2C:0x30]), util.BoolToU32(mon.Box.HasVictoryRibbon), 17, 1))
	binary.LittleEndian.PutUint32(substructData[0x2C:0x30], util.WriteBitsU32(binary.LittleEndian.Uint32(substructData[0x2C:0x30]), util.BoolToU32(mon.Box.HasArtistRibbon), 18, 1))
	binary.LittleEndian.PutUint32(substructData[0x2C:0x30], util.WriteBitsU32(binary.LittleEndian.Uint32(substructData[0x2C:0x30]), util.BoolToU32(mon.Box.HasEffortRibbon), 19, 1))
	binary.LittleEndian.PutUint32(substructData[0x2C:0x30], util.WriteBitsU32(binary.LittleEndian.Uint32(substructData[0x2C:0x30]), util.BoolToU32(mon.Box.HasGiftRibbon1), 20, 1))
	binary.LittleEndian.PutUint32(substructData[0x2C:0x30], util.WriteBitsU32(binary.LittleEndian.Uint32(substructData[0x2C:0x30]), util.BoolToU32(mon.Box.HasGiftRibbon2), 21, 1))
	binary.LittleEndian.PutUint32(substructData[0x2C:0x30], util.WriteBitsU32(binary.LittleEndian.Uint32(substructData[0x2C:0x30]), util.BoolToU32(mon.Box.HasGiftRibbon3), 22, 1))
	binary.LittleEndian.PutUint32(substructData[0x2C:0x30], util.WriteBitsU32(binary.LittleEndian.Uint32(substructData[0x2C:0x30]), util.BoolToU32(mon.Box.HasGiftRibbon4), 23, 1))
	binary.LittleEndian.PutUint32(substructData[0x2C:0x30], util.WriteBitsU32(binary.LittleEndian.Uint32(substructData[0x2C:0x30]), util.BoolToU32(mon.Box.HasGiftRibbon5), 24, 1))
	binary.LittleEndian.PutUint32(substructData[0x2C:0x30], util.WriteBitsU32(binary.LittleEndian.Uint32(substructData[0x2C:0x30]), util.BoolToU32(mon.Box.HasGiftRibbon6), 25, 1))
	binary.LittleEndian.PutUint32(substructData[0x2C:0x30], util.WriteBitsU32(binary.LittleEndian.Uint32(substructData[0x2C:0x30]), util.BoolToU32(mon.Box.HasGiftRibbon7), 26, 1))
	binary.LittleEndian.PutUint32(substructData[0x2C:0x30], util.WriteBitsU32(binary.LittleEndian.Uint32(substructData[0x2C:0x30]), uint32(mon.Box.FatefulEncounter), 27, 4))
	binary.LittleEndian.PutUint32(substructData[0x2C:0x30], util.WriteBitsU32(binary.LittleEndian.Uint32(substructData[0x2C:0x30]), util.BoolToU32(mon.Box.FollowObedienceRules), 31, 1))
	encryptSubstructs(substructData, mon.Box.Personality, mon.Box.OtID, data[0x20:0x50])
	binary.LittleEndian.PutUint16(data[0x1C:0x1E], calculateChecksum(substructData))
	binary.LittleEndian.PutUint32(data[0x50:0x54], util.WriteBitsU32(binary.LittleEndian.Uint32(data[0x50:0x54]), uint32(mon.SleepTurnsRemaining), 0, 3))
	binary.LittleEndian.PutUint32(data[0x50:0x54], util.WriteBitsU32(binary.LittleEndian.Uint32(data[0x50:0x54]), util.BoolToU32(mon.IsPoisoned), 3, 1))
	binary.LittleEndian.PutUint32(data[0x50:0x54], util.WriteBitsU32(binary.LittleEndian.Uint32(data[0x50:0x54]), util.BoolToU32(mon.IsBurned), 4, 1))
	binary.LittleEndian.PutUint32(data[0x50:0x54], util.WriteBitsU32(binary.LittleEndian.Uint32(data[0x50:0x54]), util.BoolToU32(mon.IsFrozen), 5, 1))
	binary.LittleEndian.PutUint32(data[0x50:0x54], util.WriteBitsU32(binary.LittleEndian.Uint32(data[0x50:0x54]), util.BoolToU32(mon.IsParalyzed), 6, 1))
	binary.LittleEndian.PutUint32(data[0x50:0x54], util.WriteBitsU32(binary.LittleEndian.Uint32(data[0x50:0x54]), util.BoolToU32(mon.IsToxicPoisoned), 7, 1))
	binary.LittleEndian.PutUint32(data[0x50:0x54], util.WriteBitsU32(binary.LittleEndian.Uint32(data[0x50:0x54]), uint32(mon.ToxicTurnCounter), 8, 4))
	data[0x54] = mon.Level
	data[0x55] = mon.Mail
	binary.LittleEndian.PutUint16(data[0x56:0x58], mon.Hp)
	binary.LittleEndian.PutUint16(data[0x58:0x5A], mon.MaxHp)
	binary.LittleEndian.PutUint16(data[0x5A:0x5C], mon.Attack)
	binary.LittleEndian.PutUint16(data[0x5C:0x5E], mon.Defense)
	binary.LittleEndian.PutUint16(data[0x5E:0x60], mon.Speed)
	binary.LittleEndian.PutUint16(data[0x60:0x62], mon.SpecialAttack)
	binary.LittleEndian.PutUint16(data[0x62:0x64], mon.SpecialDefense)
	return nil
}

var substructPermutations = [][]int{
	{0, 1, 2, 3},
	{0, 1, 3, 2},
	{0, 2, 1, 3},
	{0, 3, 1, 2},
	{0, 2, 3, 1},
	{0, 3, 2, 1},
	{1, 0, 2, 3},
	{1, 0, 3, 2},
	{2, 0, 1, 3},
	{3, 0, 1, 2},
	{2, 0, 3, 1},
	{3, 0, 2, 1},
	{1, 2, 0, 3},
	{1, 3, 0, 2},
	{2, 1, 0, 3},
	{3, 1, 0, 2},
	{2, 3, 0, 1},
	{3, 2, 0, 1},
	{1, 2, 3, 0},
	{1, 3, 2, 0},
	{2, 1, 3, 0},
	{3, 1, 2, 0},
	{2, 3, 1, 0},
	{3, 2, 1, 0},
}

func decryptSubstructs(data []byte, personality uint32, otID uint32) []byte {
	// First, decrypt the raw data.
	decryptedData := make([]byte, len(data))
	for i, v := range data {
		decryptedData[i] = v ^ byte((otID >> ((i % 4) * 8))) ^ byte((personality >> ((i % 4) * 8)))
	}

	// Rearrange the decrypted data in order according to the possible permutations
	// that each of the 4 substructs can be arranged.
	result := make([]byte, len(data))
	permutation := substructPermutations[personality%uint32(len(substructPermutations))]
	for i, offset := range permutation {
		srcOffset := offset * 12
		destOffset := i * 12
		copy(result[destOffset:destOffset+12], decryptedData[srcOffset:srcOffset+12])
	}
	return result
}

func encryptSubstructs(data []byte, personality uint32, otID uint32, destination []byte) {
	// First, arrange the unencrypted data in order according to the possible permutations
	// that each of the 4 substructs can be arranged.
	permutation := substructPermutations[personality%uint32(len(substructPermutations))]
	for i, offset := range permutation {
		srcOffset := i * 12
		destOffset := offset * 12
		copy(destination[destOffset:destOffset+12], data[srcOffset:srcOffset+12])
	}

	// Encrypt the data.
	for i, v := range destination {
		destination[i] = v ^ byte((otID >> ((i % 4) * 8))) ^ byte((personality >> ((i % 4) * 8)))
	}
}

func calculateChecksum(data []byte) uint16 {
	var checksum uint16
	for i := 0; i < len(data); i += 2 {
		checksum += binary.LittleEndian.Uint16(data[i : i+2])
	}
	return checksum
}
