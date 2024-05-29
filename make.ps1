param(
    [switch] $installPackagesManager,
    [switch] $installDependencies,
    [switch] $build,
    [switch] $run,
    [switch] $clean,
    [switch] $release
)

$MAKE_SCRIPT_DIR = ($PSScriptRoot)
$RELEASE_DIR = "$MAKE_SCRIPT_DIR\release"
$BINARY_NAME = "lazygit-repository-manager"
$BINARY = "$RELEASE_DIR\${BINARY_NAME}.exe"
$POWERSHELL_VENDOR_DIR = "$MAKE_SCRIPT_DIR\vendor\powershell-utils"

# IMPORT LIBS
. "$POWERSHELL_VENDOR_DIR\MainUtils.ps1"

function _exitSuccess() {
    oklog "Done."
    exit 0
}

function createDirectory($directory) {
    if (!(directoryexists "$directory")) {
        mkdir "$directory" | Out-Null
    }
}

function _copyFiles() {
    infolog "Copy necessary files..."
    $iconsDir = "$MAKE_SCRIPT_DIR\icon"
    $vendorReleaseDir = "$RELEASE_DIR\vendor"

    createDirectory "$vendorReleaseDir"
    Copy-Item "$iconsDir" -Destination "$RELEASE_DIR" -Recurse -Force
    Copy-Item "$POWERSHELL_VENDOR_DIR" -Destination "$vendorReleaseDir" -Recurse -Force
    Copy-Item "$MAKE_SCRIPT_DIR\scripts\install.ps1" -Destination "$RELEASE_DIR" -Recurse -Force
    Copy-Item "$MAKE_SCRIPT_DIR\scripts\uninstall.ps1" -Destination "$RELEASE_DIR" -Recurse -Force
    Copy-Item "$MAKE_SCRIPT_DIR\README.md" -Destination "$RELEASE_DIR" -Recurse -Force
}

function _prepareRelease($withClean) {
    if ($withClean) {
        _clean
    }
    _build
    _copyFiles
}

function _release() {
    _prepareRelease $true
    Compress-Archive "$RELEASE_DIR\*" -DestinationPath "$MAKE_SCRIPT_DIR\${BINARY_NAME}-win.zip" -Force
}

function _clean() {
    deletedirectory "$RELEASE_DIR"
    deletefile "$MAKE_SCRIPT_DIR\go.sum"
    deletefile "$MAKE_SCRIPT_DIR\go.work.sum"
    deletefile "$MAKE_SCRIPT_DIR\$BINARY_NAME.zip"
}

function _build() {
    infolog "Build WINDOWS app..."
	export GOOS=windows
	export GOARCH=amd64
    evaladvanced "go build -o `"$BINARY`" `"$MAKE_SCRIPT_DIR\src`""
}

function main() {
    if ($installPackagesManager) {
        install_winget
        install_scoop
        infolog "Please, restart terminal."
        _exitSuccess
    } elseif ($installDependencies) {
        install_golang
        install_cpp_c
        evaladvanced "go mod tidy"
        infolog "Please, restart terminal."
        _exitSuccess
    } elseif ($build) {
        export CGO_ENABLED=1
        _build
        _exitSuccess
    } elseif ($run) {
        export CGO_ENABLED=1
        _prepareRelease $false
        infolog "Run app..."
        Invoke-Expression "$BINARY"
    } elseif ($release) {
        _release
        _exitSuccess
    } elseif ($clean) {
        _clean
        _exitSuccess
    }
    log "make.ps1 -installPackagesManager|-installDependencies|-build|-run|-release|-clean"
}
main