package main

import (
	"fmt"
	"os"
)

var Debug = false // Toggle this to enable/disable debug output

func Printf(format string, a ...any) {
	if Debug {
		fmt.Printf(format, a...)
	}
}

func Println(a ...any) {
	if Debug {
		fmt.Println(a...)
	}
}

func Fatalf(format string, a ...any) {
	if Debug {
		fmt.Fprintf(os.Stderr, format, a...)
		os.Exit(1)
	}
}
