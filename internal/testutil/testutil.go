package testutil

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/glanotte/grove/pkg/worktree"
)

// TempDir creates a temporary directory for testing
func TempDir(t *testing.T) string {
	dir, err := ioutil.TempDir("", "grove-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	t.Cleanup(func() {
		os.RemoveAll(dir)
	})
	return dir
}

// CreateTestRepo creates a test Git repository
func CreateTestRepo(t *testing.T) string {
	dir := TempDir(t)
	
	// Initialize git repo
	err := RunCommand(dir, "git", "init", "--bare")
	if err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}
	
	return dir
}

// CreateTestConfig creates a test configuration
func CreateTestConfig() *worktree.Config {
	return &worktree.Config{
		Version: 1,
		Project: worktree.ProjectConfig{
			Name:   "testapp",
			Domain: "test.local",
		},
		Worktree: worktree.WorktreeConfig{
			BasePath:      "./worktrees",
			NamingPattern: "{branch}",
		},
		Docker: worktree.DockerConfig{
			Enabled:     true,
			ComposeFile: "docker-compose.yml",
			PortOffset:  10000,
			NetworkName: "{project_name}_network",
		},
		Web: worktree.WebConfig{
			Enabled:          true,
			ProxyType:        "nginx-proxy",
			SubdomainPattern: "{branch}.{project_domain}",
		},
		Templates: worktree.TemplateConfig{
			Default: "standard",
			Available: map[string]worktree.TemplateDefinition{
				"standard": {
					Files: []worktree.TemplateFile{
						{Src: "docker-compose.yml.tmpl", Dest: "docker-compose.yml"},
						{Src: ".env.tmpl", Dest: ".env"},
					},
				},
			},
		},
		Variables: map[string]interface{}{
			"db_name_prefix": "testapp",
			"redis_prefix":   "testapp",
		},
	}
}

// CreateTestWorktreeStructure creates a test worktree directory structure
func CreateTestWorktreeStructure(t *testing.T, baseDir string) {
	groveDir := filepath.Join(baseDir, ".grove")
	templatesDir := filepath.Join(groveDir, "templates")
	
	// Create directories
	err := os.MkdirAll(templatesDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create .grove/templates: %v", err)
	}
	
	// Create template files
	dockerTemplate := `version: '3.8'
services:
  app:
    container_name: {{.ProjectName}}_{{.BranchName}}_app
    image: nginx:alpine
    ports:
      - "{{.WebPort}}:80"
    networks:
      - {{.NetworkName}}
networks:
  {{.NetworkName}}:
    external: true
`
	
	envTemplate := `APP_NAME={{.ProjectName}}_{{.BranchName}}
APP_URL=https://{{.BranchName}}.{{.ProjectDomain}}
WEB_PORT={{.WebPort}}
`
	
	err = ioutil.WriteFile(filepath.Join(templatesDir, "docker-compose.yml.tmpl"), []byte(dockerTemplate), 0644)
	if err != nil {
		t.Fatalf("Failed to create docker-compose template: %v", err)
	}
	
	err = ioutil.WriteFile(filepath.Join(templatesDir, ".env.tmpl"), []byte(envTemplate), 0644)
	if err != nil {
		t.Fatalf("Failed to create .env template: %v", err)
	}
}

// CreateTestManager creates a test manager with proper setup
func CreateTestManager(t *testing.T) (*worktree.Manager, string) {
	baseDir := TempDir(t)
	CreateTestWorktreeStructure(t, baseDir)
	
	manager, err := worktree.NewManager(baseDir)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}
	
	return manager, baseDir
}

// AssertFileExists checks if a file exists
func AssertFileExists(t *testing.T, path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("Expected file %s to exist, but it doesn't", path)
	}
}

// AssertFileNotExists checks if a file doesn't exist
func AssertFileNotExists(t *testing.T, path string) {
	if _, err := os.Stat(path); err == nil {
		t.Errorf("Expected file %s to not exist, but it does", path)
	}
}

// AssertFileContains checks if a file contains specific content
func AssertFileContains(t *testing.T, path, expectedContent string) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", path, err)
	}
	
	if !contains(string(content), expectedContent) {
		t.Errorf("Expected file %s to contain %q, but it doesn't. Content:\n%s", path, expectedContent, string(content))
	}
}

// AssertDirExists checks if a directory exists
func AssertDirExists(t *testing.T, path string) {
	if stat, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("Expected directory %s to exist, but it doesn't", path)
	} else if err != nil {
		t.Errorf("Error checking directory %s: %v", path, err)
	} else if !stat.IsDir() {
		t.Errorf("Expected %s to be a directory, but it's not", path)
	}
}

// RunCommand runs a command in a specific directory
func RunCommand(dir, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	return cmd.Run()
}

// contains checks if a string contains a substring
func contains(str, substr string) bool {
	return len(str) >= len(substr) && (str == substr || contains(str[1:], substr))
}

// MockGitCommand creates a mock git command for testing
type MockGitCommand struct {
	Commands []string
	Outputs  []string
	Errors   []error
	CallCount int
}

func (m *MockGitCommand) Run(command string, args ...string) ([]byte, error) {
	defer func() { m.CallCount++ }()
	
	fullCommand := command + " " + strings.Join(args, " ")
	m.Commands = append(m.Commands, fullCommand)
	
	if m.CallCount < len(m.Outputs) {
		return []byte(m.Outputs[m.CallCount]), nil
	}
	
	if m.CallCount < len(m.Errors) {
		return nil, m.Errors[m.CallCount]
	}
	
	return []byte(""), nil
}

// TableTest represents a table-driven test case
type TableTest struct {
	Name    string
	Input   interface{}
	Want    interface{}
	WantErr bool
}

// RunTableTests runs table-driven tests
func RunTableTests(t *testing.T, tests []TableTest, testFunc func(t *testing.T, test TableTest)) {
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			testFunc(t, tt)
		})
	}
}

// SkipIfShort skips a test if running in short mode
func SkipIfShort(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}
}

// RequireCI skips a test if not running in CI environment
func RequireCI(t *testing.T) {
	if os.Getenv("CI") == "" {
		t.Skip("Skipping test - requires CI environment")
	}
}

// RequireDocker skips a test if Docker is not available
func RequireDocker(t *testing.T) {
	err := RunCommand("", "docker", "version")
	if err != nil {
		t.Skip("Skipping test - Docker not available")
	}
}

// RequireGit skips a test if Git is not available
func RequireGit(t *testing.T) {
	err := RunCommand("", "git", "version")
	if err != nil {
		t.Skip("Skipping test - Git not available")
	}
}