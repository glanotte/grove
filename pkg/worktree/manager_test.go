package worktree

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewManager(t *testing.T) {
	tests := []struct {
		name    string
		baseDir string
		wantErr bool
	}{
		{
			name:    "valid directory",
			baseDir: "/tmp/test-repo",
			wantErr: false,
		},
		{
			name:    "empty directory",
			baseDir: "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory for testing
			tempDir := t.TempDir()
			if tt.baseDir != "" {
				tt.baseDir = tempDir
			}

			manager, err := NewManager(tt.baseDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewManager() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if manager == nil {
					t.Error("NewManager() returned nil manager")
				}
				if manager.BaseDir != tt.baseDir {
					t.Errorf("NewManager() BaseDir = %v, want %v", manager.BaseDir, tt.baseDir)
				}
			}
		})
	}
}

func TestManager_sanitizeBranchName(t *testing.T) {
	manager := &Manager{}

	tests := []struct {
		name       string
		branchName string
		want       string
	}{
		{
			name:       "simple branch name",
			branchName: "main",
			want:       "main",
		},
		{
			name:       "feature branch with slash",
			branchName: "feature/user-auth",
			want:       "feature-user-auth",
		},
		{
			name:       "branch with underscores",
			branchName: "feature_user_auth",
			want:       "feature-user-auth",
		},
		{
			name:       "mixed case",
			branchName: "Feature/User-Auth",
			want:       "feature-user-auth",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := manager.sanitizeBranchName(tt.branchName)
			if got != tt.want {
				t.Errorf("sanitizeBranchName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_getWorktreePath(t *testing.T) {
	tempDir := t.TempDir()
	
	manager := &Manager{
		BaseDir: tempDir,
		Config: &Config{
			Worktree: WorktreeConfig{
				BasePath: "./worktrees",
			},
		},
	}

	tests := []struct {
		name           string
		safeBranchName string
		want           string
	}{
		{
			name:           "simple branch",
			safeBranchName: "main",
			want:           filepath.Join(tempDir, "worktrees", "main"),
		},
		{
			name:           "feature branch",
			safeBranchName: "feature-auth",
			want:           filepath.Join(tempDir, "worktrees", "feature-auth"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := manager.getWorktreePath(tt.safeBranchName)
			if got != tt.want {
				t.Errorf("getWorktreePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_calculatePort(t *testing.T) {
	manager := &Manager{
		Config: &Config{
			Docker: DockerConfig{
				PortOffset: 10000,
			},
		},
	}

	tests := []struct {
		name       string
		branchName string
		wantMin    int
		wantMax    int
	}{
		{
			name:       "consistent hashing",
			branchName: "main",
			wantMin:    10000,
			wantMax:    10999,
		},
		{
			name:       "different branch different port",
			branchName: "feature-auth",
			wantMin:    10000,
			wantMax:    10999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := manager.calculatePort(tt.branchName)
			if got < tt.wantMin || got > tt.wantMax {
				t.Errorf("calculatePort() = %v, want between %v and %v", got, tt.wantMin, tt.wantMax)
			}

			// Test consistency - same input should give same output
			got2 := manager.calculatePort(tt.branchName)
			if got != got2 {
				t.Errorf("calculatePort() not consistent: %v != %v", got, got2)
			}
		})
	}
}

func TestManager_buildTemplateContext(t *testing.T) {
	manager := &Manager{
		Config: &Config{
			Project: ProjectConfig{
				Name:   "testapp",
				Domain: "app.test",
			},
			Docker: DockerConfig{
				Enabled:     true,
				NetworkName: "{project_name}_network",
				PortOffset:  10000,
			},
			Variables: map[string]interface{}{
				"db_name_prefix": "testapp",
				"redis_prefix":   "testapp",
			},
		},
	}

	worktreePath := "/tmp/worktrees/feature-auth"
	branchName := "feature/auth"

	ctx := manager.buildTemplateContext(worktreePath, branchName)

	// Test required fields
	if ctx["BranchName"] != "feature-auth" {
		t.Errorf("Expected BranchName 'feature-auth', got %v", ctx["BranchName"])
	}
	if ctx["OriginalBranchName"] != "feature/auth" {
		t.Errorf("Expected OriginalBranchName 'feature/auth', got %v", ctx["OriginalBranchName"])
	}
	if ctx["ProjectName"] != "testapp" {
		t.Errorf("Expected ProjectName 'testapp', got %v", ctx["ProjectName"])
	}
	if ctx["ProjectDomain"] != "app.test" {
		t.Errorf("Expected ProjectDomain 'app.test', got %v", ctx["ProjectDomain"])
	}
	if ctx["WorktreePath"] != worktreePath {
		t.Errorf("Expected WorktreePath '%s', got %v", worktreePath, ctx["WorktreePath"])
	}
	if ctx["NetworkName"] != "testapp_network" {
		t.Errorf("Expected NetworkName 'testapp_network', got %v", ctx["NetworkName"])
	}

	// Test custom variables
	if ctx["db_name_prefix"] != "testapp" {
		t.Errorf("Expected db_name_prefix 'testapp', got %v", ctx["db_name_prefix"])
	}
	if ctx["redis_prefix"] != "testapp" {
		t.Errorf("Expected redis_prefix 'testapp', got %v", ctx["redis_prefix"])
	}

	// Test WebPort is calculated
	if _, ok := ctx["WebPort"]; !ok {
		t.Error("Expected WebPort to be present in context")
	}
}

// Benchmark tests
func BenchmarkManager_sanitizeBranchName(b *testing.B) {
	manager := &Manager{}
	branchName := "feature/very-long-branch-name-with-many-characters"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.sanitizeBranchName(branchName)
	}
}

func BenchmarkManager_calculatePort(b *testing.B) {
	manager := &Manager{
		Config: &Config{
			Docker: DockerConfig{
				PortOffset: 10000,
			},
		},
	}
	branchName := "feature-auth"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.calculatePort(branchName)
	}
}

// Test helper functions
func setupTestManager(t *testing.T) (*Manager, string) {
	tempDir := t.TempDir()
	
	// Create .grove directory
	groveDir := filepath.Join(tempDir, ".grove")
	if err := os.MkdirAll(groveDir, 0755); err != nil {
		t.Fatalf("Failed to create .grove directory: %v", err)
	}

	manager, err := NewManager(tempDir)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	return manager, tempDir
}

func TestMain(m *testing.M) {
	// Setup code here if needed
	code := m.Run()
	// Teardown code here if needed
	os.Exit(code)
}