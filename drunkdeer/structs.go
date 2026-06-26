package main

import (
	"encoding/json"
	"strconv"
	"strings"
)

type RapidTriggerSettings struct {
	Enabled           bool    `json:"enabled"`
	DefaultDownstroke float32 `json:"defaultDownstroke"`
	DefaultUpstroke   float32 `json:"defaultUpstroke"`
}

// ColorSetting accepts JSON int (color index 1-8), string color name, or hex string (#rrggbb)
type ColorSetting struct {
	Index int    // color index for mode select (1-8)
	Fill  string // hex fill color for custom mode
}

func (c *ColorSetting) UnmarshalJSON(data []byte) error {
	// try int first
	var i int
	if err := json.Unmarshal(data, &i); err == nil {
		c.Index = i
		return nil
	}
	// try string
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "#") {
		c.Fill = s
	} else {
		n, err := strconv.Atoi(s)
		if err == nil {
			c.Index = n
		} else {
			c.Index = colorNameToIndex(strings.ToLower(s))
		}
	}
	return nil
}

func (c ColorSetting) MarshalJSON() ([]byte, error) {
	if c.Fill != "" {
		return json.Marshal(c.Fill)
	}
	return json.Marshal(c.Index)
}

func colorNameToIndex(name string) int {
	switch name {
	case "red":
		return 1
	case "orange":
		return 2
	case "yellow":
		return 3
	case "green":
		return 4
	case "cyan":
		return 5
	case "blue":
		return 6
	case "purple":
		return 7
	case "white":
		return 8
	}
	return 1 // default red
}

type LightSettings struct {
	Enabled    bool              `json:"enabled"`
	Direction  int               `json:"direction"`
	Sequence   int               `json:"sequence"`
	Speed      int               `json:"speed"`
	Brightness int               `json:"brightness"`
	ColorIndex int               `json:"colorIndex,omitempty"`
	Color      *ColorSetting     `json:"color,omitempty"`
	Colors     map[string]string `json:"colors,omitempty"`
	TurboColor *ColorSetting     `json:"turboColor,omitempty"`
}

type RemapSettings struct {
	Default map[string]string `json:"Default,omitempty"`
	Fn      map[string]string `json:"Fn,omitempty"`
	Menu    map[string]string `json:"Menu,omitempty"`
}

type Config struct {
	Model            string                `json:"model"`
	RapidTrigger     RapidTriggerSettings  `json:"rapidTrigger"`
	Turbo            bool                  `json:"turbo"`
	DefaultActuation float32               `json:"defaultActuation"`
	ActuationPoints  map[string]float32    `json:"actuationPoints"`
	RapidTriggers    map[string][2]float32 `json:"rapidTriggers"`
	Light            LightSettings         `json:"light"`
	Remap            RemapSettings         `json:"remap,omitempty"`
}

type Args struct {
	Command  string   `arg:"positional"`
	CmdValue []string `arg:"positional"`
	Import   string   `arg:"-i,--import" help:"Import a drunkdeer webdriver profile from the specified file/url"`
	Debug    bool     `arg:"-d,--debug" help:"Enable debug mode"`
	Index    int      `arg:"-i,--index" help:"Keyboard index to use (0 for first device, 1 for second, etc.)"`
	Profiles bool   `arg:"-p,--profiles" help:"Show all available profiles"`
	Reset    bool   `arg:"-r,--reset" help:"Reset the keyboard to default settings"`
	Load     string `arg:"-L,--load" help:"Load a profile from the specified file/url"`
	Save     string `arg:"-S,--save" help:"Save URL as profile"`
	Version  bool   `arg:"-v,--version" help:"Show version information"`
	List     bool   `arg:"-l,--list" help:"List all connected devices"`
}
