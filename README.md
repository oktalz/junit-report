# JUnit Report XML Generator for Go

This repository provides a Go package and a CLI application to generate JUnit XML reports compatible with GitLab CI (and other CI systems).

## Go Package

The `junit` package allows you to programmatically create JUnit XML reports.

### Installation

```bash
go get github.com/oktalz/junit-report
```

### Usage

```go
package main

import (
	"log"

	"github.com/oktalz/junit-report"
)

func main() {
	// Create a new report
	ts := junit.NewTestSuites()

	// Add a test suite
	suite := ts.AddSuite("My Suite")

	// Add a passing test case
	suite.AddMessageOK("main.go", "Build Check", "Build successful")

	// Add an error test case
	suite.AddMessageError("main.go", "Build Check", "Build failed")

	// Add a failing test case
	suite.AddMessageFailed("linter.go", "Lint Check", "Linting failed: variable unused")

	// Write to file
	if err := ts.Write("report.xml"); err != nil {
		log.Fatal(err)
	}
}
```

You can also load an existing report and append to it:

```go
ts, err := junit.Load("report.xml")
if err != nil {
    // Handle error (e.g., create new if not found)
    ts = junit.NewTestSuites()
}

suite := ts.GetOrCreateSuite("My Suite")
suite.AddMessageOK("test.go", "New Test", "Passed")

ts.Write("report.xml")
```

## CLI Application

The `junit-report` CLI allows you to add test results to a JUnit XML file from the command line. This is useful for scripts and CI pipelines.

### Installation

```bash
go install github.com/oktalz/junit-report/cmd/junit-report@latest
```

### Usage

The `add` subcommand is used to add a test case result.

#### Arguments

- `--output`: Path to the output XML file. Can also be set via `JUNIT_FILE` environment variable.
- `--status`: Status of the check (`ok`, `error` or `failed`). Default: `ok`.
- `--file`: File path related to the check.
- `--message`: Name/Message for the check.
- `--description`: Detailed description or output.
- `--suite`: Name of the test suite. Default: `General`.

#### Examples

**Add a passing check:**

```bash
junit-report add --output=report.xml --status=ok --file=main.go --message="Build" --description="Build successful"
```

**Add an error check:**

```bash
junit-report add --output=report.xml --status=error --file=main.go --message="Build" --description="Build failed"
```	

**Add a failing check:**

```bash
junit-report add --output=report.xml --status=failed --file=linter.go --message="Linter" --description="Found 5 lint errors"
```

**Using Environment Variable:**

```bash
export JUNIT_FILE=report.xml
junit-report add --status=ok --message="Env Test"
```

**Appending to existing file:**

If the output file exists, the tool will load it and append the new test case to the specified suite (creating the suite if it doesn't exist).
