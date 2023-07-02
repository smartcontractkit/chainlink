import { ethers } from 'hardhat';
import { toWei, increaseTime5Minutes, toHex } from '../test-helpers/helpers';
import { assert, expect } from 'chai';
import { BigNumber, constants, Contract, ContractFactory } from 'ethers';
import { Roles, getUsers } from '../test-helpers/setup';
import { bigNumEquals, evmRevert } from '../test-helpers/matchers';
import { convertFufillParams, decodeRunRequest, encodeOracleRequest, RunRequest } from '../test-helpers/oracle';
import cbor from 'cbor';
import { makeDebug } from '../test-helpers/debug';

const d = makeDebug('BasicConsumer');
let basicConsumerFactory: ContractFactory;
let oracleFactory: ContractFactory;
let linkTokenFactory: ContractFactory;

let roles: Roles;

before(async () => {
roles = await getUsers().roles;
basicConsumerFactory = await ethers.getContractFactory('src/v0.6/tests/BasicConsumer.sol:BasicConsumer', roles.defaultAccount);
oracleFactory = await ethers.getContractFactory('src/v0.6/Oracle.sol:Oracle', roles.oracleNode);
linkTokenFactory = await ethers.getContractFactory('src/v0.4/LinkToken.sol:LinkToken', roles.defaultAccount);
});

describe('BasicConsumer', () => {
const specId = '0x4c7b7ffb66b344fbaa64995af81e355a'.padEnd(66, '0');
const currency = 'USD';
const payment = toWei('1');
let link: Contract;
let oc: Contract;
let cc: Contract;

beforeEach(async () => {
link = await linkTokenFactory.connect(roles.defaultAccount).deploy();
oc = await oracleFactory.connect(roles.oracleNode).deploy(link.address);
cc = await basicConsumerFactory.connect(roles.defaultAccount).deploy(link.address, oc.address, specId);
});

it('has a predictable gas price [ @skip-coverage ]', async () => {
const rec = await ethers.provider.getTransactionReceipt(cc.deployTransaction.hash ?? '');
assert.isBelow(rec.gasUsed?.toNumber() ?? -1, 1750000);
});

describe('#requestEthereumPrice', () => {
describe('without LINK', () => {
it('reverts', async () => {
await expect(cc.requestEthereumPrice(currency, payment)).to.be.reverted;
});
});
  describe('with LINK', () => {
  beforeEach(async () => {
    await link.transfer(cc.address, toWei('1'));
  });

  it('triggers a log event in the Oracle contract', async () => {
    const tx = await cc.requestEthereumPrice(currency, payment);
    const receipt = await tx.wait();

    const log = receipt?.logs?.[3];
    assert.equal(log?.address.toLowerCase(), oc.address.toLowerCase());

    const request = decodeRunRequest(log);
    const expected = {
      path: ['USD'],
      get: 'https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD,EUR,JPY',
    };

    assert.equal(toHex(specId), request.specId);
    bigNumEquals(toWei('1'), request.payment);
    assert.equal(cc.address.toLowerCase(), request.requester.toLowerCase());
    assert.equal(1, request.dataVersion);
    assert.deepEqual(expected, cbor.decodeFirstSync(request.data));
  });

  it('has a reasonable gas cost [ @skip-coverage ]', async () => {
    const tx = await cc.requestEthereumPrice(currency, payment);
    const receipt = await tx.wait();
