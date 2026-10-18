package main

import (
	"fmt"
	"io"
	"os"

	"vihkat/fleet"
)

func main() {
	if err := run(os.Stdin, os.Stdout, fleet.ServiceCenter{}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(input io.Reader, output io.Writer, inspector fleet.Inspector) error {
	container := newServiceContainer()
	container.register("inspector", func(_ *serviceContainer) interface{} {
		return inspector
	})
	container.register("fleet", func(c *serviceContainer) interface{} {
		return fleet.NewFleet(c.resolve("inspector").(fleet.Inspector))
	})
	container.register("console", func(c *serviceContainer) interface{} {
		return newConsoleApp(c.resolve("fleet").(*fleet.Fleet), input, output)
	})
	return container.resolve("console").(*consoleApp).Run()
}
