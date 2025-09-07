# Project: awsctl

awsctl is a Go-based CLI tool for managing AWS infrastructure across multiple regions with a single command. Here's what it does:

## Core libraries

- cobra - Command-line interface framework for building CLI commands
- AWS SDK for Go - Official AWS SDK for interacting with AWS services
- godotenv - Loading environment variables from .env files
- golang.org/x/net - Extended networking capabilities

## Architecture

- Command Structure: Uses Cobra's command pattern with main commands (get, set, delete, scan, create, describe, check)
- Helper Package: Contains AWS service-specific logic (ec2.go, rds.go, s3.go, ecr.go, etc.)
- Output Package: Handles different output formats (text/JSON)
- Hooks Package: Supports webhook integrations (Slack notifications)
- Multi-region Support: Can execute commands on single region or all regions with --region all
- Dry-run Safety: Commands run in dry-run mode by default, require --yes flag for execution

## Coding standards

- Use golang best practices
- Keep the code simple, and readable
- Make sure to maintain KISS concept

