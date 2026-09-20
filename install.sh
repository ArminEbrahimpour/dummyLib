#!/bin/bash


sudo apt install libayatana-appindicator-glib-dev

go mod download

go build -o dummylib ./cmd

sudo cp dummylib /usr/bin/

set -euo pipefail

# Resolve absolute path of the project dir (the dir containing this script)
projectdir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Make sure the target dirs exist
mkdir -p "$HOME/.local/share/applications"
mkdir -p "$HOME/.config/autostart"

desktop_entry="[Desktop Entry]
Type=Application
Name=dummylib
Comment=Run dummylib
Exec=/usr/bin/dummylib
Icon=$projectdir/statics/foxies.png
Terminal=false
Categories=Utility;"

# 1) Menu entry (shows in app launcher)
printf '%s\n' "$desktop_entry" > "$HOME/.local/share/applications/dummylib.desktop"

# 2) Autostart entry (runs at login) // uncomment this line if you want auto start on startup   
#printf '%s\n' "$desktop_entry" > "$HOME/.config/autostart/dummylib.desktop"

# Refresh the desktop database so the menu entry appears immediately
if command -v update-desktop-database >/dev/null 2>&1; then
    update-desktop-database "$HOME/.local/share/applications" || true
fi

echo "Installed desktop entries."
echo "  Menu:      ~/.local/share/applications/dummylib.desktop"
echo "  Autostart: ~/.config/autostart/dummylib.desktop"
