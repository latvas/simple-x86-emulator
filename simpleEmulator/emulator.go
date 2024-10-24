package simpleEmulator

import (
	"fmt"
	"strconv"
	"strings"
)

type instruction struct {
	command string
	params  map[int]string
}

// emulator struct representing the CPU state
type emulator struct {
	registers map[string]int
	flags     map[string]bool
	memory    map[int]instruction
	program   map[int]string
	pc        uint // Program Counter
	memoryEnd uint
}

// newEmulator creates a new Emulator instance with registers and flags initialized
func newEmulator() *emulator {
	return &emulator{
		// Initializing registers to 0
		registers: map[string]int{
			"rax": 0, "rbx": 0, "rcx": 0, "rdx": 0,
		},
		// Initializing flags to false
		flags: map[string]bool{
			"ZF": false, // Zero Flag
			"SF": false, // Sign Flag
			"CF": false, // Carry Flag
			"OF": false, // Overflow Flag
		},
		memory:    make(map[int]instruction),
		program:   make(map[int]string),
		pc:        0, // Starting Program Counter
		memoryEnd: 0, // Memory end limit
	}
}

// load loads the program into memory at a given address
func (e *emulator) load(addr uint, program []string) {
	for i, line := range program {
		inst := instruction{
			params: make(map[int]string),
		}
		tokens := strings.Split(line, " ")
		inst.command = strings.Trim(tokens[0], ", ")
		for i := 1; i < len(tokens); i++ {
			inst.params[i-1] = strings.Trim(tokens[i], ", ")
		}
		e.memory[int(addr)+i] = inst
	}
	e.memoryEnd = addr + uint(len(program))
}

// step executes one instruction at the current program counter
func (e *emulator) step() error {
	if e.pc >= e.memoryEnd {
		return fmt.Errorf("program counter out of bounds")
	}

	// Decode current instruction from memory
	inst := e.memory[int(e.pc)]

	// Execute based on the instruction command
	var err error
	switch inst.command {
	case "ADD":
		err = e.executeADD(inst) // Execute ADD instruction
		e.pc++                   // Move to the next instruction
	case "SBB":
		err = e.executeSBB(inst) // Execute SBB instruction
		e.pc++
	case "ADOX":
		err = e.executeADOX(inst) // Execute ADOX instruction
		e.pc++
	case "JMP":
		err = e.executeJMP(inst) // Execute JMP instruction
	case "JGE":
		err = e.executeJGE(inst) // Execute JGE instruction
	default:
		err = fmt.Errorf("unknown instruction: %s", inst.command)
	}

	return err
}

// executeADD implements the ADD operation and updates the flags
func (e *emulator) executeADD(inst instruction) error {
	// Retrieve the operands (register or immediate value)
	dest, err := e.getValueFromOperand(inst.params[0])
	if err != nil {
		return err
	}
	src, err := e.getValueFromOperand(inst.params[1])
	if err != nil {
		return err
	}

	// Perform addition
	result := dest + src

	// Set flags
	e.flags["ZF"] = (result == 0)                // Zero Flag
	e.flags["SF"] = (result < 0)                 // Sign Flag
	e.flags["CF"] = (result < dest)              // Carry Flag (unsigned overflow)
	e.flags["OF"] = ((dest < 0) == (src < 0)) && // Overflow Flag (signed overflow)
		(result < 0 != (dest < 0))

	// Store the result back into the destination register
	e.registers[inst.params[0]] = result

	return nil
}

// executeSBB implements the SBB operation (subtract with borrow) and updates the flags
func (e *emulator) executeSBB(inst instruction) error {
	// Retrieve the operands
	dest, err := e.getValueFromOperand(inst.params[0])
	if err != nil {
		return err
	}
	src, err := e.getValueFromOperand(inst.params[1])
	if err != nil {
		return err
	}

	// Perform subtraction with borrow (borrow is 1 if CF is set)
	borrow := 0
	if e.flags["CF"] {
		borrow = 1
	}
	result := dest - src - borrow

	// Set flags
	e.flags["ZF"] = (result == 0)                       // Zero Flag
	e.flags["SF"] = (result < 0)                        // Sign Flag
	e.flags["CF"] = (dest < src+borrow)                 // Carry Flag (unsigned borrow)
	e.flags["OF"] = ((dest < 0) != (src+borrow < 0)) && // Overflow Flag
		(result < 0 != (dest < 0))

	// Store the result back into the destination register
	e.registers[inst.params[0]] = result

	return nil
}

// executeADOX implements the ADOX operation (unsigned addition) and updates only the OF flag
func (e *emulator) executeADOX(inst instruction) error {
	// Retrieve the operands
	dest, err := e.getValueFromOperand(inst.params[0])
	if err != nil {
		return err
	}
	src, err := e.getValueFromOperand(inst.params[1])
	if err != nil {
		return err
	}

	// Perform addition
	result := dest + src

	// Set Overflow Flag (OF)
	e.flags["OF"] = (result < dest) // Only check unsigned overflow

	// Store the result back into the destination register
	e.registers[inst.params[0]] = result

	return nil
}

// executeJMP implements the JMP operation (unconditional jump)
func (e *emulator) executeJMP(inst instruction) error {
	// Jump to the address specified by the operand
	addr, err := e.getValueFromOperand(inst.params[0])
	if err != nil {
		return err
	}

	// Set the program counter to the new address
	e.pc = uint(addr)

	return nil
}

// executeJGE implements the JGE operation (jump if greater or equal)
func (e *emulator) executeJGE(inst instruction) error {
	// Jump if SF == OF (greater or equal condition for signed numbers)
	if e.flags["SF"] == e.flags["OF"] {
		addr, err := e.getValueFromOperand(inst.params[0])
		if err != nil {
			return err
		}
		e.pc = uint(addr) // Set PC to the target address
	} else {
		// Otherwise, proceed to the next instruction
		e.pc++
	}

	return nil
}

// getValueFromOperand processes the operand (register or integer value)
func (e *emulator) getValueFromOperand(operand string) (int, error) {
	if value, err := strconv.Atoi(operand); err == nil {
		// Operand is an integer value
		return value, nil
	} else if value, exists := e.registers[operand]; exists {
		// Operand is a register
		return value, nil
	} else {
		// Invalid operand
		return 0, fmt.Errorf("invalid operand: %s", operand)
	}
}

// printRegisters displays the current state of all registers
func (e *emulator) printRegisters() {
	fmt.Println("Registers:")
	for reg, value := range e.registers {
		fmt.Printf("%s: %d\n", reg, value)
	}
}

// printMemory prints the memory content from a specific address and size
func (e *emulator) printMemory(addr, size int) {
	fmt.Println("Memory:")
	for i := addr; i < addr+size; i++ {
		e.printInstruction(i)
	}
}

// printInstruction prints the instruction at a given memory address
func (e *emulator) printInstruction(addr int) {
	if inst, ok := e.memory[addr]; ok {
		instString := fmt.Sprintf("%s ", inst.command)
		for _, param := range inst.params {
			instString += fmt.Sprintf("%s ", param)
		}
		fmt.Printf("Instruction at %d: %s\n", addr, instString)
	} else {
		fmt.Printf("No instruction at %d\n", addr)
	}
}
