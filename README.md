# service-ocr

## Development Setup

### Git Hooks

This repository includes a pre-commit hook that automatically runs `make format` before each commit to ensure consistent code formatting.

**Enable the hook:**
```bash
git config core.hooksPath .githooks
```

**What it does:**
- Detects staged `.go` files
- Runs `make format` to apply gofmt/goimports
- If formatting changes any files, the commit is blocked
- You must review and stage the formatted files before committing again

**To disable temporarily:**
```bash
git commit --no-verify
```