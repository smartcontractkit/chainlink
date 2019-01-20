import { deploy, bigNum } from './support/helpers'
import { assertBigNum } from './support/matchers'

const wordSizeBits = 256
const wordSizeBytes = wordSizeBits / 8

const keySizeBits = 2048
assert(keySizeBits % wordSizeBits === 0, `Key size must be multiple of words`)
const keySizeWords = keySizeBits / wordSizeBits
const keySizeBytes = keySizeBits / 8

// A roughly 2048-bit prime, per https://2ton.com.au/getprimes/random/2048
const prime = bigNum(
  '0x' +
    '91d18d4420ab0cae83964b2310dc7277e61ad331d8e37de4923b355308cd2387' +
    '46dbd85833993853724b0b7048c0d331b177ecc9486ba14142a38cf292b8be6c' +
    '861852d02fa41f1a12b5c0e13716fba0887bb5568b7caf1eac255fa6fded398b' +
    'ff863ab3391450edc27ec52dc92bd66df1dc818fa58259aca354d5cbdfe427fa' +
    'ec81c497231ae625c8f3afc0a37b8fe7752ad8c0cd04fd2a1177680d0334b2a1' +
    'ee60cd49f629a8c5e71ad3cc1af7b26fc29c7112be6162604b82f0cba28cc2d3' +
    '521f09edbdf598be03adcf4797b50b948418bc01e298ae1815d5d2c7af41f795' +
    '4471f3f52b60da23e73b8e27706ea90c877071ddc20e3ad78404f352306157b7')

const toHexString = byteArray =>
  Array.prototype.map.call(
    byteArray, byte => ('0' + (byte & 0xFF).toString(16)).slice(-2)).join('')

const numToUint256Array = n => { // n as keySizeWords-length array of uint256's
  const asBytes = bigNum(n).toArray('be', keySizeBytes)
  const rv = []
  for (let bytesStart = 0; bytesStart < asBytes.length; bytesStart += wordSizeBytes) {
    const uint256AsBytes = asBytes.slice(bytesStart, bytesStart + wordSizeBytes)
    rv.push(bigNum('0x' + toHexString(uint256AsBytes)))
  }
  assert(rv.length = keySizeWords)
  rv.forEach(c => assert(c.bitLength() <= wordSizeBits, 'Should be uint256 chunks'))
  // assertBigNum(n, bigNum([].concat(rv)), 'rv should be uint256 chunking of n')
  return rv
}

const uint256ArrayToNum = a => {
  const asBytes = [].concat(...a.map(e => bigNum(e).toArray('be', wordSizeBytes)))
  return bigNum('0x' + toHexString(asBytes))
}

contract('RSAVRF', async () => {
  let RSAVRF
  beforeEach(async () => {
    RSAVRF = await deploy('RSAVRF.sol', numToUint256Array(prime))
  })
  it('Accurately computes a bigModExp', async () => {
    const exp = async n => uint256ArrayToNum(
      await RSAVRF.bigModExp(numToUint256Array(n)))
    assertBigNum(await exp(2), 8, '2³≡8 mod p')
    assertBigNum(await exp(prime), 0, 'p³≡0 mod p')
    const minusOne = prime.sub(bigNum(1))
    assertBigNum(await exp(minusOne), minusOne, '(-1)³≡-1 mod p')
  })
})
