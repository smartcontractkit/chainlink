#!/usr/bin/env node
const { execSync } = require("child_process");
const fs = require("fs");
const path = require("path");

function main() {
  const args = process.argv.slice(2);
  const useExistingDist = args.includes("--use-existing-dist");
  const chainlinkVersion = getVersion();

  if (!useExistingDist) {
    const goreleaserConfig = mustGetEnv("GORELEASER_CONFIG");
    const releaseType = mustGetEnv("RELEASE_TYPE");
    const command = constructGoreleaserCommand(
      releaseType,
      chainlinkVersion,
      goreleaserConfig
    );

    if (process.env.DRY_RUN) {
      console.log(`Generated command: ${command}`);
      console.log("Dry run enabled. Exiting without executing the command.");
      return;
    } else {
      console.log(`Executing command: ${command}`);
      execSync(command, { stdio: "inherit" });
    }
  } else {
    console.log(
      "Skipping Goreleaser command execution as '--use-existing-dist' is set."
    );
  }

  const artifacts = getArtifacts();
  const dockerImages = extractDockerImages(artifacts);
  const repoSha = execSync("git rev-parse HEAD", { encoding: "utf-8" }).trim();

  const results = dockerImages.map((image) => {
    try {
      console.log(
        `Checking version for image: ${image}, expected version: ${chainlinkVersion}, expected SHA: ${repoSha}`
      );
      const versionOutput = execSync(`docker run --rm ${image} --version`, {
        encoding: "utf-8",
      });
      console.log(`Output from image ${image}: ${versionOutput}`);

      const cleanedOutput = versionOutput
        .replace("chainlink version ", "")
        .trim();
      const [version, sha] = cleanedOutput.split("@");
      if (!version || !sha) {
        throw new Error("Version or SHA not found in output.");
      }

      if (sha.trim() !== repoSha) {
        throw new Error(`SHA mismatch: Expected ${repoSha}, got ${sha.trim()}`);
      }
      if (version.trim() !== chainlinkVersion) {
        throw new Error(
          `Version mismatch: Expected ${chainlinkVersion}, got ${version.trim()}`
        );
      }

      return { image, success: true, message: "Version check passed." };
    } catch (error) {
      return { image, success: false, message: error.message };
    }
  });

  printSummary(results);
  if (results.some((result) => !result.success)) {
    process.exit(1);
  }
}

function printSummary(results) {
  const passed = results.filter((result) => result.success);
  const failed = results.filter((result) => !result.success);

  console.log("\nSummary:");
  console.log(`Total images checked: ${results.length}`);
  console.log(`Passed: ${passed.length}`);
  console.log(`Failed: ${failed.length}`);

  if (passed.length > 0) {
    console.log("\nPassed images:");
    passed.forEach((result) =>
      console.log(`${result.image}:\n${result.message}`)
    );
  }

  if (failed.length > 0) {
    console.log("\nFailed images:");
    failed.forEach((result) =>
      console.log(`${result.image}:\n${result.message}`)
    );
  }
}

function getArtifacts() {
  const distDir = path.resolve(process.cwd(), "dist");
  const files = [];

  function findJsonFiles(dir) {
    const items = fs.readdirSync(dir, { withFileTypes: true });
    for (const item of items) {
      const fullPath = path.join(dir, item.name);
      if (item.isDirectory()) {
        // Skip child directories if an artifacts.json exists in the current directory
        const parentArtifacts = path.join(dir, "artifacts.json");
        if (fs.existsSync(parentArtifacts)) {
          console.log(
            `Skipping child directory: ${fullPath} because a parent artifacts.json exists at: ${parentArtifacts}`
          );
        } else {
          findJsonFiles(fullPath);
        }
      } else if (item.isFile() && item.name === "artifacts.json") {
        console.log(`Found artifacts.json at: ${fullPath}`);
        files.push(fullPath);
      }
    }
  }

  findJsonFiles(distDir);

  if (files.length === 0) {
    console.error("Error: No artifacts.json found in /dist.");
    process.exit(1);
  }

  // Merge all artifacts.json files into one
  let mergedArtifacts = [];

  for (const file of files) {
    const artifactsJson = JSON.parse(fs.readFileSync(file, "utf-8"));
    mergedArtifacts = mergedArtifacts.concat(artifactsJson);
  }

  // Remove duplicate Docker images based on the artifact name
  const uniqueArtifacts = Array.from(
    new Map(
      mergedArtifacts.map((artifact) => [artifact.name, artifact])
    ).values()
  );

  return uniqueArtifacts;
}

function extractDockerImages(artifacts) {
  const dockerImages = artifacts
    .filter(
      (artifact) =>
        artifact.type === "Docker Image" ||
        artifact.type === "Published Docker Image"
    )
    .map((artifact) => artifact.name);

  if (dockerImages.length === 0) {
    console.error("Error: No Docker images found in artifacts.json.");
    process.exit(1);
  }

  console.log(`Found Docker images:\n  - ${dockerImages.join("\n  - ")}`);
  return dockerImages;
}

function constructGoreleaserCommand(releaseType, version, goreleaserConfig) {
  const flags = [];
  const debugFlag = (process.env.DEBUG == 'true') ? '--verbose' : '';

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
    return `CHAINLINK_VERSION=${version} goreleaser ${debugFlag} ${subCmd} ${flagsStr}`;
  } else {
    return `CHAINLINK_VERSION=${version} goreleaser ${debugFlag} ${subCmd} --config ${goreleaserConfig} ${flagsStr}`;
  }
}

function checkReleaseType(releaseType) {
  const VALID_RELEASE_TYPES = ["nightly", "merge", "snapshot", "release"];

  if (!VALID_RELEASE_TYPES.includes(releaseType)) {
    const validReleaseTypesStr = VALID_RELEASE_TYPES.join(", ");
    console.error(
      `Error: Invalid release type: ${releaseType}. Must be one of: ${validReleaseTypesStr}`
    );
    process.exit(1);
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
    const pkgPath = path.resolve(process.cwd(), "package.json");
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

main();
