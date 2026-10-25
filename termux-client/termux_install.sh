go build

echo "Installing medivpn client binary..."
if [ -d "$HOME/../usr/bin" ]; then
    install ./medivpn-client "$HOME/../usr/bin/medivpn"
else
    echo "/usr/bin not found."
    exit 1
fi
echo "\`- Done."

CFG_DIR="$HOME/.config"

echo "Checking for existing config directory..."
if [ ! -d "$CFG_DIR" ]; then
    echo "~/.config not found. Creating..."
    mkdir $CFG_DIR
    mkdir $CFG_DIR/medivpn
    echo "Copying config..."
    cp ./internal/config/config.yaml.example $CFG_DIR/medivpn/config.yaml
    echo "\`- Done."
elif [ ! -d "$CFG_DIR/medivpn" ]; then
    echo "No directory found. Creating and copying..."
    mkdir $CFG_DIR/medivpn
    cp ./internal/config/config.yaml.example $CFG_DIR/medivpn/config.yaml
    echo "\`- Done."
else
    echo "Usual config directory already found. Skipping config setup."
fi
echo "medivpn installed."