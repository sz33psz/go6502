package go6502

import "testing"

func TestFlags(t *testing.T) {
	flags := Flags{}
	if flags.HasCarry() {
		flagShouldBeCleared(t, "C")
	}
	flags.SetCarry(true)
	if !flags.HasCarry() {
		flagShouldBeSet(t, "C")
	}
	flags.SetCarry(false)
	if flags.HasCarry() {
		flagShouldBeCleared(t, "C")
	}

	if flags.HasOverflow() {
		flagShouldBeCleared(t, "V")
	}
	flags.SetOverflow(true)
	if !flags.HasOverflow() {
		flagShouldBeSet(t, "V")
	}
	flags.SetOverflow(false)
	if flags.HasOverflow() {
		flagShouldBeCleared(t, "V")
	}

	flags.SetOverflow(true)
	flags.SetNegative(true)
	flags.SetCarry(true)
	flags.SetNegative(false)
	if !flags.HasOverflow() {
		flagShouldBeSet(t, "V")
	}
	if !flags.HasCarry() {
		flagShouldBeSet(t, "C")
	}
	if flags.HasNegative() {
		flagShouldBeCleared(t, "N")
	}
}

func flagShouldBeCleared(t *testing.T, flag string) {
	t.Fatalf("Should have %v clear, but it is set", flag)
}

func flagShouldBeSet(t *testing.T, flag string) {
	t.Fatalf("Should have %v set, but it is clear", flag)
}
