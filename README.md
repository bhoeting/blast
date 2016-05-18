# Blast
A programming language implemented in Go.

## How to use it
This project isn't packaged very well because it's not intended to be used for anything other than learning purposes.  The next steps will show you how to get the code to mess around with.

This You must have Go installed.  I recommend using [gvm](https://github.com/moovweb/gvm) for this.  Then,
	
	cd $GOPATH
	go get github.com/bhoeting/blast
	cd src/github.com/bhoeting/blast
	go get
	go install
	
Cool, now it's installed.  Now you should it.  Copy one of the example programs below into `program.blast`, then run `blast program.blast`.  If it worked, cool.  If not, this project is learning purposes, which means it's your duty to fix it.
	

## How it works
* Lexer scans each character and builds `Nodes`.  The types of nodes are
    * `nodeTypeUnkown`
        * There shouldn't be any of these.
    * `nodeTypeFuncCall`
        * `fizzbuzz`
    * `nodeTypeVariable`
        * `x`
    * `nodeTypeNumber`
        * `13`
    * `nodeTypeString`
        * `"hi"` 
    * `nodeTypeParen`
        * `)`
    * `nodeTypeBoolean`
        * `false`
        * `true`
    * `nodeTypeOperator`
    	* `+`
    	* `-`
    	* `*`
    	* `/`
    	* `^`
    	* `=`
    	* `==`
    	* `!=`
    	* `<`
    	* `<=`
    	* `>`
    	* `>=`
    	* `&&`
    	* `||`
    	* `->`
    	* `%`
    * `nodeTypeComma`
    	* `,`
    * `nodeTypeArgCount`
    	*  this is implicitly parsed and used for evaluating function calls
    * `nodeTypeReserved`
        * `for`
        * `end`
        * `function`
        * `return`
        * `else`
        * `if`
* The parser converts list of `Nodes` into reverse polish notation, then evaluates that expression.
* Each `line` belongs to a `block`, and each `block` belongs to another `block`.  Each `block` has its own `scope`.


## Example code
### Recursive FizzBuzz
```lua
function fizzbuzz(number)
	if number > 100
		return 0
	end

	fizz = number % 3 == 0
	buzz = number % 5 == 0

	if fizz
		print("Fizz")
	end

	if buzz
		print("Buzz")
	end

	if fizz == false && buzz == false
		print(number)
	end

	println()
	return fizzbuzz(number + 1)
end
```

### FizzBuzz with for loop
```lua
function fizzbuzzWithLoop()
	for 1 -> 200, number
		fizz = number % 3 == 0
		buzz = number % 5 == 0

		if fizz
			print("Fizz")
		end

		if buzz
			print("Buzz")
		end

		if fizz == false && buzz == false
			print(number)
		end

		println()
	end

	return 0
end
```

### Fibonacci with recursion
```lua
-- get the fifth fib number
println("Fibonacci with recursion =>", fib(5))

-- Return the fibonacci number at the 
-- specified index
function fib(index = 5, acc = 1, prev = 0)
	if index == 1 
		return acc
	end

	return fib(index - 1, acc + prev, acc)
end	
```
