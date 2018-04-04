var path = require("path");
var fs = require("fs");
var OS = require("os");
var prependFile = require('prepend-file');
var WebpackOnBuildPlugin = require('on-build-webpack');

var applyBaseConfig = require('./base.webpack.config')

var outputDir = path.join(__dirname, '..', 'build');
var outputFilename = 'cli.node.js';

module.exports = applyBaseConfig({
  entry: './cli.js',
  output: {
    path: outputDir,
    filename: outputFilename,
  },
  module: {
    rules: [
      { test: /\.js$/, use: "shebang-loader" }
    ]
  },
  plugins: [
    // Put the shebang back on and make sure it's executable.
    new WebpackOnBuildPlugin(function(stats) {
      var outputFile = path.join(outputDir, outputFilename);
      if (fs.existsSync(outputFile)) {
        prependFile.sync(outputFile, '#!/usr/bin/env node' + OS.EOL);
        fs.chmodSync(outputFile, '755');
      }
    })
  ]
})
