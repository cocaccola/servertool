package main

import (
	"flag"
	"fmt"
	"os"

	"slack/servertool/internal/config"
)

var (
	ConfigFile = flag.String("config", "", "path to config file")
)

func main() {
	flag.Parse()

	if *ConfigFile == "" {
		fmt.Fprintln(os.Stderr, "please supply a path to a config file")
		os.Exit(2)
	}

	r, rm, err := config.Parse(*ConfigFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse config file: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Println(r)
	fmt.Println(rm)
}
