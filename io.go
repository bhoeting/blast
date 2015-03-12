package blast

import (
	"bufio"
	"fmt"
	"os"
)

func CLI(args []string) {
	b := NewBlast()
	cr := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")
		input, err := cr.ReadString('\n')

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(b.ParseLine(input[:len(input)-1]))
	}
}
