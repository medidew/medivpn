go build

if [ -d "$HOME/.local/bin" ]; then
    install ./medivpn-client "$HOME/.local/bin/medivpn"
else
    sudo install ./medivpn-client "/usr/local/bin/medivpn"
fi

CFG_DIR="$HOME/.config/medivpn"

if [ ! -d "$CFG_DIR" ]; then
    mkdir $CFG_DIR
    cp ./internal/config/config.yaml.example $CFG_DIR/config.yaml
fi