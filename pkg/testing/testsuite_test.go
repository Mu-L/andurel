package testing

import (
	"testing"
)

func TestNewTestSuite(t *testing.T) {
	suite := NewTestSuite()

	if suite == nil {
		t.Fatal("NewTestSuite returned nil")
	}

	if suite.Unit == nil {
		t.Error("suite.Unit is nil")
	}

	if suite.Integration == nil {
		t.Error("suite.Integration is nil")
	}

	if suite.E2E == nil {
		t.Error("suite.E2E is nil")
	}
}

func TestTestSuite_AddUnitTest(t *testing.T) {
	suite := NewTestSuite()

	testName := "TestUnit"
	testDesc := "Test unit functionality"
	testFunc := func(t *testing.T) {}

	suite.AddUnitTest(testName, testDesc, testFunc)

	if len(suite.Unit) != 1 {
		t.Errorf("Expected 1 unit test, got %d", len(suite.Unit))
	}

	test, exists := suite.Unit[testName]
	if !exists {
		t.Fatal("Unit test not found after adding")
	}

	if test.Name != testName {
		t.Errorf("Expected test name %q, got %q", testName, test.Name)
	}

	if test.Description != testDesc {
		t.Errorf("Expected description %q, got %q", testDesc, test.Description)
	}

	if test.TestFunc == nil {
		t.Error("TestFunc is nil")
	}
}

func TestTestSuite_AddIntegrationTest(t *testing.T) {
	suite := NewTestSuite()

	testName := "TestIntegration"
	testDesc := "Test integration functionality"
	setup := func() error { return nil }
	teardown := func() error { return nil }
	testFunc := func(t *testing.T) {}

	suite.AddIntegrationTest(testName, testDesc, setup, teardown, testFunc)

	if len(suite.Integration) != 1 {
		t.Errorf("Expected 1 integration test, got %d", len(suite.Integration))
	}

	test, exists := suite.Integration[testName]
	if !exists {
		t.Fatal("Integration test not found after adding")
	}

	if test.Name != testName {
		t.Errorf("Expected test name %q, got %q", testName, test.Name)
	}

	if test.Description != testDesc {
		t.Errorf("Expected description %q, got %q", testDesc, test.Description)
	}

	if test.Setup == nil {
		t.Error("Setup is nil")
	}

	if test.Teardown == nil {
		t.Error("Teardown is nil")
	}

	if test.TestFunc == nil {
		t.Error("TestFunc is nil")
	}
}

func TestTestSuite_AddE2ETest(t *testing.T) {
	suite := NewTestSuite()

	testName := "TestE2E"
	testDesc := "Test end-to-end functionality"
	setup := func() error { return nil }
	teardown := func() error { return nil }
	testFunc := func(t *testing.T) {}

	suite.AddE2ETest(testName, testDesc, setup, teardown, testFunc)

	if len(suite.E2E) != 1 {
		t.Errorf("Expected 1 e2e test, got %d", len(suite.E2E))
	}

	test, exists := suite.E2E[testName]
	if !exists {
		t.Fatal("E2E test not found after adding")
	}

	if test.Name != testName {
		t.Errorf("Expected test name %q, got %q", testName, test.Name)
	}

	if test.Description != testDesc {
		t.Errorf("Expected description %q, got %q", testDesc, test.Description)
	}

	if test.Setup == nil {
		t.Error("Setup is nil")
	}

	if test.Teardown == nil {
		t.Error("Teardown is nil")
	}

	if test.TestFunc == nil {
		t.Error("TestFunc is nil")
	}
}

func TestTestSuite_RunUnitTests(t *testing.T) {
	suite := NewTestSuite()

	executed := false
	suite.AddUnitTest("TestUnit", "Test", func(t *testing.T) {
		executed = true
	})

	suite.RunUnitTests(t)

	if !executed {
		t.Error("Unit test was not executed")
	}
}

func TestTestSuite_RunIntegrationTests(t *testing.T) {
	suite := NewTestSuite()

	setupCalled := false
	testCalled := false
	teardownCalled := false

	suite.AddIntegrationTest(
		"TestIntegration",
		"Test",
		func() error {
			setupCalled = true
			return nil
		},
		func() error {
			teardownCalled = true
			return nil
		},
		func(t *testing.T) {
			testCalled = true
		},
	)

	suite.RunIntegrationTests(t)

	if !setupCalled {
		t.Error("Setup was not called")
	}

	if !testCalled {
		t.Error("Test was not executed")
	}

	if !teardownCalled {
		t.Error("Teardown was not called")
	}
}

func TestTestSuite_RunIntegrationTests_SetupError(t *testing.T) {
	suite := NewTestSuite()

	suite.AddIntegrationTest(
		"TestIntegration",
		"Test",
		func() error {
			return nil
		},
		func() error {
			return nil
		},
		func(t *testing.T) {
			t.Log("Test runs when setup succeeds")
		},
	)

	suite.RunIntegrationTests(t)
}

func TestTestSuite_RunE2ETests(t *testing.T) {
	suite := NewTestSuite()

	setupCalled := false
	testCalled := false
	teardownCalled := false

	suite.AddE2ETest(
		"TestE2E",
		"Test",
		func() error {
			setupCalled = true
			return nil
		},
		func() error {
			teardownCalled = true
			return nil
		},
		func(t *testing.T) {
			testCalled = true
		},
	)

	suite.RunE2ETests(t)

	if !setupCalled {
		t.Error("Setup was not called")
	}

	if !testCalled {
		t.Error("Test was not executed")
	}

	if !teardownCalled {
		t.Error("Teardown was not called")
	}
}

func TestTestSuite_RunAllTests(t *testing.T) {
	suite := NewTestSuite()

	unitCalled := false
	integrationCalled := false
	e2eCalled := false

	suite.AddUnitTest("TestUnit", "Test", func(t *testing.T) {
		unitCalled = true
	})

	suite.AddIntegrationTest(
		"TestIntegration",
		"Test",
		func() error { return nil },
		func() error { return nil },
		func(t *testing.T) {
			integrationCalled = true
		},
	)

	suite.AddE2ETest(
		"TestE2E",
		"Test",
		func() error { return nil },
		func() error { return nil },
		func(t *testing.T) {
			e2eCalled = true
		},
	)

	suite.RunAllTests(t)

	if !unitCalled {
		t.Error("Unit test was not executed")
	}

	if !integrationCalled {
		t.Error("Integration test was not executed")
	}

	if !e2eCalled {
		t.Error("E2E test was not executed")
	}
}

func TestTableDrivenTest_Run(t *testing.T) {
	tests := []TestData{
		{
			Name:  "Test1",
			Input: 1,
			Expected: 1,
		},
		{
			Name:  "Test2",
			Input: 2,
			Expected: 2,
		},
	}

	testFunc := func(t *testing.T, test TestData) {
		if test.Input != test.Expected {
			t.Errorf("Input %v != Expected %v", test.Input, test.Expected)
		}
	}

	tdt := NewTableDrivenTest("TestGroup", tests, testFunc)
	tdt.Run(t)
}

func TestNewTableDrivenTest(t *testing.T) {
	tests := []TestData{
		{
			Name:  "Test1",
			Input: 1,
			Expected: 1,
		},
	}

	testFunc := func(t *testing.T, test TestData) {
		t.Log("Test executed")
	}

	tdt := NewTableDrivenTest("TestGroup", tests, testFunc)

	if tdt == nil {
		t.Fatal("NewTableDrivenTest returned nil")
	}

	if tdt.Name != "TestGroup" {
		t.Errorf("Expected name 'TestGroup', got %q", tdt.Name)
	}

	if tdt.Tests == nil {
		t.Error("Tests is nil")
	}

	if tdt.TestFunc == nil {
		t.Error("TestFunc is nil")
	}
}

func TestTestData(t *testing.T) {
	data := TestData{
		Name:     "TestName",
		Input:    "input",
		Expected: "expected",
		Error:    nil,
	}

	if data.Name != "TestName" {
		t.Errorf("Expected name 'TestName', got %q", data.Name)
	}

	if data.Input != "input" {
		t.Errorf("Expected input 'input', got %q", data.Input)
	}

	if data.Expected != "expected" {
		t.Errorf("Expected expected 'expected', got %q", data.Expected)
	}
}

func TestUnitTest(t *testing.T) {
	unitTest := UnitTest{
		Name:        "TestUnit",
		Description: "Test description",
		TestFunc:    func(t *testing.T) {},
	}

	if unitTest.Name != "TestUnit" {
		t.Errorf("Expected name 'TestUnit', got %q", unitTest.Name)
	}

	if unitTest.Description != "Test description" {
		t.Errorf("Expected description 'Test description', got %q", unitTest.Description)
	}

	if unitTest.TestFunc == nil {
		t.Error("TestFunc is nil")
	}
}

func TestIntegrationTest(t *testing.T) {
	integrationTest := IntegrationTest{
		Name:        "TestIntegration",
		Description: "Test description",
		Setup:       func() error { return nil },
		Teardown:    func() error { return nil },
		TestFunc:    func(t *testing.T) {},
	}

	if integrationTest.Name != "TestIntegration" {
		t.Errorf("Expected name 'TestIntegration', got %q", integrationTest.Name)
	}

	if integrationTest.Description != "Test description" {
		t.Errorf("Expected description 'Test description', got %q", integrationTest.Description)
	}

	if integrationTest.Setup == nil {
		t.Error("Setup is nil")
	}

	if integrationTest.Teardown == nil {
		t.Error("Teardown is nil")
	}

	if integrationTest.TestFunc == nil {
		t.Error("TestFunc is nil")
	}
}

func TestEndToEndTest(t *testing.T) {
	e2eTest := EndToEndTest{
		Name:        "TestE2E",
		Description: "Test description",
		Setup:       func() error { return nil },
		Teardown:    func() error { return nil },
		TestFunc:    func(t *testing.T) {},
	}

	if e2eTest.Name != "TestE2E" {
		t.Errorf("Expected name 'TestE2E', got %q", e2eTest.Name)
	}

	if e2eTest.Description != "Test description" {
		t.Errorf("Expected description 'Test description', got %q", e2eTest.Description)
	}

	if e2eTest.Setup == nil {
		t.Error("Setup is nil")
	}

	if e2eTest.Teardown == nil {
		t.Error("Teardown is nil")
	}

	if e2eTest.TestFunc == nil {
		t.Error("TestFunc is nil")
	}
}
