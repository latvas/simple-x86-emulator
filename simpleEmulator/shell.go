package simpleEmulator

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type EmulatorShell struct {
	emu emulator
}

// LoadProgram загружает программу в эмулятор
func (s *EmulatorShell) LoadProgram() error {
	//TODO change reader to stdin
	reader := bufio.NewReader(os.Stdin)
	// file, _ := os.Open("program.txt")
	//reader := bufio.NewReader(file)

	fmt.Println("Enter starting memory address:")
	startAddrInput, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("cannot read adress: %w", err)
	}
	startAddrInput = strings.TrimSpace(startAddrInput)
	startAddr, err := strconv.ParseUint(startAddrInput, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid starting address: %w", err)
	}

	fmt.Println("Enter program instructions:")
	var program []string

	// Use a loop to read lines until "END" is encountered
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("error while reading program: %w", err)
		}

		// Trim whitespace and check for "END"
		line = strings.TrimSpace(line)
		if line == "END" {
			break // Exit the loop if "END" is encountered
		}

		if line != "" { // Optional: Ignore empty lines
			program = append(program, line)
		}
	}

	// Загрузка программы в эмулятор
	s.emu.load(uint(startAddr), program)
	fmt.Println("Program loaded successfully.")
	return nil
}

// ShellLoop реализует основной цикл для работы с эмулятором
func (s *EmulatorShell) ShellLoop() error {
	//TODO change reader to stdin
	reader := bufio.NewReader(os.Stdin)
	//file, _ := os.Open("emu_cmd.txt")
	//reader := bufio.NewReader(file)

	for {
		fmt.Print("> ") // Вывод приглашения в командной строке
		input, err := reader.ReadString('\n')
		if err == io.EOF {
			return fmt.Errorf("EOF while read instruction")
		}
		if err != nil {
			return fmt.Errorf("error while read instruction")
		}
		input = strings.TrimSpace(input)

		// Обработка команд
		switch {
		case strings.HasPrefix(input, "step"):
			err := s.emu.step()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		case strings.HasPrefix(input, "print_registers"):
			s.emu.printRegisters()
		case strings.HasPrefix(input, "print_memory"):
			// Пример команды: print_memory(10, 5)

			// Regular expression to find two groups of digits
			re := regexp.MustCompile(`\d+`)

			// Find all matches of digits
			matches := re.FindAllString(input, -1)

			// Ensure there are exactly two matches
			if len(matches) != 2 {
				fmt.Println("Usage: print_memory <addr> <size>")
			}

			// Convert the string matches to integers
			addr, err := strconv.Atoi(matches[0])
			if err != nil {
				fmt.Printf("Invalid address: %v\n", err)
			}
			size, err := strconv.Atoi(matches[1])
			if err != nil {
				fmt.Printf("Invalid size: %v\n", err)
			}
			s.emu.printMemory(addr, size)

		case strings.HasPrefix(input, "print_instruction"):
			// Пример команды: print_instruction(10)

			// Regular expression to find digits inside parentheses
			re := regexp.MustCompile(`\d+`)

			// Find the first match of digits
			match := re.FindString(input)

			// Convert the string match to an integer
			addr, err := strconv.Atoi(match)

			if err != nil {
				fmt.Printf("Invalid address: %v\n", err)
			}
			s.emu.printInstruction(addr)
		case input == "exit":
			fmt.Println("Exiting shell.")
			return nil
		default:
			fmt.Println("Unknown command. Available commands: step, print_registers, print_memory, print_instruction, exit")
		}
	}
}

// NewEmulatorShell создает новую оболочку для эмулятора
func NewEmulatorShell() *EmulatorShell {
	return &EmulatorShell{
		emu: *newEmulator(),
	}
}
