/**
 * Modify a truffle box with the given solidity version
 *
 * @param solidityVersion A tuple of alias and version of a solidity version, e.g ['v0.4', '0.4.24']
 * @param path The path to the truffle box
 * @param dryRun Whether to actually modify the files in-place or to print the modified files to stdout
 */
export declare function modifyTruffleBoxWith([solcVersionAlias, solcVersion]: [string, string], path: string, dryRun: boolean): void;
/**
 * Get a solidity version by its alias or version number
 *
 * @param versionAliasOrVersion Either a solidity version alias "v0.6" | "0.6" or its full version "0.6.2"
 * @throws error if version given could not be found
 */
export declare function getSolidityVersionBy(versionAliasOrVersion: string): [string, string];
/**
 * Get a list of available solidity versions based on what's published in @chainlink/contracts
 *
 * The returned format is [alias, version] where alias can be "v0.6" | "0.6" and full version can be "0.6.2"
 */
export declare function getSolidityVersions(): [string, string][];
/**
 * Get the path to the truffle config
 *
 * @param basePath The path to the truffle box
 */
export declare function getTruffleConfig(basePath: string): string;
/**
 * Get a list of all javascript files in the truffle box
 *
 * @param basePath The path to the truffle box
 */
export declare function getJavascriptFiles(basePath: string): string[];
/**
 * Get path to the package.json
 *
 * @param basePath The path to the truffle box
 */
export declare function getPackageJson(basePath: string): string;
/**
 * Get a list of all solidity files in the truffle box
 *
 * @param basePath The path to the truffle box
 */
export declare function getSolidityFiles(basePath: string): string[];
//# sourceMappingURL=truffle-box.d.ts.map