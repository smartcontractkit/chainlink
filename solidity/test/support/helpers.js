BigNumber = require('bignumber.js');
moment = require('moment');

(() => {
  eth = web3.eth;

  before(async function () {
    accounts = await eth.accounts;
    Accounts = accounts.slice(1);

    oracleNode = Accounts[0];
    stranger = Accounts[1];
    consumer = Accounts[2];
  });

  Eth = function sendEth(method, params) {
    params = params || [];

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

})();
