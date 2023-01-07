package go6502

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
)

const (
	ScreenWidth  = 320
	ScreenHeight = 240

	BlockWidth   = 8
	BlockHeight  = 8
	BlocksInLine = ScreenWidth / BlockWidth

	PixelsPerByte = 8

	ColorMappingsBytes = 2 * (ScreenWidth * ScreenHeight) / (BlockWidth * BlockHeight) // 2400b
	PixelColorBytes    = ScreenWidth * ScreenHeight / PixelsPerByte                    // 9600b
)

type Screen struct {
	colorMappings []uint8
	pixels        []uint8
	addressStart  uint16
}

/*
Screen memory consists of two parts:
- Color mappings - for each block of 8px x 8px, there are two colors selected
- Pixel colors - one bit for each
*/

func (s *Screen) WithinRange(address uint16) bool {
	return address >= s.addressStart && address < s.addressStart+ColorMappingsBytes+PixelColorBytes
}

func (s *Screen) Get(address uint16) uint8 {
	internalAddress := address - s.addressStart
	if internalAddress < ColorMappingsBytes {
		return s.colorMappings[internalAddress]
	} else {
		return s.pixels[internalAddress-ColorMappingsBytes]
	}
}

func (s *Screen) Set(address uint16, value uint8) {
	internalAddress := address - s.addressStart
	if internalAddress < ColorMappingsBytes {
		s.colorMappings[internalAddress] = value
	} else {
		s.pixels[internalAddress-ColorMappingsBytes] = value
	}
}

func (s *Screen) SetMapping(x int, y int, fg uint8, bg uint8) {
	blockNumber := s.getBlockNumber(x, y)
	s.colorMappings[blockNumber*2] = fg
	s.colorMappings[blockNumber*2+1] = bg
}

func (s *Screen) getColorMappings(x, y int) (uint8, uint8) {
	blockNumber := s.getBlockNumber(x, y)
	return s.colorMappings[blockNumber*2], s.colorMappings[blockNumber*2+1]
}

func (s *Screen) getBlockNumber(x int, y int) int {
	blockX := x / BlockWidth
	blockY := y / BlockHeight
	blockNumber := blockY*BlocksInLine + blockX
	return blockNumber
}

func (s *Screen) GetPixels(x, y int) uint8 {
	byteInRow := x / PixelsPerByte
	byteNumber := y*(ScreenWidth/PixelsPerByte) + byteInRow
	return s.pixels[byteNumber]
}

func (s *Screen) Update() error {
	return nil
}

func (s *Screen) Draw(screen *ebiten.Image) {
	for blockStartY := 0; blockStartY < ScreenHeight; blockStartY += BlockHeight {
		for blockStartX := 0; blockStartX < ScreenWidth; blockStartX += BlockWidth {
			s.drawBlock(screen, blockStartX, blockStartY)
		}
	}
}

func (s *Screen) drawBlock(screen *ebiten.Image, blockStartX int, blockStartY int) {
	colorFg, colorBg := s.getColorMappings(blockStartX, blockStartY)
	for line := 0; line < BlockHeight; line++ {
		lineStartY := blockStartY + line
		linePixels := s.GetPixels(blockStartX, lineStartY)
		s.drawLine(screen, blockStartX, lineStartY, linePixels, colorFg, colorBg)
	}
}

func (s *Screen) drawLine(screen *ebiten.Image, lineStartX int, y int, pixels uint8, fg uint8, bg uint8) {
	for b := 0; b < PixelsPerByte; b++ {
		pixel := pixels & (1 << b)
		pixelColor := fg
		if pixel == 0 {
			pixelColor = bg
		}
		ebitenutil.DrawRect(screen, float64(lineStartX+b), float64(y), 1, 1, mapColor(pixelColor))
	}
}

func (s *Screen) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func NewScreen(addressStart uint16) *Screen {
	return &Screen{
		colorMappings: make([]uint8, ColorMappingsBytes),
		pixels:        make([]uint8, PixelColorBytes),
		addressStart:  addressStart,
	}
}

var (
	redMappings   = []uint8{0x00, 0x20, 0x40, 0x60, 0x80, 0xA0, 0xC0, 0xFF}
	greenMappings = []uint8{0x00, 0x20, 0x40, 0x60, 0x80, 0xA0, 0xC0, 0xFF}
	blueMappings  = []uint8{0x00, 0x55, 0xAA, 0xFF}
)

func mapColor(eightBitColor uint8) color.RGBA {
	r := eightBitColor >> 5
	g := (eightBitColor & 0b00011100) >> 2
	b := eightBitColor & 0b00000011

	return color.RGBA{
		R: redMappings[r],
		G: greenMappings[g],
		B: blueMappings[b],
		A: 0xFF,
	}
}
