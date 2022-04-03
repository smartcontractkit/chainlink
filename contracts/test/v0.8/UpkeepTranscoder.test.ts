import { ethers } from 'hardhat'
import { assert } from 'chai'
import { evmRevert } from '../test-helpers/matchers'
import { UpkeepTranscoder__factory as UpkeepTranscoderFactory } from '../../typechain/factories/UpkeepTranscoder__factory'
import { UpkeepTranscoder } from '../../typechain'

let upkeepTranscoderFactory: UpkeepTranscoderFactory
let transcoder: UpkeepTranscoder

before(async () => {
  upkeepTranscoderFactory = await ethers.getContractFactory('UpkeepTranscoder')
})

describe('UpkeepTranscoder', () => {
  beforeEach(async () => {
    transcoder = await upkeepTranscoderFactory.deploy()
  })

  describe('#typeAndVersion', () => {
    it('uses the correct type and version', async () => {
      const typeAndVersion = await transcoder.typeAndVersion()
      assert.equal(typeAndVersion, 'UpkeepTranscoder 1.0.0')
    })
  })

  describe('#transcodeUpkeeps', () => {
    const encodedData = '0xc0ffee'

    it('reverts if the from type is not 1', async () => {
      await evmRevert(
        transcoder.transcodeUpkeeps(2, 1, encodedData),
        'InvalidTranscoding()',
      )
    })

    it('reverts if the from type is not 1', async () => {
      await evmRevert(
        transcoder.transcodeUpkeeps(1, 2, encodedData),
        'InvalidTranscoding()',
      )
    })

    context('when from and to versions equal 1', () => {
      it('returns the data that was passed in', async () => {
        const response = await transcoder.transcodeUpkeeps(1, 1, encodedData)
        assert.equal(encodedData, response)
      })
    })
  })
})
