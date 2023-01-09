package go6502

import "testing"

func TestPHA(t *testing.T) {
	cpu := NewDefaultMemoryCPU()
	cpu.Memory.Set(ResetVectorL, 0x00, 0x00)

	cpu.Memory.Set(0x0000, OpPHA)
	cpu.Initialize()
	cpu.A = 0xF1 // Value 0xF1 should be pushed to stack

	cpu.Advance()
	if cpu.S != 0xFE {
		t.Fatalf("Stack pointer should be decreased after PHA. Expected %x, got %x", 0xFE, cpu.S)
	}

	if cpu.Memory.Get(0x01FF) != 0xF1 {
		t.Fatalf("A register value should be pushed to stack. Expected %x, got %x", 0xF1, cpu.Memory.Get(0x01FF))
	}
}
func TestPHP(t *testing.T) {
	cpu := NewDefaultMemoryCPU()
	cpu.Memory.Set(ResetVectorL, 0x00, 0x00)

	cpu.Memory.Set(0x0000, OpPHP)
	cpu.Initialize()
	cpu.Flags.SetNegative(true)
	cpu.Flags.SetZero(false)
	cpu.Flags.SetCarry(true)

	cpu.Advance()
	if cpu.S != 0xFE {
		t.Fatalf("Stack pointer should be decreased after PHP. Expected %x, got %x", 0xFE, cpu.S)
	}

	if cpu.Memory.Get(0x01FF) != 0b10000001 { //Nv--dizC
		t.Fatalf("Flags should be pushed to stack. Expected %x, got %x", 0b10000001, cpu.Memory.Get(0x01FF))
	}
}

func TestPLA(t *testing.T) {
	cpu := NewDefaultMemoryCPU()
	cpu.Memory.Set(ResetVectorL, 0x00, 0x00)
	cpu.Memory.Set(0x01FF, 0xA7)

	cpu.Memory.Set(0x0000, OpPLA)
	cpu.Initialize()
	cpu.S = 0xFE // One value on stack

	cpu.Advance()

	if cpu.S != 0xFF {
		t.Fatalf("Stack pointer should be increased after PLA. Expected %x, got %x", 0xFF, cpu.S)
	}

	if cpu.A != 0xA7 {
		t.Fatalf("A register value should be pulled from stack. Expected %x, got %x", 0xA7, cpu.A)
	}
}

func TestPLP(t *testing.T) {
	cpu := NewDefaultMemoryCPU()
	cpu.Memory.Set(ResetVectorL, 0x00, 0x00)
	cpu.Memory.Set(0x01FF, 0b10000001) // Nv--dizC

	cpu.Memory.Set(0x0000, OpPLP)
	cpu.Initialize()
	cpu.S = 0xFE // One value on stack

	cpu.Advance()

	if cpu.S != 0xFF {
		t.Fatalf("Stack pointer should be increased after PLP. Expected %x, got %x", 0xFF, cpu.S)
	}

	if cpu.Flags.String() != "Nv--dizC" {
		t.Fatalf("Flags should be pulled from stack. Expected %v, got %v", "Nv--dizC", cpu.Flags.String())
	}
}

func TestJMP(t *testing.T) {
	cpu := NewDefaultMemoryCPU()
	cpu.Memory.Set(ResetVectorL, 0x00, 0x00)

	cpu.Memory.Set(0x23FF, 0xCA, 0x11) // Store jump address which will be used by Indirect JMP

	cpu.Memory.Set(0x0000, OpJMP_absolute, 0x40, 0x20)
	cpu.Memory.Set(0x2040, OpJMP_indirect, 0xFF, 0x23)
	cpu.Initialize()
	cpu.Advance()

	if cpu.PC != 0x2040 {
		t.Fatalf("PC after Absolute JMP wasn't changed. Expected %x, got %x", 0x2040, cpu.PC)
	}

	cpu.Advance()

	if cpu.PC != 0x11CA {
		t.Fatalf("PC after Indirect JMP wasn't changed. Expected %x, got %x", 0x11CA, cpu.PC)
	}
}

func TestJSR(t *testing.T) {
	cpu := NewDefaultMemoryCPU()
	cpu.Memory.Set(ResetVectorL, 0xC0, 0x01) // Starting at address 0x01C0

	cpu.Memory.Set(0x01C0, OpJSR_absolute, 0x40, 0x20)
	cpu.Initialize()
	cpu.Advance()

	if cpu.PC != 0x2040 {
		t.Fatalf("PC after JSR wasn't changed. Expected %x, got %x", 0x2040, cpu.PC)
	}

	if cpu.S != 0xFD {
		t.Fatalf("S after JSR wasn't changed. Expceted %x, got %x", 0xFD, cpu.S)
	}

	if cpu.Memory.Get(0x01FF) != 0x01 {
		t.Fatalf("PCH wasn't pushed to stack. Expected %x, got %x", 0x01, cpu.Memory.Get(0x1FF))
	}

	if cpu.Memory.Get(0x01FE) != 0xC2 {
		t.Fatalf("PCH wasn't pushed to stack. Expected %x, got %x", 0xC2, cpu.Memory.Get(0x1FE))
	}
}

func TestRTS(t *testing.T) {
	cpu := NewDefaultMemoryCPU()
	cpu.Memory.Set(ResetVectorL, 0x00, 0x00)

	cpu.Memory.Set(0x01FE, 0x30) // Lower bytes of return address
	cpu.Memory.Set(0x01FF, 0xC2) // Higher bytes of return address

	cpu.Memory.Set(0x0000, OpRTS)
	cpu.Initialize()
	cpu.S = 0xFD // 2 values on stack
	cpu.Advance()

	if cpu.PC != 0xC230 {
		t.Fatalf("PC wasn't loaded from stack. Expected %x, got %x", 0xC230, cpu.PC)
	}
}

// TODO: Old tests - refactor
func TestAccumulator(t *testing.T) {
	cpu := NewDefaultMemoryCPU()
	cpu.Memory.Set(ResetVectorL, 0x00, 0x00)
	cpu.Memory.Set(0xF0, 0x66) //Some initial value to load

	cpu.Memory.Set(0x0000, OpLDA_zeropage, 0xF0)
	cpu.Initialize()
	cpu.Advance()
	if cpu.A != 0x66 {
		t.Fatalf("A wasn't loaded from zeropage. Expected %x, got %x", 0x66, cpu.A)
	}

	cpu.Memory.Set(0x0002, OpTAX)
	cpu.Advance()
	if cpu.X != cpu.A {
		t.Fatalf("A wasn't copied to X. Expected X to be %v, got %v", cpu.A, cpu.X)
	}

	cpu.Memory.Set(0x0003, OpINX)
	cpu.Advance()
	if cpu.X != cpu.A+1 {
		t.Fatalf("X wasn't incremented. Expected X to be %v, got %v", cpu.A+1, cpu.X)
	}

	cpu.Memory.Set(0x0004, OpTXA)
	cpu.Advance()
	if cpu.A != cpu.X {
		t.Fatalf("X wasn't copied to A. Expected A to be %v, got %v", cpu.X, cpu.A)
	}

	cpu.Memory.Set(0x0005, OpLDA_imm, 0x80)
	cpu.Advance()
	if cpu.A != 0x80 {
		t.Fatalf("A wasn't loaded with immediate value. Expected %v, got %v", 0x80, cpu.A)
	}

	cpu.Memory.Set(0x0007, OpLDX_imm, 0x01)
	cpu.Memory.Set(0x0009, OpSTA_zeropage_x, 0xF0) //Should store at 0x00F1 (0x00F0 + X)

	cpu.Advance()
	cpu.Advance()

	if cpu.A != cpu.Memory.Get(0x00F1) {
		t.Fatalf("A wasn't stored at zero page + X (0x00F1). Expected %v, got %v", cpu.A, cpu.Memory.Get(0x00F1))
	}

	cpu.Memory.Set(0x000B, OpSTA_zeropage, 0xF2)
	cpu.Advance()

	if cpu.A != cpu.Memory.Get(0x00F2) {
		t.Fatalf("A wasn't stored at zero page (0x00F2). Expected %v, got %v", cpu.A, cpu.Memory.Get(0x00F2))
	}

	cpu.Memory.Set(0x000D, OpSTA_absolute, 0x34, 0x12)
	cpu.Advance()

	if cpu.A != cpu.Memory.Get(0x1234) {
		t.Fatalf("A wasnt stored at absolute location (0x1234). Expected %v, got %v", cpu.A, cpu.Memory.Get(0x1234))
	}

	cpu.Memory.Set(0x0010, OpCMP_imm, 0x79)
	cpu.Advance()

	if !cpu.Flags.HasCarry() {
		t.Fatalf("Carry flag was not set after CMP with lower number")
	}

	cpu.Memory.Set(0x0012, OpCMP_imm, 0x81)
	cpu.Advance()

	if cpu.Flags.HasCarry() {
		t.Fatalf("Carry flag was set after CMP with higher number")
	}

	cpu.Memory.Set(0x0014, OpCMP_imm, 0x80)
	cpu.Advance()

	if !cpu.Flags.HasCarry() {
		t.Fatalf("Carry flag was not set after CMP with equal number")
	}
}
