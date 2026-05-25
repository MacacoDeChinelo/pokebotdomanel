package commands

var typeChart = map[string]struct {
	weakTo      []string
	resistantTo []string
	immuneTo    []string
}{
	"normal": {
		weakTo:   []string{"fighting"},
		immuneTo: []string{"ghost"},
	},
	"grass": {
		weakTo:      []string{"bug", "fire", "flying", "ice", "poison"},
		resistantTo: []string{"ground", "rock", "water"},
	},
	"fire": {
		weakTo:      []string{"ground", "rock", "water"},
		resistantTo: []string{"bug", "grass", "ice", "steel"},
	},
	"water": {
		weakTo:      []string{"electric", "grass"},
		resistantTo: []string{"fire", "ground", "rock"},
	},
	"electric": {
		weakTo:      []string{"ground"},
		resistantTo: []string{"flying", "water"},
	},
	"flying": {
		weakTo:      []string{"electric", "ice", "rock"},
		resistantTo: []string{"bug", "fighting", "grass"},
		immuneTo:    []string{"ground"},
	},
	"ice": {
		weakTo:      []string{"fighting", "fire", "rock", "steel"},
		resistantTo: []string{"dragon", "flying", "grass", "ground"},
	},
	"rock": {
		weakTo:      []string{"fighting", "grass", "ground", "steel", "water"},
		resistantTo: []string{"bug", "fire", "flying", "ice"},
	},
	"ground": {
		weakTo:      []string{"ice", "grass", "water"},
		resistantTo: []string{"electric", "fire", "poison", "rock", "steel"},
		immuneTo:    []string{"electric"},
	},
	"steel": {
		weakTo:      []string{"fighting", "fire", "ground"},
		resistantTo: []string{"fairy", "ice", "rock"},
		immuneTo:    []string{"poison"},
	},
	"fighting": {
		weakTo:      []string{"fairy", "flying", "psychic"},
		resistantTo: []string{"dark", "ice", "normal", "rock", "steel"},
	},
	"dark": {
		weakTo:      []string{"bug", "fairy", "fighting"},
		resistantTo: []string{"ghost", "psychic"},
	},
	"psychic": {
		weakTo:      []string{"bug", "dark", "ghost"},
		resistantTo: []string{"fighting", "poison"},
	},
	"poison": {
		weakTo:      []string{"ground", "psychic"},
		resistantTo: []string{"fairy", "grass"},
	},
	"bug": {
		weakTo:      []string{"fire", "flying", "rock"},
		resistantTo: []string{"dark", "grass", "psychic"},
	},
	"fairy": {
		weakTo:      []string{"steel", "poison"},
		resistantTo: []string{"dark", "dragon", "fighting"},
		immuneTo:    []string{"dragon"},
	},
	"ghost": {
		weakTo:      []string{"dark", "ghost"},
		resistantTo: []string{"bug", "poison"},
		immuneTo:    []string{"normal", "fighting"},
	},
	"dragon": {
		weakTo:      []string{"dragon", "fairy", "ice"},
		resistantTo: []string{"electric", "fire", "water", "grass"},
	},
}

func checkEffectiveness(attackType string, pokemonTypes []string) string {
	multiplier := 1.0

	for _, pType := range pokemonTypes {
		multiplier *= getMultiplier(attackType, pType)
	}

	switch {
	case multiplier == 0:
		return "imune"
	case multiplier > 0 && multiplier < 0.5:
		return "nadaefetivo"
	case multiplier == 0.5:
		return "poucoefetivo"
	case multiplier == 1:
		return "neutro"
	case multiplier == 2:
		return "superefetivo"
	case multiplier >= 4:
		return "ultraefetivo"
	default:
		return "neutro"
	}
}

func getMultiplier(attack string, defenderType string) float64 {
	data := typeChart[defenderType]

	for _, immune := range data.immuneTo {
		if immune == attack {
			return 0
		}
	}

	for _, weak := range data.weakTo {
		if weak == attack {
			return 2
		}
	}

	for _, resist := range data.resistantTo {
		if resist == attack {
			return 0.5
		}
	}

	return 1
}
