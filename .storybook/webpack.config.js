const path = require('path')

module.exports = (baseConfig, env, defaultConfig) => {
  defaultConfig.resolve.modules.push(path.resolve(__dirname, '../gui/src'))

  return defaultConfig
}
