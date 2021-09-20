import { Command, flags } from '@oclif/command';
import * as Parser from '@oclif/parser';
export default class Box extends Command {
    static description: string;
    static examples: string[];
    static flags: {
        help: Parser.flags.IBooleanFlag<void>;
        interactive: Parser.flags.IBooleanFlag<boolean>;
        solVer: flags.IOptionFlag<string | undefined>;
        list: Parser.flags.IBooleanFlag<boolean>;
        dryRun: Parser.flags.IBooleanFlag<boolean>;
    };
    static args: Parser.args.IArg[];
    run(): Promise<void>;
    /**
     * Handle printing out a list of available solidity versions
     */
    private handleList;
    /**
     * Handle interactive mode.
     * Prompts user for a solidity version number then proceeds to
     * do a find-replace within their box for the selected version
     *
     * @param path The path to the truffle box
     * @param dryRun Dont replace the file contents, print the diff instead
     */
    private handleInteractive;
    /**
     * Handle non-interactive mode "--solVer".
     * solidity version number then proceeds to
     * do a find-replace within their box for the selected version
     *
     * @param path The path to the truffle box
     * @param dryRun Dont replace the file contents, print the diff instead
     * @param versionAliasOrVersion Either a solidity version alias "v0.6" | "0.6" or its full version "0.6.2"
     */
    private handleNonInteractive;
    private getFullVersion;
}
//# sourceMappingURL=box.d.ts.map