import d from 'debug';
import * as config from './config';
/**
 * Get contract versions and their directories
 */
export declare function getContractDirs(conf: config.App): {
    dir: string;
    version: string;
}[];
/**
 * Get artifact verions and their directories
 */
export declare function getArtifactDirs(conf: config.App): {
    dir: string;
    version: string;
}[];
/**
 * Create a logger specifically for debugging. The root level namespace is based on the package name.
 *
 * @see https://www.npmjs.com/package/debug
 * @param fileName The filename that the debug logger is being used in for namespacing purposes.
 */
export declare function debug(fileName: string): d.Debugger;
/**
 * Load a json file at the specified path.
 *
 * @param path The file path relative to the cwd to read in the json file from.
 */
export declare function getJsonFile(path: string): unknown;
//# sourceMappingURL=utils.d.ts.map