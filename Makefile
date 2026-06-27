.PHONY: all clean deps wayland

all: drunkdeer

drunkdeer:
	@GO_TAGS=""; \
	if [ -n "$$WAYLAND_DISPLAY" ]; then GO_TAGS="wayland"; fi; \
	CGO_ENABLED=1 go build -tags "$$GO_TAGS" -o drunkdeer ./cmd/drunkdeer

wayland:
	$(MAKE) GO_TAGS=wayland

clean:
	rm -f drunkdeer

deps:
	@echo "Debian/Ubuntu: sudo apt install libgl1-mesa-dev libxrandr-dev libxinerama-dev libxcursor-dev libxi-dev libxext-dev"
	@echo "Fedora:        sudo dnf install mesa-libGL-devel libXrandr-devel libXinerama-devel libXcursor-devel libXi-devel libXext-devel"
	@echo "Arch:          sudo pacman -S mesa libxrandr libxinerama libxcursor libxi libxext"
	@echo "Nix:           nix-shell -p mesa libxrandr libxinerama libxcursor libxi libxext"
