import {
  checkPublicABI,
  defaultAccount,
  deploy,
  lPad,
  rPad,
  newHash,
  oracleNode,
  toHex,
  toHexWithoutPrefix
} from './support/helpers'

const ensSubnodeHash = (node, name) => {
  let label = toHexWithoutPrefix(rPad(name))
  let combo = web3.sha3(node + label, {encoding: 'hex'})
  return combo.toString()
}

contract('UpdatableConsumer', () => {
  const sourcePath = 'examples/UpdatableConsumer.sol'

  const ensRoot = toHex(lPad('\x00'))
  const rootDomain = 'domainlink'
  const rootHash = ensSubnodeHash(ensRoot, rPad(rootDomain))
  const tokenDomain = 'link'
  const tokenHash = ensSubnodeHash(rootHash, tokenDomain)
  const oracleDomain = 'oracle'
  const oracleHash = ensSubnodeHash(rootHash, oracleDomain)
  const specId = newHash('0x4c7b7ffb66b344fbaa64995af81e355a')

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

    uc = await deploy(sourcePath, toHex(specId), ens.address, rootHash)
  })

  describe('constructor', () => {
    it('pulls the token contract address from the resolver', async () => {
      assert.equal(link.address, await uc.publicLinkToken.call())
    })

    it('pulls the oracle contract address from the resolver', async () => {
      assert.equal(oc.address, await uc.publicOracle.call())
    })
  })

  describe('#updateOracle', () => {
    describe('when the ENS resolver has been updated', () => {
      const newOracleAddress = '0xf000000000000000000000000000000000000ba7'

      beforeEach(async () => {
        await ensResolver.setAddr(oracleHash, newOracleAddress, {from: oracleNode})
      })

      it("updates the contract's oracle address", async () => {
        await uc.updateOracle()

        assert.equal(newOracleAddress, await uc.publicOracle.call())
      })
    })

    describe("when the ENS resolver has not been updated", () => {
      it("keeps the same oracle address", async () => {
        await uc.updateOracle()

        assert.equal(oc.address, await uc.publicOracle.call())
      })
    })
  })
})
