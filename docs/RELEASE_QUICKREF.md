# Release Quick Reference Card

Quick reference for creating releases of claude-init CLI.

## Prerequisites

```bash
# Verify working tree is clean
git status

# Verify tests pass
make test-race

# Verify linting passes
make lint
```

## Create Release (Automated)

```bash
# Run the release script
./scripts/release.sh v0.1.0
```

The script will:
1. Check working tree is clean
2. Run tests and linters
3. Build for all platforms
4. Create release archives
5. Generate checksums
6. Create git tag
7. Ask to push tag

## Create Release (Manual)

```bash
# 1. Build release
make release VERSION=v0.1.0

# 2. Create tag
git tag -a v0.1.0 -m "Release v0.1.0"

# 3. Push tag (triggers GitHub Actions)
git push origin v0.1.0
```

## Verify Release

```bash
# Using GitHub CLI
gh release view v0.1.0

# Download and test a binary
curl -LO https://github.com/danielrossellosanchez/claude-init/releases/download/v0.1.0/claude-init-v0.1.0-linux-amd64.tar.gz
curl -LO https://github.com/danielrossellosanchez/claude-init/releases/download/v0.1.0/claude-init-v0.1.0-linux-amd64.tar.gz.sha256

# Verify checksum
sha256sum -c claude-init-v0.1.0-linux-amd64.tar.gz.sha256

# Extract and test
tar -xzf claude-init-v0.1.0-linux-amd64.tar.gz
./claude-init-linux-amd64 version --verbose
```

## Rollback Release

```bash
# Delete tag locally
git tag -d v0.1.0

# Delete tag remotely
git push origin :refs/tags/v0.1.0

# Delete GitHub release
gh release delete v0.1.0
```

## Version Format

Follow Semantic Versioning: `v<major>.<minor>.<patch>[-<prerelease>]`

Examples:
- `v1.0.0` - Stable release
- `v1.2.3-beta.1` - Pre-release
- `v2.0.0-rc.1` - Release candidate

## Platform Matrix

| OS       | Arch   | Binary                        |
|----------|--------|-------------------------------|
| Linux    | amd64  | claude-init-*-linux-amd64.tar.gz |
| Linux    | arm64  | claude-init-*-linux-arm64.tar.gz |
| macOS    | amd64  | claude-init-*-darwin-amd64.tar.gz |
| macOS    | arm64  | claude-init-*-darwin-arm64.tar.gz |
| Windows  | amd64  | claude-init-*-windows-amd64.tar.gz |

## Troubleshooting

### Tag already exists
```bash
git tag -d v0.1.0
git push origin :refs/tags/v0.1.0
```

### Working tree not clean
```bash
git status
git add .
git commit -m "Prepare for release"
```

### Tests failing
```bash
make test-race TEST_FLAGS=-v
```

### Build failing
```bash
make clean-all
make build-all
```

## CI/CD Pipeline

When you push a tag:
1. CI runs tests and linters
2. Build creates binaries for all platforms
3. Release creates GitHub Release
4. Artifacts are uploaded

Monitor at: https://github.com/danielrossellosanchez/claude-init/actions

## Documentation

Full documentation: [docs/RELEASE.md](docs/RELEASE.md)
Implementation details: [docs/PHASE10_BUILD_RELEASE.md](docs/PHASE10_BUILD_RELEASE.md)
