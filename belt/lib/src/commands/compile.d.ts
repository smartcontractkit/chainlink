import { Command, flags } from '@oclif/command';
import * as Parser from '@oclif/parser';
export default class Compile extends Command {
    static description: string;
    static examples: string[];
    static flags: {
        help: Parser.flags.IBooleanFlag<void>;
        config: flags.IOptionFlag<string>;
    };
    static args: Parser.args.IArg[];
    run(): Promise<void>;
    private compileSingle;
    private compileAll;
}
//# sourceMappingURL=compile.d.ts.map