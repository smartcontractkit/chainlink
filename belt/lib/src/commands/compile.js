"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    Object.defineProperty(o, k2, { enumerable: true, get: function() { return m[k]; } });
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (k !== "default" && Object.prototype.hasOwnProperty.call(mod, k)) __createBinding(result, mod, k);
    __setModuleDefault(result, mod);
    return result;
};
Object.defineProperty(exports, "__esModule", { value: true });
const command_1 = require("@oclif/command");
class Compile extends command_1.Command {
    async run() {
        const { flags, args, argv } = this.parse(Compile);
        if (argv.length === 0) {
            this._help();
        }
        try {
            const config = await Promise.resolve().then(() => __importStar(require('../services/config')));
            const conf = config.load(flags.config);
            const compilation = args.compiler === 'all'
                ? this.compileAll(conf)
                : this.compileSingle(args.compiler, conf);
            await compilation;
        }
        catch (e) {
            this.error(e);
        }
    }
    async compileSingle(compilerName, conf) {
        const compiler = await Promise.resolve().then(() => __importStar(require(`../services/compilers/${compilerName}`)));
        await compiler.compileAll(conf);
    }
    async compileAll(conf) {
        const compilers = await Promise.resolve().then(() => __importStar(require('../services/compilers')));
        await compilers.solc.compileAll(conf);
        await Promise.all([
            compilers.truffle.compileAll(conf),
            compilers.ethers.compileAll(conf),
        ]);
    }
}
exports.default = Compile;
Compile.description = 'Run various compilers and/or codegenners that target solidity smart contracts.';
Compile.examples = [
    `$ belt compile all

Creating directory at abi/v0.4...
Creating directory at abi/v0.5...
Creating directory at abi/v0.6...
Compiling 35 contracts...
...
...
Aggregator artifact saved!
AggregatorProxy artifact saved!
Chainlink artifact saved!
...`,
];
Compile.flags = {
    help: command_1.flags.help({ char: 'h' }),
    config: command_1.flags.string({
        char: 'c',
        default: 'app.config.json',
        description: 'Location of the configuration file',
    }),
};
Compile.args = [
    {
        name: 'compiler',
        description: 'Compile solidity smart contracts and output their artifacts',
        options: ['solc', 'ethers', 'truffle', 'all'],
    },
];
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiY29tcGlsZS5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uL3NyYy9jb21tYW5kcy9jb21waWxlLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7OztBQUFBLDRDQUErQztBQUkvQyxNQUFxQixPQUFRLFNBQVEsaUJBQU87SUFxQzFDLEtBQUssQ0FBQyxHQUFHO1FBQ1AsTUFBTSxFQUFFLEtBQUssRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLEdBQUcsSUFBSSxDQUFDLEtBQUssQ0FBQyxPQUFPLENBQUMsQ0FBQTtRQUNqRCxJQUFJLElBQUksQ0FBQyxNQUFNLEtBQUssQ0FBQyxFQUFFO1lBQ3JCLElBQUksQ0FBQyxLQUFLLEVBQUUsQ0FBQTtTQUNiO1FBRUQsSUFBSTtZQUNGLE1BQU0sTUFBTSxHQUFHLHdEQUFhLG9CQUFvQixHQUFDLENBQUE7WUFDakQsTUFBTSxJQUFJLEdBQUcsTUFBTSxDQUFDLElBQUksQ0FBQyxLQUFLLENBQUMsTUFBTSxDQUFDLENBQUE7WUFDdEMsTUFBTSxXQUFXLEdBQ2YsSUFBSSxDQUFDLFFBQVEsS0FBSyxLQUFLO2dCQUNyQixDQUFDLENBQUMsSUFBSSxDQUFDLFVBQVUsQ0FBQyxJQUFJLENBQUM7Z0JBQ3ZCLENBQUMsQ0FBQyxJQUFJLENBQUMsYUFBYSxDQUFDLElBQUksQ0FBQyxRQUFRLEVBQUUsSUFBSSxDQUFDLENBQUE7WUFFN0MsTUFBTSxXQUFXLENBQUE7U0FDbEI7UUFBQyxPQUFPLENBQUMsRUFBRTtZQUNWLElBQUksQ0FBQyxLQUFLLENBQUMsQ0FBQyxDQUFDLENBQUE7U0FDZDtJQUNILENBQUM7SUFFTyxLQUFLLENBQUMsYUFBYSxDQUFDLFlBQW9CLEVBQUUsSUFBUztRQUV6RCxNQUFNLFFBQVEsR0FBK0Isd0RBQzNDLHlCQUF5QixZQUFZLEVBQUUsR0FDeEMsQ0FBQTtRQUVELE1BQU0sUUFBUSxDQUFDLFVBQVUsQ0FBQyxJQUFJLENBQUMsQ0FBQTtJQUNqQyxDQUFDO0lBRU8sS0FBSyxDQUFDLFVBQVUsQ0FBQyxJQUFTO1FBQ2hDLE1BQU0sU0FBUyxHQUFHLHdEQUFhLHVCQUF1QixHQUFDLENBQUE7UUFFdkQsTUFBTSxTQUFTLENBQUMsSUFBSSxDQUFDLFVBQVUsQ0FBQyxJQUFJLENBQUMsQ0FBQTtRQUNyQyxNQUFNLE9BQU8sQ0FBQyxHQUFHLENBQUM7WUFDaEIsU0FBUyxDQUFDLE9BQU8sQ0FBQyxVQUFVLENBQUMsSUFBSSxDQUFDO1lBQ2xDLFNBQVMsQ0FBQyxNQUFNLENBQUMsVUFBVSxDQUFDLElBQUksQ0FBQztTQUNsQyxDQUFDLENBQUE7SUFDSixDQUFDOztBQTFFSCwwQkEyRUM7QUExRVEsbUJBQVcsR0FDaEIsZ0ZBQWdGLENBQUE7QUFFM0UsZ0JBQVEsR0FBRztJQUNoQjs7Ozs7Ozs7Ozs7SUFXQTtDQUNELENBQUE7QUFFTSxhQUFLLEdBQUc7SUFDYixJQUFJLEVBQUUsZUFBSyxDQUFDLElBQUksQ0FBQyxFQUFFLElBQUksRUFBRSxHQUFHLEVBQUUsQ0FBQztJQUMvQixNQUFNLEVBQUUsZUFBSyxDQUFDLE1BQU0sQ0FBQztRQUNuQixJQUFJLEVBQUUsR0FBRztRQUNULE9BQU8sRUFBRSxpQkFBaUI7UUFDMUIsV0FBVyxFQUFFLG9DQUFvQztLQUNsRCxDQUFDO0NBQ0gsQ0FBQTtBQUVNLFlBQUksR0FBdUI7SUFDaEM7UUFDRSxJQUFJLEVBQUUsVUFBVTtRQUNoQixXQUFXLEVBQ1QsNkRBQTZEO1FBQy9ELE9BQU8sRUFBRSxDQUFDLE1BQU0sRUFBRSxRQUFRLEVBQUUsU0FBUyxFQUFFLEtBQUssQ0FBQztLQUM5QztDQUNGLENBQUEiLCJzb3VyY2VzQ29udGVudCI6WyJpbXBvcnQgeyBDb21tYW5kLCBmbGFncyB9IGZyb20gJ0BvY2xpZi9jb21tYW5kJ1xuaW1wb3J0ICogYXMgUGFyc2VyIGZyb20gJ0BvY2xpZi9wYXJzZXInXG5pbXBvcnQgeyBBcHAgfSBmcm9tICcuLi9zZXJ2aWNlcy9jb25maWcnXG5cbmV4cG9ydCBkZWZhdWx0IGNsYXNzIENvbXBpbGUgZXh0ZW5kcyBDb21tYW5kIHtcbiAgc3RhdGljIGRlc2NyaXB0aW9uID1cbiAgICAnUnVuIHZhcmlvdXMgY29tcGlsZXJzIGFuZC9vciBjb2RlZ2VubmVycyB0aGF0IHRhcmdldCBzb2xpZGl0eSBzbWFydCBjb250cmFjdHMuJ1xuXG4gIHN0YXRpYyBleGFtcGxlcyA9IFtcbiAgICBgJCBiZWx0IGNvbXBpbGUgYWxsXG5cbkNyZWF0aW5nIGRpcmVjdG9yeSBhdCBhYmkvdjAuNC4uLlxuQ3JlYXRpbmcgZGlyZWN0b3J5IGF0IGFiaS92MC41Li4uXG5DcmVhdGluZyBkaXJlY3RvcnkgYXQgYWJpL3YwLjYuLi5cbkNvbXBpbGluZyAzNSBjb250cmFjdHMuLi5cbi4uLlxuLi4uXG5BZ2dyZWdhdG9yIGFydGlmYWN0IHNhdmVkIVxuQWdncmVnYXRvclByb3h5IGFydGlmYWN0IHNhdmVkIVxuQ2hhaW5saW5rIGFydGlmYWN0IHNhdmVkIVxuLi4uYCxcbiAgXVxuXG4gIHN0YXRpYyBmbGFncyA9IHtcbiAgICBoZWxwOiBmbGFncy5oZWxwKHsgY2hhcjogJ2gnIH0pLFxuICAgIGNvbmZpZzogZmxhZ3Muc3RyaW5nKHtcbiAgICAgIGNoYXI6ICdjJyxcbiAgICAgIGRlZmF1bHQ6ICdhcHAuY29uZmlnLmpzb24nLFxuICAgICAgZGVzY3JpcHRpb246ICdMb2NhdGlvbiBvZiB0aGUgY29uZmlndXJhdGlvbiBmaWxlJyxcbiAgICB9KSxcbiAgfVxuXG4gIHN0YXRpYyBhcmdzOiBQYXJzZXIuYXJncy5JQXJnW10gPSBbXG4gICAge1xuICAgICAgbmFtZTogJ2NvbXBpbGVyJyxcbiAgICAgIGRlc2NyaXB0aW9uOlxuICAgICAgICAnQ29tcGlsZSBzb2xpZGl0eSBzbWFydCBjb250cmFjdHMgYW5kIG91dHB1dCB0aGVpciBhcnRpZmFjdHMnLFxuICAgICAgb3B0aW9uczogWydzb2xjJywgJ2V0aGVycycsICd0cnVmZmxlJywgJ2FsbCddLFxuICAgIH0sXG4gIF1cblxuICBhc3luYyBydW4oKSB7XG4gICAgY29uc3QgeyBmbGFncywgYXJncywgYXJndiB9ID0gdGhpcy5wYXJzZShDb21waWxlKVxuICAgIGlmIChhcmd2Lmxlbmd0aCA9PT0gMCkge1xuICAgICAgdGhpcy5faGVscCgpXG4gICAgfVxuXG4gICAgdHJ5IHtcbiAgICAgIGNvbnN0IGNvbmZpZyA9IGF3YWl0IGltcG9ydCgnLi4vc2VydmljZXMvY29uZmlnJylcbiAgICAgIGNvbnN0IGNvbmYgPSBjb25maWcubG9hZChmbGFncy5jb25maWcpXG4gICAgICBjb25zdCBjb21waWxhdGlvbiA9XG4gICAgICAgIGFyZ3MuY29tcGlsZXIgPT09ICdhbGwnXG4gICAgICAgICAgPyB0aGlzLmNvbXBpbGVBbGwoY29uZilcbiAgICAgICAgICA6IHRoaXMuY29tcGlsZVNpbmdsZShhcmdzLmNvbXBpbGVyLCBjb25mKVxuXG4gICAgICBhd2FpdCBjb21waWxhdGlvblxuICAgIH0gY2F0Y2ggKGUpIHtcbiAgICAgIHRoaXMuZXJyb3IoZSlcbiAgICB9XG4gIH1cblxuICBwcml2YXRlIGFzeW5jIGNvbXBpbGVTaW5nbGUoY29tcGlsZXJOYW1lOiBzdHJpbmcsIGNvbmY6IEFwcCkge1xuICAgIHR5cGUgQ29tcGlsZXJzID0gdHlwZW9mIGltcG9ydCgnLi4vc2VydmljZXMvY29tcGlsZXJzJylcbiAgICBjb25zdCBjb21waWxlcjogQ29tcGlsZXJzW2tleW9mIENvbXBpbGVyc10gPSBhd2FpdCBpbXBvcnQoXG4gICAgICBgLi4vc2VydmljZXMvY29tcGlsZXJzLyR7Y29tcGlsZXJOYW1lfWBcbiAgICApXG5cbiAgICBhd2FpdCBjb21waWxlci5jb21waWxlQWxsKGNvbmYpXG4gIH1cblxuICBwcml2YXRlIGFzeW5jIGNvbXBpbGVBbGwoY29uZjogQXBwKSB7XG4gICAgY29uc3QgY29tcGlsZXJzID0gYXdhaXQgaW1wb3J0KCcuLi9zZXJ2aWNlcy9jb21waWxlcnMnKVxuXG4gICAgYXdhaXQgY29tcGlsZXJzLnNvbGMuY29tcGlsZUFsbChjb25mKVxuICAgIGF3YWl0IFByb21pc2UuYWxsKFtcbiAgICAgIGNvbXBpbGVycy50cnVmZmxlLmNvbXBpbGVBbGwoY29uZiksXG4gICAgICBjb21waWxlcnMuZXRoZXJzLmNvbXBpbGVBbGwoY29uZiksXG4gICAgXSlcbiAgfVxufVxuIl19