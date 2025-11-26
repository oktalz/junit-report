package junit

import (
	"encoding/xml"
	"os"
	"testing"
)

func TestGitLabExample(t *testing.T) {
	// Recreate the example from GitLab docs
	// <testsuites>
	//   <testsuite name="Authentication Tests" tests="1" failures="1">
	//     <testcase classname="LoginTest" name="test_invalid_password" file="spec/auth_spec.rb" time="0.23">
	//       <failure>Expected authentication to fail</failure>
	//       <system-out>[[ATTACHMENT|screenshots/failure.png]]</system-out>
	//     </testcase>
	//   </testsuite>
	// </testsuites>

	ts := NewTestSuites()
	suite := ts.AddSuite("Authentication Tests")
	tc := suite.AddTestCase("test_invalid_password", "LoginTest")
	tc.SetFile("spec/auth_spec.rb")
	tc.SetTime("0.23")
	tc.Failed("", "", "Expected authentication to fail")
	tc.SetSystemOut("[[ATTACHMENT|screenshots/failure.png]]")

	// Update counts
	ts.UpdateCounts()

	output, err := ts.Marshal()
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// We need to be careful with comparison because attributes order might vary or omitempty might affect it.
	// Also timestamp is dynamic.
	// So let's unmarshal both back to structs and compare relevant fields, or just check if it contains expected strings.

	// For simplicity in this specific check, let's verify key elements are present.

	// Unmarshal the generated output
	var parsed TestSuites
	if err := xml.Unmarshal(output, &parsed); err != nil {
		t.Fatalf("Failed to unmarshal generated XML: %v", err)
	}

	if len(parsed.TestSuites) != 1 {
		t.Fatalf("Expected 1 test suite, got %d", len(parsed.TestSuites))
	}

	pSuite := parsed.TestSuites[0]
	if pSuite.Name != "Authentication Tests" {
		t.Errorf("Expected suite name 'Authentication Tests', got '%s'", pSuite.Name)
	}
	if pSuite.Tests != 1 {
		t.Errorf("Expected 1 test, got %d", pSuite.Tests)
	}
	if pSuite.Failures != 1 {
		t.Errorf("Expected 1 failure, got %d", pSuite.Failures)
	}

	if len(pSuite.TestCases) != 1 {
		t.Fatalf("Expected 1 test case, got %d", len(pSuite.TestCases))
	}

	pTc := pSuite.TestCases[0]
	if pTc.Name != "test_invalid_password" {
		t.Errorf("Expected test case name 'test_invalid_password', got '%s'", pTc.Name)
	}
	if pTc.ClassName != "LoginTest" {
		t.Errorf("Expected classname 'LoginTest', got '%s'", pTc.ClassName)
	}
	if pTc.File != "spec/auth_spec.rb" {
		t.Errorf("Expected file 'spec/auth_spec.rb', got '%s'", pTc.File)
	}
	if pTc.Time != "0.23" {
		t.Errorf("Expected time '0.23', got '%s'", pTc.Time)
	}
	if pTc.Failure == nil {
		t.Fatal("Expected failure, got nil")
	}
	if pTc.Failure.Content != "Expected authentication to fail" {
		t.Errorf("Expected failure content 'Expected authentication to fail', got '%s'", pTc.Failure.Content)
	}
	if pTc.SystemOut == nil {
		t.Fatal("Expected system-out, got nil")
	}
	if pTc.SystemOut.Content != "[[ATTACHMENT|screenshots/failure.png]]" {
		t.Errorf("Expected system-out content '[[ATTACHMENT|screenshots/failure.png]]', got '%s'", pTc.SystemOut.Content)
	}
}

func TestHelperMethods(t *testing.T) {
	ts := NewTestSuites()
	suite := ts.AddSuite("Helper Tests")

	suite.AddMessageOK("file1.go", "Check 1", "All good")
	suite.AddMessageFailed("file2.go", "Check 2", "Something wrong")

	ts.UpdateCounts()

	if suite.Tests != 2 {
		t.Errorf("Expected 2 tests, got %d", suite.Tests)
	}
	if suite.Failures != 1 {
		t.Errorf("Expected 1 failure, got %d", suite.Failures)
	}

	tc1 := suite.TestCases[0]
	if tc1.Name != "Check 1" {
		t.Errorf("Expected name 'Check 1', got '%s'", tc1.Name)
	}
	if tc1.File != "file1.go" {
		t.Errorf("Expected file 'file1.go', got '%s'", tc1.File)
	}
	if tc1.SystemOut.Content != "All good" {
		t.Errorf("Expected system-out 'All good', got '%s'", tc1.SystemOut.Content)
	}
	if tc1.Failure != nil {
		t.Error("Expected no failure for tc1")
	}

	tc2 := suite.TestCases[1]
	if tc2.Name != "Check 2" {
		t.Errorf("Expected name 'Check 2', got '%s'", tc2.Name)
	}
	if tc2.File != "file2.go" {
		t.Errorf("Expected file 'file2.go', got '%s'", tc2.File)
	}
	if tc2.Failure == nil {
		t.Error("Expected failure for tc2")
	}
	if tc2.Failure.Message != "Something wrong" {
		t.Errorf("Expected failure message 'Something wrong', got '%s'", tc2.Failure.Message)
	}
}

func TestWrite(t *testing.T) {
	ts := NewTestSuites()
	suite := ts.AddSuite("Write Test")
	suite.AddMessageOK("file.go", "Check", "OK")

	// Test with explicit filename
	tmpFile := "test_output.xml"
	defer os.Remove(tmpFile)

	if err := ts.Write(tmpFile); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if len(content) == 0 {
		t.Error("File content is empty")
	}

	// Test with ENV variable
	envFile := "env_output.xml"
	defer os.Remove(envFile)
	os.Setenv("JUNIT_FILE", envFile)
	defer os.Unsetenv("JUNIT_FILE")

	if err := ts.Write(""); err != nil {
		t.Fatalf("Failed to write file using env var: %v", err)
	}

	contentEnv, err := os.ReadFile(envFile)
	if err != nil {
		t.Fatalf("Failed to read env file: %v", err)
	}

	if len(contentEnv) == 0 {
		t.Error("Env file content is empty")
	}

	// Test error case
	os.Unsetenv("JUNIT_FILE")
	if err := ts.Write(""); err == nil {
		t.Error("Expected error when no filename provided, got nil")
	}
}

func TestLoadAndGetOrCreateSuite(t *testing.T) {
	tmpFile := "load_test.xml"
	defer os.Remove(tmpFile)

	// Create initial file
	ts := NewTestSuites()
	s1 := ts.AddSuite("Suite 1")
	s1.AddTestCase("Test 1", "Class 1")
	if err := ts.Write(tmpFile); err != nil {
		t.Fatalf("Failed to write initial file: %v", err)
	}

	// Load file
	loadedTs, err := Load(tmpFile)
	if err != nil {
		t.Fatalf("Failed to load file: %v", err)
	}

	if len(loadedTs.TestSuites) != 1 {
		t.Fatalf("Expected 1 suite, got %d", len(loadedTs.TestSuites))
	}
	if loadedTs.TestSuites[0].Name != "Suite 1" {
		t.Errorf("Expected suite name 'Suite 1', got '%s'", loadedTs.TestSuites[0].Name)
	}

	// Get existing suite
	s1Existing := loadedTs.GetOrCreateSuite("Suite 1")
	if s1Existing != loadedTs.TestSuites[0] {
		t.Error("Expected to get existing suite, got different object")
	}

	// Create new suite
	s2 := loadedTs.GetOrCreateSuite("Suite 2")
	if s2.Name != "Suite 2" {
		t.Errorf("Expected suite name 'Suite 2', got '%s'", s2.Name)
	}
	if len(loadedTs.TestSuites) != 2 {
		t.Errorf("Expected 2 suites, got %d", len(loadedTs.TestSuites))
	}
}
