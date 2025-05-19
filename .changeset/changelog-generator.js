/*
 * Based off of https://github.com/changesets/changesets/blob/7323704dff6e76f488370db384579b86c95c866f/packages/changelog-github/src/index.ts
 */

const ghInfo = require("@changesets/get-github-info");

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
          let { links } = await ghInfo.getInfo({
            repo: options.repo,
            commit: cs.commit,
          });
          return links.commit;
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
      let { links } = await ghInfo.getInfoFromPullRequest({
        repo: options.repo,
        pull: prFromSummary,
      });
      if (commitFromSummary) {
        const shortCommitId = commitFromSummary.slice(0, 7);
        links = {
          ...links,
          commit: `[\`${shortCommitId}\`](https://github.com/${options.repo}/commit/${commitFromSummary})`,
        };
      }
      return links;
    }
    const commitToFetchFrom = commitFromSummary || changeset.commit;
    if (commitToFetchFrom) {
      let { links } = await ghInfo.getInfo({
        repo: options.repo,
        commit: commitToFetchFrom,
      });
      return links;
    }
    return {
      commit: null,
      pull: null,
      user: null,
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
