package main

import (
	"fmt"
	"time"
	"math/rand"
	"math"
)

func count(s string) {
	for i := 0; i < 10; i++ {
		fmt.Printf("%v %v\n", s, i)
		time.Sleep(100 * time.Millisecond)
	}
}
func shots() (b bool, shotsAmount float64) {
	shotsAmount = math.Round(rand.Float64()) * 5
	shotsAmount = shotsAmount * 5

	if b != true {
		return b, nil
	}
	return nil, shotsAmount
}

func main(){
	// go count("sheep")
	// count("mule")
	fmt.Println(shots())
}