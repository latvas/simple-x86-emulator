package simpleEmulator

import (
	"testing"
)

// TestEmulator_ADD tests the ADD instruction
func TestEmulator_ADD(t *testing.T) {
	emulator := newEmulator()

	// Program to add two values in rax and rbx
	program := []string{
		"ADD rax, 10", // rax = rax + 10
		"ADD rbx, 5",  // rbx = rbx + 5
	}

	// Load the program starting at address 0
	emulator.load(0, program)

	// Execute each instruction step by step
	err := emulator.step()
	if err != nil {
		t.Fatalf("step failed: %v", err)
	}
	if emulator.registers["rax"] != 10 {
		t.Fatalf("expected rax=10, got rax=%d", emulator.registers["rax"])
	}

	err = emulator.step()
	if err != nil {
		t.Fatalf("step failed: %v", err)
	}
	if emulator.registers["rbx"] != 5 {
		t.Fatalf("expected rbx=5, got rbx=%d", emulator.registers["rbx"])
	}

	// Check flags
	if emulator.flags["ZF"] {
		t.Fatal("expected ZF to be false, got true")
	}
}

// TestEmulator_SBB tests the SBB instruction
func TestEmulator_SBB(t *testing.T) {
	emulator := newEmulator()

	// Initialize rax and set the carry flag
	emulator.registers["rax"] = 10
	emulator.flags["CF"] = true

	// Program to subtract from rax with borrow
	program := []string{
		"SBB rax, 5", // rax = rax - 5 - CF (carry flag)
	}

	// Load the program starting at address 0
	emulator.load(0, program)

	// Execute the instruction
	err := emulator.step()
	if err != nil {
		t.Fatalf("step failed: %v", err)
	}

	// Check the result of rax
	expectedRAX := 4 // 10 - 5 - 1 (borrow)
	if emulator.registers["rax"] != expectedRAX {
		t.Fatalf("expected rax=%d, got rax=%d", expectedRAX, emulator.registers["rax"])
	}

	// Check flags
	if emulator.flags["ZF"] {
		t.Fatal("expected ZF to be false, got true")
	}
}

// TestEmulator_ADOX tests the ADOX instruction
func TestEmulator_ADOX(t *testing.T) {
	emulator := newEmulator()

	// Program to perform ADOX on rax and a value
	program := []string{
		"ADOX rax, 5", // rax = rax + 5 (unsigned addition)
	}

	// Load the program starting at address 0
	emulator.load(0, program)

	// Execute the instruction
	err := emulator.step()
	if err != nil {
		t.Fatalf("step failed: %v", err)
	}

	// Check the result of rax
	expectedRAX := 5 // rax starts at 0 and is incremented by 5
	if emulator.registers["rax"] != expectedRAX {
		t.Fatalf("expected rax=%d, got rax=%d", expectedRAX, emulator.registers["rax"])
	}

	// Check flags (OF should be false, no overflow)
	if emulator.flags["OF"] {
		t.Fatal("expected OF to be false, got true")
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

	// Load the program starting at address 0
	emulator.load(0, program)

	// Execute the JMP instruction
	err := emulator.step()
	if err != nil {
		t.Fatalf("step failed: %v", err)
	}

	// Check if program counter has jumped correctly
	if emulator.pc != 2 {
		t.Fatalf("expected pc=2, got pc=%d", emulator.pc)
	}

	// Execute the next instruction (after the jump)
	err = emulator.step()
	if err != nil {
		t.Fatalf("step failed: %v", err)
	}
	if emulator.registers["rbx"] != 10 {
		t.Fatalf("expected rbx=10, got rbx=%d", emulator.registers["rbx"])
	}
}

// TestEmulator_JGE tests the JGE instruction
func TestEmulator_JGE(t *testing.T) {
	emulator := newEmulator()

	// Program to test the jump if greater or equal (JGE)
	program := []string{
		"ADD rax, -5", // rax = -5
		"JGE 3",       // Jump to instruction at address 3 if rax >= 0 (SF == OF)
		"ADD rbx, 10", // This should be skipped
		"ADD rbx, 20", // This should be executed
	}

	// Load the program starting at address 0
	emulator.load(0, program)

	// Step through the first instruction (ADD)
	err := emulator.step()
	if err != nil {
		t.Fatalf("step failed: %v", err)
	}

	// Now the SF != OF, so JGE should NOT jump
	err = emulator.step()
	if err != nil {
		t.Fatalf("step failed: %v", err)
	}

	// Check the program counter to ensure it skipped the correct instruction
	if emulator.pc != 3 {
		t.Fatalf("expected pc=3, got pc=%d", emulator.pc)
	}

	// Step through the final instruction (ADD rbx, 20)
	err = emulator.step()
	if err != nil {
		t.Fatalf("step failed: %v", err)
	}
	if emulator.registers["rbx"] != 20 {
		t.Fatalf("expected rbx=20, got rbx=%d", emulator.registers["rbx"])
	}
}
