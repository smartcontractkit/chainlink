const { createTransformer } = require('babel-jest')

module.exports = createTransformer({ rootMode: 'upward' })
