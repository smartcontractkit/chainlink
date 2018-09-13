import {
  checkPublicABI,
  defaultAccount,
  deploy,
  lPad,
  rPad,
  oracleNode,
  toHex,
  toHexWithoutPrefix,
  newUint8ArrayFromStr,
} from './support/helpers'
import utils from 'ethereumjs-util'


const ensRoot = toHex(lPad("\x00"))

const ensSubnodeHash = (node, name) => {
  let label = toHexWithoutPrefix(rPad(name))
  let combo = web3.sha3(node + label, {encoding: 'hex'})
  return combo.toString()
}

const rootDomain = 'domainlink'
const rootHash = ensSubnodeHash(ensRoot, rPad(rootDomain))
const tokenDomain = 'link'
const tokenHash = ensSubnodeHash(rootHash, tokenDomain)
const oracleDomain = 'oracle'
const oracleHash = ensSubnodeHash(rootHash, oracleDomain)

contract('UpdatableConsumer', () => {
  const sourcePath = 'examples/UpdatableConsumer.sol'

  let ens, ensResolver, link, oc, uc

  beforeEach(async () => {
    link = await deploy('LinkToken.sol')
    oc = await deploy('Oracle.sol', link.address)
    ens = await deploy('ENSRegistry.sol')
    ensResolver = await deploy('PublicResolver.sol', ens.address)

    // register domain
    await ens.setSubnodeOwner('', rootDomain, oracleNode)
    await ens.setResolver(rootHash, ensResolver.address, {from: oracleNode})
    await ensResolver.setAddr(rootHash, oc.address, {from: oracleNode})

    // register token subdomain
    await ens.setSubnodeOwner(rootHash, tokenDomain, oracleNode, {from: oracleNode})
    await ensResolver.setAddr(tokenHash, link.address, {from: oracleNode})

    // register oracle subdomain
    await ens.setSubnodeOwner(rootHash, oracleDomain, oracleNode, {from: oracleNode})
    await ensResolver.setAddr(oracleHash, oc.address, {from: oracleNode})

    uc = await deploy(sourcePath, ens.address, rootHash)
  })

  it('has a limited public interface', () => {
    checkPublicABI(artifacts.require(sourcePath), [
      'publicLinkToken',
      'publicOracle'
    ])
  })

  describe('constructor', () => {
    it('pulls the oracle contract address from the resolver', async () => {
      assert.equal(link.address, await uc.publicLinkToken.call())
      assert.equal(oc.address, await uc.publicOracle.call())
    })
  })
})
