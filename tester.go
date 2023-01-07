package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"go6502/go6502"
	"log"
	"sync"
)

const (
	ScreenMemoryStart = 0xD000
)

func main() {
	//Memory layout:
	// 0x0000 - 0xCFFF - RAM
	// 0xD000 - 0xFF00 - Screen
	// 0xFF00 - 0xFFFF - RAM
	firstRamSegment := go6502.NewRAM(0x0000, ScreenMemoryStart)
	screen := go6502.NewScreen(ScreenMemoryStart)
	secondRamSegment := go6502.NewRAM(0xFF00, 0x0100)
	cpu := go6502.NewCPU(go6502.NewMemory(
		firstRamSegment,
		screen,
		secondRamSegment))

	//Reset vector
	cpu.Memory.Set(go6502.ResetVectorL, 0x00, 0x00)

	for i := 0; i < go6502.ScreenWidth; i++ {
		screen.SetMapping(i, 0, 0b11111111, 0b00000000)
	}

	cpu.Memory.Set(0x0000, go6502.OpLDA_imm, 0b00111000)
	for i := 0; i < go6502.BlocksInLine*8; i++ {
		addr := 0x0002 + 0x03*i
		pixelAddress := 0xD960 + i
		cpu.Memory.Set(uint16(addr), go6502.OpSTA_absolute, uint8(pixelAddress%256), uint8(pixelAddress/256)) //Beginning of pixel memory
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Go6502")

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		cpu.Initialize()
		for i := 0; i < go6502.BlocksInLine*8; i++ {
			cpu.Advance()
		}
		wg.Done()
	}()

	if err := ebiten.RunGame(screen); err != nil {
		log.Fatal(err)
	}

	wg.Wait()
}
