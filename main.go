package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func main() {
	num := 10
	for j := 0; j < num; j++ {
		res := getRand(num)
		fmt.Println(res)
	}
}

func getRand(num int) float64 {
	var mu sync.Mutex
	mu.Lock()
	v := rand.Float64()
	mu.Unlock()
	return v
}
