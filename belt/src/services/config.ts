import { test } from 'shelljs'
import { getJsonFile } from './utils'

/**
 * Rudimentary type helper to recursively mark string | boolean values as unknown.
 *
 * Made specifically for validation of AppConfig
 */
type DeepUnknown<T> = {
  [K in keyof T]: T[K] extends string | boolean ? unknown : DeepUnknown<T[K]>
}

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
  contractsDir: string
  /**
   * The directory where all of the smart contract artifacts should be outputted
   */
  artifactsDir: string
  /**
   * The directory where all contract abstractions should be outputted
   */
  contractAbstractionDir: string
  /**
   * Instruct sol-compiler to use a dockerized solc instance for higher performance,
   * or to use solcjs
   */
  useDockerisedSolc: boolean
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
      [dir: string]: string
    }
  }
  /**
   * Versions to publically show in our truffle box options
   */
  publicVersions: string[]
}

/**
 * Load a validated configuration file in JSON format for app configuration purposes.
 *
 * @param path The path relative to the current working directory to load the configuration file from.
 */
export function load(path: string): App {
  let json: unknown
  try {
    json = getJsonFile(path)
  } catch (e) {
    throw Error(`Could not load config at "${path}".\n\n${e}`)
  }
  assertAppConfig(json)

  return json
}

function assertAppConfig(json: unknown): asserts json is App {
  function assertStr(val: unknown, prop: string): asserts val is string {
    if (typeof val !== 'string') {
      throw Error(
        `Expected value of config.${prop} to be a string\nGot: ${val}`,
      )
    }
  }

  function assertDir(val: unknown, prop: string): asserts val is string {
    assertStr(val, prop)
    if (!test('-d', val)) {
      throw Error(
        `Expected value of config.${prop} to be a directory\nGot: ${val}`,
      )
    }
  }

  function assertBool(val: unknown, prop: string): asserts val is boolean {
    if (typeof val !== 'boolean') {
      throw Error(
        `Expected value of config.${prop} to be a boolean\nGot: ${val}`,
      )
    }
  }

  function assertCompilerSettings(
    val: unknown,
  ): asserts val is App['compilerSettings'] {
    if (!val || typeof val !== 'object') {
      throw Error(
        `Expected value of config.compilerSettings to be an object\nGot:${val}`,
      )
    }

    const compilerSettings = val as DeepUnknown<App['compilerSettings']>
    if (
      !compilerSettings.versions ||
      typeof compilerSettings.versions !== 'object'
    ) {
      throw Error(
        `Expected value of config.compilerSettings.versions to be a dictionary\nGot:${JSON.stringify(
          compilerSettings.versions,
        )}`,
      )
    }
  }

  const appConfig = json as DeepUnknown<App>
  assertDir(appConfig.contractsDir, 'contractsDir')
  assertStr(appConfig.artifactsDir, 'artifactsDir')
  assertStr(appConfig.contractAbstractionDir, 'contractAbstractionDir')
  assertBool(appConfig.useDockerisedSolc, 'useDockerisedSolc')
  assertCompilerSettings(appConfig.compilerSettings)
}
