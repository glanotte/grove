# Grove - Git Worktree Manager

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
git clone https://github.com/glanotte/grove.git
cd grove
make build
make install
```

### Using Go

```bash
go install github.com/glanotte/grove@latest
```

## Quick Start

1. **Initialize a repository for worktree management:**
   ```bash
   grove init git@github.com:user/repo.git
   cd repo
   ```

2. **Create configuration:**
   ```bash
   make dev-init  # Creates .gitworktree/config.yaml and templates
   ```

3. **Create a new worktree:**
   ```bash
   grove create feature/new-feature
   ```

4. **List worktrees:**
   ```bash
   grove list
   ```

5. **Switch to a worktree:**
   ```bash
   grove switch feature/new-feature
   ```

## Configuration

Configuration is stored in `.gitworktree/config.yaml`. See `examples/configs/` for example configurations:

- `basic.yaml` - Simple configuration with Docker and nginx-proxy
- `advanced.yaml` - Full-featured configuration with all options

## Shell Integration

For enhanced functionality, add the shell integration to your shell:

```bash
# Add to ~/.bashrc or ~/.zshrc
source /path/to/grove/scripts/grove.sh
```

This provides:
- Auto-completion for commands and worktree names
- Enhanced `grove switch` that actually changes directory
- Helper functions like `grove-info`, `grove-up`, `grove-down`
- Aliases: `grl`, `grc`, `grs`, `grr`

## Commands

- `grove init <repo-url>` - Initialize a bare repository
- `grove create <branch>` - Create a new worktree
- `grove list` - List all worktrees
- `grove remove <worktree>` - Remove a worktree
- `grove switch <worktree>` - Switch to a worktree (with shell integration)
- `grove version` - Show version information

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