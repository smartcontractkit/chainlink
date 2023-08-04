import { execSync } from "child_process";
import {
  GolangConfig,
  GoSplits,
  GoPackageData,
  TestsBySplit,
} from "../types.mjs";
import { simpleSplit } from "../splitter.mjs";

export interface GetGoPackagesReturn {
  packages: string[];
  testsBySplit: TestsBySplit;
  splits: GoSplits;
  serializedSplits: string
}

export function getPackageList(
  config: GolangConfig,
): GetGoPackagesReturn {
  const { numOfSplits } = config;
  const rawPackages = execSync(
    "go list -json ./... | jq -s '[.[] | {ImportPath, TestGoFiles}]'",
    { encoding: "utf8" }
  );
  const packages: GoPackageData[] = JSON.parse(rawPackages.trimEnd());
  const packagePaths = packages.map((item) => item.ImportPath);
  return handleSplit(packagePaths, numOfSplits);
}

function handleSplit(
  packages: string[],
  numOfSplits: number
): GetGoPackagesReturn {
  console.log(`${packages.length} packages to split...`);
  const packagesBySplit = simpleSplit(packages, [], numOfSplits);
  const splits: GoSplits = packagesBySplit.map((pkgs, i) => ({
    idx: `${i + 1}`,
    id: `${i + 1}/${numOfSplits}`,
    pkgs: pkgs.join(" "),
  }));
  const o: GetGoPackagesReturn = {
    packages,
    testsBySplit: packagesBySplit,
    splits,
    serializedSplits: JSON.stringify(splits),
  };
  return o;
}
