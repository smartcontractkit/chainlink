import { resolve, join, parse } from 'path'
import chalk from 'chalk'
import { cp, rm, ls, test, mkdir, cat } from 'shelljs'
import { writeFileSync } from 'fs'

// when this is non-empty, no files will be written
const DRY_RUN = process.env.DRY_RUN

// logging functions with colour output
const err = (...args: string[]) => console.error(chalk.red(...args))
const warn = (...args: string[]) => console.warn(chalk.yellow(...args))
const log = (...args: string[]) => console.log(chalk.green(...args))
const info = (...args: string[]) => console.log(chalk.grey(...args))

function main() {
  const [generatedPath, distPath] = [process.argv[2], process.argv[3]].map(p =>
    resolve(p),
  )

  exportGeneratedContractFactories(generatedPath, distPath)
}
main()

/**
 * Export all generated contract factories and their associated types
 * @param generatedPath The path of the generated files
 * @param distPath The path of the post-tsc generated files
 */
export function exportGeneratedContractFactories(
  generatedPath: string,
  distPath: string,
): void {
  const dir = getGeneratedFilePaths(generatedPath)
  const exportPaths = dir
    .map(makeExportPath)
    .filter(Boolean)
    .join('\n')
  info(`Export paths:\n${exportPaths}`)

  if (!DRY_RUN) {
    makeBarrelFile(generatedPath, exportPaths)
    copyTypings(generatedPath, distPath)
  }
}

/**
 * This copies the .d.ts files from the generated phase over,
 * since the typescript compiler drops any .d.ts source files during
 * compilation
 * @param generatedPath The path of the generated files
 * @param distPath The path of the post-tsc generated files
 */
function copyTypings(generatedPath: string, distPath: string): void {
  mkdir('-p', distPath)
  cp(`${generatedPath}/*.d.ts`, distPath)
}

/**
 * Create a barrel file which contains all of the exports.
 * This will replace the existing barrel file if it already exists.
 * @path the path to create the barrel file
 * @param data The data to write to the barrel file
 */
function makeBarrelFile(path: string, data: string): void {
  const exportFilePath = join(path, 'index.ts')
  warn(`Writing barrel file to ${exportFilePath}`)
  writeFileSync(exportFilePath, data)
  mergeIndexes(path)
}

/**
 * Making a barrel file makes us end up with two index files,
 * since one already is generated (albeit with a .d.ts extension).
 *
 * This function merges both of them into one index file, and deletes the
 * .d.fs one.
 * @param path The path of the generated files
 */
function mergeIndexes(path: string): void {
  const exportFilePath = join(path, 'index.ts')
  const declarationsFilePath = join(path, 'index.d.ts')
  const declarationFile = cat(declarationsFilePath)
  const exportsFile = cat(exportFilePath)

  writeFileSync(exportFilePath, [declarationFile, exportsFile].join('\n'))
  rm(declarationsFilePath)
}

/**
 * Check if the generated directory for smart contract factories exists
 * @param path The path of the generated files
 */
function generatedDirExists(path: string): boolean {
  log(`Checking if directory: ${path} exists...`)

  return test('-d', path)
}

/**
 * Get all the generated file paths from the generated directory
 * @param path The path of the generated files
 */
function getGeneratedFilePaths(path: string): string[] {
  if (!generatedDirExists(path)) {
    err(`Directory ${path} does not exist. Exiting...`)
    process.exit(1)
  }
  log(`Directory ${path} exists, continuing...`)
  return ls(path)
}

/**
 * Create an es6 export of a filename, handles interface and index conflicts naively
 * @param fileName The filename to export
 */
function makeExportPath(fileName: string): string {
  const { name } = parse(fileName)

  if (name.endsWith('.d')) {
    return ''
  } else if (name === 'index') {
    return ''
  } else {
    return `export * from './${name}'`
  }
}
