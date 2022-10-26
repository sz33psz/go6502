package go6502

import (
	"fmt"
)

const (
	ResetVectorL = 0xFFFC
	ResetVectorH = 0xFFFD
)

const (
	OpNOOP = 0xEA

	// Register basic ops
	OpTAX = 0xAA
	OpTXA = 0x8A
	OpDEX = 0xCA
	OpINX = 0xE8
	OpTAY = 0xA8
	OpTYA = 0x98
	OpDEY = 0x88
	OpINY = 0xC8

	// Accumulator
	OpLDA_imm        = 0xA9
	OpLDA_zeropage   = 0xA5
	OpLDA_zeropage_x = 0xB5
	OpSTA_zeropage   = 0x85
	OpSTA_absolute   = 0x8D
	OpSTA_zeropage_x = 0x95

	OpCMP_imm = 0xC9

	OpLDX_imm      = 0xA2
	OpLDX_zeropage = 0xA6

	OpLDY_imm      = 0xA0
	OpLDY_zeropage = 0xA4
)

type CPU struct {
	A      uint8
	X      uint8
	Y      uint8
	PC     uint16
	S      uint8
	Flags  *Flags
	Memory *Memory
}

func (cpu *CPU) String() string {
	return fmt.Sprintf("A: %02X\tX: %02X\tY: %02X\tS: %02X\nPC: %04X\tFlags: %s",
		cpu.A, cpu.X, cpu.Y, cpu.S, cpu.PC, cpu.Flags)
}

func (cpu *CPU) Initialize() {
	resetVector := uint16(cpu.Memory.Get(ResetVectorH))<<8 + uint16(cpu.Memory.Get(ResetVectorL))
	cpu.PC = resetVector
}

func (cpu *CPU) Advance() {
	instruction := cpu.getNextInstruction()
	switch instruction {
	case OpNOOP:
		//Noop

	// Register ops
	case OpTAX:
		cpu.tax()
	case OpTXA:
		cpu.txa()
	case OpDEX:
		cpu.dex()
	case OpINX:
		cpu.inx()
	case OpTAY:
		cpu.tay()
	case OpTYA:
		cpu.tya()
	case OpDEY:
		cpu.dey()
	case OpINY:
		cpu.iny()

	//Compare Accumulator
	case OpCMP_imm:
		cpu.cmp_imm()

	//LDA
	case OpLDA_imm:
		cpu.lda_imm()
	case OpLDA_zeropage:
		cpu.lda_zeropage()
	case OpLDA_zeropage_x:
		cpu.lda_zeropage_x()
	//TODO: other addressing modes
	//STA
	case OpSTA_zeropage:
		cpu.sta_zeropage()
	case OpSTA_absolute:
		cpu.sta_absolute()
	case OpSTA_zeropage_x:
		cpu.sta_zeropage_x()
	//TODO: other addressing modes

	//LDX
	case OpLDX_imm:
		cpu.ldx_imm()
	case OpLDX_zeropage:
		cpu.ldx_zeropage()
	//TODO: other addressing modes

	//LDY
	case OpLDY_imm:
		cpu.ldy_imm()
	case OpLDY_zeropage:
		cpu.ldy_zeropage()
	}

	println(cpu.String())
}

func (cpu *CPU) getNextInstruction() uint8 {
	opcode := cpu.Memory.Get(cpu.PC)
	cpu.PC += 1
	return opcode
}

func (cpu *CPU) updateNZ(value uint8) {
	if value == 0 {
		cpu.Flags.SetZero()
	} else {
		cpu.Flags.ClearZero()
	}

	if value&(1<<7) > 0 {
		cpu.Flags.SetNegative()
	} else {
		cpu.Flags.ClearNegative()
	}

}

func (cpu *CPU) tax() {
	cpu.X = cpu.A
	cpu.updateNZ(cpu.X)
}

func (cpu *CPU) txa() {
	cpu.A = cpu.X
	cpu.updateNZ(cpu.A)
}

func (cpu *CPU) dex() {
	cpu.X -= 1
	cpu.updateNZ(cpu.X)
}

func (cpu *CPU) inx() {
	cpu.X += 1
	cpu.updateNZ(cpu.X)
}

func (cpu *CPU) tay() {
	cpu.Y = cpu.A
	cpu.updateNZ(cpu.Y)
}

func (cpu *CPU) tya() {
	cpu.A = cpu.Y
	cpu.updateNZ(cpu.A)
}

func (cpu *CPU) dey() {
	cpu.Y -= 1
	cpu.updateNZ(cpu.Y)
}

func (cpu *CPU) iny() {
	cpu.Y += 1
	cpu.updateNZ(cpu.Y)
}

func (cpu *CPU) cmp_imm() {
	imm := cpu.getNextInstruction()
	if imm <= cpu.A {
		cpu.Flags.SetCarry()
	} else {
		cpu.Flags.ClearCarry()
	}
	cpu.updateNZ(cpu.A)
}

func (cpu *CPU) lda_imm() {
	cpu.A = cpu.getNextInstruction()
	cpu.updateNZ(cpu.A)
}

func (cpu *CPU) lda_zeropage() {
	location := cpu.getNextInstruction()
	cpu.A = cpu.Memory.Get(uint16(location))
	cpu.updateNZ(cpu.A)
}

func (cpu *CPU) lda_zeropage_x() {
	location := (uint16(cpu.getNextInstruction()) + uint16(cpu.X)) & 0xFF
	cpu.A = cpu.Memory.Get(location)
	cpu.updateNZ(cpu.A)
}

func (cpu *CPU) sta_zeropage() {
	location := cpu.getNextInstruction()
	cpu.Memory.Set(uint16(location), cpu.A)
}

func (cpu *CPU) sta_zeropage_x() {
	location := (uint16(cpu.getNextInstruction()) + uint16(cpu.X)) & 0xFF
	cpu.Memory.Set(location, cpu.A)
}

func (cpu *CPU) sta_absolute() {
	location := uint16(cpu.getNextInstruction()) + (uint16(cpu.getNextInstruction()) << 8)
	cpu.Memory.Set(location, cpu.A)
}

func (cpu *CPU) ldx_imm() {
	cpu.X = cpu.getNextInstruction()
	cpu.updateNZ(cpu.X)
}

func (cpu *CPU) ldx_zeropage() {
	location := cpu.getNextInstruction()
	cpu.X = cpu.Memory.Get(uint16(location))
	cpu.updateNZ(cpu.X)
}

func (cpu *CPU) ldy_imm() {
	cpu.Y = cpu.getNextInstruction()
	cpu.updateNZ(cpu.Y)
}

func (cpu *CPU) ldy_zeropage() {
	location := cpu.getNextInstruction()
	cpu.Y = cpu.Memory.Get(uint16(location))
	cpu.updateNZ(cpu.Y)
}

func NewDefaultMemoryCPU() *CPU {
	return &CPU{
		A:      0,
		Y:      0,
		X:      0,
		PC:     0,
		S:      0,
		Flags:  &Flags{0},
		Memory: DefaultMemory(),
	}
}
