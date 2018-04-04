const assert = require('chai').assert;
const utils = require('../utils/index.js');
const BN = require('bn.js');

describe('test utilty methods', () => {
  describe('stripZeros', () => {
    it('should strip zeros from buffer', () => {
      assert.deepEqual(utils.stripZeros(new Buffer('000001', 'hex')), new Buffer('01', 'hex'));
    });

    it('should throw while stripping String', () => {
      assert.equal(utils.stripZeros('0000000'), '');
    });
  });

  describe('bnToBuffer', () => {
    it('bnToBuffer', () => {
      assert.deepEqual(utils.bnToBuffer(new BN(10)), (new Buffer(`0${(new BN(10)).toString(16)}`, 'hex')));
    });
  });

  describe('getKeys', () => {
    it('invalid type', () => {
      assert.throws(() => utils.getKeys(undefined), Error);
    });

    it('invalid abi', () => {
      assert.throws(() => utils.getKeys([{ type: 239823 }]), Error);
    });
  });

  describe('isHexString', () => {
    it('should detect invalid hex string length', () => {
      assert.deepEqual(utils.isHexString('0x3a', 5), false);
    });
  });

  describe('getParamCoder', () => {
    it('should throw when invalid sol type uint[][][]', () => {
      assert.throws(() => utils.getParamCoder('uint[][][]'), Error);
    });

    it('should throw when invalid sol type null', () => {
      assert.throws(() => utils.getParamCoder(null), Error);
    });

    it('should throw when invalid sol type booluint', () => {
      assert.throws(() => utils.getParamCoder('booluint'), Error);
    });

    it('should throw when invalid sol type stringint2', () => {
      assert.throws(() => utils.getParamCoder('stringint2'), Error);
    });

    it('should throw when invalid sol type uint[==-3-]', () => {
      assert.throws(() => utils.getParamCoder('real'), Error);
    });

    it('should throw when invalid sol type uint8bool', () => {
      assert.throws(() => utils.getParamCoder('uint8bool'), Error);
    });

    it('should throw when invalid sol type uint8string', () => {
      assert.throws(() => utils.getParamCoder('uint8string'), Error);
    });

    it('should throw when invalid sol type uint8bytes', () => {
      assert.throws(() => utils.getParamCoder('uint8bytes'), Error);
    });

    it('should throw when invalid sol type uint64dskfjk', () => {
      assert.throws(() => utils.getParamCoder('uint64uint2uint4'), Error);
    });

    it('should throw when invalid sol type bytes64', () => {
      assert.throws(() => utils.getParamCoder('bytes64'), Error);
    });

    it('should throw when invalid sol type uint8address', () => {
      assert.throws(() => utils.getParamCoder('uint8address'), Error);
    });

    it('should throw when invalid sol type false', () => {
      assert.throws(() => utils.getParamCoder('false'), Error);
    });
  });

  describe('coderAddress', () => {
    it('not valid coder address', () => {
      assert.throws(() => utils.coderAddress.encode('sfdsfd'), Error);
    });

    it('not valid coder address', () => {
      assert.throws(() => utils.coderAddress.decode('sfdsfd', 40), Error);
    });

    it('should decode nicely', () => {
      assert.deepEqual(utils.coderAddress.decode([], 40), {
        consumed: 32,
        value: '0x',
      });
    });
  });

  describe('coderFixedBytes', () => {
    it('not valid coder bytes', () => {
      assert.throws(() => utils.coderFixedBytes(10).decode('sfdsfd', 40), Error);
    });
  });

  describe('coderDynamicBytes', () => {
    it('invalid dynamic bytes decode should throw', () => {
      assert.throws(() => utils.coderDynamicBytes.decode('2', 40), Error);
    });

    it('invalid dynamic bytes decode should throw', () => {
      assert.throws(() => utils.coderDynamicBytes.decode('asddsdfsfd', 0), Error);
    });

    it('invalid dynamic bytes decode should throw', () => {
      assert.throws(() => { utils.coderDynamicBytes.decode('0a', 5000); }, Error);
    });
  });

  describe('coderArray', () => {
    it('should throw when coder array encode is fed non array', () => {
      assert.throws(() => utils.coderArray('uint', 6).encode(243), Error);
    });
  });

  describe('hexOrBuffer', () => {
    it('uneven hex should pad', () => {
      assert.deepEqual(utils.hexOrBuffer('0xa'), new Buffer('0a', 'hex'));
    });
  });

  describe('hexlify', () => {
    it('valid number', () => {
      assert.deepEqual(utils.hexlify(10), '0x0a');
    });

    it('valid bignumber', () => {
      assert.deepEqual(utils.hexlify(new BN(10)), '0x0a');
    });

    it('valid hex', () => {
      assert.deepEqual(utils.hexlify('0x0a'), '0x0a');
    });

    it('invalid hex', () => {
      assert.throws(() => utils.hexlify([]), Error);
    });
  });
});
