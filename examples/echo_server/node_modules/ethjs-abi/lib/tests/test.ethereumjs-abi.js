'use strict';

/* eslint-disable */

var assert = require('chai').assert;
// var abi = require('../index.js');
var BN = require('bn.js');
var BigNumber = require('bignumber.js');
var encodeParams = require('../index.js').encodeParams;
var decodeParams = require('../index.js').decodeParams;

// Official test vectors from https://github.com/ethereum/wiki/wiki/Ethereum-Contract-ABI

/*
describe('official test vector 1 (encoding)', function () {
  it('should equal', function () {
    var a = abi.methodID('baz', [ 'uint32', 'bool' ]).toString('hex') + abi.rawEncode([ 'uint32', 'bool' ], [ 69, 1 ]).toString('hex')
    var b = 'cdcd77c000000000000000000000000000000000000000000000000000000000000000450000000000000000000000000000000000000000000000000000000000000001'
    assert.equal(a, b)
  })
})

describe('official test vector 3 (encoding)', function () {
  it('should equal', function () {
    var a = abi.methodID('sam', [ 'bytes', 'bool', 'uint256[]' ]).toString('hex') + abi.rawEncode([ 'bytes', 'bool', 'uint256[]' ], [ 'dave', true, [ 1, 2, 3 ] ]).toString('hex')
    var b = 'a5643bf20000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000000464617665000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000003'
    assert.equal(a, b)
  })
})

describe('official test vector 4 (encoding)', function () {
  it('should equal', function () {
    var a = abi.methodID('f', [ 'uint', 'uint32[]', 'bytes10', 'bytes' ]).toString('hex') + abi.rawEncode([ 'uint', 'uint32[]', 'bytes10', 'bytes' ], [ 0x123, [ 0x456, 0x789 ], '1234567890', 'Hello, world!' ]).toString('hex')
    var b = '8be6524600000000000000000000000000000000000000000000000000000000000001230000000000000000000000000000000000000000000000000000000000000080313233343536373839300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e0000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000004560000000000000000000000000000000000000000000000000000000000000789000000000000000000000000000000000000000000000000000000000000000d48656c6c6f2c20776f726c642100000000000000000000000000000000000000'
    assert.equal(a, b)
  })
})

// Homebrew tests

describe('method signature', function () {
  it('should work with test()', function () {
    assert.equal(abi.methodID('test', []).toString('hex'), 'f8a8fd6d')
  })
  it('should work with test(uint)', function () {
    assert.equal(abi.methodID('test', [ 'uint' ]).toString('hex'), '29e99f07')
  })
  it('should work with test(uint256)', function () {
    assert.equal(abi.methodID('test', [ 'uint256' ]).toString('hex'), '29e99f07')
  })
  it('should work with test(uint, uint)', function () {
    assert.equal(abi.methodID('test', [ 'uint', 'uint' ]).toString('hex'), 'eb8ac921')
  })
})

describe('event signature', function () {
  it('should work with test()', function () {
    assert.equal(abi.eventID('test', []).toString('hex'), 'f8a8fd6dd9544ca87214e80c840685bd13ff4682cacb0c90821ed74b1d248926')
  })
  it('should work with test(uint)', function () {
    assert.equal(abi.eventID('test', [ 'uint' ]).toString('hex'), '29e99f07d14aa8d30a12fa0b0789b43183ba1bf6b4a72b95459a3e397cca10d7')
  })
  it('should work with test(uint256)', function () {
    assert.equal(abi.eventID('test', [ 'uint256' ]).toString('hex'), '29e99f07d14aa8d30a12fa0b0789b43183ba1bf6b4a72b95459a3e397cca10d7')
  })
  it('should work with test(uint, uint)', function () {
    assert.equal(abi.eventID('test', [ 'uint', 'uint' ]).toString('hex'), 'eb8ac9210327650aab0044de896b150391af3be06f43d0f74c01f05633b97a70')
  })
})
*/

function rawEncode() {
  var args = [].slice.call(arguments);

  return encodeParams(args[0], args[1]).slice(2);
}

function rawDecode() {
  var args = [].slice.call(arguments);

  return decodeParams(args[0], args[1]);
}

var abi = {
  rawEncode: rawEncode,
  rawDecode: rawDecode
};

describe('encoding negative int32', function () {
  it('should equal', function () {
    var a = abi.rawEncode(['int32'], [-2]).toString('hex');
    var b = 'fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe';
    assert.equal(a, b);
  });
});

describe('encoding negative int256', function () {
  it('should equal', function () {
    var a = abi.rawEncode(['int256'], [new BigNumber('-19999999999999999999999999999999999999999999999999999999999999', 10)]).toString('hex');
    var b = 'fffffffffffff38dd0f10627f5529bdb2c52d4846810af0ac000000000000001';
    assert.equal(a, b);
  });
});

describe('encoding string >32bytes', function () {
  it('should equal', function () {
    var a = abi.rawEncode(['string'], [' hello world hello world hello world hello world  hello world hello world hello world hello world  hello world hello world hello world hello world hello world hello world hello world hello world']).toString('hex');
    var b = '000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000c22068656c6c6f20776f726c642068656c6c6f20776f726c642068656c6c6f20776f726c642068656c6c6f20776f726c64202068656c6c6f20776f726c642068656c6c6f20776f726c642068656c6c6f20776f726c642068656c6c6f20776f726c64202068656c6c6f20776f726c642068656c6c6f20776f726c642068656c6c6f20776f726c642068656c6c6f20776f726c642068656c6c6f20776f726c642068656c6c6f20776f726c642068656c6c6f20776f726c642068656c6c6f20776f726c64000000000000000000000000000000000000000000000000000000000000';
    assert.equal(a, b);
  });
});

describe('encoding uint32 response', function () {
  it('should equal', function () {
    var a = abi.rawEncode(['uint32'], [42]).toString('hex');
    var b = '000000000000000000000000000000000000000000000000000000000000002a';
    assert.equal(a, b);
  });
});

describe('encoding string response (unsupported)', function () {
  it('should equal', function () {
    var a = abi.rawEncode(['string'], ['a response string (unsupported)']).toString('hex');
    var b = '0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000001f6120726573706f6e736520737472696e672028756e737570706f727465642900';
    assert.equal(a, b);
  });
});

describe('encoding', function () {
  it('should work for uint256', function () {
    var a = abi.rawEncode(['uint256'], [1]).toString('hex');
    var b = '0000000000000000000000000000000000000000000000000000000000000001';
    assert.equal(a, b);
  });
  it('should work for uint', function () {
    var a = abi.rawEncode(['uint'], [1]).toString('hex');
    var b = '0000000000000000000000000000000000000000000000000000000000000001';
    assert.equal(a, b);
  });
  it('should work for int256', function () {
    var a = abi.rawEncode(['int256'], [-1]).toString('hex');
    var b = 'ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff';
    assert.equal(a, b);
  });
});

describe('encoding bytes33', function () {
  it('should fail', function () {
    assert.throws(function () {
      abi.rawEncode('fail', ['bytes33'], ['']);
    }, Error);
  });
});

describe('encoding uint0', function () {
  it('should fail', function () {
    assert.throws(function () {
      abi.rawEncode('fail', ['uint0'], [1]);
    }, Error);
  });
});

describe('encoding uint257', function () {
  it('should fail', function () {
    assert.throws(function () {
      abi.rawEncode('fail', ['uint257'], [1]);
    }, Error);
  });
});

describe('encoding int0', function () {
  it('should fail', function () {
    assert.throws(function () {
      abi.rawEncode(['int0'], [1]);
    }, Error);
  });
});

describe('encoding int257', function () {
  it('should fail', function () {
    assert.throws(function () {
      abi.rawEncode(['int257'], [1]);
    }, Error);
  });
});

describe('encoding uint[2] with [1,2,3]', function () {
  it('should fail', function () {
    assert.throws(function () {
      abi.rawEncode(['uint[2]'], [[1, 2, 3]]);
    }, Error);
  });
});

/*
describe('encoding uint8 with 9bit data', function () {
  it('should fail', function () {
    assert.throws(function () {
      console.log(abi.rawEncode([ 'uint8' ], [ new BN(1).iushln(9) ]));
    }, Error)
  })
})
*/

// Homebrew decoding tests

describe('decoding uint32', function () {
  it('should equal', function () {
    var a = abi.rawDecode(['uint32'], new Buffer('000000000000000000000000000000000000000000000000000000000000002a', 'hex'));
    var b = new BigNumber(42);
    assert.equal(Object.keys(a).length, 1);
    assert.equal(a[0].toString(), b.toString());
  });
});

describe('decoding uint256[]', function () {
  it('should equal', function () {
    var a = abi.rawDecode(['uint256[]'], new Buffer('00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000003', 'hex'));
    var b = new BN(1);
    var c = new BigNumber(2);
    var d = new BN(3);

    assert.equal(Object.keys(a).length, 1);
    assert.equal(a[0].length, 3);
    assert.equal(a[0][0].toString(), b.toString());
    assert.equal(a[0][1].toString(), c.toString());
    assert.equal(a[0][2].toString(), d.toString());
  });
});

describe('decoding bytes', function () {
  it('should equal', function () {
    var a = abi.rawDecode(['bytes'], new Buffer('0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000b68656c6c6f20776f726c64000000000000000000000000000000000000000000', 'hex'));
    var b = new Buffer('68656c6c6f20776f726c64', 'hex');

    assert.equal(Object.keys(a).length, 1);
    assert.equal(new Buffer(a[0].slice(2), 'hex').toString(), b.toString());
  });
});

describe('decoding string', function () {
  it('should equal', function () {
    var a = abi.rawDecode(['string'], new Buffer('0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000b68656c6c6f20776f726c64000000000000000000000000000000000000000000', 'hex'));
    var b = 'hello world';
    assert.equal(Object.keys(a).length, 1);
    assert.equal(a[0], b);
  });
});

describe('decoding int32', function () {
  it('should equal', function () {
    var a = abi.rawDecode(['int32'], new Buffer('fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe', 'hex'));
    var b = new BigNumber(-2);
    assert.equal(Object.keys(a).length, 1);
    assert.equal(a[0].toString(), b.toString());

    a = abi.rawDecode(['int64'], new Buffer('ffffffffffffffffffffffffffffffffffffffffffffffffffffb29c26f344fe', 'hex'));
    b = new BN(-85091238591234);
    assert.equal(Object.keys(a).length, 1);
    assert.equal(a[0].toString(), b.toString());
  });
  /*
  it('should fail', function () {
    assert.throws(function () {
      abi.rawDecode([ 'int32' ], new Buffer('ffffffffffffffffffffffffffffffffffffffffffffffffffffb29c26f344fe', 'hex'))
    }, Error)
  })
  */
});

describe('decoding bool, uint32', function () {
  it('should equal', function () {
    var a = abi.rawDecode(['bool', 'uint32'], new Buffer('0000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002a', 'hex'));
    assert.equal(Object.keys(a).length, 2);
    assert.equal(a[0], true);
    assert.equal(a[1].toString(), new BN(42).toString());
  });
});

describe('decoding bool, uint256[]', function () {
  it('should equal', function () {
    var a = abi.rawDecode(['bool', 'uint256[]'], new Buffer('000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002a', 'hex'));
    assert.equal(Object.keys(a).length, 2);
    assert.equal(a[0], true);
    assert.equal(a[1].length, 1);
    assert.equal(a[1][0].toString(), new BN(42).toString());
  });
});

describe('decoding uint256[], bool', function () {
  it('should equal', function () {
    var a = abi.rawDecode(['uint256[]', 'bool'], '0x000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002a');

    assert.equal(Object.keys(a).length, 2);
    assert.equal(a[1], true);
    assert.equal(a[0].length, 1);
    assert.equal(a[0][0].toString(), new BN(42).toString());
  });
});

describe('decoding fixed-array', function () {
  it('uint[3]', function () {
    var a = abi.rawDecode(['uint[3]'], new Buffer('000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000003', 'hex'));
    assert.equal(Object.keys(a).length, 1);
    assert.equal(a[0].length, 3);
    assert.equal(a[0][0].toString(10), 1);
    assert.equal(a[0][1].toString(10), 2);
    assert.equal(a[0][2].toString(10), 3);
  });
});

describe('decoding (uint[2], uint)', function () {
  it('should work', function () {
    var a = abi.rawDecode(['uint[2]', 'uint'], new Buffer('0000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000005c0000000000000000000000000000000000000000000000000000000000000003', 'hex'));
    assert.equal(Object.keys(a).length, 2);
    assert.equal(a[0].length, 2);
    assert.equal(a[0][0].toString(10), 1);
    assert.equal(a[0][1].toString(10), 92);
    assert.equal(a[1].toString(10), 3);
  });
});

/* FIXME: should check that the whole input buffer was consumed
describe('decoding uint[2] with [1,2,3]', function () {
  it('should fail', function () {
    assert.throws(function () {
      abi.rawDecode([ 'uint[2]' ], new Buffer('00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000003', 'hex'))
    }, Error)
  })
})
*/

/*
describe('stringify', function () {
  it('should be hex prefixed for address', function () {
    assert.deepEqual(abi.stringify([ 'address' ], [ new BN('1234', 16) ]), [ '0x1234' ])
  })

  it('should be hex prefixed for bytes', function () {
    assert.deepEqual(abi.stringify([ 'bytes' ], [ new Buffer('1234', 'hex') ]), [ '0x1234' ])
  })

  it('should be hex prefixed for bytesN', function () {
    assert.deepEqual(abi.stringify([ 'bytes32' ], [ new Buffer('1234', 'hex') ]), [ '0x1234' ])
  })

  it('should be a number for uint', function () {
    assert.deepEqual(abi.stringify([ 'uint' ], [ 42 ]), [ '42' ])
  })

  it('should be a number for uintN', function () {
    assert.deepEqual(abi.stringify([ 'uint8' ], [ 42 ]), [ '42' ])
  })

  it('should be a number for int', function () {
    assert.deepEqual(abi.stringify([ 'int' ], [ -42 ]), [ '-42' ])
  })

  it('should be a number for intN', function () {
    assert.deepEqual(abi.stringify([ 'int8' ], [ -42 ]), [ '-42' ])
  })

  it('should work for bool (true)', function () {
    assert.deepEqual(abi.stringify([ 'bool' ], [ true ]), [ 'true' ])
  })

  it('should work for bool (false)', function () {
    assert.deepEqual(abi.stringify([ 'bool' ], [ false ]), [ 'false' ])
  })

  it('should work for address[]', function () {
    assert.deepEqual(abi.stringify([ 'address[]' ], [ [ new BN('1234', 16), new BN('5678', 16) ] ]), [ '0x1234, 0x5678' ])
  })

  it('should work for address[2]', function () {
    assert.deepEqual(abi.stringify([ 'address[2]' ], [ [ new BN('1234', 16), new BN('5678', 16) ] ]), [ '0x1234, 0x5678' ])
  })

  it('should work for bytes[]', function () {
    assert.deepEqual(abi.stringify([ 'bytes[]' ], [ [ new Buffer('1234', 'hex'), new Buffer('5678', 'hex') ] ]), [ '0x1234, 0x5678' ])
  })

  it('should work for bytes[2]', function () {
    assert.deepEqual(abi.stringify([ 'bytes[2]' ], [ [ new Buffer('1234', 'hex'), new Buffer('5678', 'hex') ] ]), [ '0x1234, 0x5678' ])
  })

  it('should work for uint[]', function () {
    assert.deepEqual(abi.stringify([ 'uint[]' ], [ [ 1, 2, 3 ] ]), [ '1, 2, 3' ])
  })

  it('should work for uint[3]', function () {
    assert.deepEqual(abi.stringify([ 'uint[3]' ], [ [ 1, 2, 3 ] ]), [ '1, 2, 3' ])
  })

  it('should work for int[]', function () {
    assert.deepEqual(abi.stringify([ 'int[]' ], [ [ -1, -2, -3 ] ]), [ '-1, -2, -3' ])
  })

  it('should work for int[3]', function () {
    assert.deepEqual(abi.stringify([ 'int[3]' ], [ [ -1, -2, -3 ] ]), [ '-1, -2, -3' ])
  })

  it('should work for multiple entries', function () {
    assert.deepEqual(abi.stringify([ 'bool', 'bool' ], [ true, true ]), [ 'true', 'true' ])
  })
})

// Tests for Solidity's tight packing
describe('solidity tight packing bool', function () {
  it('should equal', function () {
    var a = abi.solidityPack(
      [ 'bool' ],
      [ true ]
    )
    var b = '01'
    assert.equal(a.toString('hex'), b.toString('hex'))

    a = abi.solidityPack(
      [ 'bool' ],
      [ false ]
    )
    b = '00'
    assert.equal(a.toString('hex'), b.toString('hex'))
  })
})

describe('solidity tight packing address', function () {
  it('should equal', function () {
    var a = abi.solidityPack(
      [ 'address' ],
      [ new BN('43989fb883ba8111221e89123897538475893837', 16) ]
    )
    var b = '43989fb883ba8111221e89123897538475893837'
    assert.equal(a.toString('hex'), b.toString('hex'))
  })
})

describe('solidity tight packing string', function () {
  it('should equal', function () {
    var a = abi.solidityPack(
      [ 'string' ],
      [ 'test' ]
    )
    var b = '74657374'
    assert.equal(a.toString('hex'), b.toString('hex'))
  })
})

describe('solidity tight packing bytes', function () {
  it('should equal', function () {
    var a = abi.solidityPack(
      [ 'bytes' ],
      [ new Buffer('123456', 'hex') ]
    )
    var b = '123456'
    assert.equal(a.toString('hex'), b.toString('hex'))
  })
})

describe('solidity tight packing bytes8', function () {
  it('should equal', function () {
    var a = abi.solidityPack(
      [ 'bytes8' ],
      [ new Buffer('123456', 'hex') ]
    )
    var b = '1234560000000000'
    assert.equal(a.toString('hex'), b.toString('hex'))
  })
})

describe('solidity tight packing uint', function () {
  it('should equal', function () {
    var a = abi.solidityPack(
      [ 'uint' ],
      [ 42 ]
    )
    var b = '000000000000000000000000000000000000000000000000000000000000002a'
    assert.equal(a.toString('hex'), b.toString('hex'))
  })
})

describe('solidity tight packing uint16', function () {
  it('should equal', function () {
    var a = abi.solidityPack(
      [ 'uint16' ],
      [ 42 ]
    )
    var b = '002a'
    assert.equal(a.toString('hex'), b.toString('hex'))
  })
})

describe('solidity tight packing int', function () {
  it('should equal', function () {
    var a = abi.solidityPack(
      [ 'int' ],
      [ -42 ]
    )
    var b = 'ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffd6'
    assert.equal(a.toString('hex'), b.toString('hex'))
  })
})

describe('solidity tight packing int16', function () {
  it('should equal', function () {
    var a = abi.solidityPack(
      [ 'int16' ],
      [ -42 ]
    )
    var b = 'ffd6'
    assert.equal(a.toString('hex'), b.toString('hex'))
  })
})

describe('solidity tight packing sha3', function () {
  it('should equal', function () {
    var a = abi.soliditySHA3(
      [ 'address', 'address', 'uint', 'uint' ],
      [ new BN('43989fb883ba8111221e89123897538475893837', 16), 0, 10000, 1448075779 ]
    )
    var b = 'c3ab5ca31a013757f26a88561f0ff5057a97dfcc33f43d6b479abc3ac2d1d595'
    assert.equal(a.toString('hex'), b.toString('hex'))
  })
})

describe('solidity tight packing sha256', function () {
  it('should equal', function () {
    var a = abi.soliditySHA256(
      [ 'address', 'address', 'uint', 'uint' ],
      [ new BN('43989fb883ba8111221e89123897538475893837', 16), 0, 10000, 1448075779 ]
    )
    var b = '344d8cb0711672efbdfe991f35943847c1058e1ecf515ff63ad936b91fd16231'
    assert.equal(a.toString('hex'), b.toString('hex'))
  })
})

describe('solidity tight packing ripemd160', function () {
  it('should equal', function () {
    var a = abi.solidityRIPEMD160(
      [ 'address', 'address', 'uint', 'uint' ],
      [ new BN('43989fb883ba8111221e89123897538475893837', 16), 0, 10000, 1448075779 ]
    )
    var b = '000000000000000000000000a398cc72490f72048efa52c4e92067e8499672e7'
    assert.equal(a.toString('hex'), b.toString('hex'))
  })
})

describe('solidity tight packing with small ints', function () {
  it('should equal', function () {
    var a = abi.soliditySHA3(
      [ 'address', 'address', 'int64', 'uint192' ],
      [ new BN('43989fb883ba8111221e89123897538475893837', 16), 0, 10000, 1448075779 ]
    )
    var b = '1c34bbd3d419c05d028a9f13a81a1212e33cb21f4b96ce1310442911c62c6986'
    assert.equal(a.toString('hex'), b.toString('hex'))
  })
})

describe('converting from serpent types', function () {
  it('should equal', function () {
    assert.deepEqual(abi.fromSerpent('s'), [ 'bytes' ])
    assert.deepEqual(abi.fromSerpent('i'), [ 'int256' ])
    assert.deepEqual(abi.fromSerpent('a'), [ 'int256[]' ])
    assert.deepEqual(abi.fromSerpent('b8'), [ 'bytes8' ])
    assert.deepEqual(abi.fromSerpent('b8i'), [ 'bytes8', 'int256' ])
    assert.deepEqual(abi.fromSerpent('b32'), [ 'bytes32' ])
    assert.deepEqual(abi.fromSerpent('b32i'), [ 'bytes32', 'int256' ])
    assert.deepEqual(abi.fromSerpent('sb8ib8a'), [ 'bytes', 'bytes8', 'int256', 'bytes8', 'int256[]' ])
    assert.throws(function () {
      abi.fromSerpent('i8')
    })
    assert.throws(function () {
      abi.fromSerpent('x')
    })
  })
})

describe('converting to serpent types', function () {
  it('should equal', function () {
    assert.equal(abi.toSerpent([ 'bytes' ]), 's')
    assert.equal(abi.toSerpent([ 'int256' ]), 'i')
    assert.equal(abi.toSerpent([ 'int256[]' ]), 'a')
    assert.equal(abi.toSerpent([ 'bytes8' ]), 'b8')
    assert.equal(abi.toSerpent([ 'bytes32' ]), 'b32')
    assert.equal(abi.toSerpent([ 'bytes', 'bytes8', 'int256', 'bytes8', 'int256[]' ]), 'sb8ib8a')
    assert.throws(function () {
      abi.toSerpent('int8')
    })
    assert.throws(function () {
      abi.toSerpent('bool')
    })
  })
})
*/

describe('utf8 handling', function () {
  it('should encode latin and extensions', function () {
    var a = abi.rawEncode(['string'], ['ethereum számítógép']).toString('hex');
    var b = '00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000017657468657265756d20737ac3a16dc3ad74c3b367c3a970000000000000000000';
    assert.equal(a, b);
  });
  it('should encode non-latin characters', function () {
    var a = abi.rawEncode(['string'], ['为什么那么认真？']).toString('hex');
    var b = '00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000018e4b8bae4bb80e4b988e982a3e4b988e8aea4e79c9fefbc9f0000000000000000';
    assert.equal(a, b);
  });
  it('should decode latin and extensions', function () {
    var a = 'ethereum számítógép';
    var b = abi.rawDecode(['string'], new Buffer('00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000017657468657265756d20737ac3a16dc3ad74c3b367c3a970000000000000000000', 'hex'));
    assert.equal(a, b[0]);
  });
  it('should decode non-latin characters', function () {
    var a = '为什么那么认真？';
    var b = abi.rawDecode(['string'], new Buffer('00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000018e4b8bae4bb80e4b988e982a3e4b988e8aea4e79c9fefbc9f0000000000000000', 'hex'));
    assert.equal(a, b[0]);
  });
});

describe('encoding ufixed128x128', function () {
  it('should equal', function () {
    assert.throws(function () {
      var a = abi.rawEncode(['ufixed128x128'], [1]).toString('hex');
      var b = '0000000000000000000000000000000100000000000000000000000000000000';
      assert.equal(a, b);
    }, Error);
  });
});

describe('encoding fixed128x128', function () {
  it('should equal', function () {
    assert.throws(function () {
      var a = abi.rawEncode(['fixed128x128'], [-1]).toString('hex');
      var b = 'ffffffffffffffffffffffffffffffff00000000000000000000000000000000';
      assert.equal(a, b);
    }, Error);
  });
});

describe('decoding ufixed128x128', function () {
  it('should equal', function () {
    assert.throws(function () {
      var a = new Buffer('0000000000000000000000000000000100000000000000000000000000000000', 'hex');
      var b = abi.rawDecode(['ufixed128x128'], a);
      assert.equal(b[0].toNumber(), 1);
    }, Error);
  });
  it('decimals should fail', function () {
    assert.throws(function () {
      var a = new Buffer('0000000000000000000000000000000100000000000000000000000000000001', 'hex');
      assert.throws(function () {
        abi.rawDecode(['ufixed128x128'], a);
      }, /^Error: Decimals not supported yet/);
    }, Error);
  });
});

describe('decoding fixed128x128', function () {
  it('should equal', function () {
    assert.throws(function () {
      var a = new Buffer('ffffffffffffffffffffffffffffffff00000000000000000000000000000000', 'hex');
      var b = abi.rawDecode(['fixed128x128'], a);
      assert.equal(b[0].toNumber(), -1);
    }, Error);
  });
  it('decimals should fail', function () {
    assert.throws(function () {
      var a = new Buffer('ffffffffffffffffffffffffffffffff00000000000000000000000000000001', 'hex');
      assert.throws(function () {
        abi.rawDecode(['fixed128x128'], a);
      }, /^Error: Decimals not supported yet/);
    }, Error);
  });
});

/*
describe('encoding -1 as uint', function () {
  it('should throw', function () {
    assert.throws(function () {
      console.log(
      abi.rawEncode([ 'uint' ], [ -1 ]));
    }, /^Error: Supplied uint is negative/)
  })
})
*/

describe('encoding 256 bits as bytes', function () {
  it('should not leave trailing zeroes', function () {
    var a = abi.rawEncode(['bytes'], [new Buffer('ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff', 'hex')]);
    assert.equal(a.toString('hex'), '00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000020ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff');
  });
});