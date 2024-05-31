param(
    [switch] $installPackagesManager,
    [switch] $installDependencies,
    [switch] $build,
    [switch] $run,
    [switch] $clean,
    [switch] $release,
    [switch] $generateScoopInstaller,
    [switch] $generatePackage
)

$MAKE_SCRIPT_DIR = ($PSScriptRoot)

# APP INFORMATION AREA
$APP_INFORMATION_FILE = "$MAKE_SCRIPT_DIR\app-information"
$APP_NAME = ((Get-Content "$APP_INFORMATION_FILE" | findstr ^NAME=) -split '=')[1]
$APP_ID = ((Get-Content "$APP_INFORMATION_FILE" | findstr ^ID=) -split '=')[1]
$APP_VERSION = ((Get-Content "$APP_INFORMATION_FILE" | findstr ^VERSION=) -split '=')[1]
$APP_WIN_ICON = ((Get-Content "$APP_INFORMATION_FILE" | findstr ^WIN_ICON=) -split '=')[1]
$APP_LINUX_ICON = ((Get-Content "$APP_INFORMATION_FILE" | findstr ^LINUX_ICON=) -split '=')[1]
$APP_DISPLAY_NAME = ((Get-Content "$APP_INFORMATION_FILE" | findstr ^DISPLAY_NAME=) -split '=')[1]

# OTHERS
$RELEASE_DIR = "$MAKE_SCRIPT_DIR\release"
$POWERSHELL_VENDOR_DIR = "$MAKE_SCRIPT_DIR\vendor\powershell-utils"
$BASH_VENDOR_DIR = "$MAKE_SCRIPT_DIR\vendor\bash-utils"
$FYNE_CROSS_DIR = "$MAKE_SCRIPT_DIR\fyne-cross"
$BINARY = "$RELEASE_DIR\${APP_NAME}.exe"

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

function _copyDirectory($directory, $destination) {
    if ((directoryexists "$directory")) {
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
function _buildRelease() {
    param(
        [Parameter(Mandatory)]
        [ValidateSet("windows","linux")]
        [string] $platform
    )
    infolog "Build $platform app..."
    $icon = "$APP_WIN_ICON"
    export CGO_ENABLED=1
    if ($platform -eq "linux") {
        $icon = "$APP_LINUX_ICON"
    }
    evaladvanced "fyne-cross.exe $platform -arch=amd64 -name=`"$APP_NAME`" -app-id=`"$APP_ID`" -icon=`"/icon/$icon`" -app-version=`"$APP_VERSION`" `"$MAKE_SCRIPT_DIR\src\.`""
}

function _release() {
    if (!(confirm "Install or run docker before, please." $true)) {
        _exitSuccess
    }
    _clean
    _buildRelease -platform "windows"
    _buildRelease -platform "linux"
    _generatePackage
    _generateScoopInstaller
    deletedirectory "$FYNE_CROSS_DIR"
}

function _build() {
    _clean
    evaladvanced "go build -o `"$BINARY`" `"$MAKE_SCRIPT_DIR\src\main.go`""
    _preparePackage
}

function _preparePackage() {
    infolog "Copy necessary files..."
    $iconsDir = "$MAKE_SCRIPT_DIR\icon"
    $vendorReleaseDir = "$RELEASE_DIR\vendor"
    $releaseDate = (Get-date -Format "dd/MM/yyyy - HH:mm:ss")

    _createDirectory "$vendorReleaseDir"
    _copyDirectory -directory "$iconsDir" -destination "$RELEASE_DIR"
    _copyDirectory -directory "$POWERSHELL_VENDOR_DIR" -destination "$vendorReleaseDir"
    _copyDirectory -directory "$BASH_VENDOR_DIR" -destination "$vendorReleaseDir"
    _copyFile -file "$MAKE_SCRIPT_DIR\scripts\install.sh" -destination "$RELEASE_DIR"
    _copyFile -file "$MAKE_SCRIPT_DIR\scripts\uninstall.sh" -destination "$RELEASE_DIR"
    _copyFile -file "$MAKE_SCRIPT_DIR\README.md" -destination "$RELEASE_DIR"
    _copyFile -file "$FYNE_CROSS_DIR\bin\windows-amd64\${APP_NAME}.exe" -destination "$RELEASE_DIR"
    _copyFile -file "$FYNE_CROSS_DIR\bin\linux-amd64\src" -destination "$RELEASE_DIR\${APP_NAME}"
    _copyFile -file "$APP_INFORMATION_FILE" -destination "$RELEASE_DIR"
    writefile "$RELEASE_DIR\app-information" "`nRELEASE_DATE=${releaseDate}" -append
}

function _generatePackage() {
    _preparePackage
    # Generte package
    Compress-Archive "$RELEASE_DIR\*" -DestinationPath "$MAKE_SCRIPT_DIR\${APP_NAME}-${APP_VERSION}.zip" -Force
}

function _generateScoopInstaller() {
    $data = @"
{
  `"version`": `"$APP_VERSION`",
  `"description`": `"Lazygit repository management`",
  `"homepage`": `"https://github.com/zecarneiro/lazygit-repository-manager`",
  `"url`": `"https://github.com/zecarneiro/lazygit-repository-manager/releases/download/v${APP_VERSION}/lazygit-repository-manager-${APP_VERSION}.zip`",
  `"bin`": `"${APP_NAME}.exe`",
  `"shortcuts`": [
        [`"${APP_NAME}.exe`", `"$APP_DISPLAY_NAME`"]
	],
  `"persist`": `".data`",
  `"checkver`": `"github`",
  `"autoupdate`": {
    `"url`": `"https://github.com/zecarneiro/lazygit-repository-manager/releases/download/v`$version/lazygit-repository-manager-${APP_VERSION}.zip`"
  }
}
"@
    writefile "$MAKE_SCRIPT_DIR\${APP_NAME}.json" $data
}

function _clean() {
    deletedirectory "$RELEASE_DIR"
    deletedirectory "$FYNE_CROSS_DIR"
    deletefile "$MAKE_SCRIPT_DIR\go.sum"
    deletefile "$MAKE_SCRIPT_DIR\go.work.sum"
    deletefile "$MAKE_SCRIPT_DIR\${APP_NAME}-${APP_VERSION}.zip"
}

function main() {
    if ($installPackagesManager) {
        install_winget
        install_scoop
        infolog "Please, restart terminal."
    } elseif ($installDependencies) {
        install_golang
        install_cpp_c
        evaladvanced "go install github.com/fyne-io/fyne-cross@latest"
        evaladvanced "go install fyne.io/fyne/v2/cmd/fyne@latest"
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
    } elseif ($generateScoopInstaller) {
        _generateScoopInstaller
    }  elseif ($generatePackage) {
        _generatePackage
    } else {
        log "make.ps1 -installPackagesManager|-installDependencies|-build|-run|-release|-clean|-generateScoopInstaller|-generatePackage"
    }
}
main