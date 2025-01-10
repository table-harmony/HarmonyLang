#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "Installing HarmonyLang..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed. Please install Go first.${NC}"
    exit 1
fi

# Create directory for binary
mkdir -p $HOME/.harmony/bin

# Install the package
go install github.com/table-harmony/HarmonyLang@latest

# Create harmony runner script
cat > $HOME/.harmony/bin/harmony << 'EOF'
#!/bin/bash

if [ "$1" = "run" ] && [ -n "$2" ]; then
    if [ -f "$2" ]; then
        $(go env GOPATH)/bin/HarmonyLang "$2"
    else
        echo "Error: File $2 not found"
        exit 1
    fi
elif [ "$1" = "repl" ]; then
    $(go env GOPATH)/bin/HarmonyLang
else
    echo "Usage:"
    echo "  harmony run <file.harmony>  - Run a Harmony source file"
    echo "  harmony repl               - Start Harmony REPL"
fi
EOF

# Make the script executable
chmod +x $HOME/.harmony/bin/harmony

# Add to PATH if not already there
HARMONY_PATH="export PATH=\$PATH:\$HOME/.harmony/bin"
if ! grep -q "$HARMONY_PATH" ~/.bashrc; then
    echo $HARMONY_PATH >> ~/.bashrc
    echo $HARMONY_PATH >> ~/.zshrc 2>/dev/null || true
fi

echo -e "${GREEN}Installation complete!${NC}"
echo "Please restart your terminal or run: source ~/.bashrc"
echo -e "\nUsage:"
echo "  harmony run <file.harmony>  - Run a Harmony source file"
echo "  harmony repl               - Start Harmony REPL"