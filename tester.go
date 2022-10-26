package main

import "go6502/go6502"

func main() {
	cpu := go6502.NewDefaultMemoryCPU()

	//Reset vector
	cpu.Memory.Set(go6502.ResetVectorL, 0x00, 0x00)

	cpu.Memory.Set(0x0000, go6502.OpLDA_imm, 0x04)
	cpu.Memory.Set(0x0002, go6502.OpCMP_imm, 0x03)

	cpu.Advance()
	cpu.Advance()
}
