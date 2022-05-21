package main

import (
	"qa/cmd"
)

func main() {
	if err := cmd.RunUnsecure(); err != nil {
		panic(err)
	}
}
