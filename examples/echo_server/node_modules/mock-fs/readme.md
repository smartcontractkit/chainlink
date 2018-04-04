# `mock-fs`

The `mock-fs` module allows Node's built-in [`fs` module](http://nodejs.org/api/fs.html) to be backed temporarily by an in-memory, mock file system.  This lets you run tests against a set of mock files and directories instead of lugging around a bunch of test fixtures.

## Example

The code below makes it so the `fs` module is temporarily backed by a mock file system with a few files and directories.

```js
const mock = require('mock-fs');

mock({
  'path/to/fake/dir': {
    'some-file.txt': 'file content here',
    'empty-dir': {/** empty directory */}
  },
  'path/to/some.png': new Buffer([8, 6, 7, 5, 3, 0, 9]),
  'some/other/path': {/** another empty directory */}
});
```

When you are ready to restore the `fs` module (so that it is backed by your real file system), call [`mock.restore()`](#mockrestore). Note that calling this may be **mandatory** in some cases. See [istanbuljs/nyc#324](https://github.com/istanbuljs/nyc/issues/324#issuecomment-234018654)

```js
// after a test runs
mock.restore();
```

## Upgrading to version 4

Instead of overriding all methods of the built-in `fs` module, the library now overrides `process.binding('fs')`.  The purpose of this change is to avoid conflicts with other libraries that override `fs` methods (e.g. `graceful-fs`) and to make it possible to work with multiple Node releases without maintaining copied and slightly modified versions of Node's `fs` module.

Breaking changes:

 * The `mock.fs()` function has been removed.  This returned an object with `fs`-like methods without overriding the built-in `fs` module.
 * The object created by `fs.Stats` is no longer an instance of `fs.Stats` (though it has all the same properties and methods).
 * Lazy `require()` do not use the real filesystem.
 * Tests are no longer run in Node < 4.

Some of these breaking changes may be restored in a future release.

## Docs

### <a id='mockconfigoptions'>`mock(config, options)`</a>

Configure the `fs` module so it is backed by an in-memory file system.

Calling `mock` sets up a mock file system with two directories by default: `process.cwd()` and `os.tmpdir()` (or `os.tmpDir()` for older Node).  When called with no arguments, just these two directories are created.  When called with a `config` object, additional files, directories, and symlinks are created.  To avoid creating a directory for `process.cwd()` and `os.tmpdir()`, see the [`options`](#options) below.

Property names of the `config` object are interpreted as relative paths to resources (relative from `process.cwd()`).  Property values of the `config` object are interpreted as content or configuration for the generated resources.

*Note that paths should always use forward slashes (`/`) - even on Windows.*

### <a id='options'>`options`</a>

The second (optional) argument may include the properties below.

 * `createCwd` - `boolean` Create a directory for `process.cwd()`.  This is `true` by default.
 * `createTmp` - `boolean` Create a directory for `os.tmpdir()`.  This is `true` by default.

### Creating files

When `config` property values are a `string` or `Buffer`, a file is created with the provided content.  For example, the following configuration creates a single file with string content (in addition to the two default directories).
```js
mock({
  'path/to/file.txt': 'file content here'
});
```

To create a file with additional properties (owner, permissions, atime, etc.), use the [`mock.file()`](#mockfileproperties) function described below.

### <a id='mockfileproperties'>`mock.file(properties)`</a>

Create a factory for new files.  Supported properties:

 * **content** - `string|Buffer` File contents.
 * **mode** - `number` File mode (permission and sticky bits).  Defaults to `0666`.
 * **uid** - `number` The user id.  Defaults to `process.getuid()`.
 * **gid** - `number` The group id.  Defaults to `process.getgid()`.
 * **atime** - `Date` The last file access time.  Defaults to `new Date()`.  Updated when file contents are accessed.
 * **ctime** - `Date` The last file change time.  Defaults to `new Date()`.  Updated when file owner or permissions change.
 * **mtime** - `Date` The last file modification time.  Defaults to `new Date()`.  Updated when file contents change.
 * **birthtime** - `Date` The time of file creation.  Defaults to `new Date()`.

To create a mock filesystem with a very old file named `foo`, you could do something like this:
```js
mock({
  foo: mock.file({
    content: 'file content here',
    ctime: new Date(1),
    mtime: new Date(1)
  })
});
```

Note that if you want to create a file with the default properties, you can provide a `string` or `Buffer` directly instead of calling `mock.file()`.

### Creating directories

When `config` property values are an `Object`, a directory is created.  The structure of the object is the same as the `config` object itself.  So an empty directory can be created with a simple object literal (`{}`).  The following configuration creates a directory containing two files (in addition to the two default directories):
```js
// note that this could also be written as
// mock({'path/to/dir': { /** config */ }})
mock({
  path: {
    to: {
      dir: {
        file1: 'text content',
        file2: new Buffer([1, 2, 3, 4])
      }
    }
  }
});
```

To create a directory with additional properties (owner, permissions, atime, etc.), use the [`mock.directory()`](mockdirectoryproperties) function described below.

### <a id='mockdirectoryproperties'>`mock.directory(properties)`</a>

Create a factory for new directories.  Supported properties:

 * **mode** - `number` Directory mode (permission and sticky bits).  Defaults to `0777`.
 * **uid** - `number` The user id.  Defaults to `process.getuid()`.
 * **gid** - `number` The group id.  Defaults to `process.getgid()`.
 * **atime** - `Date` The last directory access time.  Defaults to `new Date()`.
 * **ctime** - `Date` The last directory change time.  Defaults to `new Date()`.  Updated when owner or permissions change.
 * **mtime** - `Date` The last directory modification time.  Defaults to `new Date()`.  Updated when an item is added, removed, or renamed.
 * **birthtime** - `Date` The time of directory creation.  Defaults to `new Date()`.
 * **items** - `Object` Directory contents.  Members will generate additional files, directories, or symlinks.

To create a mock filesystem with a directory with the relative path `some/dir` that has a mode of `0755` and two child files, you could do something like this:
```js
mock({
  'some/dir': mock.directory({
    mode: 0755,
    items: {
      file1: 'file one content',
      file2: new Buffer([8, 6, 7, 5, 3, 0, 9])
    }
  })
});
```

Note that if you want to create a directory with the default properties, you can provide an `Object` directly instead of calling `mock.directory()`.

### Creating symlinks

Using a `string` or a `Buffer` is a shortcut for creating files with default properties.  Using an `Object` is a shortcut for creating a directory with default properties.  There is no shortcut for creating symlinks.  To create a symlink, you need to call the [`mock.symlink()`](#mocksymlinkproperties) function described below.

### <a id='mocksymlinkproperties'>`mock.symlink(properties)`</a>

Create a factory for new symlinks.  Supported properties:

 * **path** - `string` Path to the source (required).
 * **mode** - `number` Symlink mode (permission and sticky bits).  Defaults to `0666`.
 * **uid** - `number` The user id.  Defaults to `process.getuid()`.
 * **gid** - `number` The group id.  Defaults to `process.getgid()`.
 * **atime** - `Date` The last symlink access time.  Defaults to `new Date()`.
 * **ctime** - `Date` The last symlink change time.  Defaults to `new Date()`.
 * **mtime** - `Date` The last symlink modification time.  Defaults to `new Date()`.
 * **birthtime** - `Date` The time of symlink creation.  Defaults to `new Date()`.

To create a mock filesystem with a file and a symlink, you could do something like this:
```js
mock({
  'some/dir': {
    'regular-file': 'file contents',
    'a-symlink': mock.symlink({
      path: 'regular-file'
    })
  }
});
```

### Restoring the file system

### <a id='mockrestore'>`mock.restore()`</a>

Restore the `fs` binding to the real file system.  This undoes the effect of calling `mock()`.  Typically, you would set up a mock file system before running a test and restore the original after.  Using a test runner with `beforeEach` and `afterEach` hooks, this might look like the following:

```js
beforeEach(function() {
  mock({
    'fake-file': 'file contents'
  });
});
afterEach(mock.restore);
```

## Install

Using `npm`:

```
npm install mock-fs --save-dev
```

## Caveats

When you require `mock-fs`, Node's own `fs` module is patched to allow the binding to the underlying file system to be swapped out.  If you require `mock-fs` *before* any other modules that modify `fs` (e.g. `graceful-fs`), the mock should behave as expected.

**Note** `mock-fs` is not compatible with `graceful-fs@3.x` but works with `graceful-fs@4.x`.

Mock `fs.Stats` objects have the following properties: `dev`, `ino`, `nlink`, `mode`, `size`, `rdev`, `blksize`, `blocks`, `atime`, `ctime`, `mtime`, `birthtime`, `uid`, and `gid`.  In addition, all of the `is*()` method are provided (e.g. `isDirectory()`, `isFile()`, et al.).

Mock file access is controlled based on file mode where `process.getuid()` and `process.getgid()` are available (POSIX systems).  On other systems (e.g. Windows) the file mode has no effect.

Tested on Linux, OSX, and Windows using Node 0.10 through 8.x.  Check the tickets for a list of [known issues](https://github.com/tschaub/mock-fs/issues).

[![Current Status](https://secure.travis-ci.org/tschaub/mock-fs.png?branch=master)](https://travis-ci.org/tschaub/mock-fs)
