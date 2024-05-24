import {sieveSlowTests} from "../src/sieve.mjs";
import {testArr} from "./fixtures.mjs";

describe("sieveSlowTests", () => {
  it("works", () => {
    expect(sieveSlowTests([])).toMatchSnapshot();
    expect(sieveSlowTests([], [])).toMatchSnapshot();
    expect(sieveSlowTests(["keepme"], [])).toMatchSnapshot();
    expect(sieveSlowTests(["keepme"])).toMatchSnapshot();
    expect(sieveSlowTests(testArr, [])).toMatchSnapshot();
    expect(sieveSlowTests(testArr, ["noself"])).toMatchSnapshot();
    expect(sieveSlowTests(testArr, ["ouch.test.ts"])).toMatchSnapshot();
    expect(sieveSlowTests(testArr, ["bo", "ouch.test.ts"])).toMatchSnapshot();
  });
});
