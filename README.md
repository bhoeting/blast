# Blast
A programming language implemented in Go.

## Examples
### Recursive FizzBuzz
```lua
function fizzbuzz(number)
	if number > 100
		return 0
	end

	fizz = modulus(number, 3) == 0
	buzz = modulus(number, 5) == 0

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
		fizz = modulus(number, 3) == 0
		buzz = modulus(number, 5) == 0

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
