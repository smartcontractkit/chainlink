#!/usr/bin/env node
const { execSync } = require("child_process");

function main() {
  const goreleaserConfig = mustGetEnv("GORELEASER_CONFIG");
  const releaseType = mustGetEnv("RELEASE_TYPE");
  const command = constructGoreleaserCommand(releaseType, goreleaserConfig);

  if (process.env.DRY_RUN) {
    console.log(`Generated command: ${command}`);
    console.log("Dry run enabled. Exiting without executing the command.");
    return;
  } else {
    console.log(`Executing command: ${command}`);
    execSync(command, { stdio: "inherit" });
  }
}

main();

function constructGoreleaserCommand(releaseType, goreleaserConfig) {
  const version = getVersion();
  const flags = [];

  checkReleaseType(releaseType);

  let subCmd = "release";
  const splitArgs = ["--split", "--clean"];

  switch (releaseType) {
    case "release":
      flags.push(...splitArgs);
      break;
    case "nightly":
      flags.push("--nightly", ...splitArgs);
      break;
    case "snapshot":
      flags.push("--snapshot", ...splitArgs);
      break;
    case "merge":
      flags.push("--merge");
      subCmd = "continue";
      break;
  }

  const flagsStr = flags.join(" ");
  if (releaseType === "merge") {
    return `CHAINLINK_VERSION=${version} goreleaser ${subCmd} ${flagsStr}`;
  } else {
    return `CHAINLINK_VERSION=${version} goreleaser ${subCmd} --config ${goreleaserConfig} ${flagsStr}`;
  }
}

function checkReleaseType(releaseType) {
  const VALID_RELEASE_TYPES = ["nightly", "merge", "snapshot", "release"];

  if (!VALID_RELEASE_TYPES.includes(releaseType)) {
    const validReleaseTypesStr = VALID_RELEASE_TYPES.join(", ");
    console.error(
      `Error: Invalid release type: ${releaseType}. Must be one of: ${validReleaseTypesStr}`
    );
  }
}

function mustGetEnv(key) {
  const val = process.env[key];
  if (!val || val.trim() === "") {
    console.error(`Error: Environment variable ${key} is not set or empty.`);
    process.exit(1);
  }

  return val.trim();
}

function getVersion() {
  try {
    const pkgPath = process.cwd() + "/package.json";
    console.log("Looking for chainlink version in package.json at: ", pkgPath);
    const packageJson = require(pkgPath);
    if (!packageJson.version) {
      console.error(
        'Error: "version" field is missing or empty in package.json.'
      );
      process.exit(1);
    }
    console.log("Resolved version: ", packageJson.version);

    return packageJson.version;
  } catch (err) {
    console.error(`Error reading package.json: ${err.message}`);
    process.exit(1);
  }
}
