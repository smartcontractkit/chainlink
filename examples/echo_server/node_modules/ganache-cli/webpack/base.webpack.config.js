var path = require("path");
var fs = require("fs");
var OS = require("os");

var outputDir = path.join(__dirname, '..', 'build');

var _ = require('lodash')


module.exports = function(override) {
  return _.merge({}, {
    target: 'node',
    output: {
      path: outputDir
    },
    devtool: 'source-map',
    resolve: {
      alias: {
        "electron": path.join(__dirname, "..", "./nil.js"),
        "ws": path.join(__dirname, "..", "./nil.js"),
        "scrypt": "js-scrypt",
        "secp256k1": path.join(__dirname, "..", "node_modules", "secp256k1", "elliptic.js")
      }
    },
    mode: 'production'
  }, override)
}
