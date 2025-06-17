# Homebrew Setup Guide

This guide explains how to set up Homebrew distribution for the Aphelion CLI.

## Prerequisites

1. **GitHub Repository**: `github.com/Exmplr-AI/aphelion-cli` (main project)
2. **Homebrew Tap Repository**: `github.com/Exmplr-AI/homebrew-aphelion` (needs to be created)

## Setup Steps

### 1. Create the Homebrew Tap Repository

```bash
# Create a new repository named 'homebrew-aphelion' under the Exmplr-AI organization
# GitHub URL: https://github.com/Exmplr-AI/homebrew-aphelion
```

### 2. Set up the Homebrew Tap Repository

Clone and set up the homebrew tap:

```bash
git clone git@github.com:Exmplr-AI/homebrew-aphelion.git
cd homebrew-aphelion

# Create Formula directory
mkdir Formula

# Copy the formula file from the main repo
cp ../aphelion-cli/aphelion.rb Formula/aphelion.rb

# Create initial commit
git add .
git commit -m "Initial Homebrew formula for aphelion"
git push origin main
```

### 3. Configure Deploy Key

1. **Generate SSH Key** (if not already done):
   ```bash
   ssh-keygen -t ed25519 -C "aphelion-homebrew-deploy" -f ~/.ssh/homebrew-deploy-key
   ```

2. **Add Deploy Key to homebrew-aphelion**:
   - Go to `github.com/Exmplr-AI/homebrew-aphelion` → Settings → Deploy keys
   - Add the public key (`homebrew-deploy-key.pub`) with **write access**
   - Title: "Aphelion CLI Release Automation"

3. **Add Private Key to aphelion-cli Secrets**:
   - Go to `github.com/Exmplr-AI/aphelion-cli` → Settings → Secrets
   - Create new secret: `HOMEBREW_DEPLOY_KEY`
   - Value: Contents of the private key (`homebrew-deploy-key`)

### 4. Push Main Repository

```bash
# In the aphelion-cli directory
git add .
git commit -m "Initial commit with Homebrew support"
git push origin main

# Create and push a tag to trigger the first release
git tag v1.0.0
git push origin v1.0.0
```

### 5. Verify Release Process

1. Check that GitHub Actions creates binaries and releases
2. Verify that the homebrew formula is updated automatically
3. Test installation: `brew tap exmplr-ai/aphelion && brew install aphelion`

## File Structure

### Main Repository (aphelion-cli)
```
.github/workflows/
└── release.yml              # Builds binaries, creates releases, and updates Homebrew formula

aphelion.rb                  # Template formula file
HOMEBREW_SETUP.md           # This file
```

### Homebrew Tap Repository (homebrew-aphelion)
```
Formula/
└── aphelion.rb             # Actual Homebrew formula
README.md                   # Basic documentation
```

## Installation for Users

Once set up, users can install aphelion using:

```bash
# Add the tap
brew tap exmplr-ai/aphelion

# Install aphelion
brew install aphelion

# Or in one command
brew install exmplr-ai/aphelion/aphelion
```

## Updating Releases

1. **Create a new tag**: `git tag v1.1.0 && git push origin v1.1.0`
2. **GitHub Actions will**:
   - Build binaries for all platforms
   - Create a GitHub release
   - Automatically update the Homebrew formula with new version and checksums
3. **Users update with**: `brew upgrade aphelion`

## Troubleshooting

### Common Issues

1. **Missing HOMEBREW_DEPLOY_KEY**: Add the deploy key private key as a GitHub secret
2. **Formula not updating**: Check that the deploy key has write access to the tap repo
3. **Build failures**: Verify Go version and dependencies in GitHub Actions
4. **SSH errors**: Ensure the deploy key is properly formatted (include -----BEGIN/END----- lines)

### Manual Formula Update

If automatic updates fail, manually update the formula:

1. Download the release checksums from GitHub
2. Update `Formula/aphelion.rb` with new version and checksums
3. Commit and push changes

### Testing Locally

```bash
# Test formula syntax
brew audit --strict Formula/aphelion.rb

# Test installation locally
brew install --build-from-source Formula/aphelion.rb

# Test upgrade
brew reinstall Formula/aphelion.rb
```

## Repository URLs

- **Main Project**: https://github.com/Exmplr-AI/aphelion-cli
- **Homebrew Tap**: https://github.com/Exmplr-AI/homebrew-aphelion
- **Releases**: https://github.com/Exmplr-AI/aphelion-cli/releases