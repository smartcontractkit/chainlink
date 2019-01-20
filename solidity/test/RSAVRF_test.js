import { deploy, bigNum } from './support/helpers'
import { assertBigNum } from './support/matchers'

const wordSizeBits = 256
const wordSizeBytes = wordSizeBits / 8

const keySizeBits = 4096
assert(keySizeBits % wordSizeBits === 0, `Key size must be multiple of uint256's`)
const keySizeWords = keySizeBits / wordSizeBits
const keySizeBytes = keySizeBits / 8

// A roughly 4096-bit prime, per https://2ton.com.au/getprimes/random/4096
const prime = bigNum(
  '0x' +
    'a8b0f26b02ce8f196aff5f98ec5a449fada1e4d493241397616b71221e117219' +
    'ce12fb8e9ef465b7a8ed29787313ef6dde3696dc7c36aadad4734104bde8295a' +
    'f5a9800e6576d4467b37032e45eeff3e23dc9bfc0ee02716f6b84596bb65ce11' +
    '13f046e060a5d81eb78557d665f621c8f6452dff50679caaccd8597eb18baa40' +
    'edb3b6873a431732e905dd45bfc9c86b935b7c9c56879656f0affd610a18a328' +
    'a86036a2b3af780520fa64decd008a03836b22c7752d7bc21d89a7fddd70ae4d' +
    'ed319fd1e38670444ed80af2eb7cee6a37f6cb46caa036e7b50bcb3864334e4f' +
    '66d3b3208b923116bec2a48e6749b00546f0ef7a73cee802cec690d84fd21d66' +
    '5725051fd825da5caaa7fadd158f245d0267e0836a91ccb13a5d1e48cfd75d5b' +
    'fa036336a18e4d4835229bb2716c188a299e7fc3cabbe3839d6511a52ae1b8b8' +
    '22e62330249f9f4115347b2000b58458906836c9dc54a0f2a83a788c094171f0' +
    '1fb7002f6165cde22b09973dfbf4db594c6a90a05929bb3a7522d9776b97e019' +
    '6d7049ce7cfe9469633859eb123c67849e9f80521b470cf38a532aba7985ce0a' +
    '66f8c8e4afa0f9861a819b6a91cbe0332fda69d910fb3994eb60f5db938b7485' +
    '4ed774f018ed36299668915f71520c8f8629bc8f280f36d38c4782d0b6a8232f' +
    'cf9650052caf2c3879ef7cfd3eec1f9c8b10d8aaa8103dd56aaeca80fce9441b')

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
