package simpleEmulator

import (
	"testing"
)

// TestEmulator_ADD tests the ADD instruction and checks flags
func TestEmulator_ADD(t *testing.T) {
	emulator := newEmulator()

	// Program to add two values in rax and rbx
	program := []string{
		"ADD rax, 10",  // rax = rax + 10
		"ADD rbx, -10", // rbx = rbx + (-10)
	}

	// Load the program starting at address 0
	emulator.load(0, program)

	// Execute ADD rax, 10
	err := emulator.step()
	if err != nil {
		t.Fatalf("step failed: %v", err)
	}
	if emulator.registers["rax"] != 10 {
		t.Fatalf("expected rax=10, got rax=%d", emulator.registers["rax"])
	}
	// Flags check: no overflow or zero, but sign flag is clear
	if emulator.flags["ZF"] || emulator.flags["SF"] || emulator.flags["OF"] {
		t.Fatal("expected ZF=0, SF=0, OF=0")
	}

	// Execute ADD rbx, -10
	err = emulator.step()
	if err != nil {
		t.Fatalf("step failed: %v", err)
	}
	if emulator.registers["rbx"] != -10 {
		t.Fatalf("expected rbx=-10, got rbx=%d", emulator.registers["rbx"])
	}
	// Check if the sign flag is set correctly (SF=1 for negative values)
	if !emulator.flags["SF"] {
		t.Fatal("expected SF=1 for negative result, got SF=0")
	}
	// Zero flag should be false, overflow flag should be false
	if emulator.flags["ZF"] || emulator.flags["OF"] {
		t.Fatal("expected ZF=0, OF=0")
	}
}

// TestEmulator_SBB tests the SBB instruction and checks flags
func TestEmulator_SBB(t *testing.T) {
	emulator := newEmulator()

	// Initialize rax and set the carry flag
	emulator.registers["rax"] = 10
	emulator.flags["CF"] = true

	// Program to subtract with borrow
	program := []string{
		"SBB rax, 5", // rax = rax - 5 - CF (carry flag)
	}

	// Load the program
	emulator.load(0, program)

	// Execute SBB rax, 5
	err := emulator.step()
	if err != nil {
		t.Fatalf("step failed: %v", err)
	}

	// Expected result: 10 - 5 - 1 = 4
	expectedRAX := 4
	if emulator.registers["rax"] != expectedRAX {
		t.Fatalf("expected rax=%d, got rax=%d", expectedRAX, emulator.registers["rax"])
	}

	// Flag checks: CF should be 0 (no borrow), ZF=0, SF=0, OF=0
	if emulator.flags["CF"] || emulator.flags["ZF"] || emulator.flags["SF"] || emulator.flags["OF"] {
		t.Fatal("expected CF=0, ZF=0, SF=0, OF=0")
	}
}

// TestEmulator_ADOX tests the ADOX instruction and checks flags
func TestEmulator_ADOX(t *testing.T) {
	emulator := newEmulator()

	// Program to perform ADOX on rax
	program := []string{
		"ADOX rax, 5", // rax = rax + 5 (unsigned addition)
	}

	// Load the program
	emulator.load(0, program)

	// Execute ADOX rax, 5
	err := emulator.step()
	if err != nil {
		t.Fatalf("step failed: %v", err)
	}

	// Expected result: rax = 5
	expectedRAX := 5
	if emulator.registers["rax"] != expectedRAX {
		t.Fatalf("expected rax=%d, got rax=%d", expectedRAX, emulator.registers["rax"])
	}

	// Flags check: no overflow, zero, or carry
	if emulator.flags["OF"] || emulator.flags["CF"] || emulator.flags["ZF"] {
		t.Fatal("expected OF=0, CF=0, ZF=0")
	}
}

// TestEmulator_JMP tests the JMP instruction
func TestEmulator_JMP(t *testing.T) {
	emulator := newEmulator()

	// Program to jump to an address
	program := []string{
		"JMP 2",       // Jump to instruction at address 2
		"ADD rax, 5",  // This should be skipped
		"ADD rbx, 10", // This should be executed after the jump
	}

	// Load the program
	emulator.load(0, program)

	// Execute JMP 2
	err := emulator.step()
	if err != nil {
		t.Fatalf("step failed: %v", err)
	}

	// Check if program counter has jumped to the correct location
	if emulator.pc != 2 {
		t.Fatalf("expected pc=2, got pc=%d", emulator.pc)
	}

	// Execute ADD rbx, 10 (after jump)
	err = emulator.step()
	if err != nil {
		t.Fatalf("step failed: %v", err)
	}
	if emulator.registers["rbx"] != 10 {
		t.Fatalf("expected rbx=10, got rbx=%d", emulator.registers["rbx"])
	}
}

// TestEmulator_JGE tests the JGE instruction and flag checks
func TestEmulator_JGE(t *testing.T) {
	emulator := newEmulator()

	// Program to test the JGE (jump if greater or equal) instruction
	program := []string{
		"ADD rax, -5", // rax = -5 (SF=1)
		"JGE 3",       // Jump if rax >= 0 (this shouldn't jump since rax is negative)
		"ADD rbx, 10", // Should be executed if no jump occurs
		"ADD rbx, 20", // Should be executed if jump occurs
	}

	// Load the program
	emulator.load(0, program)

	// Step 1: ADD rax, -5
	err := emulator.step()
	if err != nil {
		t.Fatalf("step failed: %v", err)
	}
	if emulator.registers["rax"] != -5 {
		t.Fatalf("expected rax=-5, got rax=%d", emulator.registers["rax"])
	}
	// Flag check: SF should be set, OF=0 (no overflow), ZF=0
	if !emulator.flags["SF"] || emulator.flags["OF"] || emulator.flags["ZF"] {
		t.Fatal("expected SF=1, OF=0, ZF=0")
	}

	// Step 2: JGE (should not jump since rax < 0)
	err = emulator.step()
	if err != nil {
		t.Fatalf("step failed: %v", err)
	}
	if emulator.pc != 2 {
		t.Fatalf("expected pc=2 (no jump), got pc=%d", emulator.pc)
	}

	// Step 3: ADD rbx, 10
	err = emulator.step()
	if err != nil {
		t.Fatalf("step failed: %v", err)
	}
	if emulator.registers["rbx"] != 10 {
		t.Fatalf("expected rbx=10, got rbx=%d", emulator.registers["rbx"])
	}

	// Check flags (no change expected in this case)
	if emulator.flags["ZF"] || emulator.flags["OF"] || emulator.flags["CF"] || emulator.flags["SF"] {
		t.Fatal("expected all flags to remain unchanged")
	}
}

// TestEmulator_JGE_JumpOccurs tests the JGE instruction where the jump occurs
func TestEmulator_JGE_JumpOccurs(t *testing.T) {
	emulator := newEmulator()

	// Program to test the JGE (jump if greater or equal) instruction
	program := []string{
		"ADD rax, 5",  // rax = 5 (SF=0, ZF=0)
		"JGE 3",       // Jump to instruction 3 if rax >= 0 (SF=0)
		"ADD rbx, 10", // Should be skipped due to jump
		"ADD rbx, 20", // Should be executed after jump
	}

	// Load the program
	emulator.load(0, program)

	// Step 1: ADD rax, 5
	err := emulator.step()
	if err != nil {
		t.Fatalf("step failed: %v", err)
	}
	if emulator.registers["rax"] != 5 {
		t.Fatalf("expected rax=5, got rax=%d", emulator.registers["rax"])
	}
	// Flag check: SF=0 (positive value), ZF=0, OF=0
	if emulator.flags["SF"] || emulator.flags["ZF"] || emulator.flags["OF"] {
		t.Fatal("expected SF=0, ZF=0, OF=0 after ADD rax, 5")
	}

	// Step 2: JGE (should jump to instruction 3 because rax >= 0)
	err = emulator.step()
	if err != nil {
		t.Fatalf("step failed: %v", err)
	}
	if emulator.pc != 3 {
		t.Fatalf("expected pc=3 (jump occurred), got pc=%d", emulator.pc)
	}

	// Step 3: ADD rbx, 20 (executed after jump)
	err = emulator.step()
	if err != nil {
		t.Fatalf("step failed: %v", err)
	}
	if emulator.registers["rbx"] != 20 {
		t.Fatalf("expected rbx=20, got rbx=%d", emulator.registers["rbx"])
	}

	// Check flags after the final step (no change expected)
	if emulator.flags["ZF"] || emulator.flags["OF"] || emulator.flags["CF"] || emulator.flags["SF"] {
		t.Fatal("expected all flags to remain unchanged after final step")
	}
}
