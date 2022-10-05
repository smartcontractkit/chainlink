import {$, cd, glob, fs} from "zx";
import path from "node:path";
import {summary, setOutput} from "@actions/core";

/**
 * An array of all tests
 */
type Tests = string[];
/**
 * An array of tests, indexed by shard
 */
type TestsByShard = string[][];
interface Shard {
  /**
   * The shard index
   * @example "4"
   */
  idx: string;

  /**
   * The shard index in the context of all shards
   * @example "4/10"
   */
  id: string;
}
interface GoShard extends Shard {
  /**
   * A space delimited list of packages to run within this shard
   */
  pkgs: string;
}
interface SolidityShard extends Shard {
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
type GoShards = GoShard[];
interface GolangConfig {
  type: "golang";
  numOfShards: number;
}
interface SolidityConfig {
  type: "solidity";
  basePath: string;
  shards: {
    parallelism: number;
    dir: string;
  }[];
}

function getConfigFrom(path?: string): GolangConfig | SolidityConfig {
  if (!path) {
    throw Error("No config path given, specify a path via $CONFIG");
  }
  try {
    const config = fs.readJsonSync(path);
    return config;
  } catch (e: unknown) {
    throw Error(
      `Could not find config file at path: ${path}. ${(e as Error).message}`
    );
  }
}

async function main() {
  $.verbose = false;
  await runAtGitRoot();
  const configPath = process.env.CONFIG;
  const config = getConfigFrom(configPath);
  if (config.type === "golang") {
    await handleGolang(config);
  } else if (config.type === "solidity") {
    await handleSolidity(config);
  } else {
    throw Error(`Invalid config given`);
  }
}
main();

async function handleGolang(config: GolangConfig) {
  const {numOfShards} = config;
  const rawPackages = await $`go list ./...`;
  const packages = rawPackages.stdout.trimEnd().split("\n");
  console.log(`${packages.length} packages to shard...`);
  const packagesByShard = simpleShard(packages, numOfShards);
  const shards: GoShards = packagesByShard.map((pkgs, i) => ({
    idx: `${i + 1}`,
    id: `${i + 1}/${numOfShards}`,
    pkgs: pkgs.join(" "),
  }));
  const serializedShards = JSON.stringify(shards);
  setOutput("shards", serializedShards);
  createSummary(packages, packagesByShard, shards);
}

async function handleSolidity(config: SolidityConfig) {
  const {basePath, shards: configByShard} = config;
  const shards = await Promise.all(
    configByShard.map(async ({dir, parallelism}) => {
      const globPath = path.join(basePath, dir, "/**/*.test.ts");
      const rawTests = await glob(globPath);
      const tests = rawTests.map((r) => r.replace("contracts/", ""));
      const testsByShard = simpleShard(tests, parallelism);

      const shards: SolidityShard[] = testsByShard.map((tests, i) => ({
        idx: `${dir}_${i + 1}`,
        id: `${dir} ${i + 1}/${parallelism}`,
        tests: tests.join(" "),
        coverageTests:
          tests.length === 1 ? tests.join(",") : `{${tests.join(",")}}`,
      }));
      return shards;
    })
  );

  const serializedShards = JSON.stringify(shards.flat());
  setOutput("shards", serializedShards);
}

/**
 * A simple sharding strategy that fills each shard by the test list order
 *
 * @param tests The tests to shard
 * @param numOfShards The number of shards to create
 * @returns
 */
function simpleShard(tests: Tests, numOfShards: number): TestsByShard {
  const maxTestsPerShard = Math.max(tests.length / numOfShards);

  const testsByShard: TestsByShard = new Array(numOfShards)
    .fill(null)
    .map(() => []);

  tests.forEach((test, i) => {
    const shardIndex = Math.floor(i / maxTestsPerShard);
    testsByShard[shardIndex].push(test);
  });

  return testsByShard;
}

function createSummary(
  packages: Tests,
  packagesByShard: TestsByShard,
  shards: GoShards
) {
  if (!process.env.CI) {
    return;
  }
  const numberOfPackages = packages.length;
  const numberOfShards = packagesByShard.length;
  const postProcessedNumberOfPackages = packagesByShard.flat().length;

  summary
    .addHeading("Sharding Summary")
    .addHeading(
      `Number of packages from "go list ./...": ${numberOfPackages}`,
      3
    )
    .addHeading(
      `Number of packages placed into shards: ${postProcessedNumberOfPackages}`,
      3
    )
    .addHeading(`Number of shards created: ${numberOfShards}`, 3)
    .addBreak()
    .addTable([
      [
        {data: "Shard Number", header: true},
        {data: "Packages Tested", header: true},
      ],
      ...shards.map((p) => {
        const mappedPackages = p.pkgs
          .split(" ")
          .map(
            (packageName) =>
              `<li> ${packageName.replace(
                "github.com/smartcontractkit/",
                ""
              )} </li>`
          )
          .join("\n");

        return [p.id, mappedPackages];
      }),
    ])
    .write();
}

async function runAtGitRoot() {
  const gitRoot = await $`git rev-parse --show-toplevel`;
  cd(gitRoot.stdout.trimEnd());
}
