import {Tests, TestsBySplit} from "./types.mjs";

/**
 * Split tests by first prioritizing slow tests being spread over each split, then filling each split by test list order.
 *
 * @example
 * Given the following arguments:
 * tests: ['foo.test', 'bar.test', 'baz.test', 'yup.test', 'nope.test']
 * slowTests: ['bonk.test', 'bop.test', 'ouch.test.ts']
 * numOfSplits: 2
 *
 * We get the following output:
 * 1. Spread slow tests across splits: [['bonk.test', 'ouch.test.ts'], ['bop.test']]
 * 2. Insert list of tests: [['bonk.test', 'ouch.test.ts', 'foo.test', 'bar.test'], ['bop.test', 'baz.test', 'yup.test', 'nope.test']]
 *
 * @param tests A list of tests to distribute across splits by the test list order
 * @param slowTests A list of slow tests, where the list of tests is evenly distributed across all splits before inserting regular tests
 * @param numOfSplits The number of splits to spread tests across
 */
export function simpleSplit(
  tests: Tests,
  slowTests: Tests,
  numOfSplits: number
): TestsBySplit {
  const maxTestsPerSplit = Math.max(tests.length / numOfSplits);

  const testsBySplit: TestsBySplit = new Array(numOfSplits)
    .fill(null)
    .map(() => []);

  // Evenly distribute slow tests over each split
  slowTests.forEach((test, i) => {
    const splitIndex = i % numOfSplits;
    testsBySplit[splitIndex].push(test);
  });

  tests.forEach((test, i) => {
    const splitIndex = Math.floor(i / maxTestsPerSplit);
    testsBySplit[splitIndex].push(test);
  });

  return testsBySplit;
}
