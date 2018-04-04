'use strict';

var counter = 0;

/**
 * Permissions.
 * @enum {number}
 */
var permissions = {
  USER_READ: 256, // 0400
  USER_WRITE: 128, // 0200
  USER_EXEC: 64, // 0100
  GROUP_READ: 32, // 0040
  GROUP_WRITE: 16, // 0020
  GROUP_EXEC: 8, // 0010
  OTHER_READ: 4, // 0004
  OTHER_WRITE: 2, // 0002
  OTHER_EXEC: 1 // 0001
};

function getUid() {
  return process.getuid && process.getuid();
}

function getGid() {
  return process.getgid && process.getgid();
}

/**
 * A filesystem item.
 * @constructor
 */
function Item() {
  var now = Date.now();

  /**
   * Access time.
   * @type {Date}
   */
  this._atime = new Date(now);

  /**
   * Change time.
   * @type {Date}
   */
  this._ctime = new Date(now);

  /**
   * Birth time.
   * @type {Date}
   */
  this._birthtime = new Date(now);

  /**
   * Modification time.
   * @type {Date}
   */
  this._mtime = new Date(now);

  /**
   * Permissions.
   */
  this._mode = 438; // 0666

  /**
   * User id.
   * @type {number}
   */
  this._uid = getUid();

  /**
   * Group id.
   * @type {number}
   */
  this._gid = getGid();

  /**
   * Item number.
   * @type {number}
   */
  this._id = ++counter;

  /**
   * Number of links to this item.
   */
  this.links = 0;
}

/**
 * Determine if the current user has read permission.
 * @return {boolean} The current user can read.
 */
Item.prototype.canRead = function() {
  var uid = getUid();
  var gid = getGid();
  var can = false;
  if (uid === 0) {
    can = true;
  } else if (uid === this._uid) {
    can = (permissions.USER_READ & this._mode) === permissions.USER_READ;
  } else if (gid === this._gid) {
    can = (permissions.GROUP_READ & this._mode) === permissions.GROUP_READ;
  } else {
    can = (permissions.OTHER_READ & this._mode) === permissions.OTHER_READ;
  }
  return can;
};

/**
 * Determine if the current user has write permission.
 * @return {boolean} The current user can write.
 */
Item.prototype.canWrite = function() {
  var uid = getUid();
  var gid = getGid();
  var can = false;
  if (uid === 0) {
    can = true;
  } else if (uid === this._uid) {
    can = (permissions.USER_WRITE & this._mode) === permissions.USER_WRITE;
  } else if (gid === this._gid) {
    can = (permissions.GROUP_WRITE & this._mode) === permissions.GROUP_WRITE;
  } else {
    can = (permissions.OTHER_WRITE & this._mode) === permissions.OTHER_WRITE;
  }
  return can;
};

/**
 * Determine if the current user has execute permission.
 * @return {boolean} The current user can execute.
 */
Item.prototype.canExecute = function() {
  var uid = getUid();
  var gid = getGid();
  var can = false;
  if (uid === 0) {
    can = true;
  } else if (uid === this._uid) {
    can = (permissions.USER_EXEC & this._mode) === permissions.USER_EXEC;
  } else if (gid === this._gid) {
    can = (permissions.GROUP_EXEC & this._mode) === permissions.GROUP_EXEC;
  } else {
    can = (permissions.OTHER_EXEC & this._mode) === permissions.OTHER_EXEC;
  }
  return can;
};

/**
 * Get access time.
 * @return {Date} Access time.
 */
Item.prototype.getATime = function() {
  return this._atime;
};

/**
 * Set access time.
 * @param {Date} atime Access time.
 */
Item.prototype.setATime = function(atime) {
  this._atime = atime;
};

/**
 * Get change time.
 * @return {Date} Change time.
 */
Item.prototype.getCTime = function() {
  return this._ctime;
};

/**
 * Set change time.
 * @param {Date} ctime Change time.
 */
Item.prototype.setCTime = function(ctime) {
  this._ctime = ctime;
};

/**
 * Get birth time.
 * @return {Date} Birth time.
 */
Item.prototype.getBirthtime = function() {
  return this._birthtime;
};

/**
 * Set change time.
 * @param {Date} birthtime Birth time.
 */
Item.prototype.setBirthtime = function(birthtime) {
  this._birthtime = birthtime;
};

/**
 * Get modification time.
 * @return {Date} Modification time.
 */
Item.prototype.getMTime = function() {
  return this._mtime;
};

/**
 * Set modification time.
 * @param {Date} mtime Modification time.
 */
Item.prototype.setMTime = function(mtime) {
  this._mtime = mtime;
};

/**
 * Get mode (permission only, e.g 0666).
 * @return {number} Mode.
 */
Item.prototype.getMode = function() {
  return this._mode;
};

/**
 * Set mode (permission only, e.g 0666).
 * @param {Date} mode Mode.
 */
Item.prototype.setMode = function(mode) {
  this.setCTime(new Date());
  this._mode = mode;
};

/**
 * Get user id.
 * @return {number} User id.
 */
Item.prototype.getUid = function() {
  return this._uid;
};

/**
 * Set user id.
 * @param {number} uid User id.
 */
Item.prototype.setUid = function(uid) {
  this.setCTime(new Date());
  this._uid = uid;
};

/**
 * Get group id.
 * @return {number} Group id.
 */
Item.prototype.getGid = function() {
  return this._gid;
};

/**
 * Set group id.
 * @param {number} gid Group id.
 */
Item.prototype.setGid = function(gid) {
  this.setCTime(new Date());
  this._gid = gid;
};

/**
 * Get item stats.
 * @return {Object} Stats properties.
 */
Item.prototype.getStats = function() {
  return {
    dev: 8675309,
    nlink: this.links,
    uid: this.getUid(),
    gid: this.getGid(),
    rdev: 0,
    blksize: 4096,
    ino: this._id,
    atime: this.getATime(),
    mtime: this.getMTime(),
    ctime: this.getCTime(),
    birthtime: this.getBirthtime()
  };
};

/**
 * Get the item's string representation.
 * @return {string} String representation.
 */
Item.prototype.toString = function() {
  return '[' + this.constructor.name + ']';
};

/**
 * Export the constructor.
 * @type {function()}
 */
exports = module.exports = Item;
