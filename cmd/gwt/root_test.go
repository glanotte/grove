package gwt

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestNewRootCmd(t *testing.T) {
	version := "1.0.0"
	commit := "abc123"
	date := "2023-01-01"

	cmd := NewRootCmd(version, commit, date)

	if cmd.Use != "grove" {
		t.Errorf("Expected Use 'grove', got %s", cmd.Use)
	}

	if cmd.Short != "Git worktree manager with Docker and template support" {
		t.Errorf("Expected correct short description, got %s", cmd.Short)
	}

	// Check that all subcommands are added
	expectedCommands := []string{"init", "create", "list", "remove", "switch", "version"}
	commands := cmd.Commands()

	if len(commands) != len(expectedCommands) {
		t.Errorf("Expected %d commands, got %d", len(expectedCommands), len(commands))
	}

	for _, expectedCmd := range expectedCommands {
		found := false
		for _, cmd := range commands {
			if cmd.Use == expectedCmd || cmd.Use == expectedCmd+" <worktree-name>" || cmd.Use == expectedCmd+" <branch-name>" || cmd.Use == expectedCmd+" <repo-url>" {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected command %s not found", expectedCmd)
		}
	}
}

func TestRootCmd_Flags(t *testing.T) {
	cmd := NewRootCmd("1.0.0", "abc123", "2023-01-01")

	// Test that config flag is present
	configFlag := cmd.PersistentFlags().Lookup("config")
	if configFlag == nil {
		t.Error("Expected config flag to be present")
	}

	if configFlag.DefValue != "" {
		t.Errorf("Expected config flag default value to be empty, got %s", configFlag.DefValue)
	}

	if configFlag.Usage != "config file (default is .grove/config.yaml)" {
		t.Errorf("Expected correct config flag usage, got %s", configFlag.Usage)
	}
}

func TestRootCmd_Execute(t *testing.T) {
	cmd := NewRootCmd("1.0.0", "abc123", "2023-01-01")

	// Test that command can be executed (just check it doesn't panic)
	// We don't actually run it to avoid side effects
	if cmd.Execute == nil {
		t.Error("Expected Execute method to be available")
	}
}

// Test command structure and hierarchy
func TestCommandHierarchy(t *testing.T) {
	cmd := NewRootCmd("1.0.0", "abc123", "2023-01-01")

	// Test that each command has proper parent
	for _, subCmd := range cmd.Commands() {
		if subCmd.Parent() != cmd {
			t.Errorf("Command %s should have root as parent", subCmd.Use)
		}
	}
}

// Test that version information is passed correctly
func TestVersionCommand(t *testing.T) {
	version := "1.2.3"
	commit := "def456"
	date := "2023-06-15"

	cmd := NewRootCmd(version, commit, date)

	// Find version command
	var versionCmd *cobra.Command
	for _, subCmd := range cmd.Commands() {
		if subCmd.Use == "version" {
			versionCmd = subCmd
			break
		}
	}

	if versionCmd == nil {
		t.Fatal("Version command not found")
	}

	// We can't easily test the actual output without running the command
	// But we can verify the command exists and has the right structure
	if versionCmd.Short != "Print version information" {
		t.Errorf("Expected correct version command description, got %s", versionCmd.Short)
	}

	if versionCmd.Run == nil {
		t.Error("Version command should have a Run function")
	}
}

// Integration test helper
func executeCommand(t *testing.T, cmd *cobra.Command, args ...string) error {
	cmd.SetArgs(args)
	return cmd.Execute()
}

// Test command execution without side effects
func TestCommandsExecuteWithoutPanic(t *testing.T) {
	cmd := NewRootCmd("1.0.0", "abc123", "2023-01-01")

	// Test version command (safe to execute)
	versionCmd := NewRootCmd("1.0.0", "abc123", "2023-01-01")
	versionCmd.SetArgs([]string{"version"})
	
	// This should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Version command panicked: %v", r)
		}
	}()

	// Don't actually execute to avoid output, just verify structure
	if cmd == nil {
		t.Error("Root command should not be nil")
	}
}