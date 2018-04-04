'use strict';

var util = require('util');

var Item = require('./item');

var constants = require('constants');

/**
 * A directory.
 * @constructor
 */
function Directory() {
  Item.call(this);

  /**
   * Items in this directory.
   * @type {Object.<string, Item>}
   */
  this._items = {};

  /**
   * Permissions.
   */
  this._mode = 511; // 0777
}
util.inherits(Directory, Item);

/**
 * Add an item to the directory.
 * @param {string} name The name to give the item.
 * @param {Item} item The item to add.
 * @return {Item} The added item.
 */
Directory.prototype.addItem = function(name, item) {
  if (this._items.hasOwnProperty(name)) {
    throw new Error('Item with the same name already exists: ' + name);
  }
  this._items[name] = item;
  ++item.links;
  if (item instanceof Directory) {
    // for '.' entry
    ++item.links;
    // for subdirectory
    ++this.links;
  }
  this.setMTime(new Date());
  return item;
};

/**
 * Get a named item.
 * @param {string} name Item name.
 * @return {Item} The named item (or null if none).
 */
Directory.prototype.getItem = function(name) {
  var item = null;
  if (this._items.hasOwnProperty(name)) {
    item = this._items[name];
  }
  return item;
};

/**
 * Remove an item.
 * @param {string} name Name of item to remove.
 * @return {Item} The orphan item.
 */
Directory.prototype.removeItem = function(name) {
  if (!this._items.hasOwnProperty(name)) {
    throw new Error('Item does not exist in directory: ' + name);
  }
  var item = this._items[name];
  delete this._items[name];
  --item.links;
  if (item instanceof Directory) {
    // for '.' entry
    --item.links;
    // for subdirectory
    --this.links;
  }
  this.setMTime(new Date());
  return item;
};

/**
 * Get list of item names in this directory.
 * @return {Array.<string>} Item names.
 */
Directory.prototype.list = function() {
  return Object.keys(this._items).sort();
};

/**
 * Get directory stats.
 * @return {Object} Stats properties.
 */
Directory.prototype.getStats = function() {
  var stats = Item.prototype.getStats.call(this);
  stats.mode = this.getMode() | constants.S_IFDIR;
  stats.size = 1;
  stats.blocks = 1;
  return stats;
};

/**
 * Export the constructor.
 * @type {function()}
 */
exports = module.exports = Directory;
