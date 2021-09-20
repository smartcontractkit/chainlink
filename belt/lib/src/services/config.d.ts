/**
 * Structure of the application configuration, paths are relative to the current working directory.
 *
 * Uses these configuration values for:
 * - sol-compiler
 * - codegenning ethers contract abstractions
 * - codegenning truffle contract abstractions
 * - running solhint
 */
export interface App {
    /**
     * The directory where all of the solidity smart contracts are held
     */
    contractsDir: string;
    /**
     * The directory where all of the smart contract artifacts should be outputted
     */
    artifactsDir: string;
    /**
     * The directory where all contract abstractions should be outputted
     */
    contractAbstractionDir: string;
    /**
     * Instruct sol-compiler to use a dockerized solc instance for higher performance,
     * or to use solcjs
     */
    useDockerisedSolc: boolean;
    /**
     * Various compiler settings for sol-compiler
     */
    compilerSettings: {
        /**
         * A mapping of directories to their solidity compiler versions that should be used.
         *
         * e.g.
         *  Given the following directory structure:
         * ```sh
         *  src
         *  ├── v0.4
         *  └── v0.5
         *  ```
         * Our versions dict would look like the following:
         * ```json
         * {
         *  "v0.4": "0.4.24",
         *  "v0.5": "0.5.0"
         * }
         * ```
         */
        versions: {
            [dir: string]: string;
        };
    };
    /**
     * Versions to publically show in our truffle box options
     */
    publicVersions: string[];
}
/**
 * Load a validated configuration file in JSON format for app configuration purposes.
 *
 * @param path The path relative to the current working directory to load the configuration file from.
 */
export declare function load(path: string): App;
//# sourceMappingURL=config.d.ts.map