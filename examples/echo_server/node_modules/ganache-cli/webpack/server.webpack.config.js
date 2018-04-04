var path = require("path");
var fs = require("fs");
var OS = require("os");

var applyBaseConfig = require('./base.webpack.config')

var outputDir = path.join(__dirname, '..', 'build');
var outputFilename = 'server.node.js';

module.exports = applyBaseConfig({
  entry: './node_modules/ganache-core/lib/server.js',
  target: 'node',
  output: {
    path: outputDir,
    filename: outputFilename,
    library: "Server",
    libraryTarget: 'umd',
    umdNamedDefine: true
  }
})
