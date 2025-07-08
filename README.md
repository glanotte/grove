# GWT - Git Worktree Manager

A CLI tool for managing git worktrees with template support, Docker integration, and automatic web serving configuration.

## Features

- **Git Worktree Management**: Create, list, remove, and switch between git worktrees
- **Template Support**: Automatically generate configuration files from templates
- **Docker Integration**: Automatic Docker Compose setup with port management
- **Web Proxy Support**: Automatic subdomain configuration (nginx-proxy, Traefik)
- **Shell Integration**: Easy switching between worktrees with shell functions

## Installation

### From Source

```bash
git clone https://github.com/glanotte/gwt.git
cd gwt
make build
make install
```

### Using Go

```bash
go install github.com/glanotte/gwt@latest
```

## Quick Start

1. **Initialize a repository for worktree management:**
   ```bash
   gwt init git@github.com:user/repo.git
   cd repo
   ```

2. **Create configuration:**
   ```bash
   make dev-init  # Creates .gitworktree/config.yaml and templates
   ```

3. **Create a new worktree:**
   ```bash
   gwt create feature/new-feature
   ```

4. **List worktrees:**
   ```bash
   gwt list
   ```

5. **Switch to a worktree:**
   ```bash
   gwt switch feature/new-feature
   ```

## Configuration

Configuration is stored in `.gitworktree/config.yaml`. See `examples/configs/` for example configurations:

- `basic.yaml` - Simple configuration with Docker and nginx-proxy
- `advanced.yaml` - Full-featured configuration with all options

## Shell Integration

For enhanced functionality, add the shell integration to your shell:

```bash
# Add to ~/.bashrc or ~/.zshrc
source /path/to/gwt/scripts/gwt.sh
```

This provides:
- Auto-completion for commands and worktree names
- Enhanced `gwt switch` that actually changes directory
- Helper functions like `gwt-info`, `gwt-up`, `gwt-down`
- Aliases: `gwl`, `gwc`, `gws`, `gwr`

## Commands

- `gwt init <repo-url>` - Initialize a bare repository
- `gwt create <branch>` - Create a new worktree
- `gwt list` - List all worktrees
- `gwt remove <worktree>` - Remove a worktree
- `gwt switch <worktree>` - Switch to a worktree (with shell integration)
- `gwt version` - Show version information

## Templates

Templates are stored in `.gitworktree/templates/` and processed when creating worktrees. See `templates/README.md` for available variables and usage.

## Development

```bash
# Install dependencies
make deps

# Run tests
make test

# Build
make build

# Run
make run

# Clean
make clean
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and ensure they pass
5. Submit a pull request

## License

MIT License - see LICENSE file for details.