# Package Statistics Analyzer (Go Version)

## Overview
This Go program analyzes Debian package repositories, specifically targeting the "Contents" index files. It downloads the compressed Contents file for a specified architecture from a Debian mirror, parses it, and outputs statistics about the top packages containing the most files.

## Usage
1. To build the Program, compile the program using the go build command:
```bash
go build package_analyzer.go
```
This will generate an executable named package_analyzer.

2. To run the program, use the following command, where `<architecture>` is your target architecture (e.g., `amd64`):
```bash
./package_analyzer <architecture> -mirror <mirror-url> -top <number>
```

## Profiling
1. Run the program as described above.
2. While the program is waiting for the input to exit, open a new terminal window and run:
```bash
go tool pprof http://localhost:6060/debug/pprof/heap
```
3. At the pprof prompt, type top to see the top consumers of memory.
4. You can also generate a graphical representation of the memory profile by typing web at the pprof prompt. This requires Graphviz to be installed. There's a sample SVG graphical representation output in this repository.
5. When you've finished profiling, return to the program and press Enter to allow it to exit.

## Memory Usage Verification
The profiling step is crucial to verify that the program does not load the entire Contents file into memory. The provided profiling data shows the memory allocation during runtime, confirming efficient memory usage consistent with the program's design to process data line by line without loading the entire file into memory.
