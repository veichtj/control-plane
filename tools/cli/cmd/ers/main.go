package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kyma-project/control-plane/tools/cli/pkg/ers/command"
)

/*
 * Run:
 * go run cmd/ers/main.go
 */
func main() {
	setupCloseHandler()
	cmd := command.New()

	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}

}

func setupCloseHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-c
		fmt.Printf("\r- Signal '%v' received from Terminal. Exiting...\n ", sig)
		os.Exit(0)
	}()
}
