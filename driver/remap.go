package driver

func BuildRemapChunk(layer, chunkNum byte, keys map[int]*RemapKey) []byte {
	report := make([]byte, 63)
	report[0] = PACKET_IDENTITY
	report[1] = 0x02
	report[2] = 0x04
	report[3] = chunkNum
	report[4] = 14
	report[5] = layer

	start := int(chunkNum-1) * 9
	r := 6
	for w := 0; w < 9; w++ {
		s := start + w
		if s >= 126 {
			break
		}
		k, ok := keys[s]
		if ok && k != nil {
			report[r] = byte(s)
			r++
			switch k.KeyType {
			case 0:
				report[r] = k.KeyCmd
				r += 2
				report[r] = k.KeyCode
				r += 3
			case 1, 2:
				report[r] = k.KeyCmd
				r++
				report[r] = k.KeyCode
				r += 4
			case 3:
				report[r] = k.KeyCmd
				r += 5
			case 4:
				report[r] = 0xF8
				r++
				report[r] = k.KeyCode
				r++
				report[r] = 0
				r++
				report[r] = 0
				r++
				report[r] = 0
				r++
			default:
				report[r] = k.KeyCmd
				r += 5
			}
		} else {
			r += 6
		}
	}
	report[60] = 0xA5

	return report
}
