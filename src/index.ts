#!/usr/bin/env node

import path from 'path';
import NodeMenu from 'node-menu';
import { IConfig } from './interface/Iconfig';
import { EOperation } from './enum/Eoperation';
import { DirectoriesUtils } from '../vendor/utils/typescript/directories-utils';
import { ConsoleUtils, ICommandInfo } from '../vendor/utils/typescript/console-utils';
import { FileUtils } from '../vendor/utils/typescript/file-utils';
import { LoggerUtils } from '../vendor/utils/typescript/logger-utils';
import { SystemUtils } from '../vendor/utils/typescript/system-utils';
import { Response } from '../vendor/utils/typescript/entities/response';
import { FunctionUtils } from '../vendor/utils/typescript/function-utils';

class LazygitRepositoryManager {
    private readonly CONFIG_FILE_NAME: string = path.resolve(DirectoriesUtils.systemConfig, 'lazygit-repository-manager.json');
    private readonly menu: NodeMenu = require('node-menu');
    private readonly delimiterWithTitle: number = 40;
    private consoleUtils: ConsoleUtils;
    private config: IConfig;

    constructor() {
        this.consoleUtils = new ConsoleUtils();
        if (!FileUtils.fileExist(this.CONFIG_FILE_NAME)) {
            FileUtils.writeJsonFile(this.CONFIG_FILE_NAME, {
                setStoreCredentials: false,
                data: []
            });
        }
        this.config = FileUtils.readJsonFile<IConfig>(this.CONFIG_FILE_NAME);
    }

    private updateConfig() {
        FileUtils.writeJsonFile(this.CONFIG_FILE_NAME, this.config);
    }

    private processAnotherCmd(cmd: ICommandInfo): Response<string> {
        return this.consoleUtils.execSync({...cmd, realTime: true, verbose: true, verboseOnlyCommand: true});
    }

    private run(directory?: string): void {
        const cmd: ICommandInfo = {
            cmd: 'lazygit',
            args: [],
        };
        if (directory) {
            cmd.args = ['-p', `'${directory}'`];
        }
        this.processAnotherCmd(cmd);
    }

    private isRepositoryValid(repositoryDir: string): boolean {
        return FileUtils.fileExist(path.resolve(repositoryDir, '.git'));
    }

    private processReadRepository(repository: string) {
        if (!this.isRepositoryValid(repository)) {
            LoggerUtils.warn(`Repository directory ${repository} is not valid`);
            const response = ConsoleUtils.readKeyboardSync(`Delete the repository: ${repository}?`, {canChoiceBeNull: true, choices: ['Y', 'n']});
            if (!response || response == 'Y') {
                this.config.data = this.config.data.filter(r => r != repository);
                this.updateConfig();
                this.processMenu(true);
            }
        } else {
            this.run(repository);
        }
    }

    private process(operation: EOperation) {
        if (operation == EOperation.openRepository) {
            const repository = ConsoleUtils.readKeyboardSync('Select repository', {choices: [], canChoiceBeNull: true});
            if (!this.isRepositoryValid(repository)) {
                LoggerUtils.warn(`Repository directory ${repository} is not valid`);
            } else {
                this.config.data.push(repository);
                this.updateConfig();
                this.run(repository);
                this.processMenu(true);
            }
        } else if (operation == EOperation.createNewRepository) {
            const fullPath = ConsoleUtils.readKeyboardSync('Insert full path of repository', {choices: [], canChoiceBeNull: true});
            if (!fullPath) {
                LoggerUtils.warn('Repository full path is not valid');
            } else {
                FileUtils.createDir(fullPath);
                this.processAnotherCmd({
                    cmd: 'git init',
                    cwd: fullPath
                });
            }
        } else if (operation == EOperation.cloneRepository) {
            const pathRepositry = ConsoleUtils.readKeyboardSync('Insert path or ENTER(to current directory) for repository', {choices: [], canChoiceBeNull: true});
            const url = ConsoleUtils.readKeyboardSync('Insert url of repository', {choices: [], canChoiceBeNull: true});
            if (!url) {
                LoggerUtils.warn('Repository url is not valid');
            } else {
                this.processAnotherCmd({
                    cmd: `git clone "${url}"`,
                    cwd: pathRepositry || undefined
                });
            }
        }
    }

    private setCredentialsHelper() {
        const title = this.config.setStoreCredentials ? 'Unset' : 'Set';
        this.menu.addItem(`${title} Credential Helper`, () => {
            let cmd: ICommandInfo;
            if (this.config.setStoreCredentials) {
                cmd = { cmd: 'git config --global --unset credential.helper' };
                FileUtils.deleteFile(path.resolve(SystemUtils.systemInfo.homeDir, '.git-credentials'));
                this.config.setStoreCredentials = false;
            } else {
                cmd = { cmd: 'git config --global credential.helper "store"' };
                this.config.setStoreCredentials = true;
            }
            this.processAnotherCmd(cmd);
            this.updateConfig();
            this.processMenu(true);
        }, this, []);
    }

    private processMenu(withReset?: boolean) {
        if (withReset) {
            this.menu.resetMenu();
        }
        this.menu.customHeader(() => {
            LoggerUtils.title('Lazygit Repository Manager');
        });
        this.menu.addDelimiter('*', this.delimiterWithTitle, 'Repositories');
        this.config.data?.forEach((repository) => {
            this.menu.addItem(repository, () => {
                this.processReadRepository(repository);
            }, this, []);
        });
        this.menu.addDelimiter('-', this.delimiterWithTitle, 'Others');
        this.setCredentialsHelper();
        this.menu.addItem('Open repository', () => {
            this.process(EOperation.openRepository);
        }, this, []);
        this.menu.addItem('Clone Repository', () => {
            this.process(EOperation.cloneRepository);
        }, this, []);
        this.menu.addItem('Create new repository', () => {
            this.process(EOperation.createNewRepository);
        }, this, []);
        this.menu.start();
    }

    private setUserInfoOnGit() {
        const maxCount = 5;
        const process = (isEmail: boolean) => {
            let count = 0;
            const cmd: ICommandInfo = { cmd: isEmail ? 'git config user.email' : 'git config user.name' };
            while(count < maxCount) {
                let consoleData = this.consoleUtils.execSync({...cmd, verbose: false});
                if (consoleData.hasError || !consoleData.data) {
                    LoggerUtils.warn(`${isEmail ? 'Email' : 'Username'} git info is not configured`);
                    const userInfo = ConsoleUtils.readKeyboardSync(`Insert ${isEmail ? 'Email' : 'Username'}`, {choices: [], canChoiceBeNull: true});
                    if (userInfo) {
                        this.processAnotherCmd({
                            cmd: isEmail? `git config --global user.email "${userInfo}"` : `git config --global user.name "${userInfo}"`
                        });
                    }
                } else {
                    break;
                }
                count++;
            }
        };
        process(false);
        process(true);
        
    }

    public static start() {
        const lazygitRepositoryManager = new LazygitRepositoryManager();
        try {
            lazygitRepositoryManager.setUserInfoOnGit();
            lazygitRepositoryManager.processMenu();
        } catch (error) {
            LoggerUtils.error(error as Error);
        }
    }
}
LazygitRepositoryManager.start();