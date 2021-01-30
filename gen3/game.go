package gen3

// GameCode represents the game code (e.g. Emerald, Fire Red).
type GameCode byte

// Game code enumeration
const (
	GameSapphire   GameCode = 1
	GameRuby       GameCode = 2
	GameEmerald    GameCode = 3
	GameFireRed    GameCode = 4
	GameLeafGreen  GameCode = 5
	GameHeartGold  GameCode = 7
	GameSoulSilver GameCode = 8
	GameDiamond    GameCode = 10
	GamePearl      GameCode = 11
	GamePlatinum   GameCode = 12
	GameGamecube   GameCode = 15
)

func (c GameCode) String() string {
	switch c {
	case GameSapphire:
		return "Sapphire"
	case GameRuby:
		return "Ruby"
	case GameEmerald:
		return "Emerald"
	case GameFireRed:
		return "Fire Red"
	case GameLeafGreen:
		return "Leaf Green"
	case GameHeartGold:
		return "Heart Gold"
	case GameSoulSilver:
		return "Soul Silver"
	case GameDiamond:
		return "Diamond"
	case GamePearl:
		return "Pearl"
	case GamePlatinum:
		return "Platinum"
	case GameGamecube:
		return "Gamecube"
	default:
		return "Invalid Game Code"
	}
}
