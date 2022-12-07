;; testing pseudo ops in assembler
.ORIG x3000
HALT
.STRINGZ "Hello, World!"
.BLKW 4 
.FILL x3939
.END