import { ethers } from 'hardhat'
import { assert } from 'chai'
//import { evmRevert } from '../test-helpers/matchers'
import { UpkeepTranscoder30__factory as UpkeepTranscoderFactory } from '../../../typechain/factories/UpkeepTranscoder30__factory'
import { UpkeepTranscoder30 as UpkeepTranscoder } from '../../../typechain/UpkeepTranscoder30'
import { evmRevert } from '../../test-helpers/matchers'

let upkeepTranscoderFactory: UpkeepTranscoderFactory
let transcoder: UpkeepTranscoder

before(async () => {
  // @ts-ignore bug in autogen file
  upkeepTranscoderFactory = await ethers.getContractFactory(
    'UpkeepTranscoder3_0',
  )
})

describe.only('UpkeepTranscoder3_0', () => {
  const formatV1 = 0
  const formatV2 = 1
  const formatV3 = 2

  beforeEach(async () => {
    transcoder = await upkeepTranscoderFactory.deploy()
  })

  describe('#typeAndVersion', () => {
    it('uses the correct type and version', async () => {
      const typeAndVersion = await transcoder.typeAndVersion()
      assert.equal(typeAndVersion, 'UpkeepTranscoder 3.0.0')
    })
  })

  describe('#transcodeUpkeeps', () => {
    const encodedData = '0xabcd'

    it('reverts if the from type is not an enum value', async () => {
      await evmRevert(
        transcoder.transcodeUpkeeps(3, 2, encodedData),
        'function was called with incorrect parameters',
      )
    })

    it('reverts if the to version is not 2', async () => {
      await evmRevert(
        transcoder.transcodeUpkeeps(1, 3, encodedData),
        'function was called with incorrect parameters',
      )
    })

    context('when from and to versions equal', () => {
      it('returns the data that was passed in', async () => {
        let response = await transcoder.transcodeUpkeeps(
          formatV1,
          formatV1,
          encodedData,
        )
        assert.equal(encodedData, response)

        response = await transcoder.transcodeUpkeeps(
          formatV2,
          formatV2,
          encodedData,
        )
        assert.equal(encodedData, response)

        response = await transcoder.transcodeUpkeeps(
          formatV3,
          formatV3,
          encodedData,
        )
        assert.equal(encodedData, response)
      })
    })
  })
})
