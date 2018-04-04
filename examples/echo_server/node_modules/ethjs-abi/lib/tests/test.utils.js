'use strict';

var assert = require('chai').assert;
var utils = require('../utils/index.js');
var BN = require('bn.js');

describe('test utilty methods', function () {
  describe('stripZeros', function () {
    it('should strip zeros from buffer', function () {
      assert.deepEqual(utils.stripZeros(new Buffer('000001', 'hex')), new Buffer('01', 'hex'));
    });

    it('should throw while stripping String', function () {
      assert.equal(utils.stripZeros('0000000'), '');
    });
  });

  describe('bnToBuffer', function () {
    it('bnToBuffer', function () {
      assert.deepEqual(utils.bnToBuffer(new BN(10)), new Buffer('0' + new BN(10).toString(16), 'hex'));
    });
  });

  describe('getKeys', function () {
    it('invalid type', function () {
      assert.throws(function () {
        return utils.getKeys(undefined);
      }, Error);
    });

    it('invalid abi', function () {
      assert.throws(function () {
        return utils.getKeys([{ type: 239823 }]);
      }, Error);
    });
  });

  describe('isHexString', function () {
    it('should detect invalid hex string length', function () {
      assert.deepEqual(utils.isHexString('0x3a', 5), false);
    });
  });

  describe('getParamCoder', function () {
    it('should throw when invalid sol type uint[][][]', function () {
      assert.throws(function () {
        return utils.getParamCoder('uint[][][]');
      }, Error);
    });

    it('should throw when invalid sol type null', function () {
      assert.throws(function () {
        return utils.getParamCoder(null);
      }, Error);
    });

    it('should throw when invalid sol type booluint', function () {
      assert.throws(function () {
        return utils.getParamCoder('booluint');
      }, Error);
    });

    it('should throw when invalid sol type stringint2', function () {
      assert.throws(function () {
        return utils.getParamCoder('stringint2');
      }, Error);
    });

    it('should throw when invalid sol type uint[==-3-]', function () {
      assert.throws(function () {
        return utils.getParamCoder('real');
      }, Error);
    });

    it('should throw when invalid sol type uint8bool', function () {
      assert.throws(function () {
        return utils.getParamCoder('uint8bool');
      }, Error);
    });

    it('should throw when invalid sol type uint8string', function () {
      assert.throws(function () {
        return utils.getParamCoder('uint8string');
      }, Error);
    });

    it('should throw when invalid sol type uint8bytes', function () {
      assert.throws(function () {
        return utils.getParamCoder('uint8bytes');
      }, Error);
    });

    it('should throw when invalid sol type uint64dskfjk', function () {
      assert.throws(function () {
        return utils.getParamCoder('uint64uint2uint4');
      }, Error);
    });

    it('should throw when invalid sol type bytes64', function () {
      assert.throws(function () {
        return utils.getParamCoder('bytes64');
      }, Error);
    });

    it('should throw when invalid sol type uint8address', function () {
      assert.throws(function () {
        return utils.getParamCoder('uint8address');
      }, Error);
    });

    it('should throw when invalid sol type false', function () {
      assert.throws(function () {
        return utils.getParamCoder('false');
      }, Error);
    });
  });

  describe('coderAddress', function () {
    it('not valid coder address', function () {
      assert.throws(function () {
        return utils.coderAddress.encode('sfdsfd');
      }, Error);
    });

    it('not valid coder address', function () {
      assert.throws(function () {
        return utils.coderAddress.decode('sfdsfd', 40);
      }, Error);
    });

    it('should decode nicely', function () {
      assert.deepEqual(utils.coderAddress.decode([], 40), {
        consumed: 32,
        value: '0x'
      });
    });
  });

  describe('coderFixedBytes', function () {
    it('not valid coder bytes', function () {
      assert.throws(function () {
        return utils.coderFixedBytes(10).decode('sfdsfd', 40);
      }, Error);
    });
  });

  describe('coderDynamicBytes', function () {
    it('invalid dynamic bytes decode should throw', function () {
      assert.throws(function () {
        return utils.coderDynamicBytes.decode('2', 40);
      }, Error);
    });

    it('invalid dynamic bytes decode should throw', function () {
      assert.throws(function () {
        return utils.coderDynamicBytes.decode('asddsdfsfd', 0);
      }, Error);
    });

    it('invalid dynamic bytes decode should throw', function () {
      assert.throws(function () {
        utils.coderDynamicBytes.decode('0a', 5000);
      }, Error);
    });
  });

  describe('coderArray', function () {
    it('should throw when coder array encode is fed non array', function () {
      assert.throws(function () {
        return utils.coderArray('uint', 6).encode(243);
      }, Error);
    });
  });

  describe('hexOrBuffer', function () {
    it('uneven hex should pad', function () {
      assert.deepEqual(utils.hexOrBuffer('0xa'), new Buffer('0a', 'hex'));
    });
  });

  describe('hexlify', function () {
    it('valid number', function () {
      assert.deepEqual(utils.hexlify(10), '0x0a');
    });

    it('valid bignumber', function () {
      assert.deepEqual(utils.hexlify(new BN(10)), '0x0a');
    });

    it('valid hex', function () {
      assert.deepEqual(utils.hexlify('0x0a'), '0x0a');
    });

    it('invalid hex', function () {
      assert.throws(function () {
        return utils.hexlify([]);
      }, Error);
    });
  });
});