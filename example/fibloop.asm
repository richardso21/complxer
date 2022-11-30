;; Return fib(N) on R0
.ORIG x3000
  AND R0, R0, 0 ;; R0 = 0
  AND R1, R1, 0
  ADD R1, R1, 1 ;; R1 = 1
  LD R2, N      ;; R2 = N
  ADD R2, R2, -1
FIB_LOOP
  BRnz FIB_END
  ADD R3, R0, R1  ;; R3 = R0 + R1
  ADD R1, R0, 0   ;; R1 = R0
  ADD R0, R3, 0   ;; R0 = R3 (R0 + R1)
  ADD R2, R2, -1  ;; R2--
  BR FIB_LOOP
FIB_END
  HALT
N .FILL 8
.END
