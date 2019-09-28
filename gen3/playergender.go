package gen3

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
