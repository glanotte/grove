package worktree

// Config represents the worktree configuration
type Config struct {
	Version   int                    `yaml:"version"`
	Project   ProjectConfig          `yaml:"project"`
	Worktree  WorktreeConfig         `yaml:"worktree"`
	Docker    DockerConfig           `yaml:"docker"`
	Web       WebConfig              `yaml:"web"`
	Templates TemplateConfig         `yaml:"templates"`
	Variables map[string]interface{} `yaml:"variables"`
}

type ProjectConfig struct {
	Name   string `yaml:"name"`
	Domain string `yaml:"domain"`
}

type WorktreeConfig struct {
	BasePath      string `yaml:"base_path"`
	NamingPattern string `yaml:"naming_pattern"`
}

type DockerConfig struct {
	Enabled     bool   `yaml:"enabled"`
	ComposeFile string `yaml:"compose_file"`
	PortOffset  int    `yaml:"port_offset"`
	NetworkName string `yaml:"network_name"`
}

type WebConfig struct {
	Enabled          bool   `yaml:"enabled"`
	ProxyType        string `yaml:"proxy_type"`
	SubdomainPattern string `yaml:"subdomain_pattern"`
}

type TemplateConfig struct {
	Default   string                        `yaml:"default"`
	Available map[string]TemplateDefinition `yaml:"available"`
}

type TemplateDefinition struct {
	Files []TemplateFile `yaml:"files"`
}

type TemplateFile struct {
	Src  string `yaml:"src"`
	Dest string `yaml:"dest"`
}

// WorktreeInfo contains information about a worktree
type WorktreeInfo struct {
	Path   string
	Branch string
	URL    string
}
