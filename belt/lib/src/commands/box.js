"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const tslib_1 = require("tslib");
const command_1 = require("@oclif/command");
const chalk_1 = tslib_1.__importDefault(require("chalk"));
const cli_ux_1 = tslib_1.__importDefault(require("cli-ux"));
const cli = tslib_1.__importStar(require("inquirer"));
const truffle_box_1 = require("../services/truffle-box");
class Box extends command_1.Command {
    async run() {
        const { flags, args } = this.parse(Box);
        if (flags.list) {
            return this.handleList();
        }
        if (flags.interactive) {
            return await this.handleInteractive(args.path, flags.dryRun);
        }
        if (flags.solVer) {
            return this.handleNonInteractive(args.path, flags.dryRun, flags.solVer);
        }
        this._help();
    }
    /**
     * Handle printing out a list of available solidity versions
     */
    handleList() {
        const versions = truffle_box_1.getSolidityVersions().map(([alias, full]) => ({
            alias,
            full,
        }));
        this.log(chalk_1.default.greenBright('Available Solidity Versions'));
        cli_ux_1.default.table(versions, {
            alias: {},
            full: {},
        });
        this.log('');
    }
    /**
     * Handle interactive mode.
     * Prompts user for a solidity version number then proceeds to
     * do a find-replace within their box for the selected version
     *
     * @param path The path to the truffle box
     * @param dryRun Dont replace the file contents, print the diff instead
     */
    async handleInteractive(path, dryRun) {
        const solidityVersions = truffle_box_1.getSolidityVersions();
        const { solcVersion } = await cli.prompt([
            {
                name: 'solcVersion',
                type: 'list',
                choices: solidityVersions.map(([, version]) => version),
                message: 'What version of solidity do you want to use with your smart contracts?',
            },
        ]);
        const fullVersion = this.getFullVersion(solcVersion);
        truffle_box_1.modifyTruffleBoxWith(fullVersion, path, dryRun);
        this.log(chalk_1.default.greenBright(`Done!\nPlease run "npm i" to install the new changes made.`));
    }
    /**
     * Handle non-interactive mode "--solVer".
     * solidity version number then proceeds to
     * do a find-replace within their box for the selected version
     *
     * @param path The path to the truffle box
     * @param dryRun Dont replace the file contents, print the diff instead
     * @param versionAliasOrVersion Either a solidity version alias "v0.6" | "0.6" or its full version "0.6.2"
     */
    handleNonInteractive(path, dryRun, versionAliasOrVersion) {
        const fullVersion = this.getFullVersion(versionAliasOrVersion);
        truffle_box_1.modifyTruffleBoxWith(fullVersion, path, dryRun);
    }
    getFullVersion(versionAliasOrVersion) {
        let fullVersion;
        try {
            fullVersion = truffle_box_1.getSolidityVersionBy(versionAliasOrVersion);
        }
        catch {
            const error = chalk_1.default.red('Could not find given solidity version\n');
            this.log(error);
            this.handleList();
            this.exit(1);
        }
        return fullVersion;
    }
}
exports.default = Box;
Box.description = 'Modify a truffle box to a specified solidity version';
Box.examples = [
    'belt box --solVer 0.6 -d path/to/box',
    'belt box --interactive path/to/box',
    'belt box -l',
];
Box.flags = {
    help: command_1.flags.help({ char: 'h' }),
    interactive: command_1.flags.boolean({
        char: 'i',
        description: 'run this command in interactive mode',
    }),
    solVer: command_1.flags.string({
        char: 's',
        description: 'the solidity version to change the truffle box to\neither a solidity version alias "v0.6" | "0.6" or its full version "0.6.2"',
    }),
    list: command_1.flags.boolean({
        char: 'l',
        description: 'list the available solidity versions',
    }),
    dryRun: command_1.flags.boolean({
        char: 'd',
        description: 'output the replaced strings, but dont change them',
    }),
};
Box.args = [
    {
        name: 'path',
        description: 'the path to the truffle box',
    },
];
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiYm94LmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vc3JjL2NvbW1hbmRzL2JveC50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7QUFBQSw0Q0FBK0M7QUFFL0MsMERBQXlCO0FBQ3pCLDREQUF1QjtBQUN2QixzREFBK0I7QUFDL0IseURBSWdDO0FBQ2hDLE1BQXFCLEdBQUksU0FBUSxpQkFBTztJQXFDdEMsS0FBSyxDQUFDLEdBQUc7UUFDUCxNQUFNLEVBQUUsS0FBSyxFQUFFLElBQUksRUFBRSxHQUFHLElBQUksQ0FBQyxLQUFLLENBQUMsR0FBRyxDQUFDLENBQUE7UUFDdkMsSUFBSSxLQUFLLENBQUMsSUFBSSxFQUFFO1lBQ2QsT0FBTyxJQUFJLENBQUMsVUFBVSxFQUFFLENBQUE7U0FDekI7UUFFRCxJQUFJLEtBQUssQ0FBQyxXQUFXLEVBQUU7WUFDckIsT0FBTyxNQUFNLElBQUksQ0FBQyxpQkFBaUIsQ0FBQyxJQUFJLENBQUMsSUFBSSxFQUFFLEtBQUssQ0FBQyxNQUFNLENBQUMsQ0FBQTtTQUM3RDtRQUNELElBQUksS0FBSyxDQUFDLE1BQU0sRUFBRTtZQUNoQixPQUFPLElBQUksQ0FBQyxvQkFBb0IsQ0FBQyxJQUFJLENBQUMsSUFBSSxFQUFFLEtBQUssQ0FBQyxNQUFNLEVBQUUsS0FBSyxDQUFDLE1BQU0sQ0FBQyxDQUFBO1NBQ3hFO1FBRUQsSUFBSSxDQUFDLEtBQUssRUFBRSxDQUFBO0lBQ2QsQ0FBQztJQUVEOztPQUVHO0lBQ0ssVUFBVTtRQUNoQixNQUFNLFFBQVEsR0FBRyxpQ0FBbUIsRUFBRSxDQUFDLEdBQUcsQ0FBQyxDQUFDLENBQUMsS0FBSyxFQUFFLElBQUksQ0FBQyxFQUFFLEVBQUUsQ0FBQyxDQUFDO1lBQzdELEtBQUs7WUFDTCxJQUFJO1NBQ0wsQ0FBQyxDQUFDLENBQUE7UUFFSCxJQUFJLENBQUMsR0FBRyxDQUFDLGVBQUssQ0FBQyxXQUFXLENBQUMsNkJBQTZCLENBQUMsQ0FBQyxDQUFBO1FBQzFE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b3VyY2VzQ29udGVudCI6WyJpbXBvcnQgeyBDb21tYW5kLCBmbGFncyB9IGZyb20gJ0BvY2xpZi9jb21tYW5kJ1xuaW1wb3J0ICogYXMgUGFyc2VyIGZyb20gJ0BvY2xpZi9wYXJzZXInXG5pbXBvcnQgY2hhbGsgZnJvbSAnY2hhbGsnXG5pbXBvcnQgdXggZnJvbSAnY2xpLXV4J1xuaW1wb3J0ICogYXMgY2xpIGZyb20gJ2lucXVpcmVyJ1xuaW1wb3J0IHtcbiAgZ2V0U29saWRpdHlWZXJzaW9uQnksXG4gIGdldFNvbGlkaXR5VmVyc2lvbnMsXG4gIG1vZGlmeVRydWZmbGVCb3hXaXRoLFxufSBmcm9tICcuLi9zZXJ2aWNlcy90cnVmZmxlLWJveCdcbmV4cG9ydCBkZWZhdWx0IGNsYXNzIEJveCBleHRlbmRzIENvbW1hbmQge1xuICBzdGF0aWMgZGVzY3JpcHRpb24gPSAnTW9kaWZ5IGEgdHJ1ZmZsZSBib3ggdG8gYSBzcGVjaWZpZWQgc29saWRpdHkgdmVyc2lvbidcblxuICBzdGF0aWMgZXhhbXBsZXMgPSBbXG4gICAgJ2JlbHQgYm94IC0tc29sVmVyIDAuNiAtZCBwYXRoL3RvL2JveCcsXG4gICAgJ2JlbHQgYm94IC0taW50ZXJhY3RpdmUgcGF0aC90by9ib3gnLFxuICAgICdiZWx0IGJveCAtbCcsXG4gIF1cblxuICBzdGF0aWMgZmxhZ3MgPSB7XG4gICAgaGVscDogZmxhZ3MuaGVscCh7IGNoYXI6ICdoJyB9KSxcbiAgICBpbnRlcmFjdGl2ZTogZmxhZ3MuYm9vbGVhbih7XG4gICAgICBjaGFyOiAnaScsXG4gICAgICBkZXNjcmlwdGlvbjogJ3J1biB0aGlzIGNvbW1hbmQgaW4gaW50ZXJhY3RpdmUgbW9kZScsXG4gICAgfSksXG4gICAgc29sVmVyOiBmbGFncy5zdHJpbmcoe1xuICAgICAgY2hhcjogJ3MnLFxuICAgICAgZGVzY3JpcHRpb246XG4gICAgICAgICd0aGUgc29saWRpdHkgdmVyc2lvbiB0byBjaGFuZ2UgdGhlIHRydWZmbGUgYm94IHRvXFxuZWl0aGVyIGEgc29saWRpdHkgdmVyc2lvbiBhbGlhcyBcInYwLjZcIiB8IFwiMC42XCIgb3IgaXRzIGZ1bGwgdmVyc2lvbiBcIjAuNi4yXCInLFxuICAgIH0pLFxuICAgIGxpc3Q6IGZsYWdzLmJvb2xlYW4oe1xuICAgICAgY2hhcjogJ2wnLFxuICAgICAgZGVzY3JpcHRpb246ICdsaXN0IHRoZSBhdmFpbGFibGUgc29saWRpdHkgdmVyc2lvbnMnLFxuICAgIH0pLFxuICAgIGRyeVJ1bjogZmxhZ3MuYm9vbGVhbih7XG4gICAgICBjaGFyOiAnZCcsXG4gICAgICBkZXNjcmlwdGlvbjogJ291dHB1dCB0aGUgcmVwbGFjZWQgc3RyaW5ncywgYnV0IGRvbnQgY2hhbmdlIHRoZW0nLFxuICAgIH0pLFxuICB9XG5cbiAgc3RhdGljIGFyZ3M6IFBhcnNlci5hcmdzLklBcmdbXSA9IFtcbiAgICB7XG4gICAgICBuYW1lOiAncGF0aCcsXG4gICAgICBkZXNjcmlwdGlvbjogJ3RoZSBwYXRoIHRvIHRoZSB0cnVmZmxlIGJveCcsXG4gICAgfSxcbiAgXVxuXG4gIGFzeW5jIHJ1bigpIHtcbiAgICBjb25zdCB7IGZsYWdzLCBhcmdzIH0gPSB0aGlzLnBhcnNlKEJveClcbiAgICBpZiAoZmxhZ3MubGlzdCkge1xuICAgICAgcmV0dXJuIHRoaXMuaGFuZGxlTGlzdCgpXG4gICAgfVxuXG4gICAgaWYgKGZsYWdzLmludGVyYWN0aXZlKSB7XG4gICAgICByZXR1cm4gYXdhaXQgdGhpcy5oYW5kbGVJbnRlcmFjdGl2ZShhcmdzLnBhdGgsIGZsYWdzLmRyeVJ1bilcbiAgICB9XG4gICAgaWYgKGZsYWdzLnNvbFZlcikge1xuICAgICAgcmV0dXJuIHRoaXMuaGFuZGxlTm9uSW50ZXJhY3RpdmUoYXJncy5wYXRoLCBmbGFncy5kcnlSdW4sIGZsYWdzLnNvbFZlcilcbiAgICB9XG5cbiAgICB0aGlzLl9oZWxwKClcbiAgfVxuXG4gIC8qKlxuICAgKiBIYW5kbGUgcHJpbnRpbmcgb3V0IGEgbGlzdCBvZiBhdmFpbGFibGUgc29saWRpdHkgdmVyc2lvbnNcbiAgICovXG4gIHByaXZhdGUgaGFuZGxlTGlzdCgpIHtcbiAgICBjb25zdCB2ZXJzaW9ucyA9IGdldFNvbGlkaXR5VmVyc2lvbnMoKS5tYXAoKFthbGlhcywgZnVsbF0pID0+ICh7XG4gICAgICBhbGlhcyxcbiAgICAgIGZ1bGwsXG4gICAgfSkpXG5cbiAgICB0aGlzLmxvZyhjaGFsay5ncmVlbkJyaWdodCgnQXZhaWxhYmxlIFNvbGlkaXR5IFZlcnNpb25zJykpXG4gICAgdXgudGFibGUodmVyc2lvbnMsIHtcbiAgICAgIGFsaWFzOiB7fSxcbiAgICAgIGZ1bGw6IHt9LFxuICAgIH0pXG4gICAgdGhpcy5sb2coJycpXG4gIH1cblxuICAvKipcbiAgICogSGFuZGxlIGludGVyYWN0aXZlIG1vZGUuXG4gICAqIFByb21wdHMgdXNlciBmb3IgYSBzb2xpZGl0eSB2ZXJzaW9uIG51bWJlciB0aGVuIHByb2NlZWRzIHRvXG4gICAqIGRvIGEgZmluZC1yZXBsYWNlIHdpdGhpbiB0aGVpciBib3ggZm9yIHRoZSBzZWxlY3RlZCB2ZXJzaW9uXG4gICAqXG4gICAqIEBwYXJhbSBwYXRoIFRoZSBwYXRoIHRvIHRoZSB0cnVmZmxlIGJveFxuICAgKiBAcGFyYW0gZHJ5UnVuIERvbnQgcmVwbGFjZSB0aGUgZmlsZSBjb250ZW50cywgcHJpbnQgdGhlIGRpZmYgaW5zdGVhZFxuICAgKi9cbiAgcHJpdmF0ZSBhc3luYyBoYW5kbGVJbnRlcmFjdGl2ZShwYXRoOiBzdHJpbmcsIGRyeVJ1bjogYm9vbGVhbikge1xuICAgIGNvbnN0IHNvbGlkaXR5VmVyc2lvbnMgPSBnZXRTb2xpZGl0eVZlcnNpb25zKClcbiAgICBjb25zdCB7IHNvbGNWZXJzaW9uIH0gPSBhd2FpdCBjbGkucHJvbXB0KFtcbiAgICAgIHtcbiAgICAgICAgbmFtZTogJ3NvbGNWZXJzaW9uJyxcbiAgICAgICAgdHlwZTogJ2xpc3QnLFxuICAgICAgICBjaG9pY2VzOiBzb2xpZGl0eVZlcnNpb25zLm1hcCgoWywgdmVyc2lvbl0pID0+IHZlcnNpb24pLFxuICAgICAgICBtZXNzYWdlOlxuICAgICAgICAgICdXaGF0IHZlcnNpb24gb2Ygc29saWRpdHkgZG8geW91IHdhbnQgdG8gdXNlIHdpdGggeW91ciBzbWFydCBjb250cmFjdHM/JyxcbiAgICAgIH0sXG4gICAgXSlcblxuICAgIGNvbnN0IGZ1bGxWZXJzaW9uID0gdGhpcy5nZXRGdWxsVmVyc2lvbihzb2xjVmVyc2lvbilcblxuICAgIG1vZGlmeVRydWZmbGVCb3hXaXRoKGZ1bGxWZXJzaW9uLCBwYXRoLCBkcnlSdW4pXG4gICAgdGhpcy5sb2coXG4gICAgICBjaGFsay5ncmVlbkJyaWdodChcbiAgICAgICAgYERvbmUhXFxuUGxlYXNlIHJ1biBcIm5wbSBpXCIgdG8gaW5zdGFsbCB0aGUgbmV3IGNoYW5nZXMgbWFkZS5gLFxuICAgICAgKSxcbiAgICApXG4gIH1cblxuICAvKipcbiAgICogSGFuZGxlIG5vbi1pbnRlcmFjdGl2ZSBtb2RlIFwiLS1zb2xWZXJcIi5cbiAgICogc29saWRpdHkgdmVyc2lvbiBudW1iZXIgdGhlbiBwcm9jZWVkcyB0b1xuICAgKiBkbyBhIGZpbmQtcmVwbGFjZSB3aXRoaW4gdGhlaXIgYm94IGZvciB0aGUgc2VsZWN0ZWQgdmVyc2lvblxuICAgKlxuICAgKiBAcGFyYW0gcGF0aCBUaGUgcGF0aCB0byB0aGUgdHJ1ZmZsZSBib3hcbiAgICogQHBhcmFtIGRyeVJ1biBEb250IHJlcGxhY2UgdGhlIGZpbGUgY29udGVudHMsIHByaW50IHRoZSBkaWZmIGluc3RlYWRcbiAgICogQHBhcmFtIHZlcnNpb25BbGlhc09yVmVyc2lvbiBFaXRoZXIgYSBzb2xpZGl0eSB2ZXJzaW9uIGFsaWFzIFwidjAuNlwiIHwgXCIwLjZcIiBvciBpdHMgZnVsbCB2ZXJzaW9uIFwiMC42LjJcIlxuICAgKi9cbiAgcHJpdmF0ZSBoYW5kbGVOb25JbnRlcmFjdGl2ZShcbiAgICBwYXRoOiBzdHJpbmcsXG4gICAgZHJ5UnVuOiBib29sZWFuLFxuICAgIHZlcnNpb25BbGlhc09yVmVyc2lvbjogc3RyaW5nLFxuICApIHtcbiAgICBjb25zdCBmdWxsVmVyc2lvbiA9IHRoaXMuZ2V0RnVsbFZlcnNpb24odmVyc2lvbkFsaWFzT3JWZXJzaW9uKVxuXG4gICAgbW9kaWZ5VHJ1ZmZsZUJveFdpdGgoZnVsbFZlcnNpb24sIHBhdGgsIGRyeVJ1bilcbiAgfVxuXG4gIHByaXZhdGUgZ2V0RnVsbFZlcnNpb24odmVyc2lvbkFsaWFzT3JWZXJzaW9uOiBzdHJpbmcpIHtcbiAgICBsZXQgZnVsbFZlcnNpb246IFJldHVyblR5cGU8dHlwZW9mIGdldFNvbGlkaXR5VmVyc2lvbkJ5PlxuICAgIHRyeSB7XG4gICAgICBmdWxsVmVyc2lvbiA9IGdldFNvbGlkaXR5VmVyc2lvbkJ5KHZlcnNpb25BbGlhc09yVmVyc2lvbilcbiAgICB9IGNhdGNoIHtcbiAgICAgIGNvbnN0IGVycm9yID0gY2hhbGsucmVkKCdDb3VsZCBub3QgZmluZCBnaXZlbiBzb2xpZGl0eSB2ZXJzaW9uXFxuJylcbiAgICAgIHRoaXMubG9nKGVycm9yKVxuICAgICAgdGhpcy5oYW5kbGVMaXN0KClcbiAgICAgIHRoaXMuZXhpdCgxKVxuICAgIH1cblxuICAgIHJldHVybiBmdWxsVmVyc2lvblxuICB9XG59XG4iXX0=