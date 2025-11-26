package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/oktalz/junit-report"
)

func main() {
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	status := addCmd.String("status", "ok", "Status of the check (ok, failed)")
	file := addCmd.String("file", "", "File path related to the check")
	message := addCmd.String("message", "", "Message for the check")
	description := addCmd.String("description", "", "Description/Output for the check")
	outputFile := addCmd.String("output", "", "Output JUnit XML file (defaults to JUNIT_FILE env var)")
	suiteName := addCmd.String("suite", "General", "Test suite name")

	if len(os.Args) < 2 {
		fmt.Println("expected 'add' subcommand")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])
	default:
		fmt.Println("expected 'add' subcommand")
		os.Exit(1)
	}

	if *outputFile == "" {
		*outputFile = os.Getenv("JUNIT_FILE")
	}
	if *outputFile == "" {
		log.Fatal("output file is required (use --output or JUNIT_FILE env var)")
	}

	var ts *junit.TestSuites

	if _, err := os.Stat(*outputFile); err == nil {
		ts, err = junit.Load(*outputFile)
		if err != nil {
			log.Fatalf("failed to load existing file: %v", err)
		}
	} else {
		ts = junit.NewTestSuites()
	}

	suite := ts.GetOrCreateSuite(*suiteName)

	switch *status {
	case "ok":
		suite.AddMessageOK(*file, *message, *description)
	case "failed":
		suite.AddMessageFailed(*file, *message, *description)
	default:
		log.Fatalf("unknown status: %s", *status)
	}

	if err := ts.Write(*outputFile); err != nil {
		log.Fatalf("failed to write file: %v", err)
	}
}
