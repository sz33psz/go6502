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
	// 0xD000 - 0xDBB8 - Screen
	// 0xE000 - 0xFFFF - RAM
	firstRamSegment := go6502.NewRAM(0x0000, ScreenMemoryStart)
	screen := go6502.NewScreen(ScreenMemoryStart)
	secondRamSegment := go6502.NewRAM(0xE000, 0x2000)
	cpu := go6502.NewCPU(go6502.NewMemory(
		firstRamSegment,
		screen,
		secondRamSegment))

	//Reset vector
	cpu.Memory.Set(go6502.ResetVectorL, 0x00, 0x00)

	screen.SetMapping(0, 0, 0b11100000, 0b00011100)

	cpu.Memory.Set(0x0000, go6502.OpLDA_imm, 0b10010110)
	cpu.Memory.Set(0x0002, go6502.OpSTA_absolute, 0x58, 0xD2) //Beginning of pixel memory

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Go6502")

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		cpu.Initialize()
		cpu.Advance()
		cpu.Advance()
		wg.Done()
	}()

	if err := ebiten.RunGame(screen); err != nil {
		log.Fatal(err)
	}

	wg.Wait()
}
