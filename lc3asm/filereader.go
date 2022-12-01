package lc3asm

// func main() {
// 	ch := make(chan byte)
// 	go func(ch chan byte) {
// 		//  Uncomment this block to actually read from stdin
// 		reader := bufio.NewScanner(os.Stdin)
// 		reader.Split(bufio.ScanBytes)
// 		for reader.Scan() {
// 			s := reader.Bytes()
// 			ch <- s[0]
// 		}
// 		close(ch)
// 		// Simulating stdin
// 		// ch <- "A line of text"
// 		// close(ch)
// 	}(ch)

// stdinloop:
// 	for {
// 		select {
// 		case stdin, ok := <-ch:
// 			if !ok {
// 				break stdinloop
// 			} else {
// 				fmt.Println("Read input from stdin:", stdin)
// 			}
// 		case <-time.After(1 * time.Second):
// 			// Do something when there is nothing read from stdin
// 		}
// 	}
// 	fmt.Println("Done, stdin must be closed")
// }
