/**
 * An array of all tests
 */
export type Tests = string[];

/**
 * An array of tests, indexed by split
 */
export type TestsBySplit = string[][];

export interface Split {
  /**
   * The split index
   * @example "4"
   */
  idx: string;

  /**
   * The split index in the context of all splits
   * @example "4/10"
   */
  id: string;
}

export interface GoSplit extends Split {
  /**
   * A space delimited list of packages to run within this split
   */
  pkgs: string;
}

export interface SoliditySplit extends Split {
  /**
   * A string that contains a whitespace delimited list of tests to run
   *
   * This format is to support the `hardhat test` command.
   * @example test/foo.test.ts test/bar.test.ts
   */
  tests: string;

  /**
   * A string that contains a glob that expresses the list of tests to run.
   *
   * This format is used to conform to the --testfiles flag of solidity-coverage
   * @example {test/foo.test.ts,test/bar.test.ts}
   */
  coverageTests: string;
}

export type GoSplits = GoSplit[];

/**
 * Configuration file for golang tests
 */
export interface GolangConfig {
  type: "golang";
  /**
   * The number of splits to run tests across
   */
  numOfSplits: number;
}

/**
 * Configuration file for solidity tests
 */
export interface SolidityConfig {
  type: "solidity";
  /**
   * The path to the contracts tests directory, relative to the git root
   */
  basePath: string;
  splits: {
    /**
     * The number of sub-splits to run across
     */
    numOfSplits: number;
    /**
     * The directory of the tests to create sub-splits across, relative to the basePath
     */
    dir: string;
    /**
     * An array of known slow tests, to better distribute across sub-splits
     *
     * Each string is a case-sensitive matcher that will match against any substring within the list of test file paths within the `dir` configuration.
     *
     * @example
     * Given the dir `v0.8`, we get the following tests: ['v0.8/Foo1.test.ts','v0.8/bar.test.ts','v0.8/dev/eolpe/Foo.test.ts']
     *
     * If we supply the following `slowTests` argument: ['Foo']
     *
     * Then it'll match against both 'v0.8/Foo1.test.ts' and 'v0.8/dev/eolpe/Foo.test.ts'.
     */
    slowTests?: string[];
  }[];
}
