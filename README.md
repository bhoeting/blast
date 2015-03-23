# blast
A "programming language" written in Go.  This project was made for fun, experimental purposes, and the ability to say "I wrote a programming language".  I don't actually know how to write a programming language, so this project doesn't reflect how I actually write code.

## Usage
`go get github.com/bhoeting/blast`

```go
package main

import (
	"github.com/bhoeting/blast"
)

func main() {
	blast.CLI()
}
```

```
>4+5
9
>4-5
-1
>4*5
20
>4/5  
0
>4.0/5
0.8
>4*"string"
"stringstringstringstring"
>4*"string"+5
"stringstringstringstring5"
>5+5*"string"
"5stringstringstringstringstring"
>(5+5)*"string"
"stringstringstringstringstringstringstringstringstringstring"
```
