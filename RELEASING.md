# Release Process

This document describes how to create a new release of semantic-markdown.

## Prerequisites

- Git with push access to the repository
- Go 1.25 or later installed
- GitHub CLI (`gh`) installed: `brew install gh`
- Homebrew tap repository cloned locally (optional, for Homebrew distribution)
- All tests passing: `make ci`

## Release Checklist

### 1. Prepare the Release

1. Ensure all changes are merged to `main` branch
2. Update CHANGELOG.md:
   - Move items from `[Unreleased]` to new version section
   - Add release date in format `[X.Y.Z] - YYYY-MM-DD`
   - Document any breaking changes
3. Update version in any documentation if needed
4. Commit changes:
   ```bash
   git add CHANGELOG.md
   git commit -m "Prepare vX.Y.Z release"
   git push origin main
   ```

### 2. Create Version Tag

```bash
# Create annotated tag
git tag -a vX.Y.Z -m "Release vX.Y.Z"

# Push tag to repository
git push origin vX.Y.Z
```

### 3. Build Release Artifacts

```bash
# Run complete release workflow
# This will:
# - Run pre-release checks (lint, test, coverage)
# - Build binaries for all platforms
# - Create distribution archives
# - Generate checksums
# - Extract release notes from CHANGELOG
# - Generate Homebrew formula
make release
```

The release workflow creates:
- `dist/release/` - Platform-specific binaries
- `dist/tarballs/` - Compressed archives with checksums
- `dist/release-notes.md` - Extracted release notes
- `dist/homebrew/semantic-md.rb` - Homebrew formula

### 4. Create GitHub Release

```bash
# Create GitHub release with all artifacts
make github-release
```

Or manually with gh CLI:
```bash
gh release create vX.Y.Z dist/tarballs/* \
  --title "Release vX.Y.Z" \
  --notes-file dist/release-notes.md
```

### 5. Update Homebrew Tap (Optional)

If you have a Homebrew tap repository:

```bash
# Update the formula in your local tap repository
make homebrew-update TAP_DIR=/path/to/homebrew-tap

# Commit and push to tap repository
cd /path/to/homebrew-tap
git add Formula/semantic-md.rb
git commit -m "Update semantic-md to vX.Y.Z"
git push origin main
```

Users can then install with:
```bash
brew tap thorstenpfister/tap
brew install semantic-md
```

### 6. Verify Installation

Test that the release works:

```bash
# Download and test binary
curl -LO https://github.com/thorstenpfister/semantic-markdown/releases/download/vX.Y.Z/semantic-md-vX.Y.Z-darwin-arm64.tar.gz
tar -xzf semantic-md-vX.Y.Z-darwin-arm64.tar.gz
./semantic-md-darwin-arm64 version

# Test Homebrew installation (if available)
brew update
brew upgrade semantic-md
semantic-md version

# Or fresh install
brew uninstall semantic-md
brew install thorstenpfister/tap/semantic-md
```

### 7. Post-Release

1. Announce the release (social media, discussions, etc.)
2. Monitor for issues
3. Update documentation website if applicable
4. Start new `[Unreleased]` section in CHANGELOG.md for future changes

## Release Targets Reference

| Target | Description |
|--------|-------------|
| `make pre-release` | Run all checks before releasing (lint, test, coverage, verify tag) |
| `make release-build` | Build binaries for all platforms with version information |
| `make release-archives` | Create distribution archives (.tar.gz, .zip) |
| `make release-checksums` | Generate SHA256 checksums for all archives |
| `make release-validate` | Validate all artifacts and verify checksums |
| `make release-notes` | Extract release notes from CHANGELOG for current version |
| `make homebrew-formula` | Generate Homebrew formula with checksums |
| `make homebrew-update TAP_DIR=...` | Update local Homebrew tap repository |
| `make release` | Run complete release workflow (recommended) |
| `make github-release` | Create GitHub release with artifacts (requires `gh`) |
| `make clean-release` | Clean release build artifacts |

## Supported Platforms

The release process builds binaries for:

- **Linux**: AMD64, ARM64
- **macOS**: AMD64 (Intel), ARM64 (Apple Silicon)
- **Windows**: AMD64

All binaries include version information via ldflags:
- Version (from git tag)
- Git commit hash
- Build date

## Troubleshooting

### Version Detection Issues

The Makefile uses `git describe --tags` to determine the version. Ensure:
- You're on a tagged commit
- The working directory is clean (no uncommitted changes)
- Tags are pushed to the remote

Check current version:
```bash
git describe --tags --always
```

### Dirty Tag Error

If you see "Warning: Not on a clean tagged commit" during `make pre-release`:
```bash
# Check for uncommitted changes
git status

# Commit or stash changes
git add .
git commit -m "Your message"

# Or stash temporarily
git stash
```

### Checksum Mismatches

If Homebrew reports checksum mismatches:
1. Verify the release artifacts on GitHub match your local build
2. Download the actual release tarball and compute its checksum:
   ```bash
   curl -LO https://github.com/thorstenpfister/semantic-markdown/releases/download/vX.Y.Z/semantic-md-vX.Y.Z-darwin-arm64.tar.gz
   sha256sum semantic-md-vX.Y.Z-darwin-arm64.tar.gz
   ```
3. Update the formula with the correct checksum
4. Commit and push the updated formula

### GitHub CLI Authentication

If `gh release create` fails:
```bash
# Authenticate with GitHub
gh auth login

# Verify authentication
gh auth status

# Test release creation
gh release list
```

### Missing Dependencies

If builds fail:
```bash
# Verify Go version
go version  # Should be 1.25 or later

# Update dependencies
go mod tidy
go mod verify

# Run tests
make test
```

## Version Numbering

We follow [Semantic Versioning](https://semver.org/):

- **MAJOR** version (X.0.0): Incompatible API changes
- **MINOR** version (0.Y.0): Add functionality in backwards compatible manner
- **PATCH** version (0.0.Z): Backwards compatible bug fixes

Examples:
- `1.0.0` - First stable release
- `1.1.0` - Added new feature (backwards compatible)
- `1.1.1` - Bug fix (backwards compatible)
- `2.0.0` - Breaking API change

## Hotfix Releases

For urgent bugfixes:

1. Create a hotfix branch from the release tag:
   ```bash
   git checkout -b hotfix/v1.0.1 v1.0.0
   ```

2. Apply the fix and commit:
   ```bash
   # Make changes
   git add .
   git commit -m "Fix critical bug"
   ```

3. Follow the normal release process with a patch version bump:
   ```bash
   # Update CHANGELOG.md
   git add CHANGELOG.md
   git commit -m "Prepare v1.0.1 release"

   # Tag and release
   git tag -a v1.0.1 -m "Release v1.0.1"
   git push origin hotfix/v1.0.1
   git push origin v1.0.1
   make release
   make github-release
   ```

4. Merge back to main:
   ```bash
   git checkout main
   git merge hotfix/v1.0.1
   git push origin main
   ```

## Release Candidates

For testing releases before going stable:

1. Create a release candidate tag:
   ```bash
   git tag -a v1.1.0-rc1 -m "Release candidate 1.1.0-rc1"
   git push origin v1.1.0-rc1
   ```

2. Build and test:
   ```bash
   make release
   # Test the binaries
   ```

3. If issues found, fix and create rc2, rc3, etc.

4. When stable, create final release:
   ```bash
   git tag -a v1.1.0 -m "Release v1.1.0"
   git push origin v1.1.0
   make release
   make github-release
   ```

## Example: Complete Release Workflow

Here's a complete example for releasing version 1.1.0:

```bash
# 1. Prepare
cd /path/to/semantic-markdown
git checkout main
git pull origin main

# 2. Update CHANGELOG
# Edit CHANGELOG.md: move Unreleased → [1.1.0] - 2026-01-15
git add CHANGELOG.md
git commit -m "Prepare v1.1.0 release"
git push origin main

# 3. Create and push tag
git tag -a v1.1.0 -m "Release v1.1.0"
git push origin v1.1.0

# 4. Build release
make release

# Review artifacts
ls -lh dist/tarballs/
cat dist/release-notes.md

# 5. Create GitHub release
make github-release

# 6. Update Homebrew tap (if available)
make homebrew-update TAP_DIR=../homebrew-tap
cd ../homebrew-tap
git add Formula/semantic-md.rb
git commit -m "Update semantic-md to v1.1.0"
git push origin main

# 7. Verify
brew update
brew upgrade semantic-md
semantic-md version  # Should show v1.1.0

# 8. Test
echo "<h1>Test</h1>" | semantic-md convert
```

## Support

For questions or issues with the release process:
- Open an issue: https://github.com/thorstenpfister/semantic-markdown/issues
- Check existing releases: https://github.com/thorstenpfister/semantic-markdown/releases
