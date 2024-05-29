$APP_ROOT_DIR = ($PSScriptRoot)
$APP_NAME = "Lazygit Repo Manager"
$SCRIPT_SHELL_VENDOR_DIR = "$APP_ROOT_DIR\vendor\powershell-utils"

# IMPORT LIBS
. "$SCRIPT_SHELL_VENDOR_DIR\MainUtils.ps1"

function exitSuccess() {
    oklog "Done."
    exit 0
}

function uninstall {
    del_shortcut_file "$APP_NAME"
}

function main() {
    infolog "Uninstall ${APP_NAME}..."
    uninstall
    exitSuccess
}
main
