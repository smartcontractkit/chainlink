import { $, cd, glob, fs } from "zx";
import path from "node:path";
import { setOutput } from "@actions/core";
import { SolidityConfig, SoliditySplit } from "./types.mjs";
import { sieveSlowTests } from "./sieve.mjs";
import { simpleSplit } from "./splitter.mjs";

/**
 * Get a JSON formatted config file
 *
 * @param path The path to the config relative to the git root
 */
function getConfigFrom(path?: string): SolidityConfig {
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
  if (config.type === "solidity") {
    await handleSolidity(config);
  } else {
    throw Error(`Invalid config given`);
  }
}
main();

async function handleSolidity(config: SolidityConfig) {
  const { basePath, splits: configBySplit } = config;
  const splits = await Promise.all(
    configBySplit.map(
      async ({ dir, numOfSplits, slowTests: slowTestMatchers }) => {
        const globPath = path.join(basePath, dir, "/**/*.test.ts");
        const rawTests = await glob(globPath);
        const pathMappedTests = rawTests.map((r) =>
          r.replace("contracts/", "")
        );
        const { filteredTests, slowTests } = sieveSlowTests(
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

async function runAtGitRoot() {
  const gitRoot = await $`git rev-parse --show-toplevel`;
  cd(gitRoot.stdout.trimEnd());
}
