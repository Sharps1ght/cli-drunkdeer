package widget

type KeyDef struct {
	Name  string
	Value int
	Width float64
	Row   int
}

type LayoutDef struct {
	Rows      [][]KeyDef
	UnitSize  float64
	RowHeight float64
	Gap       float64
}

func G60Layout() LayoutDef {
	return LayoutDef{
		Rows: [][]KeyDef{
			{
				{Name: "Esc", Value: 21, Width: 1},
				{Name: "1", Value: 22, Width: 1}, {Name: "2", Value: 23, Width: 1},
				{Name: "3", Value: 24, Width: 1}, {Name: "4", Value: 25, Width: 1},
				{Name: "5", Value: 26, Width: 1}, {Name: "6", Value: 27, Width: 1},
				{Name: "7", Value: 28, Width: 1}, {Name: "8", Value: 29, Width: 1},
				{Name: "9", Value: 30, Width: 1}, {Name: "0", Value: 31, Width: 1},
				{Name: "-", Value: 32, Width: 1}, {Name: "=", Value: 33, Width: 1},
				{Name: "Back", Value: 34, Width: 2},
			},
			{
				{Name: "Tab", Value: 42, Width: 1.5},
				{Name: "Q", Value: 43, Width: 1}, {Name: "W", Value: 44, Width: 1},
				{Name: "E", Value: 45, Width: 1}, {Name: "R", Value: 46, Width: 1},
				{Name: "T", Value: 47, Width: 1}, {Name: "Y", Value: 48, Width: 1},
				{Name: "U", Value: 49, Width: 1}, {Name: "I", Value: 50, Width: 1},
				{Name: "O", Value: 51, Width: 1}, {Name: "P", Value: 52, Width: 1},
				{Name: "[", Value: 53, Width: 1}, {Name: "]", Value: 54, Width: 1},
				{Name: "\\", Value: 55, Width: 1},
			},
			{
				{Name: "Caps", Value: 63, Width: 1.75},
				{Name: "A", Value: 64, Width: 1}, {Name: "S", Value: 65, Width: 1},
				{Name: "D", Value: 66, Width: 1}, {Name: "F", Value: 67, Width: 1},
				{Name: "G", Value: 68, Width: 1}, {Name: "H", Value: 69, Width: 1},
				{Name: "J", Value: 70, Width: 1}, {Name: "K", Value: 71, Width: 1},
				{Name: "L", Value: 72, Width: 1},
				{Name: ";", Value: 73, Width: 1}, {Name: "'", Value: 74, Width: 1},
				{Name: "Enter", Value: 76, Width: 2.25},
			},
			{
				{Name: "Shift", Value: 84, Width: 2.25},
				{Name: "Z", Value: 86, Width: 1}, {Name: "X", Value: 87, Width: 1},
				{Name: "C", Value: 88, Width: 1}, {Name: "V", Value: 89, Width: 1},
				{Name: "B", Value: 90, Width: 1}, {Name: "N", Value: 91, Width: 1},
				{Name: "M", Value: 92, Width: 1},
				{Name: ",", Value: 93, Width: 1}, {Name: ".", Value: 94, Width: 1},
				{Name: "/", Value: 95, Width: 1},
				{Name: "Shift", Value: 97, Width: 2.75},
			},
			{
				{Name: "Ctrl", Value: 105, Width: 1.25},
				{Name: "Win", Value: 106, Width: 1.25},
				{Name: "Alt", Value: 107, Width: 1.25},
				{Name: "Space", Value: 111, Width: 6.25},
				{Name: "Alt", Value: 115, Width: 1.25},
				{Name: "Fn1", Value: 116, Width: 1.25},
				{Name: "Fn2", Value: 117, Width: 1.25},
				{Name: "Ctrl", Value: 118, Width: 1.25},
			},
		},
		UnitSize:  60,
		RowHeight: 60,
		Gap:       4,
	}
}

func A75Layout() LayoutDef {
	return LayoutDef{
		Rows: [][]KeyDef{
			{
				{Name: "Esc", Value: 0, Width: 1},
				{Name: "F1", Value: 2, Width: 1}, {Name: "F2", Value: 3, Width: 1},
				{Name: "F3", Value: 4, Width: 1}, {Name: "F4", Value: 5, Width: 1},
				{Name: "F5", Value: 6, Width: 1}, {Name: "F6", Value: 7, Width: 1},
				{Name: "F7", Value: 8, Width: 1}, {Name: "F8", Value: 9, Width: 1},
				{Name: "F9", Value: 10, Width: 1}, {Name: "F10", Value: 11, Width: 1},
				{Name: "F11", Value: 12, Width: 1}, {Name: "F12", Value: 13, Width: 1},
			},
			{
				{Name: "~", Value: 21, Width: 1},
				{Name: "1", Value: 22, Width: 1}, {Name: "2", Value: 23, Width: 1},
				{Name: "3", Value: 24, Width: 1}, {Name: "4", Value: 25, Width: 1},
				{Name: "5", Value: 26, Width: 1}, {Name: "6", Value: 27, Width: 1},
				{Name: "7", Value: 28, Width: 1}, {Name: "8", Value: 29, Width: 1},
				{Name: "9", Value: 30, Width: 1}, {Name: "0", Value: 31, Width: 1},
				{Name: "-", Value: 32, Width: 1}, {Name: "=", Value: 33, Width: 1},
				{Name: "Back", Value: 34, Width: 2},
			},
			{
				{Name: "Tab", Value: 42, Width: 1.5},
				{Name: "Q", Value: 43, Width: 1}, {Name: "W", Value: 44, Width: 1},
				{Name: "E", Value: 45, Width: 1}, {Name: "R", Value: 46, Width: 1},
				{Name: "T", Value: 47, Width: 1}, {Name: "Y", Value: 48, Width: 1},
				{Name: "U", Value: 49, Width: 1}, {Name: "I", Value: 50, Width: 1},
				{Name: "O", Value: 51, Width: 1}, {Name: "P", Value: 52, Width: 1},
				{Name: "[", Value: 53, Width: 1}, {Name: "]", Value: 54, Width: 1},
				{Name: "\\", Value: 55, Width: 1},
			},
			{
				{Name: "Caps", Value: 63, Width: 1.75},
				{Name: "A", Value: 64, Width: 1}, {Name: "S", Value: 65, Width: 1},
				{Name: "D", Value: 66, Width: 1}, {Name: "F", Value: 67, Width: 1},
				{Name: "G", Value: 68, Width: 1}, {Name: "H", Value: 69, Width: 1},
				{Name: "J", Value: 70, Width: 1}, {Name: "K", Value: 71, Width: 1},
				{Name: "L", Value: 72, Width: 1},
				{Name: ";", Value: 73, Width: 1}, {Name: "'", Value: 74, Width: 1},
				{Name: "Enter", Value: 76, Width: 2.25},
			},
			{
				{Name: "Shift", Value: 84, Width: 2.25},
				{Name: "Z", Value: 86, Width: 1}, {Name: "X", Value: 87, Width: 1},
				{Name: "C", Value: 88, Width: 1}, {Name: "V", Value: 89, Width: 1},
				{Name: "B", Value: 90, Width: 1}, {Name: "N", Value: 91, Width: 1},
				{Name: "M", Value: 92, Width: 1},
				{Name: ",", Value: 93, Width: 1}, {Name: ".", Value: 94, Width: 1},
				{Name: "/", Value: 95, Width: 1},
				{Name: "Shift", Value: 97, Width: 2.75},
			},
			{
				{Name: "Ctrl", Value: 105, Width: 1.25},
				{Name: "Win", Value: 106, Width: 1.25},
				{Name: "Alt", Value: 107, Width: 1.25},
				{Name: "Space", Value: 111, Width: 6.25},
				{Name: "Alt", Value: 115, Width: 1.25},
				{Name: "Fn1", Value: 116, Width: 1.25},
				{Name: "Fn2", Value: 117, Width: 0},
				{Name: "Ctrl", Value: 118, Width: 1.25},
			},
		},
		UnitSize:  56,
		RowHeight: 56,
		Gap:       4,
	}
}

var GetLayoutDef = func(model string) LayoutDef {
	if model == "G60" {
		return G60Layout()
	}
	return A75Layout()
}
