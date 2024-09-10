import * as core from "@actions/core";
import { join } from "path";
import { getGitTopLevel } from "./lib";
import { promises as fs } from "fs";

export async function appendIssueNumberToChangesetFile(
    prefix: string,
    changesetFile: string,
    issueNumber: string
  ) {
    const gitTopLevel = await getGitTopLevel();
    const fullChangesetPath = join(gitTopLevel, changesetFile);
    const changesetContents = await fs.readFile(fullChangesetPath, "utf-8");
    // Check if the issue number is already in the changeset file
    if (changesetContents.includes(issueNumber)) {
      core.info("Issue number already exists in changeset file, skipping...");
      return;
    }

    const updatedChangesetContents = `${changesetContents}\n\n${prefix}${issueNumber}`;
    await fs.writeFile(fullChangesetPath, updatedChangesetContents);
  }

  export function extractChangesetFile() {
    const changesetFiles = process.env.CHANGESET_FILES;
    if (!changesetFiles) {
      throw Error("Missing required environment variable CHANGESET_FILES");
    }

    const parsedChangesetFiles = JSON.parse(changesetFiles);
    if (parsedChangesetFiles.length !== 1) {
      throw Error(
        "This action only supports one changeset file per pull request."
      );
    }
    const [changesetFile] = parsedChangesetFiles;

    return { changesetFile };
  }