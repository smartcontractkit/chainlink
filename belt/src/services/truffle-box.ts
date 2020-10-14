import { join } from 'path'
import { ls, sed } from 'shelljs'
import { App } from './config'

/**
 * Modify a truffle box with the given solidity version
 *
 * @param solidityVersion A tuple of alias and version of a solidity version, e.g ['v0.4', '0.4.24']
 * @param path The path to the truffle box
 * @param dryRun Whether to actually modify the files in-place or to print the modified files to stdout
 */
export function modifyTruffleBoxWith(
  [solcVersionAlias, solcVersion]: [string, string],
  path: string,
  dryRun: boolean,
) {
  const solVersionToOzversion: Record<string, string> = {
    '0.5.0': '2.3.0',
    '0.4.24': '2.0.0',
  }
  const ozVersion = solVersionToOzversion[solcVersion]

  convertPackageJson(
    path,
    ozVersion,
    solVersionToOzversion,
    solcVersion,
    dryRun,
  )
  convertTruffleConfig(path, solcVersion, dryRun)
  convertSolidityFiles(path, ozVersion, dryRun, solcVersionAlias, solcVersion)
  convertJavascriptFiles(path, solcVersionAlias, dryRun)
}

/**
 * Get a solidity version by its alias or version number
 *
 * @param versionAliasOrVersion Either a solidity version alias "v0.6" | "0.6" or its full version "0.6.2"
 * @throws error if version given could not be found
 */
export function getSolidityVersionBy(versionAliasOrVersion: string) {
  const versions = getSolidityVersions()
  const version = versions.find(
    ([alias, full]) =>
      alias.replace('v', '') === versionAliasOrVersion.replace('v', '') ||
      full === versionAliasOrVersion,
  )
  if (!version) {
    throw Error(
      `Could not find given version, here are the available versions: ${versions}`,
    )
  }

  return version
}

/**
 * Get a list of available solidity versions based on what's published in @chainlink/contracts
 *
 * The returned format is [alias, version] where alias can be "v0.6" | "0.6" and full version can be "0.6.2"
 */
export function getSolidityVersions(): [string, string][] {
  // eslint-disable-next-line @typescript-eslint/no-var-requires
  const config: App = require('@chainlink/contracts/app.config.json')

  return Object.entries(config.compilerSettings.versions).filter(([, v]) =>
    config.publicVersions.find((pv) => pv === v),
  )
}

/**
 * Get the path to the truffle config
 *
 * @param basePath The path to the truffle box
 */
export function getTruffleConfig(basePath: string): string {
  return join(basePath, 'truffle-config.js')
}

/**
 * Get a list of all javascript files in the truffle box
 *
 * @param basePath The path to the truffle box
 */
export function getJavascriptFiles(basePath: string): string[] {
  const directories = ['scripts', 'test', 'migrations']

  return directories
    .map((d) => ls(join(basePath, d, '**', '*.js')))
    .reduce<string[]>((prev, next) => {
      return prev.concat(next)
    }, [])
}

/**
 * Get path to the package.json
 *
 * @param basePath The path to the truffle box
 */
export function getPackageJson(basePath: string): string {
  return join(basePath, 'package.json')
}

/**
 * Get a list of all solidity files in the truffle box
 *
 * @param basePath The path to the truffle box
 */
export function getSolidityFiles(basePath: string): string[] {
  return [...ls(join(basePath, 'contracts', '**', '*.sol'))]
}

function convertPackageJson(
  path: string,
  ozVersion: string,
  solVersionToOzversion: Record<string, string>,
  solcVersion: string,
  dryRun: boolean,
) {
  const packageJson = getPackageJson(path)
  if (ozVersion !== '2.0.0') {
    replaceInFile(
      `"openzeppelin-solidity": "1.12.0"`,
      `"@openzeppelin/contracts": "${solVersionToOzversion[solcVersion]}"`,
      [packageJson],
      dryRun,
    )
  }
}

function convertTruffleConfig(
  path: string,
  solcVersion: string,
  dryRun: boolean,
) {
  const truffleConfig = getTruffleConfig(path)
  replaceInFile(
    "version: '0.4.24'",
    `version: '${solcVersion}'`,
    [truffleConfig],
    dryRun,
  )
}

function convertSolidityFiles(
  path: string,
  ozVersion: string,
  dryRun: boolean,
  solcVersionAlias: string,
  solcVersion: string,
) {
  const solFiles = getSolidityFiles(path)
  if (ozVersion !== '2.0.0') {
    replaceInFile(
      'import "openzeppelin-solidity',
      `import "@openzeppelin`,
      solFiles,
      dryRun,
    )
  }
  replaceInFile(
    '@chainlink/contracts/src/v0.4',
    `@chainlink/contracts/src/${solcVersionAlias}`,
    solFiles,
    dryRun,
  )
  replaceInFile(
    'pragma solidity 0.4.24;',
    `pragma solidity ${solcVersion};`,
    solFiles,
    dryRun,
  )
}

function convertJavascriptFiles(
  path: string,
  solcVersionAlias: string,
  dryRun: boolean,
) {
  const jsFiles = getJavascriptFiles(path)
  replaceInFile(
    '@chainlink/contracts/truffle/v0.4',
    `@chainlink/contracts/truffle/${solcVersionAlias}`,
    jsFiles,
    dryRun,
  )
  // replace linktoken back to v0.4
  replaceInFile(
    `@chainlink/contracts/truffle/${solcVersionAlias}/LinkToken`,
    '@chainlink/contracts/truffle/v0.4/LinkToken',
    jsFiles,
    dryRun,
  )
}

function replaceInFile(
  regex: string | RegExp,
  replacement: string,
  files: string[],
  dryRun: boolean,
) {
  if (dryRun) {
    const { stderr, stdout } = sed(regex, replacement, files)
    if (stdout) {
      console.log(stdout)
    }
    if (stderr) {
      console.error(stderr)
    }
  } else {
    sed('-i', regex, replacement, files)
  }
}
