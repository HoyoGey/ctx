package main

import (
	"fmt"
	"time"

	"github.com/HoyoGey/ctx"
)

func main() {
	// Current time example
	now := time.Now()
	ct := ctx.NewCTX(now)
	fmt.Printf("Current time: %v\n", now)
	fmt.Printf("CTX bytes: % X\n", ct.Bytes())
	
	// Future time example
	future := time.Now().AddDate(10, 0, 0) // 10 years in the future
	futureCt := ctx.NewCTX(future)
	fmt.Printf("\nFuture time: %v\n", future)
	fmt.Printf("CTX bytes: % X\n", futureCt.Bytes())
	
	// Binary storage example
	bytes := ct.Bytes()
	restored := ctx.FromBytes(bytes)
	fmt.Printf("\nRestored time: %v\n", restored.Time())
}
