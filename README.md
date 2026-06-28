# Custom DrunkDeer Driver CLI

Please, visit [original repository](https://github.com/2xxn/cli-drunkdeer) for additional info.

## Build Dependencies

You need `gcc`, `CGO_ENABLED=1`, and X11 development libraries:

```bash
# Debian/Ubuntu
sudo apt install libgl1-mesa-dev libxrandr-dev libxinerama-dev libxcursor-dev libxi-dev libxext-dev
# Fedora
sudo dnf install mesa-libGL-devel libXrandr-devel libXinerama-devel libXcursor-devel libXi-devel libXext-devel
# Arch
sudo pacman -S mesa libxrandr libxinerama libxcursor libxi libxext
# Nix
nix-shell -p mesa libxrandr libxinerama libxcursor libxi libxext
```

## Building

```bash
git clone https://github.com/Sharps1ght/opendrunkdeer
cd opendrunkdeer
make
```

This produces a single `drunkdeer` binary with both CLI and GUI modes. On Windows, starting the .exe just opens the GUI. On Wayland, ```make``` auto-detects your display server and builds with Wayland support. To force a Wayland build: ```make wayland```.

## udev Rules

By default, HID devices require root access. Install the udev rules once to let your user account talk to the keyboard directly:

```bash
# Run these from the project root (where the etc/ directory lives)
sudo cp etc/udev/rules.d/99-drunkdeer.rules /etc/udev/rules.d/99-drunkdeer.rules
sudo udevadm control --reload-rules && sudo udevadm trigger
```

Create the `plugdev` group if it doesn't exist, then add your user and **log out and back in**:

```bash
sudo groupadd -f plugdev
sudo usermod -aG plugdev $USER
```

To make sure the utility works, **unplug and replug the keyboard** for the rule to take effect.

## Usage

```bash
drunkdeer list                      			# list connected devices
drunkdeer profiles                  			# list saved profiles
drunkdeer load <MyProfile>          			# load a profile
drunkdeer import path/to/config.json			# import a config
drunkdeer gui                       			# launch the graphical configurator
drunkdeer gui -model G60/G70/G75/A75			# GUI with a specific keyboard model
drunkdeer set <section> <subsection> <value>    # apply a setting according to .json (e.g. 'drunkdeer set light sequence 19)
```

**NOTE:** ```color(s)``` are NOT RECOMMENDED to be set via ```drunkdeer set``` due to the way configuration is handled.

## Config File structure

Speed and brightness must be between 0 and 9 inside of .json (where 9 is max), direction must be 0, 1 or 2. Actuation point should be between 0.2mm and 3.8mm for the best experience.

```json
{
    "model": "G60",
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
        "sequence": 5, 
        "color": "#FFFFFF",
        "colorTurbo": "#FF00FF",
        "colors":   {
            "W": "#FF0000", 
            "A": "#FF0000", 
            "S": "#FF0000", 
            "D": "#FF0000"
        }
    },
    "actuationPoints": {
        "W": 0.2,
        "A": 0.2,
        "S": 0.2,
        "D": 0.2,
        "CAPS": 3.8
    },
    "rapidTriggers": {
        "W": [0.2, 0.2],
        "A": [0.2, 0.2],
        "S": [0.2, 0.2],
        "D": [0.2, 0.2]
    },
    "remap":    {
        "Default":  {
            "W": "S",
            "A": "D",
            "S": "W",
            "D": "A"
        },
        "Fn":   {
            "1": "F1",
            "2": "F2",
            [...]
        },
        "Menu": {
            "1": "NUM1",
            "2": "NUM2",
            [...]
        }
}
```

For tables of actions and keys refer to [mapping tables](/mapping.md). It is only made for G60 for now, but it should overlap with the rest of Drunkdeer keyboards.
List of character names and color sequences can be found in [this file](https://github.com/Sharps1ght/opendrunkdeer/blob/main/driver/consts.go)
