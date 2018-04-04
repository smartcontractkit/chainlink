var path = require("path");
var fs = require("fs");
var OS = require("os");

var applyBaseConfig = require('./base.webpack.config')

var outputDir = path.join(__dirname, '..', 'build');
var outputFilename = 'lib.node.js';

module.exports = applyBaseConfig({
  entry: './lib.js',
  target: 'node',
  output: {
    path: outputDir,
    filename: outputFilename,
    library: "ganache",
    libraryTarget: 'umd',
    umdNamedDefine: true
  }
})
