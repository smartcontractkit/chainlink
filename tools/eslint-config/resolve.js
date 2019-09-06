// https://github.com/eslint/eslint/issues/3458#issuecomment-516716165

const path = require('path')

let currentModule = module
while (
  !/[\\/]eslint[\\/]lib[\\/]cli-engine[\\/]config-array-factory\.js/i.test(
    currentModule.filename,
  )
) {
  if (!currentModule.parent) {
    // This was tested with ESLint 6.1.0; other versions may not work
    throw new Error(
      'Failed to patch ESLint because the calling module was not recognized',
    )
  }
  currentModule = currentModule.parent
}
const eslintFolder = path.join(path.dirname(currentModule.filename), '../..')

const configArrayFactoryPath = path.join(
  eslintFolder,
  'lib/cli-engine/config-array-factory',
)
const configArrayFactoryModule = require(configArrayFactoryPath)

const moduleResolverPath = path.join(
  eslintFolder,
  'lib/shared/relative-module-resolver',
)
const ModuleResolver = require(moduleResolverPath)

const originalLoadPlugin =
  configArrayFactoryModule.ConfigArrayFactory.prototype._loadPlugin
configArrayFactoryModule.ConfigArrayFactory.prototype._loadPlugin = function(
  _,
  importerPath,
) {
  const originalResolve = ModuleResolver.resolve
  try {
    ModuleResolver.resolve = function(moduleName) {
      // resolve using importerPath instead of relativeToPath
      return originalResolve.call(this, moduleName, importerPath)
    }
    // eslint-disable-next-line prefer-rest-params
    return originalLoadPlugin.apply(this, arguments)
  } finally {
    ModuleResolver.resolve = originalResolve
  }
}
