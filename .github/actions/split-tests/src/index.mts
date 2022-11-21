import {$, cd, glob, fs} from "zx";
import path from "node:path";
import {summary, setOutput} from "@actions/core";
import {
  GolangConfig,
  SolidityConfig,
  GoSplits,
  Tests,
  SoliditySplit,
  TestsBySplit,
} from "./types.mjs";
import {sieveSlowTests} from "./sieve.mjs";
import {simpleSplit} from "./splitter.mjs";

/**
 * Get a JSON formatted config file
 *
 * @param path The path to the config relative to the git root
 */
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
  const {numOfSplits} = config;
  const rawPackages = await $`go list ./...`;
  const packages = rawPackages.stdout.trimEnd().split("\n");
  console.log(`${packages.length} packages to split...`);
  const packagesBySplit = simpleSplit(packages, [], numOfSplits);
  const splits: GoSplits = packagesBySplit.map((pkgs, i) => ({
    idx: `${i + 1}`,
    id: `${i + 1}/${numOfSplits}`,
    pkgs: pkgs.join(" "),
  }));
  const serializedSplits = JSON.stringify(splits);
  setOutput("splits", serializedSplits);
  createSummary(packages, packagesBySplit, splits);
}

async function handleSolidity(config: SolidityConfig) {
  const {basePath, splits: configBySplit} = config;
  const splits = await Promise.all(
    configBySplit.map(
      async ({dir, numOfSplits, slowTests: slowTestMatchers}) => {
        const globPath = path.join(basePath, dir, "/**/*.test.ts");
        const rawTests = await glob(globPath);
        const pathMappedTests = rawTests.map((r) =>
          r.replace("contracts/", "")
        );
        const {filteredTests, slowTests} = sieveSlowTests(
          pathMappedTests,
          slowTestMatchers
        );
        const testsBySplit = simpleSplit(filteredTests, slowTests, numOfSplits);
        const splits: SoliditySplit[] = testsBySplit.map((tests, i) => ({
          idx: `${dir}_${i + 1}`,
          id: `${dir} ${i + 1}/${numOfSplits}`,
          tests: tests.join(" "),
          coverageTests:
            tests.length === 1 ? tests.join(",") : `{${tests.join(",")}}`,
        }));
        return splits;
      }
    )
  );

  const serializedSplits = JSON.stringify(splits.flat());
  setOutput("splits", serializedSplits);
}

function createSummary(
  packages: Tests,
  packagesBySplit: TestsBySplit,
  splits: GoSplits
) {
  if (!process.env.CI) {
    return;
  }
  const numberOfPackages = packages.length;
  const numberOfSplits = packagesBySplit.length;
  const postProcessedNumberOfPackages = packagesBySplit.flat().length;

  summary
    .addHeading("Spliting Summary")
    .addHeading(
      `Number of packages from "go list ./...": ${numberOfPackages}`,
      3
    )
    .addHeading(
      `Number of packages placed into splits: ${postProcessedNumberOfPackages}`,
      3
    )
    .addHeading(`Number of splits created: ${numberOfSplits}`, 3)
    .addBreak()
    .addTable([
      [
        {data: "Split Number", header: true},
        {data: "Packages Tested", header: true},
      ],
      ...splits.map((p) => {
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
