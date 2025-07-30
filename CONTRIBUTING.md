# Contributing to GitHub PR Metrics

## Requirements

* mise (>=2025.6.8)

## Setup

1. Install mise
2. Create `mise.local.toml` in the project root:

```toml
[env]
GITHUB_URL = 'https://api.github.com'
GITHUB_TOKEN = '<your-personal-access-token>'
REPOSITORY = 'owner/repo'
```

## Development

Run the application:

```bash
mise run
```

Format code:

```bash
mise run format
```

Lint code:

```bash
mise run lint
```

## Contributing Process

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run format and lint
5. Submit a pull request

I'll review your PR and merge it if approved.
