package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("hello")

	// Keep the application running until interrupted
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Composer API is running. Press Ctrl+C to exit.")
	<-sigChan

	fmt.Println("Shutting down gracefully...")
}
