/* eslint-disable @typescript-eslint/no-var-requires */
const { createTransformer } = require('babel-jest')

module.exports = createTransformer({ rootMode: 'upward' })
