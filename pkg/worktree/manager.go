package worktree

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "text/template"
)

// Manager handles git worktree operations
type Manager struct {
    BaseDir    string
    ConfigPath string
    Config     *Config
}


// NewManager creates a new worktree manager
func NewManager(baseDir string) (*Manager, error) {
    configPath := filepath.Join(baseDir, ".gitworktree", "config.yaml")
    
    m := &Manager{
        BaseDir:    baseDir,
        ConfigPath: configPath,
    }
    
    if err := m.loadConfig(); err != nil {
        return nil, fmt.Errorf("failed to load config: %w", err)
    }
    
    return m, nil
}

// CreateWorktree creates a new git worktree with templates
func (m *Manager) CreateWorktree(branchName, baseBranch, templateName string) error {
    // Sanitize branch name for use in paths and subdomains
    safeBranchName := m.sanitizeBranchName(branchName)
    
    // Calculate worktree path
    worktreePath := m.getWorktreePath(safeBranchName)
    
    // Create git worktree
    if err := m.createGitWorktree(worktreePath, branchName, baseBranch); err != nil {
        return fmt.Errorf("failed to create git worktree: %w", err)
    }
    
    // Process templates
    if err := m.processTemplates(worktreePath, branchName, templateName); err != nil {
        return fmt.Errorf("failed to process templates: %w", err)
    }
    
    // Setup Docker if enabled
    if m.Config.Docker.Enabled {
        if err := m.setupDocker(worktreePath, safeBranchName); err != nil {
            return fmt.Errorf("failed to setup Docker: %w", err)
        }
    }
    
    // Setup web proxy if enabled
    if m.Config.Web.Enabled {
        if err := m.setupWebProxy(safeBranchName); err != nil {
            return fmt.Errorf("failed to setup web proxy: %w", err)
        }
    }
    
    return nil
}

// loadConfig loads the configuration file
func (m *Manager) loadConfig() error {
    // In a real implementation, this would use viper or yaml.Unmarshal
    // For now, we'll use a default config
    m.Config = &Config{
        Version: 1,
        Project: ProjectConfig{
            Name:   "myapp",
            Domain: "app.lvh.me",
        },
        Worktree: WorktreeConfig{
            BasePath:      "./worktrees",
            NamingPattern: "{branch}",
        },
        Docker: DockerConfig{
            Enabled:     true,
            ComposeFile: "docker-compose.yml",
            PortOffset:  10000,
            NetworkName: "{project_name}_network",
        },
        Web: WebConfig{
            Enabled:          true,
            ProxyType:        "nginx-proxy",
            SubdomainPattern: "{branch}.{project_domain}",
        },
        Variables: make(map[string]interface{}),
    }
    return nil
}

// sanitizeBranchName makes a branch name safe for use in URLs and paths
func (m *Manager) sanitizeBranchName(branchName string) string {
    // Replace slashes with dashes
    safe := strings.ReplaceAll(branchName, "/", "-")
    // Replace other problematic characters
    safe = strings.ReplaceAll(safe, "_", "-")
    // Convert to lowercase for consistency
    safe = strings.ToLower(safe)
    return safe
}

// getWorktreePath returns the full path for a worktree
func (m *Manager) getWorktreePath(safeBranchName string) string {
    basePath := m.Config.Worktree.BasePath
    if !filepath.IsAbs(basePath) {
        basePath = filepath.Join(m.BaseDir, basePath)
    }
    return filepath.Join(basePath, safeBranchName)
}

// createGitWorktree creates the actual git worktree
func (m *Manager) createGitWorktree(path, branchName, baseBranch string) error {
    // Check if branch exists
    checkCmd := exec.Command("git", "show-ref", "--verify", "--quiet", fmt.Sprintf("refs/heads/%s", branchName))
    checkCmd.Dir = m.BaseDir
    
    branchExists := checkCmd.Run() == nil
    
    if branchExists {
        // Add worktree for existing branch
        cmd := exec.Command("git", "worktree", "add", path, branchName)
        cmd.Dir = m.BaseDir
        if output, err := cmd.CombinedOutput(); err != nil {
            return fmt.Errorf("git worktree add failed: %s", output)
        }
    } else {
        // Create new branch and worktree
        cmd := exec.Command("git", "worktree", "add", "-b", branchName, path, baseBranch)
        cmd.Dir = m.BaseDir
        if output, err := cmd.CombinedOutput(); err != nil {
            return fmt.Errorf("git worktree add failed: %s", output)
        }
    }
    
    return nil
}

// processTemplates processes all template files for the worktree
func (m *Manager) processTemplates(worktreePath, branchName, templateName string) error {
    if templateName == "" {
        templateName = m.Config.Templates.Default
    }
    
    templateDef, ok := m.Config.Templates.Available[templateName]
    if !ok {
        // No templates to process
        return nil
    }
    
    // Build template context
    ctx := m.buildTemplateContext(worktreePath, branchName)
    
    for _, file := range templateDef.Files {
        if err := m.processTemplateFile(worktreePath, file, ctx); err != nil {
            return fmt.Errorf("failed to process template %s: %w", file.Src, err)
        }
    }
    
    return nil
}

// buildTemplateContext creates the context for template processing
func (m *Manager) buildTemplateContext(worktreePath, branchName string) map[string]interface{} {
    safeBranchName := m.sanitizeBranchName(branchName)
    
    ctx := make(map[string]interface{})
    
    // Standard variables
    ctx["BranchName"] = safeBranchName
    ctx["OriginalBranchName"] = branchName
    ctx["WorktreePath"] = worktreePath
    ctx["ProjectName"] = m.Config.Project.Name
    ctx["ProjectDomain"] = m.Config.Project.Domain
    
    // Docker variables
    if m.Config.Docker.Enabled {
        ctx["NetworkName"] = strings.ReplaceAll(m.Config.Docker.NetworkName, "{project_name}", m.Config.Project.Name)
        ctx["WebPort"] = m.calculatePort(safeBranchName)
    }
    
    // Custom variables
    for k, v := range m.Config.Variables {
        ctx[k] = v
    }
    
    return ctx
}

// processTemplateFile processes a single template file
func (m *Manager) processTemplateFile(worktreePath string, file TemplateFile, ctx map[string]interface{}) error {
    templatePath := filepath.Join(m.BaseDir, ".gitworktree", "templates", file.Src)
    destPath := filepath.Join(worktreePath, file.Dest)
    
    // Read template
    content, err := os.ReadFile(templatePath)
    if err != nil {
        return fmt.Errorf("failed to read template: %w", err)
    }
    
    // Parse and execute template
    tmpl, err := template.New(file.Src).Parse(string(content))
    if err != nil {
        return fmt.Errorf("failed to parse template: %w", err)
    }
    
    // Create destination directory
    destDir := filepath.Dir(destPath)
    if err := os.MkdirAll(destDir, 0755); err != nil {
        return fmt.Errorf("failed to create directory: %w", err)
    }
    
    // Write processed template
    destFile, err := os.Create(destPath)
    if err != nil {
        return fmt.Errorf("failed to create destination file: %w", err)
    }
    defer destFile.Close()
    
    if err := tmpl.Execute(destFile, ctx); err != nil {
        return fmt.Errorf("failed to execute template: %w", err)
    }
    
    return nil
}

// calculatePort generates a unique port based on branch name
func (m *Manager) calculatePort(branchName string) int {
    // Simple hash-based port assignment
    hash := 0
    for _, c := range branchName {
        hash = (hash * 31 + int(c)) % 1000
    }
    return m.Config.Docker.PortOffset + hash
}

// setupDocker sets up Docker containers for the worktree
func (m *Manager) setupDocker(worktreePath, branchName string) error {
    // Ensure Docker network exists
    networkName := strings.ReplaceAll(m.Config.Docker.NetworkName, "{project_name}", m.Config.Project.Name)
    
    checkCmd := exec.Command("docker", "network", "inspect", networkName)
    if err := checkCmd.Run(); err != nil {
        // Create network if it doesn't exist
        createCmd := exec.Command("docker", "network", "create", networkName)
        if output, err := createCmd.CombinedOutput(); err != nil {
            return fmt.Errorf("failed to create Docker network: %s", output)
        }
    }
    
    // Start containers (optional - could be manual)
    composeFile := filepath.Join(worktreePath, m.Config.Docker.ComposeFile)
    if _, err := os.Stat(composeFile); err == nil {
        fmt.Printf("Docker Compose file created at: %s\n", composeFile)
        fmt.Printf("Run 'docker-compose up -d' in the worktree to start containers\n")
    }
    
    return nil
}

// setupWebProxy configures the web proxy for the worktree
func (m *Manager) setupWebProxy(branchName string) error {
    subdomain := strings.ReplaceAll(m.Config.Web.SubdomainPattern, "{branch}", branchName)
    subdomain = strings.ReplaceAll(subdomain, "{project_domain}", m.Config.Project.Domain)
    
    fmt.Printf("Web URL: https://%s\n", subdomain)
    
    // In a real implementation, this would:
    // - Update nginx-proxy configuration
    // - Or update Traefik labels
    // - Or update /etc/hosts
    
    return nil
}

// ListWorktrees lists all active worktrees
func (m *Manager) ListWorktrees() ([]WorktreeInfo, error) {
    cmd := exec.Command("git", "worktree", "list", "--porcelain")
    cmd.Dir = m.BaseDir
    
    output, err := cmd.Output()
    if err != nil {
        return nil, fmt.Errorf("failed to list worktrees: %w", err)
    }
    
    // Parse the output and build WorktreeInfo list
    // This is simplified - real implementation would properly parse porcelain output
    var worktrees []WorktreeInfo
    
    lines := strings.Split(string(output), "\n")
    for i := 0; i < len(lines); i += 4 {
        if i+2 >= len(lines) {
            break
        }
        
        if strings.HasPrefix(lines[i], "worktree ") {
            path := strings.TrimPrefix(lines[i], "worktree ")
            branch := ""
            if strings.HasPrefix(lines[i+2], "branch ") {
                branch = strings.TrimPrefix(lines[i+2], "branch refs/heads/")
            }
            
            worktrees = append(worktrees, WorktreeInfo{
                Path:   path,
                Branch: branch,
            })
        }
    }
    
    return worktrees, nil
}


// RemoveWorktree removes a worktree and cleans up resources
func (m *Manager) RemoveWorktree(name string, force bool) error {
    // Find worktree path
    worktrees, err := m.ListWorktrees()
    if err != nil {
        return err
    }
    
    var worktreePath string
    for _, wt := range worktrees {
        if filepath.Base(wt.Path) == name || wt.Branch == name {
            worktreePath = wt.Path
            break
        }
    }
    
    if worktreePath == "" {
        return fmt.Errorf("worktree '%s' not found", name)
    }
    
    // Stop Docker containers if running
    if m.Config.Docker.Enabled {
        composeFile := filepath.Join(worktreePath, m.Config.Docker.ComposeFile)
        if _, err := os.Stat(composeFile); err == nil {
            fmt.Printf("Stopping Docker containers...\n")
            stopCmd := exec.Command("docker-compose", "down")
            stopCmd.Dir = worktreePath
            stopCmd.Run() // Ignore errors
        }
    }
    
    // Remove git worktree
    args := []string{"worktree", "remove"}
    if force {
        args = append(args, "--force")
    }
    args = append(args, worktreePath)
    
    cmd := exec.Command("git", args...)
    cmd.Dir = m.BaseDir
    
    if output, err := cmd.CombinedOutput(); err != nil {
        return fmt.Errorf("failed to remove worktree: %s", output)
    }
    
    fmt.Printf("Worktree '%s' removed successfully\n", name)
    return nil
}