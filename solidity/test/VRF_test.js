import { deploy, bigNum } from './support/helpers'
import { assertBigNum } from './support/matchers'

// Group elements are {(x,y) in GF(fieldSize)^2 | y^2=x^3+3}, where
// GF(fieldSize) is arithmetic modulo fieldSize on {0, 1, ..., fieldSize-1}
const fieldSize = bigNum(
  "21888242871839275222246405745257275088696311157297823662689037894645226208583"
);
const groupOrder = bigNum( // #{(x,y) in GF(fieldSize)^2 | y^2=x^3+3}
  "21888242871839275222246405745257275088548364400416034343698204186575808495617"
);
const generator = [1, 2].map(bigNum) // Point in EC group
const gFPNeg = n => fieldSize.sub(bigNum(n)) // Additive inverse in GF(fieldSize)
const minusGenerator = [1, gFPNeg(2)].map(bigNum) // (1,-2)
const big2 = bigNum(2)
const big3 = bigNum(3)


contract('VRF', () => {
  let VRF
  beforeEach(async () => { VRF = await deploy('VRF.sol') })
  it('Accurately calculates simple and obvious bigModExpTest inputs',
     async () => { assertBigNum(await VRF.bigModExp(3, 2, 5), (3**2) % 5) })
  it('Accurately calculates the sum of P and -P', async () => {
    const result = await VRF.addPoints(generator, minusGenerator)
    assertBigNum(0, result[0])
    assertBigNum(0, result[1])
  })
  it('Accurately multiplies a point by a scalar', async () => {
    const oneLessThanGroupSize = groupOrder.sub(bigNum(1)).toString()
    const negative = await VRF.scalarMul(generator, oneLessThanGroupSize)
    assertBigNum(minusGenerator[0], negative[0], "(p-1)*g = -g")
    assertBigNum(minusGenerator[1], negative[1], "(p-1)*g = -g")
  })
  it('Recognizes basic squares and non-squares in GF(fieldSize)', async () => {
    assert(!(await VRF.isSquare(gFPNeg(1).toString())), "-1 is not a square")
    assert(await VRF.isSquare(4), "4=2^2")
  })
  it('Can compute square roots', async () => {
    assertBigNum(2, await VRF.squareRoot(4), "4=2^2") // 4**((fieldSize-1)/2)
  })
  it('Can compute the square of the y ordinate given the x ordinate', async () => {
    assertBigNum(4, await(VRF.ySquared(1)), "4=1^3+3")
  })
  it('Can recognize valid and invalid x ordinates', async () => {
    assert(await VRF.isCurveXOrdinate(1), "2^2=1^3+3")
    assert(!(await VRF.isCurveXOrdinate(4)), "∄ y s.t. y^2=4^3+1")
  })
  it('Hashes to the curve with the same results as the golang code', async () => {
    let result = await VRF.hashToCurve(generator, 0)
    assertBigNum(bigNum(result[0]).pow(big3).add(big3).umod(fieldSize),
                 bigNum(result[1]).pow(big2).umod(fieldSize), "y^2=x^3+3")
    // See 
    result = await VRF.hashToCurve(generator, 5)
    assertBigNum(result[0], "0x247154f2ce523897365341b03669e1061049e801e8750ae708e1cb02f36cb225",
                 "mismatch with output from services/vrf/vrf_test.go/TestVRF_HashToCurve")
    assertBigNum(result[1], "0x16e1157d5b94324127e094abe222a05a5c47be3124254a6aa047d5e1f2d864ea",
                 "mismatch with output from services/vrf/vrf_test.go/TestVRF_HashToCurve")
  })
  it('Correctly computes linear combinations', async () => {
    const zero = await VRF.linearCombination(50, generator, 50, minusGenerator)
    assertBigNum(zero[0], 0, "50*g + 50*(-g) = 0")
    assertBigNum(zero[1], 0, "50*g + 50*(-g) = 0")
  })
  it('Computes the same hashed scalar from curve points as the golang code', async () => {
    assertBigNum("0x57bf013147ceec913f17ef97d3bcfad8315d99752af81f8913ad1c88493e669",
                 await VRF.scalarFromCurve(generator, generator, generator, generator, generator),
                 "mismatch with output from services/vrf/vrf_test.go/TestVRF_ScalarFromCurve")
  })
  it('Rejects an invalid VRF proof', async () => {
    assert(!(await VRF.isValidVRFOutput(generator, generator, 1, 1, 1, 0)))
  })
  it('Accepts a valid VRF output', async () => {
    const secretKey = 2
    const publicKey = await VRF.scalarMul(generator, secretKey)
    const seed = 0
    const gamma = ["0x26feb384a4a3f28742d0e0e0f5458474ba54ef9816d4d31f3bf538dfcf67cf3f",
                   "0x1eaed2431dd78ad75dd0c9f013cabff4f1d8c4c83cda79fff3855c988a3606d8"]
    const c = "0x1826029ee5a1a03cc4c58e78085c9f4daff1c7474f78f01443e463c658d85368"
    const s = "0x2588663566d33936a96b6f09bbd5c669172249d220d05c038584df07bff6f777"
    assert(await VRF.verifyVRFProof(publicKey, gamma, c, s, seed),
           `Could not validate a proof output by services/vrf/vrf.go/GenerateProof.
These proof values correspond to a blinding value ("m" in vrf.go/GenerateProof) of
0x25701d0050e4d9867aa646434b0daca74ed1f0184608cb9ac96bb10081a79e46`)
  })
})
