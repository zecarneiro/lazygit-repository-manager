#!/bin/bash
# AUTHOR: Jose M. Noronha

declare OTHERAPPS_DIR="$HOME/.otherapps"
declare INSTALL_DIR="$OTHERAPPS_DIR/{APP_NAME}"
declare SHORTCUT="$HOME/.local/share/applications/{APP_NAME}.desktop"
declare ZIP_FILE="$OTHERAPPS_DIR/lazygit-repository-manager-{APP_VERSION}.zip"
declare SYMBOLIC_SYSTEM_FILE="/usr/bin/{APP_NAME}"
declare URL="https://github.com/zecarneiro/lazygit-repository-manager/releases/download/v{APP_VERSION}/lazygit-repository-manager-{APP_VERSION}.zip"

function _printInfo() {
    local operation="$1"
    local message="$2"
    echo "[INFO] ${operation}: ${message}"
}

function _printError() {
    local message="$1"
    echo "[ERROR] ${message}"
}

function _install() {
    local data="[Desktop Entry]
Version=1.0
Type=Application
Terminal=true
Exec=$INSTALL_DIR/{APP_NAME}
Name={APP_DISPLAY_NAME}
Comment={APP_DISPLAY_NAME}
Icon=$INSTALL_DIR/linux.png"
    
    # Start
    if [ ! "$(command -v wget)" ]; then
        _printError "Please install wget first"
        exit 1
    fi
    if [ ! "$(command -v unzip)" ]; then
        _printError "Please install unzip first"
        exit 1
    fi
    if [ -f "$ZIP_FILE" ]; then
        _printInfo "Remove" "$ZIP_FILE"
        rm "$ZIP_FILE"
    fi

    _printInfo "Create" "$INSTALL_DIR"
    mkdir -p "$INSTALL_DIR"
    
    _printInfo "Download" "{APP_DISPLAY_NAME}"
    wget -O "$ZIP_FILE" "$URL" -q --show-progress || exit 1

    _printInfo "Install" "{APP_DISPLAY_NAME}"
    unzip "$ZIP_FILE" -d "$INSTALL_DIR" || exit 1
    echo -e "$data" | tee "$SHORTCUT" >/dev/null
    chmod +x "$SHORTCUT" || exit 1
    chmod +x "$INSTALL_DIR/{APP_NAME}" || exit 1
    if [ -f "$ZIP_FILE" ]; then
        _printInfo "Remove" "$ZIP_FILE"
        rm "$ZIP_FILE"
    fi

    # TODO: Install command on system level
    #_printInfo "Install" "$SYMBOLIC_SYSTEM_FILE"
    #sudo ln "$INSTALL_DIR/{APP_NAME}" "$SYMBOLIC_SYSTEM_FILE"
}

function _uninstall() {
    if [ -h "$SYMBOLIC_SYSTEM_FILE" ]; then
        _printInfo "Remove" "$SYMBOLIC_SYSTEM_FILE"
        sudo rm "$SYMBOLIC_SYSTEM_FILE"
    fi
    if [ -d "$INSTALL_DIR" ]; then
        _printInfo "Remove" "$INSTALL_DIR"
        rm -rf "$INSTALL_DIR"
    fi
    if [ -f "$SHORTCUT" ]; then
        _printInfo "Remove" "$SHORTCUT"
        rm "$SHORTCUT"
    fi
}

function main() {
    echo "Processing installation of {APP_DISPLAY_NAME}"
    echo "1. Install"
    echo "2. Uninstall"
    echo "3. Exit"
    read -p "Insert an option: " option
    case $option in
        1)
            _uninstall
            _install
        ;;
        2) _uninstall ;;
        3) exit 0 ;;
        *) echo "Invalid option!"
    esac
}
main
