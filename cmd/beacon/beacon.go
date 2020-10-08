package main

import "github.com/chen-keinan/beacon/internal/startup"

func main() {
	startup.InitCLI(startup.ArgsSanitizer)
}
