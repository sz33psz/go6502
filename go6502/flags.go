package go6502

import "strings"

type Flags struct {
	val uint8
}

func (f *Flags) HasFlag(index int) bool {
	return f.val&(1<<index) > 0
}

func (f *Flags) SetFlag(index int) {
	f.val = f.val | (1 << index)
}

func (f *Flags) ClearFlag(index int) {
	f.val = f.val & (0xFF - (1 << index))
}

func (f *Flags) HasCarry() bool {
	return f.HasFlag(0)
}

func (f *Flags) SetCarry() {
	f.SetFlag(0)
}

func (f *Flags) ClearCarry() {
	f.ClearFlag(0)
}

func (f *Flags) HasZero() bool {
	return f.HasFlag(1)
}

func (f *Flags) SetZero() {
	f.SetFlag(1)
}

func (f *Flags) ClearZero() {
	f.ClearFlag(1)
}

func (f *Flags) HasInterruptDisable() bool {
	return f.HasFlag(2)
}

func (f *Flags) SetInterruptDisable() {
	f.SetFlag(2)
}

func (f *Flags) ClearInterruptDisable() {
	f.ClearFlag(2)
}

func (f *Flags) HasDecimal() bool {
	return f.HasFlag(3)
}

func (f *Flags) SetDecimal() {
	f.SetFlag(3)
}

func (f *Flags) ClearDecimal() {
	f.ClearFlag(3)
}

func (f *Flags) HasOverflow() bool {
	return f.HasFlag(6)
}

func (f *Flags) SetOverflow() {
	f.SetFlag(6)
}

func (f *Flags) ClearOverflow() {
	f.ClearFlag(6)
}

func (f *Flags) HasNegative() bool {
	return f.HasFlag(7)
}

func (f *Flags) SetNegative() {
	f.SetFlag(7)
}

func (f *Flags) ClearNegative() {
	f.ClearFlag(7)
}

func (f *Flags) String() string {
	builder := strings.Builder{}
	if f.HasNegative() {
		builder.WriteString("N")
	} else {
		builder.WriteString("n")
	}

	if f.HasOverflow() {
		builder.WriteString("V")
	} else {
		builder.WriteString("v")
	}

	builder.WriteString("--")

	if f.HasDecimal() {
		builder.WriteString("D")
	} else {
		builder.WriteString("d")
	}

	if f.HasInterruptDisable() {
		builder.WriteString("I")
	} else {
		builder.WriteString("i")
	}

	if f.HasZero() {
		builder.WriteString("Z")
	} else {
		builder.WriteString("z")
	}

	if f.HasCarry() {
		builder.WriteString("C")
	} else {
		builder.WriteString("c")
	}
	return builder.String()
}
