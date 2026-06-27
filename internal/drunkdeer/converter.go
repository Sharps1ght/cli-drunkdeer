package drunkdeer

import (
	"encoding/json"
)

type DDConfigKey struct {
	Keyname    string  `json:"keyname"`
	Actuation  float32 `json:"action_point"`
	Downstroke float32 `json:"downstroke"`
	Upstroke   float32 `json:"upstroke"`
}

type DDConfig struct {
	StorageName string        `json:"storagename"`
	Showname    string        `json:"showname"`
	Keys        []DDConfigKey `json:"keys_array"`
}

func (c *DDConfig) mostUsedValues() (actuation, downstroke, upstroke float32) {
	countsAct := make(map[float32]int)
	countsDS := make(map[float32]int)
	countsUS := make(map[float32]int)
	for _, key := range c.Keys {
		countsAct[key.Actuation]++
		countsDS[key.Downstroke]++
		countsUS[key.Upstroke]++
	}

	var maxAct, maxDS, maxUS int
	for v, n := range countsAct {
		if n > maxAct {
			maxAct, actuation = n, v
		}
	}
	for v, n := range countsDS {
		if n > maxDS {
			maxDS, downstroke = n, v
		}
	}
	for v, n := range countsUS {
		if n > maxUS {
			maxUS, upstroke = n, v
		}
	}
	return
}

func (c *DDConfig) getModelFromStorageName() string {
	sub := c.StorageName[5:8]

	return sub
}

func (c *DDConfig) convertToCLIConfig() *Config {
	var config Config

	config.Model = c.getModelFromStorageName()
	config.DefaultActuation, config.RapidTrigger.DefaultDownstroke, config.RapidTrigger.DefaultUpstroke = c.mostUsedValues()
	config.RapidTrigger.Enabled = true

	config.Turbo = false

	config.ActuationPoints = make(map[string]float32)
	config.RapidTriggers = make(map[string][2]float32)

	config.Light = LightSettings{}
	config.Light.Enabled = false

	for _, key := range c.Keys {
		var createRtEntry bool = false

		if key.Actuation != 0 && key.Actuation != config.DefaultActuation {
			config.ActuationPoints[key.Keyname] = key.Actuation
		}

		if key.Downstroke != 0 && key.Downstroke != config.RapidTrigger.DefaultDownstroke {
			createRtEntry = true
		}

		if key.Upstroke != 0 && key.Upstroke != config.RapidTrigger.DefaultUpstroke {
			createRtEntry = true
		}

		if createRtEntry {
			config.RapidTriggers[key.Keyname] = [2]float32{key.Downstroke, key.Upstroke}
		}
	}

	return &config
}

func parseDrunkDeerConfig(config []byte) *DDConfig {
	var ddConfig DDConfig
	err := json.Unmarshal(config, &ddConfig)
	handleError("Error parsing JSON", err)

	return &ddConfig
}
