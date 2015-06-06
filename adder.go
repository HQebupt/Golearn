package main

import "fmt"

func adder() func(int) int {
	sum := 0
	return func(x int) int {
		sum += x
		fmt.Println("Sum:",sum)
		return sum
	}
}

func main() {
	pos := adder()
	for i := 0; i < 4; i++ {
		fmt.Println(pos(i))
	}
}	