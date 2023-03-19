;; testing pseudo ops in assembler
.ORIG x3000
LEA R0, S
PUTS
HALT
S 
.STRINGZ "Hello, World!"
.BLKW 4 
.FILL x3939
.END