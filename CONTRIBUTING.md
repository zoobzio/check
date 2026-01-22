# Contributing to check

Thank you for considering contributing to check.

## Development Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/zoobzio/check.git
   cd check
   ```

2. Install development tools:
   ```bash
   make install-tools
   ```

3. Install git hooks:
   ```bash
   make install-hooks
   ```

## Development Workflow

Run `make help` to see all available commands:

```bash
make help
```

### Common Commands

| Command | Purpose |
|---------|---------|
| `make test` | Run all tests with race detector |
| `make test-unit` | Run unit tests only (short mode) |
| `make lint` | Run linters |
| `make lint-fix` | Run linters with auto-fix |
| `make coverage` | Generate coverage report |
| `make check` | Quick validation (test + lint) |
| `make ci` | Full CI simulation |

### Before Submitting

1. Ensure tests pass:
   ```bash
   make test
   ```

2. Ensure linting passes:
   ```bash
   make lint
   ```

3. Or run the full check:
   ```bash
   make check
   ```

## Pull Request Process

1. Create a feature branch from `main`
2. Make your changes
3. Ensure `make check` passes
4. Submit a pull request

## Code Style

- Follow standard Go conventions
- Run `gofmt` (handled by linters)
- Write tests for new functionality
- Keep functions focused and small

## Commit Messages

Use clear, descriptive commit messages:

- `feat: add UUID validation`
- `fix: handle empty string in Required`
- `docs: update README examples`
- `test: add edge case tests for Email`

## Questions?

Open an issue for any questions or concerns.
