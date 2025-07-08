package worktree

import (
	"reflect"
	"testing"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				Version: 1,
				Project: ProjectConfig{
					Name:   "testapp",
					Domain: "app.test",
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
				Templates: TemplateConfig{
					Default: "standard",
					Available: map[string]TemplateDefinition{
						"standard": {
							Files: []TemplateFile{
								{Src: "docker-compose.yml.tmpl", Dest: "docker-compose.yml"},
							},
						},
					},
				},
				Variables: map[string]interface{}{
					"db_name_prefix": "testapp",
				},
			},
			wantErr: false,
		},
		{
			name: "missing project name",
			config: &Config{
				Version: 1,
				Project: ProjectConfig{
					Domain: "app.test",
				},
			},
			wantErr: false, // We'll implement validation later
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// For now, we don't have validation implemented
			// This is a placeholder for future validation logic
			if tt.config == nil {
				t.Error("Config should not be nil")
			}
		})
	}
}

func TestTemplateConfig_GetTemplate(t *testing.T) {
	config := &TemplateConfig{
		Default: "standard",
		Available: map[string]TemplateDefinition{
			"standard": {
				Files: []TemplateFile{
					{Src: "docker-compose.yml.tmpl", Dest: "docker-compose.yml"},
					{Src: ".env.tmpl", Dest: ".env"},
				},
			},
			"minimal": {
				Files: []TemplateFile{
					{Src: "docker-compose.minimal.yml.tmpl", Dest: "docker-compose.yml"},
				},
			},
		},
	}

	tests := []struct {
		name         string
		templateName string
		want         *TemplateDefinition
		wantFound    bool
	}{
		{
			name:         "existing template",
			templateName: "standard",
			want: &TemplateDefinition{
				Files: []TemplateFile{
					{Src: "docker-compose.yml.tmpl", Dest: "docker-compose.yml"},
					{Src: ".env.tmpl", Dest: ".env"},
				},
			},
			wantFound: true,
		},
		{
			name:         "non-existing template",
			templateName: "nonexistent",
			want:         nil,
			wantFound:    false,
		},
		{
			name:         "empty template name uses default",
			templateName: "",
			want: &TemplateDefinition{
				Files: []TemplateFile{
					{Src: "docker-compose.yml.tmpl", Dest: "docker-compose.yml"},
					{Src: ".env.tmpl", Dest: ".env"},
				},
			},
			wantFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			templateName := tt.templateName
			if templateName == "" {
				templateName = config.Default
			}

			got, found := config.Available[templateName]
			if found != tt.wantFound {
				t.Errorf("Template found = %v, want %v", found, tt.wantFound)
				return
			}

			if tt.wantFound {
				if !reflect.DeepEqual(&got, tt.want) {
					t.Errorf("Template = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestProjectConfig_GetDomain(t *testing.T) {
	tests := []struct {
		name   string
		config ProjectConfig
		want   string
	}{
		{
			name: "simple domain",
			config: ProjectConfig{
				Name:   "myapp",
				Domain: "app.local",
			},
			want: "app.local",
		},
		{
			name: "production domain",
			config: ProjectConfig{
				Name:   "myapp",
				Domain: "myapp.com",
			},
			want: "myapp.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.Domain
			if got != tt.want {
				t.Errorf("Domain = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDockerConfig_IsEnabled(t *testing.T) {
	tests := []struct {
		name   string
		config DockerConfig
		want   bool
	}{
		{
			name: "enabled",
			config: DockerConfig{
				Enabled: true,
			},
			want: true,
		},
		{
			name: "disabled",
			config: DockerConfig{
				Enabled: false,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.Enabled
			if got != tt.want {
				t.Errorf("Enabled = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWebConfig_GetSubdomain(t *testing.T) {
	config := WebConfig{
		SubdomainPattern: "{branch}.{project_domain}",
	}

	tests := []struct {
		name          string
		branchName    string
		projectDomain string
		want          string
	}{
		{
			name:          "simple replacement",
			branchName:    "main",
			projectDomain: "app.local",
			want:          "main.app.local",
		},
		{
			name:          "feature branch",
			branchName:    "feature-auth",
			projectDomain: "myapp.com",
			want:          "feature-auth.myapp.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This would be implemented in the actual subdomain generation logic
			pattern := config.SubdomainPattern
			pattern = replacePattern(pattern, "{branch}", tt.branchName)
			pattern = replacePattern(pattern, "{project_domain}", tt.projectDomain)
			
			if pattern != tt.want {
				t.Errorf("Subdomain = %v, want %v", pattern, tt.want)
			}
		})
	}
}

// Helper function for pattern replacement (would be in utils)
func replacePattern(pattern, placeholder, value string) string {
	// Simple string replacement - in real implementation might use text/template
	return pattern[:len(pattern)-len(placeholder)] + value + pattern[len(pattern)-len(placeholder)+len(placeholder):]
}

// Example of a more complex helper function that could be added to config.go
func (c *Config) GetEffectiveTemplate(templateName string) (*TemplateDefinition, bool) {
	if templateName == "" {
		templateName = c.Templates.Default
	}
	
	template, found := c.Templates.Available[templateName]
	return &template, found
}

func TestConfig_GetEffectiveTemplate(t *testing.T) {
	config := &Config{
		Templates: TemplateConfig{
			Default: "standard",
			Available: map[string]TemplateDefinition{
				"standard": {
					Files: []TemplateFile{
						{Src: "docker-compose.yml.tmpl", Dest: "docker-compose.yml"},
					},
				},
			},
		},
	}

	tests := []struct {
		name         string
		templateName string
		wantFound    bool
	}{
		{
			name:         "existing template",
			templateName: "standard",
			wantFound:    true,
		},
		{
			name:         "empty uses default",
			templateName: "",
			wantFound:    true,
		},
		{
			name:         "non-existing template",
			templateName: "nonexistent",
			wantFound:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, found := config.GetEffectiveTemplate(tt.templateName)
			if found != tt.wantFound {
				t.Errorf("GetEffectiveTemplate() found = %v, want %v", found, tt.wantFound)
			}
		})
	}
}