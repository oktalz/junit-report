package junit

import (
	"encoding/xml"
	"fmt"
	"os"
	"time"
)

// TestSuites is the root element of the JUnit XML report.
type TestSuites struct {
	XMLName    xml.Name     `xml:"testsuites"`
	TestSuites []*TestSuite `xml:"testsuite"`
}

// TestSuite represents a test suite.
type TestSuite struct {
	XMLName   xml.Name    `xml:"testsuite"`
	Name      string      `xml:"name,attr"`
	Tests     int         `xml:"tests,attr"`
	Failures  int         `xml:"failures,attr"`
	Errors    int         `xml:"errors,attr"`
	Time      string      `xml:"time,attr,omitempty"`
	Timestamp string      `xml:"timestamp,attr,omitempty"`
	TestCases []*TestCase `xml:"testcase"`
}

// TestCase represents a single test case.
type TestCase struct {
	XMLName   xml.Name   `xml:"testcase"`
	Name      string     `xml:"name,attr"`
	ClassName string     `xml:"classname,attr"`
	File      string     `xml:"file,attr,omitempty"`
	Time      string     `xml:"time,attr,omitempty"`
	Failure   *Failure   `xml:"failure,omitempty"`
	Error     *Error     `xml:"error,omitempty"`
	Skipped   *Skipped   `xml:"skipped,omitempty"`
	SystemOut *SystemOut `xml:"system-out,omitempty"`
	SystemErr *SystemErr `xml:"system-err,omitempty"`
}

// Failure represents a failed test case.
type Failure struct {
	Message string `xml:"message,attr,omitempty"`
	Type    string `xml:"type,attr,omitempty"`
	Content string `xml:",chardata"`
}

// Error represents an errored test case.
type Error struct {
	Message string `xml:"message,attr,omitempty"`
	Type    string `xml:"type,attr,omitempty"`
	Content string `xml:",chardata"`
}

// Skipped represents a skipped test case.
type Skipped struct {
	Message string `xml:"message,attr,omitempty"`
}

// SystemOut represents standard output.
type SystemOut struct {
	Content string `xml:",chardata"`
}

// SystemErr represents standard error.
type SystemErr struct {
	Content string `xml:",chardata"`
}

// NewTestSuites creates a new TestSuites object.
func NewTestSuites() *TestSuites {
	return &TestSuites{
		TestSuites: []*TestSuite{},
	}
}

// AddSuite adds a new test suite to the report.
func (ts *TestSuites) AddSuite(name string) *TestSuite {
	suite := &TestSuite{
		Name:      name,
		TestCases: []*TestCase{},
		Timestamp: time.Now().Format(time.RFC3339),
	}
	ts.TestSuites = append(ts.TestSuites, suite)
	return suite
}

// AddTestCase adds a new test case to the suite.
func (ts *TestSuite) AddTestCase(name, classname string) *TestCase {
	tc := &TestCase{
		Name:      name,
		ClassName: classname,
	}
	ts.TestCases = append(ts.TestCases, tc)
	ts.Tests++
	return tc
}

// SetTime sets the execution time for the test case.
func (tc *TestCase) SetTime(duration string) {
	tc.Time = duration
}

// SetFile sets the file path for the test case.
func (tc *TestCase) SetFile(file string) {
	tc.File = file
}

// Failed marks the test case as failed.
func (tc *TestCase) Failed(message, failureType, content string) {
	tc.Failure = &Failure{
		Message: message,
		Type:    failureType,
		Content: content,
	}
}

// Errored marks the test case as errored.
func (tc *TestCase) Errored(message, errorType, content string) {
	tc.Error = &Error{
		Message: message,
		Type:    errorType,
		Content: content,
	}
}

// Skipped marks the test case as skipped.
func (tc *TestCase) Skip(message string) {
	tc.Skipped = &Skipped{
		Message: message,
	}
}

// SetSystemOut sets the system-out content.
func (tc *TestCase) SetSystemOut(content string) {
	tc.SystemOut = &SystemOut{
		Content: content,
	}
}

// SetSystemErr sets the system-err content.
func (tc *TestCase) SetSystemErr(content string) {
	tc.SystemErr = &SystemErr{
		Content: content,
	}
}

// AddMessageOK adds a passing test case with a message and description.
// findTestCase checks if a test case with the given file and name already exists.
func (ts *TestSuite) findTestCase(file, name string) *TestCase {
	for _, tc := range ts.TestCases {
		if tc.File == file && tc.Name == name {
			return tc
		}
	}
	return nil
}

// AddMessageOK adds a passing test case with a message and description.
// If a test case with the same file and message already exists, it does nothing.
func (ts *TestSuite) AddMessageOK(file, message, description string) {
	if existing := ts.findTestCase(file, message); existing != nil {
		// Update description if needed
		if existing.SystemOut != nil && existing.SystemOut.Content != description {
			existing.SystemOut.Content = description
		}
		return
	}
	tc := ts.AddTestCase(message, "Test")
	tc.SetFile(file)
	tc.SetSystemOut(description)
}

// AddMessageFailed adds a failing test case with a message and description.
// If a test case with the same file and message already exists, it updates it to failed.
func (ts *TestSuite) AddMessageFailed(file, message, description string) {
	if existing := ts.findTestCase(file, message); existing != nil {
		// Ensure it is marked as failed
		existing.Failed(description, "Failure", description)
		return
	}
	tc := ts.AddTestCase(message, "Test")
	tc.SetFile(file)
	tc.Failed(description, "Failure", description)
}

// AddMessageError adds an errored test case with a message and description.
// If a test case with the same file and message exists, it updates it to errored.
func (ts *TestSuite) AddMessageError(file, message, description string) {
	if existing := ts.findTestCase(file, message); existing != nil {
		// Ensure it is marked as errored
		existing.Errored(description, "Error", description)
		return
	}
	tc := ts.AddTestCase(message, "Test")
	tc.SetFile(file)
	tc.Errored(description, "Error", description)
}

// Marshal returns the XML encoding of the TestSuites.
func (ts *TestSuites) Marshal() ([]byte, error) {
	return xml.MarshalIndent(ts, "", "  ")
}

// Write writes the XML report to a file.
// If fileName is empty, it checks the JUNIT_FILE environment variable.
func (ts *TestSuites) Write(fileName string) error {
	if fileName == "" {
		fileName = os.Getenv("JUNIT_FILE")
	}
	if fileName == "" {
		return fmt.Errorf("no filename provided and JUNIT_FILE env var is empty")
	}

	ts.UpdateCounts()
	data, err := ts.Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal xml: %w", err)
	}

	// Add XML header
	data = append([]byte(xml.Header), data...)

	return os.WriteFile(fileName, data, 0o644)
}

// Load loads the XML report from a file.
func Load(fileName string) (*TestSuites, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var ts TestSuites
	if err := xml.Unmarshal(data, &ts); err != nil {
		return nil, err
	}
	return &ts, nil
}

// GetOrCreateSuite returns an existing suite by name or creates a new one.
func (ts *TestSuites) GetOrCreateSuite(name string) *TestSuite {
	for _, suite := range ts.TestSuites {
		if suite.Name == name {
			return suite
		}
	}
	return ts.AddSuite(name)
}

// UpdateCounts updates the failure and error counts for the suite.
// This should be called before marshaling if counts are not manually managed.
func (ts *TestSuite) UpdateCounts() {
	ts.Failures = 0
	ts.Errors = 0
	for _, tc := range ts.TestCases {
		if tc.Failure != nil {
			ts.Failures++
		}
		if tc.Error != nil {
			ts.Errors++
		}
	}
}

// UpdateCounts updates the counts for all suites.
func (ts *TestSuites) UpdateCounts() {
	for _, suite := range ts.TestSuites {
		suite.UpdateCounts()
	}
}

func (ts *TestSuites) String() string {
	b, err := ts.Marshal()
	if err != nil {
		return fmt.Sprintf("error marshaling: %s", err)
	}
	return string(b)
}
