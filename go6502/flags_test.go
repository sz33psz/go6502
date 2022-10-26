package go6502

import "testing"

func TestFlags(t *testing.T) {
	flags := Flags{}
	if flags.HasCarry() {
		flagShouldBeCleared(t, "C")
	}
	flags.SetCarry()
	if !flags.HasCarry() {
		flagShouldBeSet(t, "C")
	}
	flags.ClearCarry()
	if flags.HasCarry() {
		flagShouldBeCleared(t, "C")
	}

	if flags.HasOverflow() {
		flagShouldBeCleared(t, "V")
	}
	flags.SetOverflow()
	if !flags.HasOverflow() {
		flagShouldBeSet(t, "V")
	}
	flags.ClearOverflow()
	if flags.HasOverflow() {
		flagShouldBeCleared(t, "V")
	}

	flags.SetOverflow()
	flags.SetNegative()
	flags.SetCarry()
	flags.ClearNegative()
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
