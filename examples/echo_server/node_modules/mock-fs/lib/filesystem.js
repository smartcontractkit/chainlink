'use strict';

var os = require('os');
var path = require('path');

var Directory = require('./directory');
var File = require('./file');
var FSError = require('./error');
var SymbolicLink = require('./symlink');

var isWindows = process.platform === 'win32';

function getPathParts(filepath) {
  var parts = path._makeLong(path.resolve(filepath)).split(path.sep);
  parts.shift();
  if (isWindows) {
    // parts currently looks like ['', '?', 'c:', ...]
    parts.shift();
    var q = parts.shift(); // should be '?'
    var base = '\\\\' + q + '\\' + parts.shift().toLowerCase();
    parts.unshift(base);
  }
  if (parts[parts.length - 1] === '') {
    parts.pop();
  }
  return parts;
}

/**
 * Create a new file system.
 * @param {Object} options Any filesystem options.
 * @param {boolean} options.createCwd Create a directory for `process.cwd()`
 *     (defaults to `true`).
 * @param {boolean} options.createTmp Create a directory for `os.tmpdir()`
 *     (defaults to `true`).
 * @constructor
 */
function FileSystem(options) {
  options = options || {};

  var createCwd = 'createCwd' in options ? options.createCwd : true;
  var createTmp = 'createTmp' in options ? options.createTmp : true;

  var root = new Directory();

  // populate with default directories
  var defaults = [];
  if (createCwd) {
    defaults.push(process.cwd());
  }

  if (createTmp) {
    defaults.push((os.tmpdir && os.tmpdir()) || os.tmpDir());
  }

  defaults.forEach(function(dir) {
    var parts = getPathParts(dir);
    var directory = root;
    var i, ii, name, candidate;
    for (i = 0, ii = parts.length; i < ii; ++i) {
      name = parts[i];
      candidate = directory.getItem(name);
      if (!candidate) {
        directory = directory.addItem(name, new Directory());
      } else if (candidate instanceof Directory) {
        directory = candidate;
      } else {
        throw new Error('Failed to create directory: ' + dir);
      }
    }
  });

  /**
   * Root directory.
   * @type {Directory}
   */
  this._root = root;
}

/**
 * Get the root directory.
 * @return {Directory} The root directory.
 */
FileSystem.prototype.getRoot = function() {
  return this._root;
};

/**
 * Get a file system item.
 * @param {string} filepath Path to item.
 * @return {Item} The item (or null if not found).
 */
FileSystem.prototype.getItem = function(filepath) {
  var parts = getPathParts(filepath);
  var currentParts = getPathParts(process.cwd());
  var item = this._root;
  var itemPath = '/';
  var name;
  for (var i = 0, ii = parts.length; i < ii; ++i) {
    name = parts[i];
    while (item instanceof SymbolicLink) {
      // Symbolic link being traversed as a directory --- If link targets
      // another symbolic link, resolve target's path relative to the original
      // link's target, otherwise relative to the current item.
      itemPath = path.resolve(path.dirname(itemPath), item.getPath());
      item = this.getItem(itemPath);
    }
    if (item) {
      if (item instanceof Directory && name !== currentParts[i]) {
        // make sure traversal is allowed
        if (!item.canExecute()) {
          throw new FSError('EACCES', filepath);
        }
      }
      item = item.getItem(name);
    }
    if (!item) {
      break;
    }
    itemPath = path.resolve(itemPath, name);
  }
  return item;
};

/**
 * Populate a directory with an item.
 * @param {Directory} directory The directory to populate.
 * @param {string} name The name of the item.
 * @param {string|Buffer|function|Object} obj Instructions for creating the
 *     item.
 */
function populate(directory, name, obj) {
  var item;
  if (typeof obj === 'string' || Buffer.isBuffer(obj)) {
    // contents for a file
    item = new File();
    item.setContent(obj);
  } else if (typeof obj === 'function') {
    // item factory
    item = obj();
  } else if (typeof obj === 'object') {
    // directory with more to populate
    item = new Directory();
    for (var key in obj) {
      populate(item, key, obj[key]);
    }
  } else {
    throw new Error('Unsupported type: ' + typeof obj + ' of item ' + name);
  }
  /**
   * Special exception for redundant adding of empty directories.
   */
  if (
    item instanceof Directory &&
    item.list().length === 0 &&
    directory.getItem(name) instanceof Directory
  ) {
    // pass
  } else {
    directory.addItem(name, item);
  }
}

/**
 * Configure a mock file system.
 * @param {Object} paths Config object.
 * @param {Object} options Any filesystem options.
 * @param {boolean} options.createCwd Create a directory for `process.cwd()`
 *     (defaults to `true`).
 * @param {boolean} options.createTmp Create a directory for `os.tmpdir()`
 *     (defaults to `true`).
 * @return {FileSystem} Mock file system.
 */
FileSystem.create = function(paths, options) {
  var system = new FileSystem(options);

  for (var filepath in paths) {
    var parts = getPathParts(filepath);
    var directory = system._root;
    var i, ii, name, candidate;
    for (i = 0, ii = parts.length - 1; i < ii; ++i) {
      name = parts[i];
      candidate = directory.getItem(name);
      if (!candidate) {
        directory = directory.addItem(name, new Directory());
      } else if (candidate instanceof Directory) {
        directory = candidate;
      } else {
        throw new Error('Failed to create directory: ' + filepath);
      }
    }
    populate(directory, parts[i], paths[filepath]);
  }

  return system;
};

/**
 * Generate a factory for new files.
 * @param {Object} config File config.
 * @return {function():File} Factory that creates a new file.
 */
FileSystem.file = function(config) {
  config = config || {};
  return function() {
    var file = new File();
    if (config.hasOwnProperty('content')) {
      file.setContent(config.content);
    }
    if (config.hasOwnProperty('mode')) {
      file.setMode(config.mode);
    } else {
      file.setMode(438); // 0666
    }
    if (config.hasOwnProperty('uid')) {
      file.setUid(config.uid);
    }
    if (config.hasOwnProperty('gid')) {
      file.setGid(config.gid);
    }
    if (config.hasOwnProperty('atime')) {
      file.setATime(config.atime);
    }
    if (config.hasOwnProperty('ctime')) {
      file.setCTime(config.ctime);
    }
    if (config.hasOwnProperty('mtime')) {
      file.setMTime(config.mtime);
    }
    if (config.hasOwnProperty('birthtime')) {
      file.setBirthtime(config.birthtime);
    }
    return file;
  };
};

/**
 * Generate a factory for new symbolic links.
 * @param {Object} config File config.
 * @return {function():File} Factory that creates a new symbolic link.
 */
FileSystem.symlink = function(config) {
  config = config || {};
  return function() {
    var link = new SymbolicLink();
    if (config.hasOwnProperty('mode')) {
      link.setMode(config.mode);
    } else {
      link.setMode(438); // 0666
    }
    if (config.hasOwnProperty('uid')) {
      link.setUid(config.uid);
    }
    if (config.hasOwnProperty('gid')) {
      link.setGid(config.gid);
    }
    if (config.hasOwnProperty('path')) {
      link.setPath(config.path);
    } else {
      throw new Error('Missing "path" property');
    }
    if (config.hasOwnProperty('atime')) {
      link.setATime(config.atime);
    }
    if (config.hasOwnProperty('ctime')) {
      link.setCTime(config.ctime);
    }
    if (config.hasOwnProperty('mtime')) {
      link.setMTime(config.mtime);
    }
    if (config.hasOwnProperty('birthtime')) {
      link.setBirthtime(config.birthtime);
    }
    return link;
  };
};

/**
 * Generate a factory for new directories.
 * @param {Object} config File config.
 * @return {function():Directory} Factory that creates a new directory.
 */
FileSystem.directory = function(config) {
  config = config || {};
  return function() {
    var dir = new Directory();
    if (config.hasOwnProperty('mode')) {
      dir.setMode(config.mode);
    }
    if (config.hasOwnProperty('uid')) {
      dir.setUid(config.uid);
    }
    if (config.hasOwnProperty('gid')) {
      dir.setGid(config.gid);
    }
    if (config.hasOwnProperty('items')) {
      for (var name in config.items) {
        populate(dir, name, config.items[name]);
      }
    }
    if (config.hasOwnProperty('atime')) {
      dir.setATime(config.atime);
    }
    if (config.hasOwnProperty('ctime')) {
      dir.setCTime(config.ctime);
    }
    if (config.hasOwnProperty('mtime')) {
      dir.setMTime(config.mtime);
    }
    if (config.hasOwnProperty('birthtime')) {
      dir.setBirthtime(config.birthtime);
    }
    return dir;
  };
};

/**
 * Module exports.
 * @type {function}
 */
var exports = (module.exports = FileSystem);
exports.getPathParts = getPathParts;
