package main

import (
	"flag"
	"fmt"
	"go-fmu/pkg/fmi2"
	"os"
)

func Run() error {
	if len(os.Args) < 2 {
		fmt.Println("expected 'dump' or 'bar' subcommands")
		os.Exit(1)
	}

	dumpCmd := flag.NewFlagSet("dump", flag.ExitOnError)
	dumpFilename := dumpCmd.String("filename", "", "filename")

	simulateCmd := flag.NewFlagSet("simulate", flag.ExitOnError)
	simulateFilename := simulateCmd.String("filename", "", "filename")

	switch os.Args[1] {
	case "dump":
		dumpCmd.Parse(os.Args[2:])
		if *dumpFilename == "" {
			dumpCmd.Usage()
			os.Exit(1)
		}

		err := fmi2.Dump(*dumpFilename, []string{"input", "output", "independent"})
		if err != nil {
			return err
		}

	case "simulate":
		simulateCmd.Parse(os.Args[2:])
		if *simulateFilename == "" {
			simulateCmd.Usage()
			os.Exit(1)
		}

		fmi2.SimulateFmu(*simulateFilename, fmi2.SimulationOptions{
			Initialize:   true,
			DebugLogging: true,
		})

	default:
		fmt.Println("expected 'dump' or 'bar' subcommands")
		os.Exit(1)
	}

	return nil
}

func main() {

	if err := Run(); err != nil {
		panic(err)
	}
}
