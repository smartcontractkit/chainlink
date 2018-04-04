'use strict';

var path = require('path');

var File = require('./file');
var FileDescriptor = require('./descriptor');
var Directory = require('./directory');
var SymbolicLink = require('./symlink');
var FSError = require('./error');
var constants = require('constants');
var getPathParts = require('./filesystem').getPathParts;

/** Workaround for optimizations in node 8 */
var fsBinding = process.binding('fs');
var statValues = fsBinding.getStatValues ? fsBinding.getStatValues() : [];

/**
 * Call the provided function and either return the result or call the callback
 * with it (depending on if a callback is provided).
 * @param {function()} callback Optional callback.
 * @param {Object} thisArg This argument for the following function.
 * @param {function()} func Function to call.
 * @return {*} Return (if callback is not provided).
 */
function maybeCallback(callback, thisArg, func) {
  if (callback && typeof callback === 'function') {
    var err = null;
    var val;
    try {
      val = func.call(thisArg);
    } catch (e) {
      err = e;
    }
    process.nextTick(function() {
      if (val === undefined) {
        callback(err);
      } else {
        callback(err, val);
      }
    });
  } else {
    return func.call(thisArg);
  }
}

/**
 * Handle FSReqWrap oncomplete.
 * @param {Function} callback The callback.
 * @return {Function} The normalized callback.
 */
function normalizeCallback(callback) {
  if (callback && typeof callback.oncomplete === 'function') {
    // Unpack callback from FSReqWrap
    callback = callback.oncomplete.bind(callback);
  }
  return callback;
}

/**
 * Handle stat optimizations introduced in Node 8.
 * See https://github.com/nodejs/node/pull/11665.
 * @param {Function} callback The callback.
 * @return {Function} The wrapped callback.
 */
function wrapStatsCallback(callback) {
  if (callback && typeof callback.oncomplete === 'function') {
    // Unpack callback from FSReqWrap
    callback = callback.oncomplete.bind(callback);
  }
  if (typeof callback === 'function') {
    return function(err, stats) {
      if (stats) {
        fillStatsArray(stats, statValues);
      }
      callback.apply(this, arguments);
    };
  } else {
    return callback;
  }
}

function notImplemented() {
  throw new Error('Method not implemented');
}

/**
 * Create a new stats object.
 * @param {Object} config Stats properties.
 * @constructor
 */
function Stats(config) {
  for (var key in config) {
    this[key] = config[key];
  }
}

/**
 * Check if mode indicates property.
 * @param {number} property Property to check.
 * @return {boolean} Property matches mode.
 */
Stats.prototype._checkModeProperty = function(property) {
  return (this.mode & constants.S_IFMT) === property;
};

/**
 * @return {Boolean} Is a directory.
 */
Stats.prototype.isDirectory = function() {
  return this._checkModeProperty(constants.S_IFDIR);
};

/**
 * @return {Boolean} Is a regular file.
 */
Stats.prototype.isFile = function() {
  return this._checkModeProperty(constants.S_IFREG);
};

/**
 * @return {Boolean} Is a block device.
 */
Stats.prototype.isBlockDevice = function() {
  return this._checkModeProperty(constants.S_IFBLK);
};

/**
 * @return {Boolean} Is a character device.
 */
Stats.prototype.isCharacterDevice = function() {
  return this._checkModeProperty(constants.S_IFCHR);
};

/**
 * @return {Boolean} Is a symbolic link.
 */
Stats.prototype.isSymbolicLink = function() {
  return this._checkModeProperty(constants.S_IFLNK);
};

/**
 * @return {Boolean} Is a named pipe.
 */
Stats.prototype.isFIFO = function() {
  return this._checkModeProperty(constants.S_IFIFO);
};

/**
 * @return {Boolean} Is a socket.
 */
Stats.prototype.isSocket = function() {
  return this._checkModeProperty(constants.S_IFSOCK);
};

/**
 * Create a new binding with the given file system.
 * @param {FileSystem} system Mock file system.
 * @constructor
 */
function Binding(system) {
  /**
   * Mock file system.
   * @type {FileSystem}
   */
  this._system = system;

  /**
   * Stats constructor.
   * @type {function}
   */
  this.Stats = Stats;

  /**
   * Lookup of open files.
   * @type {Object.<number, FileDescriptor>}
   */
  this._openFiles = {};

  /**
   * Counter for file descriptors.
   * @type {number}
   */
  this._counter = 0;
}

/**
 * Get the file system underlying this binding.
 * @return {FileSystem} The underlying file system.
 */
Binding.prototype.getSystem = function() {
  return this._system;
};

/**
 * Reset the file system underlying this binding.
 * @param {FileSystem} system The new file system.
 */
Binding.prototype.setSystem = function(system) {
  this._system = system;
};

/**
 * Get a file descriptor.
 * @param {number} fd File descriptor identifier.
 * @return {FileDescriptor} File descriptor.
 */
Binding.prototype._getDescriptorById = function(fd) {
  if (!this._openFiles.hasOwnProperty(fd)) {
    throw new FSError('EBADF');
  }
  return this._openFiles[fd];
};

/**
 * Keep track of a file descriptor as open.
 * @param {FileDescriptor} descriptor The file descriptor.
 * @return {number} Identifier for file descriptor.
 */
Binding.prototype._trackDescriptor = function(descriptor) {
  var fd = ++this._counter;
  this._openFiles[fd] = descriptor;
  return fd;
};

/**
 * Stop tracking a file descriptor as open.
 * @param {number} fd Identifier for file descriptor.
 */
Binding.prototype._untrackDescriptorById = function(fd) {
  if (!this._openFiles.hasOwnProperty(fd)) {
    throw new FSError('EBADF');
  }
  delete this._openFiles[fd];
};

/**
 * Resolve the canonicalized absolute pathname.
 * @param {string|Buffer} filepath The file path.
 * @param {string} encoding The encoding for the return.
 * @return {string|Buffer} The real path.
 */
Binding.prototype.realpath = function(filepath, encoding, callback) {
  return maybeCallback(normalizeCallback(callback), this, function() {
    var realPath;
    if (Buffer.isBuffer(filepath)) {
      filepath = filepath.toString();
    }
    var resolved = path.resolve(filepath);
    var parts = getPathParts(resolved);
    var item = this._system.getRoot();
    var itemPath = '/';
    var name, i, ii;
    for (i = 0, ii = parts.length; i < ii; ++i) {
      name = parts[i];
      while (item instanceof SymbolicLink) {
        itemPath = path.resolve(path.dirname(itemPath), item.getPath());
        item = this._system.getItem(itemPath);
      }
      if (!item) {
        throw new FSError('ENOENT', filepath);
      }
      if (item instanceof Directory) {
        itemPath = path.resolve(itemPath, name);
        item = item.getItem(name);
      } else {
        throw new FSError('ENOTDIR', filepath);
      }
    }
    if (item) {
      while (item instanceof SymbolicLink) {
        itemPath = path.resolve(path.dirname(itemPath), item.getPath());
        item = this._system.getItem(itemPath);
      }
      realPath = itemPath;
    } else {
      throw new FSError('ENOENT', filepath);
    }

    if (encoding === 'buffer') {
      realPath = new Buffer(realPath);
    }
    return realPath;
  });
};

/**
 * Fill a Float64Array with stat information
 * This is based on the internal FillStatsArray function in Node.
 * https://github.com/nodejs/node/blob/4e05952a8a75af6df625415db612d3a9a1322682/src/node_file.cc#L533
 * @param {object} stats An object with file stats
 * @param {Float64Array} statValues A Float64Array where stat values should be inserted
 * @returns {void}
 */
function fillStatsArray(stats, statValues) {
  statValues[0] = stats.dev;
  statValues[1] = stats.mode;
  statValues[2] = stats.nlink;
  statValues[3] = stats.uid;
  statValues[4] = stats.gid;
  statValues[5] = stats.rdev;
  statValues[6] = stats.blksize;
  statValues[7] = stats.ino;
  statValues[8] = stats.size;
  statValues[9] = stats.blocks;
  statValues[10] = +stats.atime;
  statValues[11] = +stats.mtime;
  statValues[12] = +stats.ctime;
  statValues[13] = +stats.birthtime;
}

/**
 * Stat an item.
 * @param {string} filepath Path.
 * @param {function(Error, Stats)|Float64Array} callback Callback (optional). In Node 7.7.0+ this will be a Float64Array
 * that should be filled with stat values.
 * @return {Stats|undefined} Stats or undefined (if sync).
 */
Binding.prototype.stat = function(filepath, callback) {
  return maybeCallback(wrapStatsCallback(callback), this, function() {
    var item = this._system.getItem(filepath);
    if (item instanceof SymbolicLink) {
      item = this._system.getItem(
        path.resolve(path.dirname(filepath), item.getPath())
      );
    }
    if (!item) {
      throw new FSError('ENOENT', filepath);
    }
    var stats = item.getStats();

    // In Node 7.7.0+, binding.stat accepts a Float64Array as the second argument,
    // which should be filled with stat values.
    // In prior versions of Node, binding.stat simply returns a Stats instance.
    if (callback instanceof Float64Array) {
      fillStatsArray(stats, callback);
    } else {
      fillStatsArray(stats, statValues);
      return new Stats(stats);
    }
  });
};

/**
 * Stat an item.
 * @param {number} fd File descriptor.
 * @param {function(Error, Stats)|Float64Array} callback Callback (optional). In Node 7.7.0+ this will be a Float64Array
 * that should be filled with stat values.
 * @return {Stats|undefined} Stats or undefined (if sync).
 */
Binding.prototype.fstat = function(fd, callback) {
  return maybeCallback(wrapStatsCallback(callback), this, function() {
    var descriptor = this._getDescriptorById(fd);
    var item = descriptor.getItem();
    var stats = item.getStats();

    // In Node 7.7.0+, binding.stat accepts a Float64Array as the second argument,
    // which should be filled with stat values.
    // In prior versions of Node, binding.stat simply returns a Stats instance.
    if (callback instanceof Float64Array) {
      fillStatsArray(stats, callback);
    } else {
      fillStatsArray(stats, statValues);
      return new Stats(stats);
    }
  });
};

/**
 * Close a file descriptor.
 * @param {number} fd File descriptor.
 * @param {function(Error)} callback Callback (optional).
 */
Binding.prototype.close = function(fd, callback) {
  maybeCallback(normalizeCallback(callback), this, function() {
    this._untrackDescriptorById(fd);
  });
};

/**
 * Open and possibly create a file.
 * @param {string} pathname File path.
 * @param {number} flags Flags.
 * @param {number} mode Mode.
 * @param {function(Error, string)} callback Callback (optional).
 * @return {string} File descriptor (if sync).
 */
Binding.prototype.open = function(pathname, flags, mode, callback) {
  return maybeCallback(normalizeCallback(callback), this, function() {
    var descriptor = new FileDescriptor(flags);
    var item = this._system.getItem(pathname);
    while (item instanceof SymbolicLink) {
      item = this._system.getItem(
        path.resolve(path.dirname(pathname), item.getPath())
      );
    }
    if (descriptor.isExclusive() && item) {
      throw new FSError('EEXIST', pathname);
    }
    if (descriptor.isCreate() && !item) {
      var parent = this._system.getItem(path.dirname(pathname));
      if (!parent) {
        throw new FSError('ENOENT', pathname);
      }
      if (!(parent instanceof Directory)) {
        throw new FSError('ENOTDIR', pathname);
      }
      item = new File();
      if (mode) {
        item.setMode(mode);
      }
      parent.addItem(path.basename(pathname), item);
    }
    if (descriptor.isRead()) {
      if (!item) {
        throw new FSError('ENOENT', pathname);
      }
      if (!item.canRead()) {
        throw new FSError('EACCES', pathname);
      }
    }
    if (descriptor.isWrite() && !item.canWrite()) {
      throw new FSError('EACCES', pathname);
    }
    if (descriptor.isTruncate()) {
      item.setContent('');
    }
    if (descriptor.isTruncate() || descriptor.isAppend()) {
      descriptor.setPosition(item.getContent().length);
    }
    descriptor.setItem(item);
    return this._trackDescriptor(descriptor);
  });
};

/**
 * Read from a file descriptor.
 * @param {string} fd File descriptor.
 * @param {Buffer} buffer Buffer that the contents will be written to.
 * @param {number} offset Offset in the buffer to start writing to.
 * @param {number} length Number of bytes to read.
 * @param {?number} position Where to begin reading in the file.  If null,
 *     data will be read from the current file position.
 * @param {function(Error, number, Buffer)} callback Callback (optional) called
 *     with any error, number of bytes read, and the buffer.
 * @return {number} Number of bytes read (if sync).
 */
Binding.prototype.read = function(
  fd,
  buffer,
  offset,
  length,
  position,
  callback
) {
  return maybeCallback(normalizeCallback(callback), this, function() {
    var descriptor = this._getDescriptorById(fd);
    if (!descriptor.isRead()) {
      throw new FSError('EBADF');
    }
    var file = descriptor.getItem();
    if (!(file instanceof File)) {
      // deleted or not a regular file
      throw new FSError('EBADF');
    }
    if (typeof position !== 'number' || position < 0) {
      position = descriptor.getPosition();
    }
    var content = file.getContent();
    var start = Math.min(position, content.length);
    var end = Math.min(position + length, content.length);
    var read = start < end ? content.copy(buffer, offset, start, end) : 0;
    descriptor.setPosition(position + read);
    return read;
  });
};

/**
 * Write to a file descriptor given a buffer.
 * @param {string} fd File descriptor.
 * @param {Array<Buffer>} buffers Array of buffers with contents to write.
 * @param {?number} position Where to begin writing in the file.  If null,
 *     data will be written to the current file position.
 * @param {function(Error, number, Buffer)} callback Callback (optional) called
 *     with any error, number of bytes written, and the buffer.
 * @return {number} Number of bytes written (if sync).
 */
Binding.prototype.writeBuffers = function(fd, buffers, position, callback) {
  return maybeCallback(normalizeCallback(callback), this, function() {
    var descriptor = this._getDescriptorById(fd);
    if (!descriptor.isWrite()) {
      throw new FSError('EBADF');
    }
    var file = descriptor.getItem();
    if (!(file instanceof File)) {
      // not a regular file
      throw new FSError('EBADF');
    }
    if (typeof position !== 'number' || position < 0) {
      position = descriptor.getPosition();
    }
    var content = file.getContent();
    var newContent = Buffer.concat(buffers);
    var newLength = position + newContent.length;
    if (content.length < newLength) {
      var tempContent = new Buffer(newLength);
      content.copy(tempContent);
      content = tempContent;
    }
    var written = newContent.copy(content, position);
    file.setContent(content);
    descriptor.setPosition(newLength);
    return written;
  });
};

/**
 * Write to a file descriptor given a buffer.
 * @param {string} fd File descriptor.
 * @param {Buffer} buffer Buffer with contents to write.
 * @param {number} offset Offset in the buffer to start writing from.
 * @param {number} length Number of bytes to write.
 * @param {?number} position Where to begin writing in the file.  If null,
 *     data will be written to the current file position.
 * @param {function(Error, number, Buffer)} callback Callback (optional) called
 *     with any error, number of bytes written, and the buffer.
 * @return {number} Number of bytes written (if sync).
 */
Binding.prototype.writeBuffer = function(
  fd,
  buffer,
  offset,
  length,
  position,
  callback
) {
  return maybeCallback(normalizeCallback(callback), this, function() {
    var descriptor = this._getDescriptorById(fd);
    if (!descriptor.isWrite()) {
      throw new FSError('EBADF');
    }
    var file = descriptor.getItem();
    if (!(file instanceof File)) {
      // not a regular file
      throw new FSError('EBADF');
    }
    if (typeof position !== 'number' || position < 0) {
      position = descriptor.getPosition();
    }
    var content = file.getContent();
    var newLength = position + length;
    if (content.length < newLength) {
      var newContent = new Buffer(newLength);
      content.copy(newContent);
      content = newContent;
    }
    var sourceEnd = Math.min(offset + length, buffer.length);
    var written = buffer.copy(content, position, offset, sourceEnd);
    file.setContent(content);
    descriptor.setPosition(newLength);
    return written;
  });
};

/**
 * Alias for writeBuffer (used in Node <= 0.10).
 * @param {string} fd File descriptor.
 * @param {Buffer} buffer Buffer with contents to write.
 * @param {number} offset Offset in the buffer to start writing from.
 * @param {number} length Number of bytes to write.
 * @param {?number} position Where to begin writing in the file.  If null,
 *     data will be written to the current file position.
 * @param {function(Error, number, Buffer)} callback Callback (optional) called
 *     with any error, number of bytes written, and the buffer.
 * @return {number} Number of bytes written (if sync).
 */
Binding.prototype.write = Binding.prototype.writeBuffer;

/**
 * Write to a file descriptor given a string.
 * @param {string} fd File descriptor.
 * @param {string} string String with contents to write.
 * @param {number} position Where to begin writing in the file.  If null,
 *     data will be written to the current file position.
 * @param {string} encoding String encoding.
 * @param {function(Error, number, string)} callback Callback (optional) called
 *     with any error, number of bytes written, and the string.
 * @return {number} Number of bytes written (if sync).
 */
Binding.prototype.writeString = function(
  fd,
  string,
  position,
  encoding,
  callback
) {
  var buffer = new Buffer(string, encoding);
  var wrapper;
  if (callback) {
    if (callback.oncomplete) {
      callback = callback.oncomplete.bind(callback);
    }
    wrapper = function(err, written, returned) {
      callback(err, written, returned && string);
    };
  }
  return this.writeBuffer(fd, buffer, 0, string.length, position, wrapper);
};

/**
 * Rename a file.
 * @param {string} oldPath Old pathname.
 * @param {string} newPath New pathname.
 * @param {function(Error)} callback Callback (optional).
 * @return {undefined}
 */
Binding.prototype.rename = function(oldPath, newPath, callback) {
  return maybeCallback(normalizeCallback(callback), this, function() {
    var oldItem = this._system.getItem(oldPath);
    if (!oldItem) {
      throw new FSError('ENOENT', oldPath);
    }
    var oldParent = this._system.getItem(path.dirname(oldPath));
    var oldName = path.basename(oldPath);
    var newItem = this._system.getItem(newPath);
    var newParent = this._system.getItem(path.dirname(newPath));
    var newName = path.basename(newPath);
    if (newItem) {
      // make sure they are the same type
      if (oldItem instanceof File) {
        if (newItem instanceof Directory) {
          throw new FSError('EISDIR', newPath);
        }
      } else if (oldItem instanceof Directory) {
        if (!(newItem instanceof Directory)) {
          throw new FSError('ENOTDIR', newPath);
        }
        if (newItem.list().length > 0) {
          throw new FSError('ENOTEMPTY', newPath);
        }
      }
      newParent.removeItem(newName);
    } else {
      if (!newParent) {
        throw new FSError('ENOENT', newPath);
      }
      if (!(newParent instanceof Directory)) {
        throw new FSError('ENOTDIR', newPath);
      }
    }
    oldParent.removeItem(oldName);
    newParent.addItem(newName, oldItem);
  });
};

/**
 * Read a directory.
 * @param {string} dirpath Path to directory.
 * @param {string} encoding The encoding ('utf-8' or 'buffer').
 * @param {function(Error, (Array.<string>|Array.<Buffer>)} callback Callback
 *     (optional) called with any error or array of items in the directory.
 * @return {Array.<string>|Array.<Buffer>} Array of items in directory (if sync).
 */
Binding.prototype.readdir = function(dirpath, encoding, callback) {
  if (encoding && typeof encoding !== 'string') {
    callback = encoding;
    encoding = 'utf-8';
  }
  return maybeCallback(normalizeCallback(callback), this, function() {
    var dpath = dirpath;
    var dir = this._system.getItem(dirpath);
    while (dir instanceof SymbolicLink) {
      dpath = path.resolve(path.dirname(dpath), dir.getPath());
      dir = this._system.getItem(dpath);
    }
    if (!dir) {
      throw new FSError('ENOENT', dirpath);
    }
    if (!(dir instanceof Directory)) {
      throw new FSError('ENOTDIR', dirpath);
    }
    var list = dir.list();
    if (encoding === 'buffer') {
      list = list.map(function(item) {
        return new Buffer(item);
      });
    }
    return list;
  });
};

/**
 * Create a directory.
 * @param {string} pathname Path to new directory.
 * @param {number} mode Permissions.
 * @param {function(Error)} callback Optional callback.
 */
Binding.prototype.mkdir = function(pathname, mode, callback) {
  maybeCallback(normalizeCallback(callback), this, function() {
    var item = this._system.getItem(pathname);
    if (item) {
      throw new FSError('EEXIST', pathname);
    }
    var parent = this._system.getItem(path.dirname(pathname));
    if (!parent) {
      throw new FSError('ENOENT', pathname);
    }
    this.access(path.dirname(pathname), parseInt('0002', 8));
    var dir = new Directory();
    if (mode) {
      dir.setMode(mode);
    }
    parent.addItem(path.basename(pathname), dir);
  });
};

/**
 * Remove a directory.
 * @param {string} pathname Path to directory.
 * @param {function(Error)} callback Optional callback.
 */
Binding.prototype.rmdir = function(pathname, callback) {
  maybeCallback(normalizeCallback(callback), this, function() {
    var item = this._system.getItem(pathname);
    if (!item) {
      throw new FSError('ENOENT', pathname);
    }
    if (!(item instanceof Directory)) {
      throw new FSError('ENOTDIR', pathname);
    }
    if (item.list().length > 0) {
      throw new FSError('ENOTEMPTY', pathname);
    }
    this.access(path.dirname(pathname), parseInt('0002', 8));
    var parent = this._system.getItem(path.dirname(pathname));
    parent.removeItem(path.basename(pathname));
  });
};

var PATH_CHARS =
  'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';

var MAX_ATTEMPTS = 62 * 62 * 62;

/**
 * Create a directory based on a template.
 * See http://web.mit.edu/freebsd/head/lib/libc/stdio/mktemp.c
 * @param {string} template Path template (trailing Xs will be replaced).
 * @param {string} encoding The encoding ('utf-8' or 'buffer').
 * @param {function(Error, string)} callback Optional callback.
 */
Binding.prototype.mkdtemp = function(prefix, encoding, callback) {
  if (encoding && typeof encoding !== 'string') {
    callback = encoding;
    encoding = 'utf-8';
  }
  return maybeCallback(normalizeCallback(callback), this, function() {
    prefix = prefix.replace(/X{0,6}$/, 'XXXXXX');
    var parentPath = path.dirname(prefix);
    var parent = this._system.getItem(parentPath);
    if (!parent) {
      throw new FSError('ENOENT', prefix);
    }
    if (!(parent instanceof Directory)) {
      throw new FSError('ENOTDIR', prefix);
    }
    this.access(parentPath, parseInt('0002', 8));
    var template = path.basename(prefix);
    var unique = false;
    var count = 0;
    var name;
    while (!unique && count < MAX_ATTEMPTS) {
      var position = template.length - 1;
      var replacement = '';
      while (template.charAt(position) === 'X') {
        replacement += PATH_CHARS.charAt(
          Math.floor(PATH_CHARS.length * Math.random())
        );
        position -= 1;
      }
      var candidate = template.slice(0, position + 1) + replacement;
      if (!parent.getItem(candidate)) {
        name = candidate;
        unique = true;
      }
      count += 1;
    }
    if (!name) {
      throw new FSError('EEXIST', prefix);
    }
    var dir = new Directory();
    parent.addItem(name, dir);
    var uniquePath = path.join(parentPath, name);
    if (encoding === 'buffer') {
      uniquePath = new Buffer(uniquePath);
    }
    return uniquePath;
  });
};

/**
 * Truncate a file.
 * @param {number} fd File descriptor.
 * @param {number} len Number of bytes.
 * @param {function(Error)} callback Optional callback.
 */
Binding.prototype.ftruncate = function(fd, len, callback) {
  maybeCallback(normalizeCallback(callback), this, function() {
    var descriptor = this._getDescriptorById(fd);
    if (!descriptor.isWrite()) {
      throw new FSError('EINVAL');
    }
    var file = descriptor.getItem();
    if (!(file instanceof File)) {
      throw new FSError('EINVAL');
    }
    var content = file.getContent();
    var newContent = new Buffer(len);
    content.copy(newContent);
    file.setContent(newContent);
  });
};

/**
 * Legacy support.
 * @param {number} fd File descriptor.
 * @param {number} len Number of bytes.
 * @param {function(Error)} callback Optional callback.
 */
Binding.prototype.truncate = Binding.prototype.ftruncate;

/**
 * Change user and group owner.
 * @param {string} pathname Path.
 * @param {number} uid User id.
 * @param {number} gid Group id.
 * @param {function(Error)} callback Optional callback.
 */
Binding.prototype.chown = function(pathname, uid, gid, callback) {
  maybeCallback(normalizeCallback(callback), this, function() {
    var item = this._system.getItem(pathname);
    if (!item) {
      throw new FSError('ENOENT', pathname);
    }
    item.setUid(uid);
    item.setGid(gid);
  });
};

/**
 * Change user and group owner.
 * @param {number} fd File descriptor.
 * @param {number} uid User id.
 * @param {number} gid Group id.
 * @param {function(Error)} callback Optional callback.
 */
Binding.prototype.fchown = function(fd, uid, gid, callback) {
  maybeCallback(normalizeCallback(callback), this, function() {
    var descriptor = this._getDescriptorById(fd);
    var item = descriptor.getItem();
    item.setUid(uid);
    item.setGid(gid);
  });
};

/**
 * Change permissions.
 * @param {string} pathname Path.
 * @param {number} mode Mode.
 * @param {function(Error)} callback Optional callback.
 */
Binding.prototype.chmod = function(pathname, mode, callback) {
  maybeCallback(normalizeCallback(callback), this, function() {
    var item = this._system.getItem(pathname);
    if (!item) {
      throw new FSError('ENOENT', pathname);
    }
    item.setMode(mode);
  });
};

/**
 * Change permissions.
 * @param {number} fd File descriptor.
 * @param {number} mode Mode.
 * @param {function(Error)} callback Optional callback.
 */
Binding.prototype.fchmod = function(fd, mode, callback) {
  maybeCallback(normalizeCallback(callback), this, function() {
    var descriptor = this._getDescriptorById(fd);
    var item = descriptor.getItem();
    item.setMode(mode);
  });
};

/**
 * Delete a named item.
 * @param {string} pathname Path to item.
 * @param {function(Error)} callback Optional callback.
 */
Binding.prototype.unlink = function(pathname, callback) {
  maybeCallback(normalizeCallback(callback), this, function() {
    var item = this._system.getItem(pathname);
    if (!item) {
      throw new FSError('ENOENT', pathname);
    }
    if (item instanceof Directory) {
      throw new FSError('EPERM', pathname);
    }
    var parent = this._system.getItem(path.dirname(pathname));
    parent.removeItem(path.basename(pathname));
  });
};

/**
 * Update timestamps.
 * @param {string} pathname Path to item.
 * @param {number} atime Access time (in seconds).
 * @param {number} mtime Modification time (in seconds).
 * @param {function(Error)} callback Optional callback.
 */
Binding.prototype.utimes = function(pathname, atime, mtime, callback) {
  maybeCallback(normalizeCallback(callback), this, function() {
    var item = this._system.getItem(pathname);
    if (!item) {
      throw new FSError('ENOENT', pathname);
    }
    item.setATime(new Date(atime * 1000));
    item.setMTime(new Date(mtime * 1000));
  });
};

/**
 * Update timestamps.
 * @param {number} fd File descriptor.
 * @param {number} atime Access time (in seconds).
 * @param {number} mtime Modification time (in seconds).
 * @param {function(Error)} callback Optional callback.
 */
Binding.prototype.futimes = function(fd, atime, mtime, callback) {
  maybeCallback(normalizeCallback(callback), this, function() {
    var descriptor = this._getDescriptorById(fd);
    var item = descriptor.getItem();
    item.setATime(new Date(atime * 1000));
    item.setMTime(new Date(mtime * 1000));
  });
};

/**
 * Synchronize in-core state with storage device.
 * @param {number} fd File descriptor.
 * @param {function(Error)} callback Optional callback.
 */
Binding.prototype.fsync = function(fd, callback) {
  maybeCallback(normalizeCallback(callback), this, function() {
    this._getDescriptorById(fd);
  });
};

/**
 * Synchronize in-core metadata state with storage device.
 * @param {number} fd File descriptor.
 * @param {function(Error)} callback Optional callback.
 */
Binding.prototype.fdatasync = function(fd, callback) {
  maybeCallback(normalizeCallback(callback), this, function() {
    this._getDescriptorById(fd);
  });
};

/**
 * Create a hard link.
 * @param {string} srcPath The existing file.
 * @param {string} destPath The new link to create.
 * @param {function(Error)} callback Optional callback.
 */
Binding.prototype.link = function(srcPath, destPath, callback) {
  maybeCallback(normalizeCallback(callback), this, function() {
    var item = this._system.getItem(srcPath);
    if (!item) {
      throw new FSError('ENOENT', srcPath);
    }
    if (item instanceof Directory) {
      throw new FSError('EPERM', srcPath);
    }
    if (this._system.getItem(destPath)) {
      throw new FSError('EEXIST', destPath);
    }
    var parent = this._system.getItem(path.dirname(destPath));
    if (!parent) {
      throw new FSError('ENOENT', destPath);
    }
    if (!(parent instanceof Directory)) {
      throw new FSError('ENOTDIR', destPath);
    }
    parent.addItem(path.basename(destPath), item);
  });
};

/**
 * Create a symbolic link.
 * @param {string} srcPath Path from link to the source file.
 * @param {string} destPath Path for the generated link.
 * @param {string} type Ignored (used for Windows only).
 * @param {function(Error)} callback Optional callback.
 */
Binding.prototype.symlink = function(srcPath, destPath, type, callback) {
  maybeCallback(normalizeCallback(callback), this, function() {
    if (this._system.getItem(destPath)) {
      throw new FSError('EEXIST', destPath);
    }
    var parent = this._system.getItem(path.dirname(destPath));
    if (!parent) {
      throw new FSError('ENOENT', destPath);
    }
    if (!(parent instanceof Directory)) {
      throw new FSError('ENOTDIR', destPath);
    }
    var link = new SymbolicLink();
    link.setPath(srcPath);
    parent.addItem(path.basename(destPath), link);
  });
};

/**
 * Read the contents of a symbolic link.
 * @param {string} pathname Path to symbolic link.
 * @param {string} encoding The encoding ('utf-8' or 'buffer').
 * @param {function(Error, (string|Buffer))} callback Optional callback.
 * @return {string|Buffer} Symbolic link contents (path to source).
 */
Binding.prototype.readlink = function(pathname, encoding, callback) {
  if (encoding && typeof encoding !== 'string') {
    callback = encoding;
    encoding = 'utf-8';
  }
  return maybeCallback(normalizeCallback(callback), this, function() {
    var link = this._system.getItem(pathname);
    if (!(link instanceof SymbolicLink)) {
      throw new FSError('EINVAL', pathname);
    }
    var linkPath = link.getPath();
    if (encoding === 'buffer') {
      linkPath = new Buffer(linkPath);
    }
    return linkPath;
  });
};

/**
 * Stat an item.
 * @param {string} filepath Path.
 * @param {function(Error, Stats)|Float64Array} callback Callback (optional). In Node 7.7.0+ this will be a Float64Array
 * that should be filled with stat values.
 * @return {Stats|undefined} Stats or undefined (if sync).
 */
Binding.prototype.lstat = function(filepath, callback) {
  return maybeCallback(wrapStatsCallback(callback), this, function() {
    var item = this._system.getItem(filepath);
    if (!item) {
      throw new FSError('ENOENT', filepath);
    }
    var stats = item.getStats();

    // In Node 7.7.0+, binding.stat accepts a Float64Array as the second argument,
    // which should be filled with stat values.
    // In prior versions of Node, binding.stat simply returns a Stats instance.
    if (callback instanceof Float64Array) {
      fillStatsArray(stats, callback);
    } else {
      fillStatsArray(stats, statValues);
      return new Stats(item.getStats());
    }
  });
};

/**
 * Tests user permissions.
 * @param {string} filepath Path.
 * @param {number} mode Mode.
 * @param {function(Error)} callback Callback (optional).
 */
Binding.prototype.access = function(filepath, mode, callback) {
  maybeCallback(normalizeCallback(callback), this, function() {
    var item = this._system.getItem(filepath);
    if (!item) {
      throw new FSError('ENOENT', filepath);
    }
    if (mode && process.getuid && process.getgid) {
      var itemMode = item.getMode();
      if (item.getUid() === process.getuid()) {
        if ((itemMode & (mode * 64)) !== mode * 64) {
          throw new FSError('EACCES', filepath);
        }
      } else if (item.getGid() === process.getgid()) {
        if ((itemMode & (mode * 8)) !== mode * 8) {
          throw new FSError('EACCES', filepath);
        }
      } else {
        if ((itemMode & mode) !== mode) {
          throw new FSError('EACCES', filepath);
        }
      }
    }
  });
};

/**
 * Not yet implemented.
 * @type {function()}
 */
Binding.prototype.StatWatcher = notImplemented;

/**
 * Export the binding constructor.
 * @type {function()}
 */
exports = module.exports = Binding;
