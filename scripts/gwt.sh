#!/usr/bin/env bash
# Git Worktree Manager Shell Integration
# Add this to your .bashrc or .zshrc:
# source /path/to/gwt.sh

# Main gwt function that wraps the gwt binary
gwt() {
    case "$1" in
        switch|cd)
            # Special handling for switch command to change directory
            if [ -z "$2" ]; then
                echo "Usage: gwt switch <worktree-name>"
                return 1
            fi
            
            # Get the worktree path from the gwt binary
            local worktree_path
            worktree_path=$(command gwt switch "$2" 2>/dev/null)
            
            if [ $? -eq 0 ] && [ -n "$worktree_path" ] && [ -d "$worktree_path" ]; then
                cd "$worktree_path" || return 1
                echo "Switched to worktree: $2"
                
                # Optional: Activate virtual environment if exists
                if [ -f "venv/bin/activate" ]; then
                    source venv/bin/activate
                elif [ -f ".venv/bin/activate" ]; then
                    source .venv/bin/activate
                fi
                
                # Optional: Source .envrc if using direnv
                if [ -f ".envrc" ] && command -v direnv &> /dev/null; then
                    direnv allow .
                fi
                
                # Optional: Show worktree status
                echo "Branch: $(git branch --show-current)"
                echo "Status: $(git status --short | wc -l) uncommitted changes"
                
                # Show Docker status if compose file exists
                if [ -f "docker-compose.yml" ]; then
                    echo "Docker: Run 'docker-compose up -d' to start services"
                fi
            else
                echo "Error: Worktree '$2' not found"
                return 1
            fi
            ;;
            
        create)
            # Run the create command and auto-switch if successful
            command gwt "$@"
            if [ $? -eq 0 ] && [ -n "$2" ]; then
                echo ""
                echo "Worktree created successfully!"
                read -p "Switch to new worktree? [Y/n] " -n 1 -r
                echo
                if [[ $REPLY =~ ^[Yy]$ ]] || [[ -z $REPLY ]]; then
                    gwt switch "$2"
                fi
            fi
            ;;
            
        *)
            # Pass through all other commands
            command gwt "$@"
            ;;
    esac
}

# Bash completion for gwt
_gwt_completion() {
    local cur prev opts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    
    # Main commands
    local commands="init create list remove switch version help"
    
    case "${prev}" in
        gwt)
            COMPREPLY=( $(compgen -W "${commands}" -- ${cur}) )
            return 0
            ;;
        switch|remove|cd)
            # Get worktree names for completion
            if command -v gwt &> /dev/null; then
                local worktrees=$(gwt list --format=names 2>/dev/null | grep -v "^$")
                COMPREPLY=( $(compgen -W "${worktrees}" -- ${cur}) )
            fi
            return 0
            ;;
        create)
            # Suggest branch names from git
            local branches=$(git branch -r 2>/dev/null | sed 's/origin\///' | grep -v HEAD)
            COMPREPLY=( $(compgen -W "${branches}" -- ${cur}) )
            return 0
            ;;
    esac
}

# ZSH completion for gwt
if [ -n "$ZSH_VERSION" ]; then
    _gwt_zsh_completion() {
        local -a commands worktrees branches
        
        commands=(
            'init:Initialize a bare repository for worktree management'
            'create:Create a new worktree from a branch'
            'list:List all worktrees with their status'
            'remove:Remove a worktree and its associated resources'
            'switch:Switch to a different worktree'
            'version:Print version information'
            'help:Show help'
        )
        
        case $words[2] in
            switch|remove|cd)
                if command -v gwt &> /dev/null; then
                    worktrees=($(gwt list --format=names 2>/dev/null | grep -v "^$"))
                    _describe 'worktree' worktrees
                fi
                ;;
            create)
                branches=($(git branch -r 2>/dev/null | sed 's/origin\///' | grep -v HEAD))
                _describe 'branch' branches
                ;;
            *)
                _describe 'command' commands
                ;;
        esac
    }
    
    compdef _gwt_zsh_completion gwt
else
    # Bash completion
    complete -F _gwt_completion gwt
fi

# Aliases for common operations
alias gwl='gwt list'
alias gwc='gwt create'
alias gws='gwt switch'
alias gwr='gwt remove'

# Helper function to show current worktree info
gwt-info() {
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        echo "Not in a git repository"
        return 1
    fi
    
    local worktree_path=$(git rev-parse --show-toplevel)
    local branch=$(git branch --show-current)
    local uncommitted=$(git status --short | wc -l)
    
    echo "Worktree: $(basename "$worktree_path")"
    echo "Branch: $branch"
    echo "Path: $worktree_path"
    echo "Uncommitted changes: $uncommitted"
    
    # Show Docker status
    if [ -f "$worktree_path/docker-compose.yml" ]; then
        local running_containers=$(docker-compose ps -q 2>/dev/null | wc -l)
        echo "Docker containers: $running_containers running"
    fi
    
    # Show web URL if configured
    if [ -f "$worktree_path/.env" ]; then
        local app_url=$(grep "APP_URL=" "$worktree_path/.env" | cut -d'=' -f2)
        if [ -n "$app_url" ]; then
            echo "URL: $app_url"
        fi
    fi
}

# Function to start Docker services in current worktree
gwt-up() {
    if [ -f "docker-compose.yml" ]; then
        echo "Starting Docker services..."
        docker-compose up -d
        echo ""
        docker-compose ps
    else
        echo "No docker-compose.yml found in current directory"
        return 1
    fi
}

# Function to stop Docker services in current worktree
gwt-down() {
    if [ -f "docker-compose.yml" ]; then
        echo "Stopping Docker services..."
        docker-compose down
    else
        echo "No docker-compose.yml found in current directory"
        return 1
    fi
}

# Function to show logs for current worktree
gwt-logs() {
    if [ -f "docker-compose.yml" ]; then
        docker-compose logs -f "$@"
    else
        echo "No docker-compose.yml found in current directory"
        return 1
    fi
}

# Export functions for use in subshells
export -f gwt gwt-info gwt-up gwt-down gwt-logs

# Optional: Add worktree info to prompt
# For bash PS1 or zsh PROMPT, you can add:
# $(git rev-parse --show-toplevel 2>/dev/null | xargs basename | sed 's/^/[/' | sed 's/$/]/')