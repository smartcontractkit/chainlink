import {Tests} from "./types.mjs";

export function sieveSlowTests(tests: Tests, slowTestMatchers?: string[]) {
  const slowTests: Tests = [];
  const filteredTests: Tests = [];

  if (!slowTestMatchers) {
    return {slowTests, filteredTests: tests};
  }

  // If the user supplies slow test matchers
  // then we go through each test to see if we get a case sensitive match

  tests.forEach((t) => {
    const isSlow = slowTestMatchers.reduce(
      (isSlow, matcher) => t.includes(matcher) || isSlow,
      false
    );
    if (isSlow) {
      slowTests.push(t);
    } else {
      filteredTests.push(t);
    }
  });

  return {slowTests, filteredTests};
}
