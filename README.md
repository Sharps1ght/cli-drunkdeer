# Custom DrunkDeer Driver CLI
This CLI is a custom driver for the DrunkDeer keyboard. It allows you to configure the keyboard's settings, including actuation points, light settings, turbo mode and rapid trigger.

## Reason for this project
This project was created out of frustration, DrunkDeer webdriver's servers are so awful that getting into the WebDriver can sometimes take up to 10 minutes (especially uncached, I have tendencies to reload using CTRL+SHIFT+R). This project is a workaround for that, it allows you to configure the keyboard without the need for the web driver, additionally "preventing" DrunkDeer from exit-scamming.

## REMAPPING IS NOT SUPPORTED YET

## Installation
### You will need gcc before installing and CGO_ENABLED=1
```bash
set CGO_ENABLED=1
go install github.com/2xxn/cli-drunkdeer/drunkdeer@latest
```
### You may simply move ./go/bin/drunkdeer* to /usr/bin under root on Linux, BUT ONLY AT YOUR OWN RISK!!!
This allows you to do this (for example):
```bash
sudo drunkdeer -h
```

```sudo``` is required cuz of permissions stuff.
This tool works, but not for everything.

## Usage
```bash
drunkdeer [command] [value?] [options?]
```

### You will generally need those 4 commands
```bash
drunkdeer list - list all available devices
drunkdeer profiles - list all available profiles
drunkdeer import [path-to-config-file] - import a DRUNKDEER ANTLER CONFIG FILE
drunkdeer load [profile-name] - load a profile into the keyboard
```
#### To import someone's CLI config file you can do `drunkdeer load [url/relative or absolute path]`


### To learn more
```bash
drunkdeer
drunkdeer -h
```

## Config File structure
### This entry is for myself and the more advanced users
### The config file is a JSON file that contains the following structure:
List of character names and color sequences can be found in [this file](https://github.com/2xxn/cli-drunkdeer/blob/main/driver/consts.go)<br>
Actuation point should be between 0.1mm and 3.9mm (although both are unadvised, you should do 0.2mm at lowest)
#### Speed and brightness must be between 0 and 9 (where 9 is max)
```json
{
    "model": "A75",
    "turbo": false,
    "defaultActuation": 2.0,
    "rapidTrigger": {
        "enabled": false,
        "defaultDownstroke": 0.0,
        "defaultUpstroke": 0.0
    },
    "light": {
        "enabled": true,
        "direction": 0,
        "speed": 5,
        "brightness": 9,
        "sequence": 5
    },
    "actuationPoints": {
        "W": 0.2,
        "A": 0.2,
        "S": 0.2,
        "D": 0.2,
        "TAB": 3.8
    },
    "rapidTriggers": {
        "A": [0.2, 0.2],
        "S": [0.2, 0.2]
    }
}
```

# WORD FROM ME, SHARPSIGHT
I (OpenCode Zen, i dunno) made color profiles work on Linux. It still requires ```sudo``` to work. ONLY TESTED WITH G60 AND EXPLICITLY MADE TO WORK WITH IT! PLEASE, don't try it on other keyboards. Or do, but at your own risk.
## Usage
There is a precomplied executable at *cli-drunkdeer/drunkdeer*. To run, use
```bash
sudo /path/to/executable/drunkdeer *command*
```
or compile yourself: go inside the cloned repo's directory and run ```go build -o drunkdeer ./drunkdeer```.
## List of sequences
Or light modes, however you like:
| Value | Name |
|-------|------|
| 0  | Off              |
| 1  | Rotating Chase   |
| 2  | Spectrum Wave    |
| 3  | Right Surfing    |
| 4  | Breathing        |
| 5  | Center Surfing   |
| 6  | Spectrum Cycle   |
| 7  | Key Ripple       |
| 8  | Always On        |
| 9  | Press To Light   |
| 10 | Center Snake     |
| 11 | Color Fountain   |
| 12 | Key Laser        |
| 13 | Glowing Fish     |
| 14 | Cross Surfing    |
| 15 | Heart            |
| 16 | Traffic          |
| 17 | Snake            |
| 18 | Raindrop         |
| 19 | Custom Colors    |
## Custom colors
If you set ```sequence``` to ```19``` you can set your colors with ```color(s)```.
This utility accepts HEX code. ```color``` decides the default color of unspecified keys, while ```colors``` is for custom color for specific key, overriding the ```color```.
```
"light":    {
[...]
"color": "#00FF00",
"colors":   {
    "ESC": "#FF0000",
    "SPACE": "#FF0000",
    "RETURN": "#FF0000"
    }
}
```
This exact setup will make every key green, except Escape, Space and Enter, these three will be red.
