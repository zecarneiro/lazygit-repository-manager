$APP_ROOT_DIR = ($PSScriptRoot)
$APP_ID = "lazygit-repo-manager"
$APP_NAME = "Lazygit Repo Manager"
$SCRIPT_SHELL_VENDOR_DIR = "$APP_ROOT_DIR\vendor\powershell-utils"

# IMPORT LIBS
. "$SCRIPT_SHELL_VENDOR_DIR\MainUtils.ps1"

function exitSuccess() {
    oklog "Done."
    exit 0
}

function install() {
    create_script_to_run_cmd_hidden "$APP_ROOT_DIR\$APP_ID" "$APP_ROOT_DIR\${APP_ID}.exe"
    create_shortcut_file -name "$APP_NAME" -target "$APP_ROOT_DIR\${APP_ID}.vbs" -icon "$APP_ROOT_DIR\icon\win-icon.ico"
}

function main() {
    infolog "Install ${APP_NAME}..."
    install
    exitSuccess
}
main
