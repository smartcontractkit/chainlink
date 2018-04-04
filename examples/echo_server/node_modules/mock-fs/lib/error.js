'use strict';

/**
 * Error codes from libuv.
 * @enum {number}
 */
var codes = {
  UNKNOWN: {
    errno: -1,
    message: 'unknown error'
  },
  OK: {
    errno: 0,
    message: 'success'
  },
  EOF: {
    errno: 1,
    message: 'end of file'
  },
  EADDRINFO: {
    errno: 2,
    message: 'getaddrinfo error'
  },
  EACCES: {
    errno: 3,
    message: 'permission denied'
  },
  EAGAIN: {
    errno: 4,
    message: 'resource temporarily unavailable'
  },
  EADDRINUSE: {
    errno: 5,
    message: 'address already in use'
  },
  EADDRNOTAVAIL: {
    errno: 6,
    message: 'address not available'
  },
  EAFNOSUPPORT: {
    errno: 7,
    message: 'address family not supported'
  },
  EALREADY: {
    errno: 8,
    message: 'connection already in progress'
  },
  EBADF: {
    errno: 9,
    message: 'bad file descriptor'
  },
  EBUSY: {
    errno: 10,
    message: 'resource busy or locked'
  },
  ECONNABORTED: {
    errno: 11,
    message: 'software caused connection abort'
  },
  ECONNREFUSED: {
    errno: 12,
    message: 'connection refused'
  },
  ECONNRESET: {
    errno: 13,
    message: 'connection reset by peer'
  },
  EDESTADDRREQ: {
    errno: 14,
    message: 'destination address required'
  },
  EFAULT: {
    errno: 15,
    message: 'bad address in system call argument'
  },
  EHOSTUNREACH: {
    errno: 16,
    message: 'host is unreachable'
  },
  EINTR: {
    errno: 17,
    message: 'interrupted system call'
  },
  EINVAL: {
    errno: 18,
    message: 'invalid argument'
  },
  EISCONN: {
    errno: 19,
    message: 'socket is already connected'
  },
  EMFILE: {
    errno: 20,
    message: 'too many open files'
  },
  EMSGSIZE: {
    errno: 21,
    message: 'message too long'
  },
  ENETDOWN: {
    errno: 22,
    message: 'network is down'
  },
  ENETUNREACH: {
    errno: 23,
    message: 'network is unreachable'
  },
  ENFILE: {
    errno: 24,
    message: 'file table overflow'
  },
  ENOBUFS: {
    errno: 25,
    message: 'no buffer space available'
  },
  ENOMEM: {
    errno: 26,
    message: 'not enough memory'
  },
  ENOTDIR: {
    errno: 27,
    message: 'not a directory'
  },
  EISDIR: {
    errno: 28,
    message: 'illegal operation on a directory'
  },
  ENONET: {
    errno: 29,
    message: 'machine is not on the network'
  },
  ENOTCONN: {
    errno: 31,
    message: 'socket is not connected'
  },
  ENOTSOCK: {
    errno: 32,
    message: 'socket operation on non-socket'
  },
  ENOTSUP: {
    errno: 33,
    message: 'operation not supported on socket'
  },
  ENOENT: {
    errno: 34,
    message: 'no such file or directory'
  },
  ENOSYS: {
    errno: 35,
    message: 'function not implemented'
  },
  EPIPE: {
    errno: 36,
    message: 'broken pipe'
  },
  EPROTO: {
    errno: 37,
    message: 'protocol error'
  },
  EPROTONOSUPPORT: {
    errno: 38,
    message: 'protocol not supported'
  },
  EPROTOTYPE: {
    errno: 39,
    message: 'protocol wrong type for socket'
  },
  ETIMEDOUT: {
    errno: 40,
    message: 'connection timed out'
  },
  ECHARSET: {
    errno: 41,
    message: 'invalid Unicode character'
  },
  EAIFAMNOSUPPORT: {
    errno: 42,
    message: 'address family for hostname not supported'
  },
  EAISERVICE: {
    errno: 44,
    message: 'servname not supported for ai_socktype'
  },
  EAISOCKTYPE: {
    errno: 45,
    message: 'ai_socktype not supported'
  },
  ESHUTDOWN: {
    errno: 46,
    message: 'cannot send after transport endpoint shutdown'
  },
  EEXIST: {
    errno: 47,
    message: 'file already exists'
  },
  ESRCH: {
    errno: 48,
    message: 'no such process'
  },
  ENAMETOOLONG: {
    errno: 49,
    message: 'name too long'
  },
  EPERM: {
    errno: 50,
    message: 'operation not permitted'
  },
  ELOOP: {
    errno: 51,
    message: 'too many symbolic links encountered'
  },
  EXDEV: {
    errno: 52,
    message: 'cross-device link not permitted'
  },
  ENOTEMPTY: {
    errno: 53,
    message: 'directory not empty'
  },
  ENOSPC: {
    errno: 54,
    message: 'no space left on device'
  },
  EIO: {
    errno: 55,
    message: 'i/o error'
  },
  EROFS: {
    errno: 56,
    message: 'read-only file system'
  },
  ENODEV: {
    errno: 57,
    message: 'no such device'
  },
  ESPIPE: {
    errno: 58,
    message: 'invalid seek'
  },
  ECANCELED: {
    errno: 59,
    message: 'peration canceled'
  }
};

/**
 * Create an error.
 * @param {string} code Error code.
 * @param {string} path Path (optional).
 * @constructor
 */
function FSError(code, path) {
  if (!codes.hasOwnProperty(code)) {
    throw new Error('Programmer error, invalid error code: ' + code);
  }
  Error.call(this);
  var details = codes[code];
  var message = code + ', ' + details.message;
  if (path) {
    message += " '" + path + "'";
  }
  this.message = message;
  this.code = code;
  this.errno = details.errno;
  Error.captureStackTrace(this, FSError);
}
FSError.prototype = new Error();

/**
 * Error constructor.
 */
exports = module.exports = FSError;
