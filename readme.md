# Complxer

[Complx](https://github.com/TricksterGuy/complx-tools) is a multi-platform Little
Computer 3 (LC-3) simulator and assembler on C++.

Complxer is an ambitious recreation of complx using Golang, albeit some
implementation details for the LC-3 are skimmed over to keep my own sanity.
I don't expect this to replace complx at any time, but I had a great time in
CS 2110 and wondered if I could reverse engineer this tool used so much
throughout the course.

## TODO

- Stepping/Debug mode for LC3
- **GUI**
- Interrupts, Async I/O
- Iron out assembler bugs

### Done

- LC3 Assembler (usable for most programs)
- Mock Registers & Memory
- Support for all op-Codes & Traps
- Memory-Mapped I/O
- Read and Load Object Files
