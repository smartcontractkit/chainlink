// SPDX-License-Identifier: MIT
// Adapted from https://github.com/OpenZeppelin/openzeppelin-contracts/blob/c9630526e24ba53d9647787588a19ffaa3dd65e1/test/math/SignedSafeMath.test.js

import { ethers } from "hardhat";
import { assert } from "chai";
import { BigNumber, constants, Contract, ContractFactory } from "ethers";
import { Personas, getUsers } from "../test-helpers/setup";
import { bigNumEquals } from "../test-helpers/matchers";

let mathFactory: ContractFactory;
let personas: Personas;

before(async () => {
  personas = (await getUsers()).personas;
  mathFactory = await ethers.getContractFactory("CheckedMathTestHelper", personas.Default);
});

const int256Max = constants.MaxInt256;
const int256Min = constants.MinInt256;

describe("CheckedMath", () => {
  let math: Contract;

  beforeEach(async () => {
    math = await mathFactory.connect(personas.Default).deploy();
  });

  describe("#add", () => {
    const a = BigNumber.from("1234");
    const b = BigNumber.from("5678");

    it("is commutative", async () => {
      const c1 = await math.add(a, b);
      const c2 = await math.add(b, a);

      bigNumEquals(c1.result, c2.result);
      assert.isTrue(c1.ok);
      assert.isTrue(c2.ok);
    });

    it("is commutative with big numbers", async () => {
      const c1 = await math.add(int256Max, int256Min);
      const c2 = await math.add(int256Min, int256Max);

      bigNumEquals(c1.result, c2.result);
      assert.isTrue(c1.ok);
      assert.isTrue(c2.ok);
    });

    it("returns false when overflowing", async () => {
      const c1 = await math.add(int256Max, 1);
      const c2 = await math.add(1, int256Max);

      bigNumEquals(0, c1.result);
      bigNumEquals(0, c2.result);
      assert.isFalse(c1.ok);
      assert.isFalse(c2.ok);
    });

    it("returns false when underflowing", async () => {
      const c1 = await math.add(int256Min, -1);
      const c2 = await math.add(-1, int256Min);

      bigNumEquals(0, c1.result);
      bigNumEquals(0, c2.result);
      assert.isFalse(c1.ok);
      assert.isFalse(c2.ok);
    });
  });

  describe("#sub", () => {
    const a = BigNumber.from("1234");
    const b = BigNumber.from("5678");

    it("subtracts correctly if it does not overflow and the result is negative", async () => {
      const c = await math.sub(a, b);
      const expected = a.sub(b);

      bigNumEquals(expected, c.result);
      assert.isTrue(c.ok);
    });

    it("subtracts correctly if it does not overflow and the result is positive", async () => {
      const c = await math.sub(b, a);
      const expected = b.sub(a);

      bigNumEquals(expected, c.result);
      assert.isTrue(c.ok);
    });

    it("returns false on overflow", async () => {
      const c = await math.sub(int256Max, -1);

      bigNumEquals(0, c.result);
      assert.isFalse(c.ok);
    });

    it("returns false on underflow", async () => {
      const c = await math.sub(int256Min, 1);

      bigNumEquals(0, c.result);
      assert.isFalse(c.ok);
    });
  });

  describe("#mul", () => {
    const a = BigNumber.from("5678");
    const b = BigNumber.from("-1234");

    it("is commutative", async () => {
      const c1 = await math.mul(a, b);
      const c2 = await math.mul(b, a);

      bigNumEquals(c1.result, c2.result);
      assert.isTrue(c1.ok);
      assert.isTrue(c2.ok);
    });

    it("multiplies by 0 correctly", async () => {
      const c = await math.mul(a, 0);

      bigNumEquals(0, c.result);
      assert.isTrue(c.ok);
    });

    it("returns false on multiplication overflow", async () => {
      const c = await math.mul(int256Max, 2);

      bigNumEquals(0, c.result);
      assert.isFalse(c.ok);
    });

    it("returns false when the integer minimum is negated", async () => {
      const c = await math.mul(int256Min, -1);

      bigNumEquals(0, c.result);
      assert.isFalse(c.ok);
    });
  });

  describe("#div", () => {
    const a = BigNumber.from("5678");
    const b = BigNumber.from("-5678");

    it("divides correctly", async () => {
      const c = await math.div(a, b);

      bigNumEquals(a.div(b), c.result);
      assert.isTrue(c.ok);
    });

    it("divides a 0 numerator correctly", async () => {
      const c = await math.div(0, a);

      bigNumEquals(0, c.result);
      assert.isTrue(c.ok);
    });

    it("returns complete number result on non-even division", async () => {
      const c = await math.div(7000, 5678);

      bigNumEquals(1, c.result);
      assert.isTrue(c.ok);
    });

    it("reverts when 0 is the denominator", async () => {
      const c = await math.div(a, 0);

      bigNumEquals(0, c.result);
      assert.isFalse(c.ok);
    });

    it("reverts on underflow with a negative denominator", async () => {
      const c = await math.div(int256Min, -1);

      bigNumEquals(0, c.result);
      assert.isFalse(c.ok);
    });
  });
});
