package go6502

import "testing"

func TestAccumulator(t *testing.T) {
	cpu := NewDefaultMemoryCPU()
	cpu.Memory.Set(ResetVectorL, 0x00, 0x00)
	cpu.Memory.Set(0xF0, 0x66) //Some initial value to load

	cpu.Memory.Set(0x0000, OpLDA_zeropage, 0xF0)
	cpu.Advance()
	if cpu.A != 0x66 {
		t.Fatalf("A wasn't loaded from zeropage. Expected %v, got %v", 0x66, cpu.A)
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
