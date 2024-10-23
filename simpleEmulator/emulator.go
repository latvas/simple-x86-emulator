package simpleEmulator

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
