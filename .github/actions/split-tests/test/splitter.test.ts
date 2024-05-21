import {simpleSplit} from "../src/splitter.mjs";
import {testArr, testSievedArr, testSlowArr} from "./fixtures.mjs";

describe("simpleSplit", () => {
  it("doesn't error on empty arrays", () => {
    expect(simpleSplit([], [], 1)).toMatchSnapshot();
    expect(simpleSplit([], [], 5)).toMatchSnapshot();
  });

  it("handles no slow test splitting", () => {
    expect(simpleSplit(testArr, [], 1)).toMatchSnapshot();
    expect(simpleSplit(testArr, [], 2)).toMatchSnapshot();
    expect(simpleSplit(testArr, [], 3)).toMatchSnapshot();
  });

  it("handles slow test splitting", () => {
    expect(simpleSplit(testSievedArr, testSlowArr, 1)).toMatchSnapshot();
    expect(simpleSplit(testSievedArr, testSlowArr, 2)).toMatchSnapshot();
    expect(simpleSplit(testSievedArr, testSlowArr, 3)).toMatchSnapshot();
  });
});
