"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.compileAll = void 0;
const sol_compiler_1 = require("@0x/sol-compiler");
const path_1 = require("path");
const utils_1 = require("../utils");
const d = utils_1.debug('solc');
/**
 * Generate solidity artifacts for all of the solidity versions under a specified contract
 * directory.
 *
 * @param conf The application configuration, e.g. where to read solidity files, where to output, etc..
 */
async function compileAll(conf) {
    return Promise.all(utils_1.getContractDirs(conf).map(async ({ dir, version }) => {
        const c = compiler(conf, dir, version);
        // Compiler#getCompilerOutputsAsync throws on compilation errors
        // this method prints any errors and warnings for us
        await c.compileAsync();
    }));
}
exports.compileAll = compileAll;
/**
 * Create a sol-compiler instance that reads in a subdirectory of smart contracts e.g. (src/v0.4, src/v0.5, ..)
 * and outputs their respective compiler artifacts e.g. (abi/v0.4, abi/v0.5)
 *
 * @param config The application specific configuration to use for sol-compiler
 * @param subDir The subdirectory to use as a namespace when reading .sol files and outputting
 * their respective artifacts
 * @param solcVersion The solidity compiler version to use with sol-compiler
 */
function compiler({ artifactsDir, useDockerisedSolc, contractsDir, compilerSettings, }, subDir, solcVersion) {
    const _d = d.extend('compiler');
    // remove our custom versions property
    const compilerSettingCopy = JSON.parse(JSON.stringify(compilerSettings));
    // @ts-expect-error
    delete compilerSettingCopy.versions;
    const settings = {
        artifactsDir: path_1.join(artifactsDir, subDir),
        compilerSettings: {
            outputSelection: {
                '*': {
                    '*': [
                        'abi',
                        'devdoc',
                        'userdoc',
                        'evm.bytecode.object',
                        'evm.bytecode.sourceMap',
                        'evm.deployedBytecode.object',
                        'evm.deployedBytecode.sourceMap',
                        'evm.methodIdentifiers',
                        'metadata',
                    ],
                },
            },
            ...compilerSettingCopy,
        },
        contracts: '*',
        contractsDir: path_1.join(contractsDir, subDir),
        isOfflineMode: false,
        shouldSaveStandardInput: false,
        solcVersion,
        useDockerisedSolc,
    };
    _d('Settings: %o', settings);
    return new sol_compiler_1.Compiler(settings);
}
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoic29sYy5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uLy4uL3NyYy9zZXJ2aWNlcy9jb21waWxlcnMvc29sYy50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7QUFBQSxtREFBNEQ7QUFDNUQsK0JBQTJCO0FBRTNCLG9DQUFpRDtBQUNqRCxNQUFNLENBQUMsR0FBRyxhQUFLLENBQUMsTUFBTSxDQUFDLENBQUE7QUFFdkI7Ozs7O0dBS0c7QUFDSSxLQUFLLFVBQVUsVUFBVSxDQUFDLElBQWdCO0lBQy9DLE9BQU8sT0FBTyxDQUFDLEdBQUcsQ0FDaEIsdUJBQWUsQ0FBQyxJQUFJLENBQUMsQ0FBQyxHQUFHLENBQUMsS0FBSyxFQUFFLEVBQUUsR0FBRyxFQUFFLE9BQU8sRUFBRSxFQUFFLEVBQUU7UUFDbkQsTUFBTSxDQUFDLEdBQUcsUUFBUSxDQUFDLElBQUksRUFBRSxHQUFHLEVBQUUsT0FBTyxDQUFDLENBQUE7UUFFdEMsZ0VBQWdFO1FBQ2hFLG9EQUFvRDtRQUNwRCxNQUFNLENBQUMsQ0FBQyxZQUFZLEVBQUUsQ0FBQTtJQUN4QixDQUFDLENBQUMsQ0FDSCxDQUFBO0FBQ0gsQ0FBQztBQVZELGdDQVVDO0FBRUQ7Ozs7Ozs7O0dBUUc7QUFDSCxTQUFTLFFBQVEsQ0FDZixFQUNFLFlBQVksRUFDWixpQkFBaUIsRUFDakIsWUFBWSxFQUNaLGdCQUFnQixHQUNMLEVBQ2IsTUFBYyxFQUNkLFdBQW1CO0lBRW5CLE1BQU0sRUFBRSxHQUFHLENBQUMsQ0FBQyxNQUFNLENBQUMsVUFBVSxDQUFDLENBQUE7SUFDL0Isc0NBQXNDO0lBQ3RDLE1BQU0sbUJBQW1CLEdBQW1DLElBQUksQ0FBQyxLQUFLLENBQ3BFLElBQUksQ0FBQyxTQUFTLENBQUMsZ0JBQWdCLENBQUMsQ0FDakMsQ0FBQTtJQUNELG1CQUFtQjtJQUNuQixPQUFPLG1CQUFtQixDQUFDLFFBQVEsQ0FBQTtJQUVuQyxNQUFNLFFBQVEsR0FBb0I7UUFDaEMsWUFBWSxFQUFFLFdBQUksQ0FBQyxZQUFZLEVBQUUsTUFBTSxDQUFDO1FBQ3hDLGdCQUFnQixFQUFFO1lBQ2hCLGVBQWUsRUFBRTtnQkFDZixHQUFHLEVBQUU7b0JBQ0gsR0FBRyxFQUFFO3dCQUNILEtBQUs7d0JBQ0wsUUFBUTt3QkFDUixTQUFTO3dCQUNULHFCQUFxQjt3QkFDckIsd0JBQXdCO3dCQUN4Qiw2QkFBNkI7d0JBQzdCLGdDQUFnQzt3QkFDaEMsdUJBQXVCO3dCQUN2QixVQUFVO3FCQUNYO2lCQUNGO2FBQ0Y7WUFDRCxHQUFHLG1CQUFtQjtTQUN2QjtRQUNELFNBQVMsRUFBRSxHQUFHO1FBQ2QsWUFBWSxFQUFFLFdBQUksQ0FBQyxZQUFZLEVBQUUsTUFBTSxDQUFDO1FBQ3hDLGFBQWEsRUFBRSxLQUFLO1FBQ3BCLHVCQUF1QixFQUFFLEtBQUs7UUFDOUIsV0FBVztRQUNYLGlCQUFpQjtLQUNsQixDQUFBO0lBQ0QsRUFBRSxDQUFDLGNBQWMsRUFBRSxRQUFRLENBQUMsQ0FBQTtJQUU1QixPQUFPLElBQUksdUJBQVEsQ0FBQyxRQUFRLENBQUMsQ0FBQTtBQUMvQixDQUFDIiwic291cmNlc0NvbnRlbnQiOlsiaW1wb3J0IHsgQ29tcGlsZXIsIENvbXBpbGVyT3B0aW9ucyB9IGZyb20gJ0AweC9zb2wtY29tcGlsZXInXG5pbXBvcnQgeyBqb2luIH0gZnJvbSAncGF0aCdcbmltcG9ydCAqIGFzIGNvbmZpZyBmcm9tICcuLi9jb25maWcnXG5pbXBvcnQgeyBkZWJ1ZywgZ2V0Q29udHJhY3REaXJzIH0gZnJvbSAnLi4vdXRpbHMnXG5jb25zdCBkID0gZGVidWcoJ3NvbGMnKVxuXG4vKipcbiAqIEdlbmVyYXRlIHNvbGlkaXR5IGFydGlmYWN0cyBmb3IgYWxsIG9mIHRoZSBzb2xpZGl0eSB2ZXJzaW9ucyB1bmRlciBhIHNwZWNpZmllZCBjb250cmFjdFxuICogZGlyZWN0b3J5LlxuICpcbiAqIEBwYXJhbSBjb25mIFRoZSBhcHBsaWNhdGlvbiBjb25maWd1cmF0aW9uLCBlLmcuIHdoZXJlIHRvIHJlYWQgc29saWRpdHkgZmlsZXMsIHdoZXJlIHRvIG91dHB1dCwgZXRjLi5cbiAqL1xuZXhwb3J0IGFzeW5jIGZ1bmN0aW9uIGNvbXBpbGVBbGwoY29uZjogY29uZmlnLkFwcCkge1xuICByZXR1cm4gUHJvbWlzZS5hbGwoXG4gICAgZ2V0Q29udHJhY3REaXJzKGNvbmYpLm1hcChhc3luYyAoeyBkaXIsIHZlcnNpb24gfSkgPT4ge1xuICAgICAgY29uc3QgYyA9IGNvbXBpbGVyKGNvbmYsIGRpciwgdmVyc2lvbilcblxuICAgICAgLy8gQ29tcGlsZXIjZ2V0Q29tcGlsZXJPdXRwdXRzQXN5bmMgdGhyb3dzIG9uIGNvbXBpbGF0aW9uIGVycm9yc1xuICAgICAgLy8gdGhpcyBtZXRob2QgcHJpbnRzIGFueSBlcnJvcnMgYW5kIHdhcm5pbmdzIGZvciB1c1xuICAgICAgYXdhaXQgYy5jb21waWxlQXN5bmMoKVxuICAgIH0pLFxuICApXG59XG5cbi8qKlxuICogQ3JlYXRlIGEgc29sLWNvbXBpbGVyIGluc3RhbmNlIHRoYXQgcmVhZHMgaW4gYSBzdWJkaXJlY3Rvcnkgb2Ygc21hcnQgY29udHJhY3RzIGUuZy4gKHNyYy92MC40LCBzcmMvdjAuNSwgLi4pXG4gKiBhbmQgb3V0cHV0cyB0aGVpciByZXNwZWN0aXZlIGNvbXBpbGVyIGFydGlmYWN0cyBlLmcuIChhYmkvdjAuNCwgYWJpL3YwLjUpXG4gKlxuICogQHBhcmFtIGNvbmZpZyBUaGUgYXBwbGljYXRpb24gc3BlY2lmaWMgY29uZmlndXJhdGlvbiB0byB1c2UgZm9yIHNvbC1jb21waWxlclxuICogQHBhcmFtIHN1YkRpciBUaGUgc3ViZGlyZWN0b3J5IHRvIHVzZSBhcyBhIG5hbWVzcGFjZSB3aGVuIHJlYWRpbmcgLnNvbCBmaWxlcyBhbmQgb3V0cHV0dGluZ1xuICogdGhlaXIgcmVzcGVjdGl2ZSBhcnRpZmFjdHNcbiAqIEBwYXJhbSBzb2xjVmVyc2lvbiBUaGUgc29saWRpdHkgY29tcGlsZXIgdmVyc2lvbiB0byB1c2Ugd2l0aCBzb2wtY29tcGlsZXJcbiAqL1xuZnVuY3Rpb24gY29tcGlsZXIoXG4gIHtcbiAgICBhcnRpZmFjdHNEaXIsXG4gICAgdXNlRG9ja2VyaXNlZFNvbGMsXG4gICAgY29udHJhY3RzRGlyLFxuICAgIGNvbXBpbGVyU2V0dGluZ3MsXG4gIH06IGNvbmZpZy5BcHAsXG4gIHN1YkRpcjogc3RyaW5nLFxuICBzb2xjVmVyc2lvbjogc3RyaW5nLFxuKSB7XG4gIGNvbnN0IF9kID0gZC5leHRlbmQoJ2NvbXBpbGVyJylcbiAgLy8gcmVtb3ZlIG91ciBjdXN0b20gdmVyc2lvbnMgcHJvcGVydHlcbiAgY29uc3QgY29tcGlsZXJTZXR0aW5nQ29weTogY29uZmlnLkFwcFsnY29tcGlsZXJTZXR0aW5ncyddID0gSlNPTi5wYXJzZShcbiAgICBKU09OLnN0cmluZ2lmeShjb21waWxlclNldHRpbmdzKSxcbiAgKVxuICAvLyBAdHMtZXhwZWN0LWVycm9yXG4gIGRlbGV0ZSBjb21waWxlclNldHRpbmdDb3B5LnZlcnNpb25zXG5cbiAgY29uc3Qgc2V0dGluZ3M6IENvbXBpbGVyT3B0aW9ucyA9IHtcbiAgICBhcnRpZmFjdHNEaXI6IGpvaW4oYXJ0aWZhY3RzRGlyLCBzdWJEaXIpLFxuICAgIGNvbXBpbGVyU2V0dGluZ3M6IHtcbiAgICAgIG91dHB1dFNlbGVjdGlvbjoge1xuICAgICAgICAnKic6IHtcbiAgICAgICAgICAnKic6IFtcbiAgICAgICAgICAgICdhYmknLFxuICAgICAgICAgICAgJ2RldmRvYycsXG4gICAgICAgICAgICAndXNlcmRvYycsXG4gICAgICAgICAgICAnZXZtLmJ5dGVjb2RlLm9iamVjdCcsXG4gICAgICAgICAgICAnZXZtLmJ5dGVjb2RlLnNvdXJjZU1hcCcsXG4gICAgICAgICAgICAnZXZtLmRlcGxveWVkQnl0ZWNvZGUub2JqZWN0JyxcbiAgICAgICAgICAgICdldm0uZGVwbG95ZWRCeXRlY29kZS5zb3VyY2VNYXAnLFxuICAgICAgICAgICAgJ2V2bS5tZXRob2RJZGVudGlmaWVycycsXG4gICAgICAgICAgICAnbWV0YWRhdGEnLFxuICAgICAgICAgIF0sXG4gICAgICAgIH0sXG4gICAgICB9LFxuICAgICAgLi4uY29tcGlsZXJTZXR0aW5nQ29weSxcbiAgICB9LFxuICAgIGNvbnRyYWN0czogJyonLFxuICAgIGNvbnRyYWN0c0Rpcjogam9pbihjb250cmFjdHNEaXIsIHN1YkRpciksXG4gICAgaXNPZmZsaW5lTW9kZTogZmFsc2UsXG4gICAgc2hvdWxkU2F2ZVN0YW5kYXJkSW5wdXQ6IGZhbHNlLFxuICAgIHNvbGNWZXJzaW9uLFxuICAgIHVzZURvY2tlcmlzZWRTb2xjLFxuICB9XG4gIF9kKCdTZXR0aW5nczogJW8nLCBzZXR0aW5ncylcblxuICByZXR1cm4gbmV3IENvbXBpbGVyKHNldHRpbmdzKVxufVxuIl19