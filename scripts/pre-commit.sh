#!/bin/bash

# Pre-commit hook to scan for secrets
# Install this by running: cp scripts/pre-commit.sh .git/hooks/pre-commit && chmod +x .git/hooks/pre-commit

set -e

echo "🔍 Scanning for secrets before commit..."

# Check if gitleaks is installed
if command -v gitleaks &> /dev/null; then
    echo "Running gitleaks scan..."
    gitleaks detect --config .gitleaks.toml --verbose --no-git
    if [ $? -ne 0 ]; then
        echo "❌ Secrets detected! Please remove them before committing."
        echo "💡 Use 'git reset HEAD~1' to undo the last commit if needed."
        exit 1
    fi
    echo "✅ No secrets detected by gitleaks"
else
    echo "⚠️  gitleaks not installed. Install with: brew install gitleaks"
fi

# Check for common secret patterns manually
echo "Running manual secret checks..."

# Check for .env files (should not be committed)
echo "Checking for .env files in git index..."
if git diff --cached --name-only | grep -E "\.env$|\.env\.local$|\.env\.production$" > /dev/null; then
    echo "❌ .env files detected in commit! These should not be committed."
    echo "Files found:"
    git diff --cached --name-only | grep -E "\.env$|\.env\.local$|\.env\.production$"
    echo "💡 .env files should be local only and are in .gitignore"
    exit 1
fi

# Check for potential API keys in code (only in files being committed)
echo "Checking for potential secrets in files being committed..."
SUSPICIOUS_PATTERNS=(
    "WEBEX_PUBLIC_WORKSPACE_API_KEY.*=.*[a-zA-Z0-9]{20,}"
    "Bearer [a-zA-Z0-9]{20,}"
    "api[_-]?key.*=.*['\"][a-zA-Z0-9]{20,}['\"]"
    "secret.*=.*['\"][a-zA-Z0-9]{20,}['\"]"
    "token.*=.*['\"][a-zA-Z0-9]{20,}['\"]"
    "sk-[a-zA-Z0-9-_]{40,}"
)

for pattern in "${SUSPICIOUS_PATTERNS[@]}"; do
    if git diff --cached | grep -E "$pattern" >/dev/null 2>&1; then
        echo "❌ Potential secret found in commit matching pattern: $pattern"
        echo "Please remove secrets before committing."
        echo "💡 Use environment variables or .env files (which are gitignored)"
        exit 1
    fi
done

echo "✅ Manual secret checks passed"
echo "🎉 All security checks passed! Proceeding with commit..."
