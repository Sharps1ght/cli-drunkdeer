GOAMD64=v1
GO_TAGS=

.PHONY: all cli gui wayland clean

all: cli gui

cli:
	GOAMD64=$(GOAMD64) go build -o drunkdeer-cli ./drunkdeer

gui:
	GOAMD64=$(GOAMD64) go build -tags '$(GO_TAGS)' -o drunkdeer-gui ./gui

wayland:
	$(MAKE) gui GO_TAGS=wayland

clean:
	rm -f drunkdeer-cli drunkdeer-gui
