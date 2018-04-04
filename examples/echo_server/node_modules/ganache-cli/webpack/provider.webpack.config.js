var path = require("path");
var fs = require("fs");
var OS = require("os");

var applyBaseConfig = require('./base.webpack.config')

var outputDir = path.join(__dirname, '..', 'build');
var outputFilename = 'provider.node.js';

module.exports = applyBaseConfig({
  entry: './node_modules/ganache-core/lib/provider.js',
  target: 'node',
  output: {
    path: outputDir,
    filename: outputFilename,
    library: "Provider",
    libraryTarget: 'umd',
    umdNamedDefine: true
  }
})
