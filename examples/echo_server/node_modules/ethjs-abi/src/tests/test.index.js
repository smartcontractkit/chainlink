const assert = require('chai').assert;
const abi = require('../index.js');
const contracts = require('./contracts.json');

describe('test basic encoding and decoding functionality', () => {
  it('should encode and decode contract data nicely', () => {
    const BalanceClaimInterface = JSON.parse(contracts.BalanceClaim.interface);
    const encodeBalanceClaimMethod1 = abi.encodeMethod(BalanceClaimInterface[0], []);
    assert.equal(encodeBalanceClaimMethod1, '0x30509bca');
    const interfaceABI = [{'constant':false,'inputs':[{'name':'_value','type':'uint256'}],'name':'set','outputs':[{'name':'','type':'bool'}],'payable':false,'type':'function'},{'constant':false,'inputs':[],'name':'get','outputs':[{'name':'storeValue','type':'uint256'}],'payable':false,'type':'function'},{'anonymous':false,'inputs':[{'indexed':false,'name':'_newValue','type':'uint256'},{'indexed':false,'name':'_sender','type':'address'}],'name':'SetComplete','type':'event'}]; // eslint-disable-line

    const setMethodInputBytecode = abi.encodeMethod(interfaceABI[0], [24000]);
    abi.decodeMethod(interfaceABI[0], '0x0000000000000000000000000000000000000000000000000000000000000001');

    abi.encodeMethod(interfaceABI[1], []);
    abi.decodeMethod(interfaceABI[1], '0x000000000000000000000000000000000000000000000000000000000000b26e');

    abi.encodeEvent(interfaceABI[2], [24000, '0xca35b7d915458ef540ade6068dfe2f44e8fa733c']);
    abi.decodeEvent(interfaceABI[2], '0x0000000000000000000000000000000000000000000000000000000000000d7d000000000000000000000000ca35b7d915458ef540ade6068dfe2f44e8fa733c');

    assert.equal(setMethodInputBytecode, '0x60fe47b10000000000000000000000000000000000000000000000000000000000005dc0');
  });
});
