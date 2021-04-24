const { assert, web3 } = require('hardhat')
const { constants } = require('@openzeppelin/test-helpers')

describe('ValidatorProxy', () => {
  let accounts
  let ValidatorProxyArtifact
  let validatorProxy

  beforeEach(async () => {
    ValidatorProxyArtifact = artifacts.require('ValidatorProxy')
    accounts = await web3.eth.getAccounts()
  })

  describe('#constructor', () => {
    let owner, aggregator, validator

    beforeEach(async () => {
      owner = accounts[0]
      aggregator = accounts[1]
      validator = accounts[2]
      validatorProxy = await ValidatorProxyArtifact.new(aggregator, validator, {
        from: owner,
      })
    })

    it('should set the aggregator addresses correctly', async () => {
      const response = await validatorProxy.getAggregators()
      assert.equal(response.current, aggregator)
      assert.equal(response.proposed, constants.ZERO_ADDRESS)
    })

    it('should set the validator addresses conrrectly', async () => {
      const response = await validatorProxy.getValidators()
      assert.equal(response.current, validator)
      assert.equal(response.proposed, constants.ZERO_ADDRESS)
    })

    it('should set the owner correctly', async () => {
      const response = await validatorProxy.owner()
      assert.equal(response, owner)
    })
  })

  describe('#validate', () => {
    // TODO
  })
  describe('#proposeNewAggregator', () => {
    // TODO
  })
  describe('#upgradeAggregator', () => {
    // TODO
  })
  describe('#proposeNewValidator', () => {
    // TODO
  })
  describe('#upgradeValidator', () => {
    // TODO
  })
})
