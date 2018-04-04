'use strict';

var Binding = require('./binding');
var FSError = require('./error');
var FileSystem = require('./filesystem');
var realBinding = process.binding('fs');
var path = require('path');

var realBindingProps = Object.assign({}, realBinding);
var realProcessProps = {
  cwd: process.cwd,
  chdir: process.chdir
};

function overrideBinding(binding) {
  for (var key in binding) {
    if (typeof binding[key] === 'function') {
      realBinding[key] = binding[key].bind(binding);
    } else {
      realBinding[key] = binding[key];
    }
  }
}

function overrideProcess(cwd, chdir) {
  process.cwd = cwd;
  process.chdir = chdir;
}

function restoreBinding() {
  var key;
  for (key in realBindingProps) {
    realBinding[key] = realBindingProps[key];
  }
  // Delete excess keys that came in when the binding was originally applied.
  for (key in realBinding) {
    if (typeof realBindingProps[key] === 'undefined') {
      delete realBinding[key];
    }
  }
}

function restoreProcess() {
  for (var key in realProcessProps) {
    process[key] = realProcessProps[key];
  }
}

/**
 * Swap out the fs bindings for a mock file system.
 * @param {Object} config Mock file system configuration.
 * @param {Object} options Any filesystem options.
 * @param {boolean} options.createCwd Create a directory for `process.cwd()`
 *     (defaults to `true`).
 * @param {boolean} options.createTmp Create a directory for `os.tmpdir()`
 *     (defaults to `true`).
 */
var exports = (module.exports = function mock(config, options) {
  var system = FileSystem.create(config, options);
  var binding = new Binding(system);

  overrideBinding(binding);

  var currentPath = process.cwd();
  overrideProcess(
    function cwd() {
      return currentPath;
    },
    function chdir(directory) {
      if (!binding.stat(path._makeLong(directory)).isDirectory()) {
        throw new FSError('ENOTDIR');
      }
      currentPath = path.resolve(currentPath, directory);
    }
  );
});

/**
 * Get hold of the mocked filesystem's 'root'
 * If fs hasn't currently been replaced, this will return an empty object
 */
exports.getMockRoot = function() {
  if (typeof realBinding.getSystem === 'undefined') {
    return {};
  } else {
    return realBinding.getSystem().getRoot();
  }
};

/**
 * Restore the fs bindings for the real file system.
 */
exports.restore = function() {
  restoreBinding();
  restoreProcess();
};

/**
 * Create a file factory.
 */
exports.file = FileSystem.file;

/**
 * Create a directory factory.
 */
exports.directory = FileSystem.directory;

/**
 * Create a symbolic link factory.
 */
exports.symlink = FileSystem.symlink;
