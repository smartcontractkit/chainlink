/*
 * Based off of https://github.com/changesets/changesets/blob/7323704dff6e76f488370db384579b86c95c866f/packages/changelog-github/src/index.ts
 */
const { getCommitInfo, getPullRequestInfo } = require("@changesets/get-github-info");

const getDependencyReleaseLine = async (changesets, dependenciesUpdated, options) => {
  if (dependenciesUpdated.length === 0) return "";
  if (!options || !options.repo) {
    throw new Error(
      'Please provide a repo to this changelog generator like this:\n"changelog": ["@changesets/changelog-github", { "repo": "org/repo" }]'
    );
  }

  const changesetLink = `- Updated dependencies [${(
    await Promise.all(
      changesets.map(async (cs) => {
        if (cs.commit) {
          const info = await getCommitInfo({
            repo: options.repo,
            commit: cs.commit,
          });
          if (!info) {
            throw new Error(
              `Could not get GitHub info for commit ${cs.commit} in ${options.repo}`
            );
          }
          return info.commit.markdownLink;
        }
      })
    )
  )
    .filter((_) => _)
    .join(", ")}]:`;


  const updatedDepsList = dependenciesUpdated.map(
    (dependency) => `  - ${dependency.name}@${dependency.newVersion}`
  );

  return [changesetLink, ...updatedDepsList].join("\n");
};

const getReleaseLine = async (changeset, _, options) => {
  if (!options || !options.repo) {
    throw new Error(
      'Please provide a repo to this changelog generator like this:\n"changelog": ["@changesets/changelog-github", { "repo": "org/repo" }]'
    );
  }

  let prFromSummary;
  let commitFromSummary;

  const replacedChangelog = changeset.summary
    .replace(/^\s*(?:pr|pull|pull\s+request):\s*#?(\d+)/im, (_, pr) => {
      let num = Number(pr);
      if (!isNaN(num)) prFromSummary = num;
      return "";
    })
    .replace(/^\s*commit:\s*([^\s]+)/im, (_, commit) => {
      commitFromSummary = commit;
      return "";
    })
    .trim();

  const [firstLine, ...futureLines] = replacedChangelog
    .split("\n")
    .map((l) => l.trimRight());

  const links = await (async () => {
    if (prFromSummary !== undefined) {
      const info = await getPullRequestInfo({
        repo: options.repo,
        pull: prFromSummary,
      });
      if (!info) {
        throw new Error(
          `Could not get GitHub info for PR #${prFromSummary} in ${options.repo}`
        );
      }
      const links = {
        commit: info.commit?.markdownLink ?? null,
        pull: info.pull.markdownLink,
      };
      if (commitFromSummary) {
        const shortCommitId = commitFromSummary.slice(0, 7);
        links.commit = `[\`${shortCommitId}\`](https://github.com/${options.repo}/commit/${commitFromSummary})`;
      }
      return links;
    }
    const commitToFetchFrom = commitFromSummary || changeset.commit;
    if (commitToFetchFrom) {
      const info = await getCommitInfo({
        repo: options.repo,
        commit: commitToFetchFrom,
      });
      if (!info) {
        throw new Error(
          `Could not get GitHub info for commit ${commitToFetchFrom} in ${options.repo}`
        );
      }
      return {
        commit: info.commit.markdownLink,
        pull: info.pull?.markdownLink ?? null,
      };
    }
    return {
      commit: null,
      pull: null,
    };
  })();

  const prefix = [
    links.pull === null ? "" : ` ${links.pull}`,
    links.commit === null ? "" : ` ${links.commit}`,
  ].join("");

  return `\n\n-${prefix ? `${prefix} -` : ""} ${firstLine}\n${futureLines
    .map((l) => `  ${l}`)
    .join("\n")}`;
};

module.exports = { getReleaseLine, getDependencyReleaseLine };
