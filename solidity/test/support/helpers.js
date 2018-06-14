const Wallet = require('../../app/wallet.js');
const Utils = require('../../app/utils.js');
const Deployer = require('../../app/deployer.js');

BigNumber = require('bignumber.js');
moment = require('moment');
abi = require('ethereumjs-abi');
util = require('ethereumjs-util');
cbor = require('cbor');

(() => {
  let utils, wallet, deployer
  eth = web3.eth;

  process.env.SOLIDITY_INCLUDE = '../../contracts/:../../contracts/lib/:../../node_modules/'

  before(async function () {
    // Default hard coded truffle accounts:
    // ==================
    // (0) 0x627306090abab3a6e1400e9345bc60c78a8bef57
    // (1) 0xf17f52151ebef6c7334fad080c5704d77216b732
    // (2) 0xc5fdf4076b8f3a5357c5e395ab970b5b54098fef
    // (3) 0x821aea9a577a9b44299b9c15c88cf3087f3b5544
    // (4) 0x0d1d4e623d10f9fba5db95830f7d3839406c6af2
    // (5) 0x2932b7a2355d6fecc4b5c0b6bd44cc31df247a2e
    // (6) 0x2191ef87e392377ec08e7c08eb105ef5448eced5
    // (7) 0x0f4f2ac550a1b4e2280d04c21cea7ebd822934b5
    // (8) 0x6330a553fc93768f612722bb8c2ec78ac90b3bbc
    // (9) 0x5aeda56215b167893e80b4fe645ba6d5bab767de

    // Private Keys
    // ==================
    // (0) c87509a1c067bbde78beb793e6fa76530b6382a4c0241e5e4a9ec0a0f44dc0d3
    // (1) ae6ae8e5ccbfb04590405997ee2d52d2b330726137b875053c36d94e974d162f
    // (2) 0dbbe8e4ae425a6d2687f1a7e3ba17bc98c673636790f1b8ad91193c05875ef1
    // (3) c88b703fb08cbea894b6aeff5a544fb92e78a18e19814cd85da83b71f772aa6c
    // (4) 388c684f0ba1ef5017716adb5d21a053ea8e90277d0868337519f97bede61418
    // (5) 659cbb0e2411a44db63778987b1e22153c086a95eb6b18bdf89de078917abc63
    // (6) 82d052c865f5763aad42add438569276c00d3d88a2d062d36b2bae914d58b8c8
    // (7) aa3680d5d48a8283413f7a108367c7299ca73f553735860a87b08f39395618b7
    // (8) 0f62d96d6675f32685bbdb8ac13cda7c23436f63efbb9d07700d8669ff12b7c4
    // (9) 8d5366123cb560bb606379f90a0bfd4769eecc0557f1b362dcae9012b548b1e5

    // HD Wallet
    // ==================
    // Mnemonic:      candy maple cake sugar pudding cream honey rich smooth crumble sweet treat
    // Base HD Path:  m/44'/60'/0'/0/{account_index}
    accounts = eth.accounts

    defaultAccount = accounts[0]
    oracleNode = accounts[1]
    stranger = accounts[2]
    consumer = accounts[3]

    const privateKey = 'c87509a1c067bbde78beb793e6fa76530b6382a4c0241e5e4a9ec0a0f44dc0d3'
    utils = new Utils(web3.currentProvider)
    wallet = new Wallet(privateKey, utils)
    deployer = new Deployer(wallet, utils)
  });

  deploy = (filePath, ...args) => {
    return deployer.perform(filePath, ...args)
  }

  Eth = function sendEth(method, params) {
    params = params || []

    return new Promise((resolve, reject) => {
      web3.currentProvider.sendAsync({
        jsonrpc: "2.0",
        method: method,
        params: params || [],
        id: new Date().getTime()
      }, function sendEthResponse(error, response) {
        if (error) {
          reject(error);
        } else {
          resolve(response.result);
        };
      }, () => {}, () => {});
    });
  };

  emptyAddress = '0x0000000000000000000000000000000000000000';

  sealBlock = async function sealBlock() {
    return Eth('evm_mine');
  };

  sendTransaction = async function sendTransaction(params) {
    return await eth.sendTransaction(params);
  }

  getBalance = async function getBalance(account) {
    return bigNum(await eth.getBalance(account));
  }

  bigNum = function bigNum(number) {
    return new BigNumber(number);
  }

  toWei = function toWei(number) {
    return bigNum(web3.toWei(number));
  }

  tokens = function tokens(number) {
    return bigNum(number * 10**18);
  }

  intToHex = function intToHex(number) {
    return '0x' + bigNum(number).toString(16);
  }

  intToHexNoPrefix = function intToHex(number) {
    return bigNum(number).toString(16);
  }

  hexToInt = function hexToInt(string) {
    return web3.toBigNumber(string);
  }

  hexToAddress = function hexToAddress(string) {
    return '0x' + string.slice(string.length - 40);
  }

  unixTime = function unixTime(time) {
    return moment(time).unix();
  }

  seconds = function seconds(number) {
    return number;
  };

  minutes = function minutes(number) {
    return number * 60;
  };

  hours = function hours(number) {
    return number * minutes(60);
  };

  days = function days(number) {
    return number * hours(24);
  };

  keccak256 = function keccak256(string) {
    return web3.sha3(string);
  }

  logTopic = function logTopic(string) {
    let hash = keccak256(string);
    return '0x' + hash.slice(26);
  }

  getLatestBlock = async function getLatestBlock() {
    return await eth.getBlock('latest', false);
  };

  getLatestTimestamp = async function getLatestTimestamp () {
    let latestBlock = await getLatestBlock()
    return web3.toDecimal(latestBlock.timestamp);
  };

  fastForwardTo = async function fastForwardTo(target) {
    let now = await getLatestTimestamp();
    assert.isAbove(target, now, "Cannot fast forward to the past");
    let difference = target - now;
    await Eth("evm_increaseTime", [difference]);
    await sealBlock();
  };

  getEvents = function getEvents(contract) {
    return new Promise((resolve, reject) => {
      contract.allEvents().get((error, events) => {
        if (error) {
          reject(error);
        } else {
          resolve(events);
        };
      });
    });
  };

  eventsOfType = function eventsOfType(events, type) {
    let filteredEvents = [];
    for (event of events) {
      if (event.event === type) filteredEvents.push(event);
    }
    return filteredEvents;
  };

  getEventsOfType = async function getEventsOfType(contract, type) {
    return eventsOfType(await getEvents(contract), type);
  };

  getLatestEvent = async function getLatestEvent(contract) {
    let events = await getEvents(contract);
    return events[events.length - 1];
  };

  assertActionThrows = function assertActionThrows(action) {
    return Promise.resolve().then(action)
      .catch(error => {
        assert(error, "Expected an error to be raised");
        assert(error.message, "Expected an error to be raised");
        return error.message;
      })
      .then(errorMessage => {
        assert(errorMessage, "Expected an error to be raised");
        invalidOpcode = errorMessage.includes("invalid opcode")
        reverted = errorMessage.includes("VM Exception while processing transaction: revert")
        assert.isTrue(invalidOpcode || reverted, 'expected error message to include "invalid JUMP" or "revert"');
        // see https://github.com/ethereumjs/testrpc/issues/39
        // for why the "invalid JUMP" is the throw related error when using TestRPC
      })
  };

  encodeUint256 = function encodeUint256(int) {
    let zeros = "0000000000000000000000000000000000000000000000000000000000000000";
    let payload = int.toString(16);
    return (zeros + payload).slice(payload.length);
  }

  encodeAddress = function encodeAddress(address) {
    return '000000000000000000000000' + address.slice(2);
  }

  encodeBytes = function encodeBytes(bytes) {
    let zeros = "0000000000000000000000000000000000000000000000000000000000000000";
    let padded = bytes.padEnd(64, 0);
    let length = encodeUint256(bytes.length / 2);
    return length + padded;
  }

  checkPublicABI = function checkPublicABI(contract, expectedPublic) {
    let actualPublic = [];
    for (method of contract.abi) {
      if (method.type == 'function') actualPublic.push(method.name);
    };

    for (method of actualPublic) {
      let index = expectedPublic.indexOf(method);
      assert.isAtLeast(index, 0, (`#${method} is NOT expected to be public`))
    }

    for (method of expectedPublic) {
      let index = actualPublic.indexOf(method);
      assert.isAtLeast(index, 0, (`#${method} is expected to be public`))
    }
  };

  functionSelector = function functionSelector(signature) {
    return "0x" + web3.sha3(signature).slice(2).slice(0, 8);
  };

  rPad = function rPad(string) {
    let wordLen = parseInt((string.length + 31) / 32) * 32;
    for (let i = string.length; i < wordLen; i++) {
      string = string + "\x00";
    }
    return string
  };

  lPad = function lPad(string) {
    let wordLen = parseInt((string.length + 31) / 32) * 32;
    for (let i = string.length; i < wordLen; i++) {
      string = "\x00" + string;
    }
    return string
  };

  lPadHex = function lPadHex(string) {
    let wordLen = parseInt((string.length + 63) / 64) * 64;
    for (let i = string.length; i < wordLen; i++) {
      string = "0" + string;
    }
    return string
  };

  toHex = function toHex(arg) {
    if (arg instanceof Buffer) {
      return arg.toString("hex");
    } else {
      return Buffer.from(arg, "ascii").toString("hex");
    }
  };

  decodeRunABI = function decodeRunABI(log) {
    let runABI = util.toBuffer(log.data);
    let types = ["bytes32", "address", "bytes4", "bytes"];
    return abi.rawDecode(types, runABI);
  };

  decodeRunRequest = function decodeRunRequest(log) {
    let runABI = util.toBuffer(log.data);
    let types = ["uint256", "bytes"];
    let [version, data] = abi.rawDecode(types, runABI);
    return [log.topics[1], log.topics[2], log.topics[3], version, data];
  };

  decodeSpecAndRunRequest = function(log) {
    let runABI = util.toBuffer(log.data);
    let types = ["uint256", "bytes"];
    let [version, data] = abi.rawDecode(types, runABI);
    return [log.topics[1], log.topics[2], version, data];
  }

  requestDataBytes = function requestDataBytes(specId, to, fHash, runId, data) {
    let types = ["uint256", "bytes32", "address", "bytes4", "bytes32", "bytes"];
    let values = [1, specId, to, fHash, runId, data];
    let encoded = abi.rawEncode(types, values);
    let funcSelector = functionSelector("requestData(uint256,bytes32,address,bytes4,bytes32,bytes)");
    return funcSelector + encoded.toString("hex");
  };

  requestDataFrom = function requestDataFrom(oc, link, amount, args) {
    return link.transferAndCall(oc.address, amount, args);
  };

  specAndRunBytes = function(to, fHash, requestId, data) {
    let types = ["uint256", "address", "bytes4", "bytes32", "bytes"];
    let values = [1, to, fHash, requestId, data];
    let encoded = abi.rawEncode(types, values);
    let funcSelector = functionSelector("specAndRun(uint256,address,bytes4,bytes32,bytes)");
    return funcSelector + encoded.toString("hex");
  };

  _0x = (val) => {
    return `0x${toHex(val)}`
  }

})();
