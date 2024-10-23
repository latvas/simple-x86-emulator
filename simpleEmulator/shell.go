package simpleEmulator

// implement interactive shell for emulator

type EmulatorShell struct {
	emu emulator
}

func (s *EmulatorShell) LoadProgram() error {
	return nil
}

func (s *EmulatorShell) ShellLoop() error {
	return nil
}

func NewEmulatorShell() *EmulatorShell {
	return &EmulatorShell{
		emu: *newEmulator(),
	}
}
