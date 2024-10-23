package main

import (
	"log"

	"github.com/latvas/simple-x86-emulator/simpleEmulator"
)

func main() {
	shell := simpleEmulator.NewEmulatorShell()
	err := shell.LoadProgram()
	if err != nil {
		log.Fatalf("Error while read program: %v", err)
	}
	for {
		err = shell.ShellLoop()
		if err != nil {
			log.Fatalf("Error while execute emulator instruction: %v", err)
		}
	}
}
