package driver

import (
	"sync"

	"github.com/sstallion/go-hid"
)

type DrunkDeerController struct {
	device   *hid.Device
	identity *DDKeyboardIdentity

	actuations  []byte
	downstrokes []byte
	upstrokes   []byte

	turbo        bool
	rapidTrigger bool
	debug        bool

	mu          sync.Mutex
	wg          sync.WaitGroup
	packetWg    sync.WaitGroup
	packetChan  chan DDPacket
	packetQueue chan []byte

	Light *DDLight

	shouldClose bool
	closeOnce   sync.Once
}

type DDPacket struct {
	Packet uint8
	Data   []byte
}

type DDLight struct {
	Direction    byte
	Speed        byte
	Sequence     byte
	Brightness   byte
	Color        byte
	Colors       map[int][3]byte // per-key custom colors (key index → RGB)
	DefaultColor [3]byte         // fallback color for unset keys in custom mode

	TurboDefaultColor [3]byte // uniform color for turbo custom mode (no per-key)
}

type DDKeyboardIdentity struct {
	KeyboardModel string
	KeyboardType  uint8

	FirmwareVersion string
	Turbo           bool
	RapidTrigger    bool
}
