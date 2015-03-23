package blast

import (
	"bufio"
	"fmt"
	"os"
)

// CLI runs the command
// line interface for the
// Blast interpreter
func CLI() {
	cr := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")
		input, err := cr.ReadString('\n')

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(newTokenStream(input[:len(input)-1]).combine().toRPN().evaluate().string())
	}
}
