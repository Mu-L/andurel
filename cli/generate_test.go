package cli

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestGenerateCommands(t *testing.T) {
	rootCmd := NewRootCommand("test", "test-date")

	tests := []struct {
		name string
		args []string
	}{
		{"generate help", []string{"generate", "--help"}},
		{"model help", []string{"generate", "model", "--help"}},
		{"controller help", []string{"generate", "controller", "--help"}},
		{"resource help", []string{"generate", "resource", "--help"}},
		{"queries help", []string{"generate", "queries", "--help"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rootCmd.SetArgs(tt.args)
			err := rootCmd.Execute()
			if err != nil {
				t.Errorf("Command %v failed: %v", tt.args, err)
			}
		})
	}
}

func TestGenerateCommandStructure(t *testing.T) {
	rootCmd := NewRootCommand("test", "test-date")

	var generateCmd *cobra.Command
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "generate" {
			generateCmd = cmd
			break
		}
	}

	if generateCmd == nil {
		t.Fatal("generate command not found")
	}

	expectedCommands := []string{"model", "controller", "view", "resource", "queries"}
	foundCommands := make(map[string]bool)

	for _, cmd := range generateCmd.Commands() {
		cmdName := strings.Fields(cmd.Use)[0]
		foundCommands[cmdName] = true
	}

	for _, expectedCmd := range expectedCommands {
		if !foundCommands[expectedCmd] {
			t.Errorf(
				"Expected command '%s' not found. Available commands: %v",
				expectedCmd,
				getCommandNames(generateCmd.Commands()),
			)
		}
	}
}

func getCommandNames(commands []*cobra.Command) []string {
	var names []string
	for _, cmd := range commands {
		cmdName := strings.Fields(cmd.Use)[0]
		names = append(names, cmdName)
	}
	return names
}

func TestQueriesCommandFlags(t *testing.T) {
	cmd := newQueriesCommand()

	if cmd.Use != "queries [name]" {
		t.Errorf("Expected Use 'queries [name]', got '%s'", cmd.Use)
	}

	refreshFlag := cmd.Flags().Lookup("refresh")
	if refreshFlag == nil {
		t.Error("Expected --refresh flag to exist")
	}
	if refreshFlag.DefValue != "false" {
		t.Errorf("Expected --refresh default to be 'false', got '%s'", refreshFlag.DefValue)
	}

	tableNameFlag := cmd.Flags().Lookup("table-name")
	if tableNameFlag == nil {
		t.Error("Expected --table-name flag to exist")
	}
	if tableNameFlag.DefValue != "" {
		t.Errorf("Expected --table-name default to be empty, got '%s'", tableNameFlag.DefValue)
	}
}

func TestGenerateModelFlags(t *testing.T) {
	cmd := newModelCommand()

	if cmd.Use != "model [name]" {
		t.Errorf("Expected Use 'model [name]', got '%s'", cmd.Use)
	}

	refreshFlag := cmd.Flags().Lookup("refresh")
	if refreshFlag == nil {
		t.Error("Expected --refresh flag to exist")
	}
	if refreshFlag.DefValue != "false" {
		t.Errorf("Expected --refresh default to be 'false', got '%s'", refreshFlag.DefValue)
	}

	tableNameFlag := cmd.Flags().Lookup("table-name")
	if tableNameFlag == nil {
		t.Error("Expected --table-name flag to exist")
	}
	if tableNameFlag.DefValue != "" {
		t.Errorf("Expected --table-name default to be empty, got '%s'", tableNameFlag.DefValue)
	}
}

func TestGenerateControllerFlags(t *testing.T) {
	cmd := newControllerCommand()

	if cmd.Use != "controller [model_name]" {
		t.Errorf("Expected Use 'controller [model_name]', got '%s'", cmd.Use)
	}

	withViewsFlag := cmd.Flags().Lookup("with-views")
	if withViewsFlag == nil {
		t.Error("Expected --with-views flag to exist")
	}
	if withViewsFlag.DefValue != "false" {
		t.Errorf("Expected --with-views default to be 'false', got '%s'", withViewsFlag.DefValue)
	}
}

func TestGenerateResourceFlags(t *testing.T) {
	cmd := newResourceCommand()

	if cmd.Use != "resource [name]" {
		t.Errorf("Expected Use 'resource [name]', got '%s'", cmd.Use)
	}

	tableNameFlag := cmd.Flags().Lookup("table-name")
	if tableNameFlag == nil {
		t.Error("Expected --table-name flag to exist")
	}
	if tableNameFlag.DefValue != "" {
		t.Errorf("Expected --table-name default to be empty, got '%s'", tableNameFlag.DefValue)
	}
}

func TestGenerateViewFlags(t *testing.T) {
	cmd := newViewCommand()

	if cmd.Use != "view [model_name]" {
		t.Errorf("Expected Use 'view [model_name]', got '%s'", cmd.Use)
	}

	withControllerFlag := cmd.Flags().Lookup("with-controller")
	if withControllerFlag == nil {
		t.Error("Expected --with-controller flag to exist")
	}
	if withControllerFlag.DefValue != "false" {
		t.Errorf("Expected --with-controller default to be 'false', got '%s'", withControllerFlag.DefValue)
	}
}

func TestGenerateModelFunction(t *testing.T) {
	tests := []struct {
		name              string
		args              []string
		refresh           bool
		tableNameOverride string
		expectError       bool
	}{
		{
			name:              "valid model without flags",
			args:              []string{"User"},
			refresh:           false,
			tableNameOverride: "",
			expectError:       true,
		},
		{
			name:              "valid model with table name",
			args:              []string{"User"},
			refresh:           false,
			tableNameOverride: "accounts",
			expectError:       true,
		},
		{
			name:              "valid model with refresh",
			args:              []string{"User"},
			refresh:           true,
			tableNameOverride: "",
			expectError:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			cmd.Flags().Bool("refresh", tt.refresh, "")
			cmd.Flags().String("table-name", tt.tableNameOverride, "")

			err := generateModel(cmd, tt.args)
			if tt.expectError && err == nil {
				t.Error("Expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestGenerateControllerFunction(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		withViews   bool
		expectError bool
	}{
		{
			name:        "valid controller without views",
			args:        []string{"User"},
			withViews:   false,
			expectError: true,
		},
		{
			name:        "valid controller with views",
			args:        []string{"User"},
			withViews:   true,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			cmd.Flags().Bool("with-views", tt.withViews, "")

			err := generateController(cmd, tt.args)
			if tt.expectError && err == nil {
				t.Error("Expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestGenerateViewFunction(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		withController bool
		expectError    bool
	}{
		{
			name:           "valid view without controller",
			args:           []string{"User"},
			withController: false,
			expectError:    true,
		},
		{
			name:           "valid view with controller",
			args:           []string{"User"},
			withController: true,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			cmd.Flags().Bool("with-controller", tt.withController, "")

			err := generateView(cmd, tt.args)
			if tt.expectError && err == nil {
				t.Error("Expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestGenerateResourceFunction(t *testing.T) {
	tests := []struct {
		name              string
		args              []string
		tableNameOverride string
		expectError       bool
	}{
		{
			name:              "valid resource without table name",
			args:              []string{"Product"},
			tableNameOverride: "",
			expectError:       true,
		},
		{
			name:              "valid resource with table name",
			args:              []string{"Feedback"},
			tableNameOverride: "user_feedback",
			expectError:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			cmd.Flags().String("table-name", tt.tableNameOverride, "")

			err := generateResource(cmd, tt.args)
			if tt.expectError && err == nil {
				t.Error("Expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestGenerateQueriesFunction(t *testing.T) {
	tests := []struct {
		name              string
		args              []string
		refresh           bool
		tableNameOverride string
		expectError       bool
	}{
		{
			name:              "valid queries without flags",
			args:              []string{"UserRole"},
			refresh:           false,
			tableNameOverride: "",
			expectError:       true,
		},
		{
			name:              "valid queries with table name",
			args:              []string{"UserRole"},
			refresh:           false,
			tableNameOverride: "users_roles",
			expectError:       true,
		},
		{
			name:              "valid queries with refresh",
			args:              []string{"UserRole"},
			refresh:           true,
			tableNameOverride: "",
			expectError:       true,
		},
		{
			name:              "valid queries with refresh and table name",
			args:              []string{"UserRole"},
			refresh:           true,
			tableNameOverride: "users_roles",
			expectError:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			cmd.Flags().Bool("refresh", tt.refresh, "")
			cmd.Flags().String("table-name", tt.tableNameOverride, "")

			err := generateQueries(cmd, tt.args)
			if tt.expectError && err == nil {
				t.Error("Expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
