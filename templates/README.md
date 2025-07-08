# Templates

This directory contains template files that will be processed when creating new worktrees.

## Available Templates

- `docker-compose.yml.tmpl` - Docker Compose configuration with application, database, and Redis services
- `.env.tmpl` - Environment variables for the application

## Template Variables

The following variables are available in templates:

- `{{.ProjectName}}` - Project name from config
- `{{.BranchName}}` - Sanitized branch name (safe for URLs/paths)
- `{{.OriginalBranchName}}` - Original branch name
- `{{.ProjectDomain}}` - Project domain from config
- `{{.WorktreePath}}` - Full path to the worktree directory
- `{{.WebPort}}` - Calculated web port for the service
- `{{.NetworkName}}` - Docker network name
- `{{.DbNamePrefix}}` - Database name prefix from config
- `{{.RedisPrefix}}` - Redis key prefix from config

## Template Functions

- `{{add .WebPort 1}}` - Add numbers (useful for calculating port offsets)
- `{{.BranchName | lower}}` - Convert to lowercase
- `{{.BranchName | upper}}` - Convert to uppercase

## Custom Variables

You can add custom variables in your config.yaml under the `variables` section:

```yaml
variables:
  custom_var: "value"
  another_var: 123
```

These will be available as `{{.CustomVar}}` and `{{.AnotherVar}}` in templates.