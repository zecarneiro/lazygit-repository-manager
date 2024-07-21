param(
    [switch] $installPackagesManager,
    [switch] $installDependencies,
    [switch] $build,
    [switch] $run,
    [switch] $clean,
    [switch] $release,
    [switch] $generateInstaller,
    [switch] $generatePackage
)

$MAKE_SCRIPT_DIR = ($PSScriptRoot)

# APP INFORMATION AREA
$APP_INFORMATION_FILE = $(Resolve-Path "$MAKE_SCRIPT_DIR\app-information")
$APP_NAME = ((Get-Content "$APP_INFORMATION_FILE" | Select-String -Pattern '^NAME=') -split '=')[1]
$APP_ID = ((Get-Content "$APP_INFORMATION_FILE" | Select-String -Pattern '^ID=') -split '=')[1]
$APP_VERSION = ((Get-Content "$APP_INFORMATION_FILE" | Select-String -Pattern '^VERSION=') -split '=')[1]
$APP_WIN_ICON = ((Get-Content "$APP_INFORMATION_FILE" | Select-String -Pattern '^WIN_ICON=') -split '=')[1]
$APP_LINUX_ICON = ((Get-Content "$APP_INFORMATION_FILE" | Select-String -Pattern '^LINUX_ICON=') -split '=')[1]
$APP_DISPLAY_NAME = ((Get-Content "$APP_INFORMATION_FILE" | Select-String -Pattern '^DISPLAY_NAME=') -split '=')[1]

# OTHERS
$RELEASE_DIR = "$MAKE_SCRIPT_DIR\release"
$POWERSHELL_VENDOR_DIR = "$MAKE_SCRIPT_DIR\vendor\powershell-utils"
$BINARY = "$RELEASE_DIR\${APP_NAME}.exe"
$BINARY_LINUX = "$RELEASE_DIR\${APP_NAME}"

# IMPORT LIBS
. "$POWERSHELL_VENDOR_DIR\MainUtils.ps1"

# ---------------------------------------------------------------------------- #
#                               GENERIC FUNCTIONS                              #
# ---------------------------------------------------------------------------- #
function _exitSuccess() {
    oklog "Done."
    exit 0
}

function _createDirectory($directory) {
    if (!(directoryexists "$directory")) {
        mkdir "$directory" | Out-Null
    }
}

function _copyDirectory($directory, $destination, $onlyFiles) {
    if ((directoryexists "$directory")) {
        if ($onlyFiles) {
            $directory = "$directory\*"
        }
        Copy-Item "$directory" -Destination "$destination" -Recurse -Force
    } else {
        warnlog "Not found directory: $directory"
    }
}

function _copyFile($file, $destination) {
    if ((fileexists "$file")) {
        Copy-Item "$file" -Destination "$destination" -Recurse -Force
    } else {
        warnlog "Not found file: $file"
    }
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
	go build -o "$BINARY" "$MAKE_SCRIPT_DIR\src\main.go"

	infolog "Build LINUX app..."
    export GOOS=linux
    go build -o "$BINARY_LINUX" "$MAKE_SCRIPT_DIR\src\main.go"
    _preparePackage
}

function _preparePackage() {
    infolog "Copy necessary files..."
    $iconsDir = "$MAKE_SCRIPT_DIR\icon"
    $releaseDate = (Get-date -Format "dd/MM/yyyy - HH:mm:ss")
    _copyDirectory -directory "$iconsDir" -destination "$RELEASE_DIR" -onlyFiles $true
    _copyFile -file "$MAKE_SCRIPT_DIR\README.md" -destination "$RELEASE_DIR"
    _copyFile -file "$APP_INFORMATION_FILE" -destination "$RELEASE_DIR"
    writefile "$RELEASE_DIR\app-information" "`nRELEASE_DATE=${releaseDate}" -append
}

function _generatePackage() {
    _preparePackage
    # Generte package
    Compress-Archive "$RELEASE_DIR\*" -DestinationPath "$MAKE_SCRIPT_DIR\${APP_NAME}-${APP_VERSION}.zip" -Force
}

function _generateInstaller() {
    $installerDir = (resolvePath "$MAKE_SCRIPT_DIR\installers")
    $replacer = @{ "{APP_VERSION}" = "${APP_VERSION}"; "{APP_NAME}" = "${APP_NAME}"; "{APP_DISPLAY_NAME}" = "${APP_DISPLAY_NAME}"}

    infolog "Generate SCOOP Installer..."
    $scoopInstallerFile = (resolvePath "$installerDir\scoop.json")
    $scoopInstallerDestFile = (resolvePath "$MAKE_SCRIPT_DIR\${APP_NAME}.json")
    _copyFile "$scoopInstallerFile" "$scoopInstallerDestFile"
    foreach ($key in $replacer.Keys) {
        writefile "$scoopInstallerDestFile" ((Get-Content "$scoopInstallerDestFile" -Raw) -replace "$key", $($replacer.$key))
    }
}

function _clean() {
    deletedirectory "$RELEASE_DIR"
    deletefile "$MAKE_SCRIPT_DIR\src\go.sum"
    deletefile "$MAKE_SCRIPT_DIR\${APP_NAME}-${APP_VERSION}.zip"
    deletefile "$MAKE_SCRIPT_DIR\${APP_NAME}.json"
    deletefile "$MAKE_SCRIPT_DIR\${APP_NAME}.sh"
}

function main() {
    if ($installPackagesManager) {
        install_winget
        install_scoop
        infolog "Please, restart terminal."
    } elseif ($installDependencies) {
        install_golang
        evaladvanced "go clean -cache -modcache -testcache"
        evaladvanced "go get -u github.com/rivo/tview@master"
        evaladvanced "go get -u github.com/gdamore/tcell/v2"
        evaladvanced "go get -u github.com/zecarneiro/golangutils"
        evaladvanced "go get -u github.com/zecarneiro/simpleconsoleui"
        evaladvanced "go mod tidy"
        infolog "Please, restart terminal."
    } elseif ($build) {
        _build
    } elseif ($run) {
        infolog "Run app..."
        Invoke-Expression "$BINARY"
    } elseif ($release) {
        _release
    } elseif ($clean) {
        _clean
    } elseif ($generateInstaller) {
        _generateInstaller
    }  elseif ($generatePackage) {
        _generatePackage
    } else {
        log "make.ps1 -installPackagesManager|-installDependencies|-build|-run|-release|-clean|-generateScoopInstaller|-generatePackage"
    }
}
main