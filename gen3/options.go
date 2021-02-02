package gen3

import (
	"fmt"
)

// ButtonMode is the configuration for what buttons can be used like the A button.
type ButtonMode byte

// ButtonMode values
const (
	ButtonModeNormal   ButtonMode = 0
	ButtonModeLR       ButtonMode = 1
	ButtonModeLEqualsA ButtonMode = 2
)

func getButtonModeString(mode ButtonMode) string {
	switch mode {
	case ButtonModeNormal:
		return "Normal"
	case ButtonModeLR:
		return "LR"
	case ButtonModeLEqualsA:
		return "L Equals A"
	default:
		return "Invalid"
	}
}

// TextSpeed is the configuration for how fast text is rendered in-game.
type TextSpeed byte

// TextSpeed values
const (
	TextSpeedSlow   TextSpeed = 0
	TextSpeedMedium TextSpeed = 1
	TextSpeedFast   TextSpeed = 2
)

func getTextSpeedString(speed TextSpeed) string {
	switch speed {
	case TextSpeedSlow:
		return "Slow"
	case TextSpeedMedium:
		return "Medium"
	case TextSpeedFast:
		return "Fast"
	default:
		return "Invalid"
	}
}

// SoundMode is the configuration for how sound is played on the system.
type SoundMode bool

// SoundMode values
const (
	SoundModeMono   SoundMode = false
	SoundModeStereo SoundMode = true
)

func getSoundModeString(mode SoundMode) string {
	switch mode {
	case SoundModeMono:
		return "Mono"
	case SoundModeStereo:
		return "Stereo"
	default:
		return "Invalid"
	}
}

// BattleStyle is the configuration for the battle style.
type BattleStyle bool

// BattleStyle values
const (
	BattleStyleShift BattleStyle = false
	BattleStyleSet   BattleStyle = true
)

func getBattleStyleString(mode BattleStyle) string {
	switch mode {
	case BattleStyleShift:
		return "Shift"
	case BattleStyleSet:
		return "Set"
	default:
		return "Invalid"
	}
}

// Options is the representation of the player's option settings.
type Options struct {
	ButtonMode       ButtonMode
	TextSpeed        TextSpeed
	FrameStyle       uint8
	SoundMode        SoundMode
	BattleStyle      BattleStyle
	BattleAnimations bool
}

func (o Options) String() string {
	return fmt.Sprintf(
		"{ButtonMode{%s} TextSpeed{%s} FrameStyle{%d} SoundMode{%s} BattleStyle{%s} BattleAnimations{%t}}",
		getButtonModeString(o.ButtonMode),
		getTextSpeedString(o.TextSpeed),
		o.FrameStyle,
		getSoundModeString(o.SoundMode),
		getBattleStyleString(o.BattleStyle),
		o.BattleAnimations,
	)
}

// GetOptions gets the player's option settings.
func (s *SaveData) GetOptions() Options {
	options := Options{}
	options.ButtonMode = ButtonMode(s.readU8(0x13))
	options.TextSpeed = TextSpeed(s.readBitsU8(0x14, 0, 3))
	options.FrameStyle = s.readBitsU8(0x14, 3, 5)
	options.SoundMode = SoundMode(s.readBit(0x15, 0))
	options.BattleStyle = BattleStyle(s.readBit(0x15, 1))
	options.BattleAnimations = !s.readBit(0x15, 2)
	return options
}

// SetOptions gets the player's option settings.
func (s *SaveData) SetOptions(options Options) error {
	if options.ButtonMode > ButtonModeLEqualsA {
		return fmt.Errorf("Invalid options button mode %d", options.ButtonMode)
	}
	if options.TextSpeed > TextSpeedFast {
		return fmt.Errorf("Invalid options text speed %d", options.TextSpeed)
	}
	if options.FrameStyle > 19 {
		return fmt.Errorf("Invalid options frame style %d. Must be in range 0-19", options.FrameStyle)
	}
	s.writeU8(byte(options.ButtonMode), 0x13)
	s.writeBitsU8(byte(options.TextSpeed), 0x14, 0, 3)
	s.writeBitsU8(options.FrameStyle, 0x14, 3, 5)
	s.writeBit(bool(options.SoundMode), 0x15, 0)
	s.writeBit(bool(options.BattleStyle), 0x15, 1)
	s.writeBit(!options.BattleAnimations, 0x15, 2)
	return nil
}
