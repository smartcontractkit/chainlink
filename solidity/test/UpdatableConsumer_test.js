import {
  checkPublicABI,
  defaultAccount,
  deploy,
  lPad,
  rPad,
  oracleNode,
  newUint8ArrayFromStr,
} from './support/helpers'
import utils from 'ethereumjs-util'

contract('UpdatableConsumer', () => {
  const sourcePath = 'examples/UpdatableConsumer.sol'
  const domainName = 'domainLink'
  const domainHash = web3.sha3(lPad("\x00") + rPad(domainName))

  let ens, ensResolver, link, oc, uc

  beforeEach(async () => {
    link = await deploy('LinkToken.sol')
    oc = await deploy('Oracle.sol', link.address)
    ens = await deploy('examples/ENSRegistry.sol')
    ensResolver = await deploy('examples/PublicResolver.sol', ens.address)

    await ens.setSubnodeOwner('', domainName, oracleNode)
    await ens.setResolver(domainHash.toString(), ensResolver.address, {from: oracleNode})
    await ens.setSubnodeOwner('', domainName, oracleNode)
    await ensResolver.setAddr(domainHash, oc.address, {from: oracleNode})

    uc = await deploy(sourcePath, link.address, ens.address, domainHash)
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
