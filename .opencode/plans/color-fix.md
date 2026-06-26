# Light Color Index Fix Plan

## Fix 1 — GUI: Correct colorOptions FirmwareIdx (gui/main.go:179-192)

**Change:** Replace the `colorOptions` array with corrected firmware indices.

```go
var colorOptions = []struct {
	Name        string
	FirmwareIdx byte
	DisplayHex  string
}{
	{"Rainbow", 0, ""},
	{"Red", 1, "#FF0000"},
	{"Green", 2, "#00FF00"},
	{"Blue", 3, "#0000FF"},
	{"Yellow", 4, "#FFFF00"},
	{"Magenta", 5, "#FF00FF"},
	{"Cyan", 6, "#00FFFF"},
	{"White", 7, "#FFFFFF"},
}
```

Changes: Green 4→2, Blue 6→3, Yellow 3→4, Magenta 7→5, Cyan 5→6, White 8→7 (dropped unused index 8).

## Fix 2 — GUI: Add startingUp guard + async to colorIndexSelect.OnChanged (gui/main.go:735-747)

**Change:** Replace the `colorIndexSelect.OnChanged` handler with a guarded + async version matching the sequence handler's pattern.

```go
colorIndexSelect := fynetool.NewSelect(colorLabels, func(label string) {
    for i, l := range colorLabels {
        if l == label {
            profile.Light.ColorIndex = int(colorOptions[i].FirmwareIdx)
            break
        }
    }
    applyKeyboardColors(kb, profile, model)
    if !startingUp {
        go func() {
            cmd := exec.Command("sudo", "-E", "drunkdeer-cli", "set", "light", "colorIndex", strconv.Itoa(profile.Light.ColorIndex))
            if out, err := cmd.CombinedOutput(); err != nil {
                log.Printf("Error setting light color index: %v\n%s", err, out)
            }
        }()
    }
})
```

## Fix 3 — CLI: Preserve colorIndex in set light sequence (drunkdeer/app.go:512)

**Change:** Line 512, replace `0x00` with `byte(state.Light.ColorIndex)`.

```go
// Before:
a.controller.SendLEDModeSelect(0x00, byte(v), byte(state.Light.Speed), br, 0x00)
// After:
a.controller.SendLEDModeSelect(0x00, byte(v), byte(state.Light.Speed), br, byte(state.Light.ColorIndex))
```
