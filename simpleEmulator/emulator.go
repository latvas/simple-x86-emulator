package simpleEmulator

import (
	"errors"
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

// newEmulator creates a new Emulator instance
func newEmulator() *emulator {
	return &emulator{
		registers: map[string]int{
			"rax": 0, "rbx": 0, "rcx": 0, "rdx": 0,
		},
		memory:    make(map[int]instruction),
		program:   make(map[int]string),
		pc:        0,
		memoryEnd: 0,
	}
}

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

func (e *emulator) step() error {
	if e.pc >= e.memoryEnd {
		return fmt.Errorf("program counter out of bounds")
	}

	// Decode current instruction
	inst := e.memory[int(e.pc)]

	// Execute based on instruction command
	var err error
	switch inst.command {
	case "ADD":
		err = e.executeADD(inst)
		e.pc++
	case "SBB":
		err = e.executeSBB(inst)
		e.pc++
	case "ADOX":
		err = e.executeADOX(inst)
		e.pc++
	case "JMP":
		err = e.executeJMP(inst)
	case "JGE":
		err = e.executeJGE(inst)
	default:
		err = fmt.Errorf("unknown instruction: %s", inst.command)
	}

	return err
}

func (e *emulator) executeADD(inst instruction) error {
	// ADD rax, rbx  OR  ADD rax, 5
	if len(inst.params) != 2 {
		return errors.New("ADD requires 2 parameters")
	}

	dstReg := inst.params[0]
	valueStr := inst.params[1]

	value, err := e.getValueFromOperand(valueStr)
	if err != nil {
		return fmt.Errorf("ADD failed: %v", err)
	}

	if _, exists := e.registers[dstReg]; !exists {
		return fmt.Errorf("unknown register: %s", dstReg)
	}

	e.registers[dstReg] += value
	return nil
}

func (e *emulator) executeSBB(inst instruction) error {
	// SBB rax, rbx  OR  SBB rax, 5
	if len(inst.params) != 2 {
		return errors.New("SBB requires 2 parameters")
	}

	dstReg := inst.params[0]
	valueStr := inst.params[1]

	value, err := e.getValueFromOperand(valueStr)
	if err != nil {
		return fmt.Errorf("SBB failed: %v", err)
	}

	if _, exists := e.registers[dstReg]; !exists {
		return fmt.Errorf("unknown register: %s", dstReg)
	}

	e.registers[dstReg] -= value
	return nil
}

func (e *emulator) executeADOX(inst instruction) error {
	// ADOX rax, rbx  OR  ADOX rax, 5
	if len(inst.params) != 2 {
		return errors.New("ADOX requires 2 parameters")
	}

	dstReg := inst.params[0]
	valueStr := inst.params[1]

	value, err := e.getValueFromOperand(valueStr)
	if err != nil {
		return fmt.Errorf("ADOX failed: %v", err)
	}

	if _, exists := e.registers[dstReg]; !exists {
		return fmt.Errorf("unknown register: %s", dstReg)
	}

	// ADOX - like ADD but also considers the overflow flag.
	e.registers[dstReg] += value
	// Overflow handling can be implemented here if necessary.
	return nil
}

func (e *emulator) executeJMP(inst instruction) error {
	// JMP 10 (Jump to address 10)
	if len(inst.params) != 1 {
		return errors.New("JMP requires 1 parameter")
	}

	targetAddrStr := inst.params[0]
	targetAddr, err := strconv.Atoi(targetAddrStr)
	if err != nil {
		return fmt.Errorf("invalid jump address: %w", err)
	}

	e.pc = uint(targetAddr)
	return nil
}

func (e *emulator) executeJGE(inst instruction) error {
	// JGE 10 (Jump to address 10 if rax >= rbx)
	if len(inst.params) != 1 {
		return errors.New("JGE requires 1 parameter")
	}

	if e.registers["rax"] >= e.registers["rbx"] {
		targetAddrStr := inst.params[0]
		targetAddr, err := strconv.Atoi(targetAddrStr)
		if err != nil {
			return fmt.Errorf("invalid jump address: %w", err)
		}

		e.pc = uint(targetAddr)
	}
	return nil
}

// getValueFromOperand обрабатывает операнд (регистры или целочисленные значения)
func (e *emulator) getValueFromOperand(operand string) (int, error) {
	if value, err := strconv.Atoi(operand); err == nil {
		// Это целое число
		return value, nil
	} else if value, exists := e.registers[operand]; exists {
		// Это регистр
		return value, nil
	} else {
		return 0, fmt.Errorf("invalid operand: %s", operand)
	}
}

func (e *emulator) printRegisters() {
	fmt.Println("Registers:")
	for reg, value := range e.registers {
		fmt.Printf("%s: %d\n", reg, value)
	}
}

func (e *emulator) printMemory(addr, size int) {
	fmt.Println("Memory:")
	for i := addr; i < addr+size; i++ {
		e.printInstruction(i)
	}
}

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
