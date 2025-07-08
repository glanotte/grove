# Examples

This directory contains example configurations and usage scenarios for GWT.

## Configuration Examples

### basic.yaml
A simple configuration suitable for most projects:
- Basic Docker setup with nginx-proxy
- Simple template processing
- Standard port allocation

### advanced.yaml
A comprehensive configuration showing all available options:
- Multiple template inheritance
- Advanced Docker configuration
- Database and Redis integration
- SSL/TLS configuration
- Monitoring and security features

## Usage Examples

### Basic Workflow

```bash
# Initialize repository
gwt init git@github.com:user/myapp.git
cd myapp

# Setup configuration
cp examples/configs/basic.yaml .gitworktree/config.yaml
cp -r templates .gitworktree/

# Create feature branch worktree
gwt create feature/user-auth

# List worktrees
gwt list

# Switch to worktree (with shell integration)
gwt switch feature/user-auth

# Start services
gwt-up

# Your app is now running at https://user-auth.app.lvh.me
```

### Advanced Setup

```bash
# Use advanced configuration
cp examples/configs/advanced.yaml .gitworktree/config.yaml

# Create worktree with custom template
gwt create feature/api-refactor --template microservices

# The worktree now has:
# - Custom docker-compose.yml for microservices
# - Database with automatic seeding
# - Redis with branch-specific databases
# - SSL certificates via Let's Encrypt
```

## Template Customization

Create your own templates by:

1. Adding template files to `.gitworktree/templates/`
2. Defining template sets in your config.yaml
3. Using template variables like `{{.BranchName}}` and `{{.ProjectName}}`

See `templates/README.md` for full documentation on available variables.