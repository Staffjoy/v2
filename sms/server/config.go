package main

type sendingConfig struct {
	WhitelistOnly bool
	Concurrency   int
	// Country code to national numbers
	Numbers map[int32]string
}

var sendingConfigs = map[string]sendingConfig{
	"development": {
		WhitelistOnly: true,
		Numbers:       map[int32]string{1: "4157693528"},
		Concurrency:   1,
	},
	"staging": {
		WhitelistOnly: true,
		Numbers:       map[int32]string{1: "4152129829"},
		Concurrency:   1,
	},
	"production": {
		WhitelistOnly: false,
		Numbers: map[int32]string{
			1:  "4152129952", // USA
			44: "1429450010", // UK
		},
		Concurrency: 1,
	},
}
