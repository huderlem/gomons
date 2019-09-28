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
