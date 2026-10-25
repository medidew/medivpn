go build

echo "Installing medivpn client binary..."
if [ -d "$HOME/.local/bin" ]; then
    install ./medivpn-client "$HOME/.local/bin/medivpn"
else
    sudo install ./medivpn-client "/usr/local/bin/medivpn"
fi
echo "\`- Done."

CFG_DIR="$HOME/.config/medivpn"

echo "Checking for existing config directory..."
if [ ! -d "$CFG_DIR" ]; then
    echo "No directory found. Creating and copying."
    mkdir $CFG_DIR
    cp ./internal/config/config.yaml.example $CFG_DIR/config.yaml
    echo "\`- Done."
else
    echo "Usual config directory already found. Skipping config setup."
fi
echo "medivpn installed."