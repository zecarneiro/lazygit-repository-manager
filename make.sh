#!/bin/bash
# AUTHOR: José M. Noronha

declare MAKE_SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"

# APP INFORMATION AREA
declare APP_INFORMATION_FILE="$MAKE_SCRIPT_DIR/app-information"
declare APP_NAME=$(cat "$APP_INFORMATION_FILE" | grep ^NAME= | cut -d '=' -f2)
declare APP_ID=$(cat "$APP_INFORMATION_FILE" | grep ^ID= | cut -d '=' -f2)
declare APP_VERSION=$(cat "$APP_INFORMATION_FILE" | grep ^VERSION= | cut -d '=' -f2)
declare APP_WIN_ICON=$(cat "$APP_INFORMATION_FILE" | grep ^WIN_ICON= | cut -d '=' -f2)
declare APP_LINUX_ICON=$(cat "$APP_INFORMATION_FILE" | grep ^LINUX_ICON= | cut -d '=' -f2)
declare APP_DISPLAY_NAME=$(cat "$APP_INFORMATION_FILE" | grep ^DISPLAY_NAME= | cut -d '=' -f2)

# OTHERS
declare RELEASE_DIR="$MAKE_SCRIPT_DIR/release"
declare POWERSHELL_VENDOR_DIR="$MAKE_SCRIPT_DIR/vendor/powershell-utils"
declare BASH_VENDOR_DIR="$MAKE_SCRIPT_DIR/vendor/bash-utils"
declare FYNE_CROSS_DIR="$MAKE_SCRIPT_DIR/fyne-cross"
declare BINARY="$RELEASE_DIR/${APP_NAME}"
declare BINARY_WINDOWS="$RELEASE_DIR/${APP_NAME}.exe"

# IMPORT LIBS
. "$BASH_VENDOR_DIR/main-utils.sh"

# ---------------------------------------------------------------------------- #
#                               GENERIC FUNCTIONS                              #
# ---------------------------------------------------------------------------- #
function _exitSuccess {
    oklog "Done."
    exit 0
}

function _copyDirectory {
    local directory="$1"
    local destination="$2"
    local onlyFiles="$3"
    if [ $(directoryexists "$directory") == true ]; then
        if [ $onlyFiles == true ]; then
            directory="$directory/."
        fi
        evaladvanced "cp -r \"$directory\" \"$destination\""
    else
        warnlog "Not found directory: $directory"
    fi
}

function _copyFile {
    local file="$1"
    local destination="$2"
    if [ $(fileexists "$file") == true ]; then
        cp "$file" "$destination"
    else
        warnlog "Not found file: $file"
    fi
}

# ---------------------------------------------------------------------------- #
#                                     MAIN                                     #
# ---------------------------------------------------------------------------- #
function _release() {
    _build
    _generatePackage
}

function _build() {
    export GOARCH=amd64
    _clean

    infolog "Build WINDOWS app..."
    export GOOS=windows
	go build -o "$BINARY_WINDOWS" "$MAKE_SCRIPT_DIR/src/main.go"

    infolog "Build LINUX app..."
    export GOOS=linux
    go build -o "$BINARY" "$MAKE_SCRIPT_DIR/src/main.go"
    _preparePackage
}

function _preparePackage() {
    infolog "Copy necessary files..."
    local iconsDir="$MAKE_SCRIPT_DIR/icon"
    local releaseDate=$(date '+%d/%m/%Y %H:%M:%S')
    _copyDirectory "$iconsDir" "$RELEASE_DIR" true
    _copyFile "$MAKE_SCRIPT_DIR/README.md" "$RELEASE_DIR"
    _copyFile "$APP_INFORMATION_FILE" "$RELEASE_DIR"
    writefile "$RELEASE_DIR/app-information" "\nRELEASE_DATE=${releaseDate}" -append
}

function _generatePackage() {
    _preparePackage
    # Generte package
    pushd .
    cd "$RELEASE_DIR"
    zip -rq "$MAKE_SCRIPT_DIR/${APP_NAME}-${APP_VERSION}.zip" .
    popd
}

function _generateInstaller() {
    local installerDir="$MAKE_SCRIPT_DIR/installers"
    local -A replacer=(
        [{APP_VERSION}]="${APP_VERSION}"
        [{APP_NAME}]="${APP_NAME}"
        [{APP_DISPLAY_NAME}]="${APP_DISPLAY_NAME}"
    )

    infolog "Generate LINUX Installer..."
    local linuxInstallerFile="$installerDir/linux.sh"
    local linuxInstallerDestFile="$MAKE_SCRIPT_DIR/${APP_NAME}.sh"
    _copyFile "$linuxInstallerFile" "$linuxInstallerDestFile"
    for key in "${!replacer[@]}"; do
        sed -i "s#$key#${replacer[$key]}#g" "$linuxInstallerDestFile"
    done
}

function _clean() {
    deletedirectory "$RELEASE_DIR"
    deletefile "$MAKE_SCRIPT_DIR/go.sum"
    deletefile "$MAKE_SCRIPT_DIR/go.work.sum"
    deletefile "$MAKE_SCRIPT_DIR/${APP_NAME}-${APP_VERSION}.zip"
    deletefile "$MAKE_SCRIPT_DIR/${APP_NAME}.json"
    deletefile "$MAKE_SCRIPT_DIR/${APP_NAME}.sh"
}

function main() {
    case "${1}" in
        -installDependencies)
            install_go
            evaladvanced "go clean -cache -modcache -testcache"
            evaladvanced "go get -u github.com/rivo/tview@master"
            evaladvanced "go get -u github.com/gdamore/tcell/v2"
            evaladvanced "go get -u github.com/zecarneiro/golangutils"
            evaladvanced "go get -u github.com/zecarneiro/simpleconsoleui"
            evaladvanced "go mod tidy"
            infolog "Please, restart terminal."
        ;;
        -build) _build ;;
        -run)
            infolog "Run app..."
            eval "$BINARY"
        ;;
        -release) _release ;;
        -clean) _clean ;;
        -generateInstaller) _generateInstaller ;;
        -generatePackage) _generatePackage ;;
        *) log "make.sh -installPackagesManager|-installDependencies|-build|-run|-release|-clean|-generateInstaller|-generatePackage" ;;
    esac
}
main "$@"
