'use strict';

var util = require('util');

var Item = require('./item');

var constants = require('constants');

/**
 * A directory.
 * @constructor
 */
function SymbolicLink() {
  Item.call(this);

  /**
   * Relative path to source.
   * @type {string}
   */
  this._path = undefined;
}
util.inherits(SymbolicLink, Item);

/**
 * Set the path to the source.
 * @param {string} pathname Path to source.
 */
SymbolicLink.prototype.setPath = function(pathname) {
  this._path = pathname;
};

/**
 * Get the path to the source.
 * @return {string} Path to source.
 */
SymbolicLink.prototype.getPath = function() {
  return this._path;
};

/**
 * Get symbolic link stats.
 * @return {Object} Stats properties.
 */
SymbolicLink.prototype.getStats = function() {
  var size = this._path.length;
  var stats = Item.prototype.getStats.call(this);
  stats.mode = this.getMode() | constants.S_IFLNK;
  stats.size = size;
  stats.blocks = Math.ceil(size / 512);
  return stats;
};

/**
 * Export the constructor.
 * @type {function()}
 */
exports = module.exports = SymbolicLink;
