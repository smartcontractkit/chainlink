const path = require('path')

module.exports = (baseConfig, env, defaultConfig) => {
  defaultConfig.resolve.modules.push(path.resolve(__dirname, '../operator_ui/src'))

  return defaultConfig
}
