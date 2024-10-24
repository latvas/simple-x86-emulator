# simple-x86-emulator

This project is a simplified emulator of a single-core x86-64 processor. It can load and execute a program consisting of a few specific instructions and provides an interactive shell to control the emulator. The project is implemented in Go.

## Features

- Supports a limited set of x86-64 instructions: `ADD`, `SBB`, `ADOX`, `JMP`, `JGE`.
- Emulates basic CPU components: registers, memory, and program counter.
- Provides an interactive shell to control the emulator.
- Instructions are decoded from memory and executed step by step.
- The emulator interface supports loading programs, executing one instruction at a time, and viewing the state of registers and memory.

## Supported Instructions

- **ADD**: Adds values from two registers.
- **SBB**: Subtracts values from two registers.
- **ADOX**: Performs addition with an overflow flag.
- **JMP**: Unconditionally jumps to a specific program address.
- **JGE**: Jumps to a specific address if a condition (greater or equal) is met.

## Usage

### Requirements

- Go 1.20 or later

### Running the Emulator

1. Clone the repository:

    ```bash
    git clone https://example.com/pablo-emulator.git
    cd pablo-emulator
    ```

2. Build the project:

    ```bash
    go build -o emulator .
    ```

3. Run the emulator:

    ```bash
    ./emulator
    ```

4. Once running, you can use the interactive shell to load programs and step through instructions.

### Interactive Shell

After starting the emulator, the following commands are available in the interactive shell:

- **`load(addr, program)`**: Load a program into memory at the specified address.
- **`step()`**: Execute the next instruction in the program.
- **`print_registers()`**: Display the values of all registers.
- **`print_memory(addr, size)`**: Display memory contents from the specified address.
- **`print_instruction(addr)`**: Display the instruction at the specified address.

### Example

```bash
$ ./emulator
Enter starting memory address:
0
Enter program instructions:
ADD rax, 78         
SBB rcx, rdx        
ADOX rbx, 10        
JMP 4              
JGE 2               
END
Program loaded successfully.
> print_memory(0, 5)
Memory:
Instruction at 0: ADD rax 78 
Instruction at 1: SBB rcx rdx 
Instruction at 2: ADOX rbx 10 
Instruction at 3: JMP 4
Instruction at 4: JGE 2

> step()       # Execute ADD rax, 78 (rax should now be 78)
> print_registers()
Registers:
rax: 78
rbx: 0
rcx: 0
rdx: 0
> step()       # Execute SBB rcx, rdx (rcx remains unchanged)
> print_registers()
Registers:
rax: 78
rbx: 0
rcx: 0
rdx: 0

> step()       # Execute ADOX rbx, 10 (rbx should now be 10)
> print_registers()
Registers:
rax: 78
rbx: 10
rcx: 0
rdx: 0

> step()       # Execute JMP 4 (program counter jumps to instruction 4)
> step()       # Execute JGE 2 (if condition met, jumps to instruction 2)

> print_registers()
Registers:
rax: 78
rbx: 10
rcx: 0
rdx: 0

```