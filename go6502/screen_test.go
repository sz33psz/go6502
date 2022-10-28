package go6502

import (
	"image/color"
	"testing"
)

func TestColor(t *testing.T) {
	allWhite := mapColor(uint8(0b11111111))
	assertColor(t, allWhite, 0xFF, 0xFF, 0xFF)
	allBlue := mapColor(uint8(0b00000011))
	assertColor(t, allBlue, 0x00, 0x00, 0xFF)
	redAndBlue := mapColor(uint8(0b11100011))
	assertColor(t, redAndBlue, 0xFF, 0x00, 0xFF)
	greenAndBlue := mapColor(uint8(0b00011111))
	assertColor(t, greenAndBlue, 0x00, 0xFF, 0xFF)

	colorsMix := mapColor(uint8(0b01010110))
	assertColor(t, colorsMix, 0x40, 0xA0, 0xAA)
}

func assertColor(t *testing.T, rgba color.RGBA, r uint8, g uint8, b uint8) {
	if rgba.R != r {
		t.Fatalf("R component is wrong. Expected: %v, got: %v", r, rgba.R)
	}
	if rgba.G != g {
		t.Fatalf("G component is wrong. Expected: %v, got: %v", g, rgba.G)
	}
	if rgba.B != b {
		t.Fatalf("B component is wrong. Expected: %v, got: %v", b, rgba.B)
	}
	if rgba.A != 255 {
		t.Fatalf("Alpha component is wrong. Expected: %v, got: %v", 255, rgba.A)
	}
}
