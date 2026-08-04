#!/bin/bash

APP_NAME="$1"
APP_VERSION="$2"
APP_DISPLAY_NAME="$3"
SO_TYPE="$4"
ROOT_DIR="$PWD"
RELEASE_DIR="$ROOT_DIR/release"
DEPLOY_DIR="$ROOT_DIR/$APP_NAME"
BINARY_DIR="$ROOT_DIR/bin"
IMAGES_DIR="$ROOT_DIR/docs/images"
INSTALLERS_DIR="$ROOT_DIR/scripts/installers"
echo ">>> Create release directory: $RELEASE_DIR"
rm -rf "$RELEASE_DIR"
mkdir -p "$RELEASE_DIR"

echo ">>> Copy binaries..."
cp "$BINARY_DIR/$APP_NAME" "$RELEASE_DIR/$APP_NAME"
cp "$BINARY_DIR/${APP_NAME}.exe" "$RELEASE_DIR/${APP_NAME}.exe"

echo ">>> Copy images..."
cp "$IMAGES_DIR/linux.png" "$RELEASE_DIR/linux.png"
cp "$IMAGES_DIR/win.ico" "$RELEASE_DIR/win.ico"

echo ">>> Create deploy directory: $DEPLOY_DIR"
rm -rf "$DEPLOY_DIR"
mkdir -p "$DEPLOY_DIR"

echo ">>> Copy installer and uninstaller..."
SCOOP_INSTALLER="$DEPLOY_DIR/$APP_NAME.json"
LINUX_INSTALLER="$DEPLOY_DIR/${APP_NAME}-install.sh"
LINUX_UNINSTALLER="$DEPLOY_DIR/${APP_NAME}-uninstall.sh"
cp "$INSTALLERS_DIR/scoop.json" "$SCOOP_INSTALLER"
cp "$INSTALLERS_DIR/linux-install.sh" "$LINUX_INSTALLER"
cp "$INSTALLERS_DIR/linux-uninstall.sh" "$LINUX_UNINSTALLER"
declare -A REPLACER=(
    [{APP_VERSION}]="${APP_VERSION}"
    [{APP_NAME}]="${APP_NAME}"
    [{APP_DISPLAY_NAME}]="${APP_DISPLAY_NAME}"
)
for key in "${!REPLACER[@]}"; do
    sed -i "s#$key#${REPLACER[$key]}#g" "$SCOOP_INSTALLER"
    sed -i "s#$key#${REPLACER[$key]}#g" "$LINUX_INSTALLER"
    sed -i "s#$key#${REPLACER[$key]}#g" "$LINUX_UNINSTALLER"
done

echo ">>> Generate package file..."
if [[ "$SO_TYPE" == "windows" ]]; then
    powershell.exe -Command "Compress-Archive '$RELEASE_DIR\*' -DestinationPath '$DEPLOY_DIR\${APP_NAME}-${APP_VERSION}.zip'" -Force
else
    cd "$RELEASE_DIR" || exit 1
    zip -rq "$DEPLOY_DIR/${APP_NAME}-${APP_VERSION}.zip" .
    cd "$ROOT_DIR" || exit 1
fi

echo ">>> Delete unnecessary files and directories..."
rm -rf "$RELEASE_DIR"
