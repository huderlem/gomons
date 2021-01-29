package gen3

import (
	"fmt"

	"github.com/huderlem/gomons/util"
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
type SoundMode byte

// SoundMode values
const (
	SoundModeMono   SoundMode = 0
	SoundModeStereo SoundMode = 1
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
type BattleStyle byte

// BattleStyle values
const (
	BattleStyleShift BattleStyle = 0
	BattleStyleSet   BattleStyle = 1
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
func (s *SaveData) GetOptions() (Options, error) {
	options := Options{}
	section, err := s.getGameSaveSection(0)
	if err != nil {
		return options, err
	}
	options.ButtonMode = ButtonMode(section.data[0x13])
	options.TextSpeed = TextSpeed(section.data[0x14] & 0x7)
	options.FrameStyle = section.data[0x14] >> 3
	options.SoundMode = SoundMode(section.data[0x15] & 0x1)
	options.BattleStyle = BattleStyle((section.data[0x15] & 0x2) >> 1)
	options.BattleAnimations = (section.data[0x15] & 0x4) == 0
	return options, nil
}

// SetOptions gets the player's option settings.
func (s *SaveData) SetOptions(options Options) error {
	section, err := s.getGameSaveSection(0)
	if err != nil {
		return err
	}
	if options.ButtonMode > ButtonModeLEqualsA {
		return fmt.Errorf("Invalid options button mode %d", options.ButtonMode)
	}
	if options.TextSpeed > TextSpeedFast {
		return fmt.Errorf("Invalid options text speed %d", options.TextSpeed)
	}
	if options.FrameStyle > 19 {
		return fmt.Errorf("Invalid options frame style %d. Must be in range 0-19", options.FrameStyle)
	}
	if options.SoundMode > SoundModeStereo {
		return fmt.Errorf("Invalid options sound mode %d. Must be in range 0-1", options.FrameStyle)
	}
	if options.BattleStyle > BattleStyleSet {
		return fmt.Errorf("Invalid options battle style %d. Must be in range 0-1", options.BattleStyle)
	}
	section.data[0x13] = byte(options.ButtonMode)
	section.data[0x14] = (section.data[0x14] &^ 0x7) | byte(options.TextSpeed)
	section.data[0x14] = (section.data[0x14] &^ 0xF8) | (options.FrameStyle << 3)
	section.data[0x15] = (section.data[0x15] &^ 0x1) | byte(options.SoundMode)
	section.data[0x15] = (section.data[0x15] &^ 0x2) | (byte(options.BattleStyle) << 1)
	section.data[0x15] = (section.data[0x15] &^ 0x4) | (util.BoolToByte(!options.BattleAnimations) << 2)
	return nil
}
