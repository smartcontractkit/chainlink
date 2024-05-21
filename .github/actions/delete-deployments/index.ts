import { Octokit } from "@octokit/action";
import { info, warning, isDebug } from "@actions/core";
import { throttling } from "@octokit/plugin-throttling";
import { retry } from "@octokit/plugin-retry";

async function main() {
  const {
    dryRun,
    environment,
    numOfPages,
    owner,
    ref,
    repo,
    debug,
    startingPage,
  } = getInputs();
  const octokit = getOctokit(debug);

  const deployments = await getDeployments({
    octokit,
    owner,
    repo,
    environment,
    ref,
    paginateOptions: {
      numOfPages,
      startingPage,
    },
  });
  const deploymentIds = deployments.map((d) => d.id);
  if (dryRun) {
    info(`Dry run: would delete deployments (${deploymentIds.length})`);
    return;
  }

  info(`Deleting deployments (${deploymentIds.length})`);
  const deleteDeployments = deploymentIds.map(async (id) => {
    const sharedArgs = {
      owner,
      repo,
      deployment_id: id,
      request: {
        retries: 0,
      },
    };

    const setStatus = await octokit.repos
      .createDeploymentStatus({
        ...sharedArgs,
        state: "inactive",
      })
      .then(() => true)
      .catch((e) => {
        warning(
          `Marking deployment id ${id} to "inactive" failed: ${e.message}`
        );
        return false;
      });
    if (!setStatus) return false;

    return octokit.repos
      .deleteDeployment({
        ...sharedArgs,
      })
      .then(() => true)
      .catch((e) => {
        warning(`Deleting deployment id ${id} failed: ${e.message}`);
        return false;
      });
  });

  const processed = await Promise.all(deleteDeployments);
  const succeeded = processed.filter((p) => !!p);
  info(
    `Successfully deleted ${succeeded.length}/${processed.length} deployments`
  );
}
main();

function getInputs() {
  const debug = !!(process.env.DEBUG || isDebug());

  const dryRun = process.env.DRY_RUN === "true";

  const environment = process.env.ENVIRONMENT;
  if (!environment) throw new Error("ENVIRONMENT not set");

  const ref = process.env.REF;

  const repository = process.env.REPOSITORY;
  if (!repository) throw new Error("REPOSITORY not set");
  const [owner, repo] = repository.split("/");

  const rawStartingPage = process.env.STARTING_PAGE;

  let startingPage: number | undefined;
  if (rawStartingPage) {
    startingPage = parseInt(rawStartingPage);
    if (isNaN(startingPage)) {
      throw new Error(`STARTING_PAGE is not a number: ${rawStartingPage}`);
    }
    if (startingPage < 0) {
      throw new Error(
        `STARTING_PAGE must be a positive integer or zero: ${rawStartingPage}`
      );
    }
    info(`Starting from page ${startingPage}`);
  }

  const rawNumOfPages = process.env.NUM_OF_PAGES;
  let numOfPages: "all" | number = "all";
  if (rawNumOfPages === "all") {
    info("Fetching all pages of deployments");
  } else {
    const parsedPages = parseInt(rawNumOfPages || "");
    if (isNaN(parsedPages)) {
      throw new Error(`NUM_OF_PAGES is not a number: ${rawNumOfPages}`);
    }
    if (parsedPages < 1) {
      throw new Error(`NUM_OF_PAGES must be greater than 0: ${rawNumOfPages}`);
    }
    numOfPages = parsedPages;
  }

  if (numOfPages === "all" && startingPage) {
    throw new Error(`Cannot use STARTING_PAGE with NUM_OF_PAGES=all`);
  }

  const parsedInputs = {
    environment,
    ref,
    owner,
    repo,
    numOfPages,
    startingPage,
    dryRun,
    debug,
  };
  info(`Configuration: ${JSON.stringify(parsedInputs)}`);
  return parsedInputs;
}

function getOctokit(debug: boolean) {
  const OctokitAPI = Octokit.plugin(throttling, retry);
  const octokit = new OctokitAPI({
    log: debug ? console : undefined,
    throttle: {
      onRateLimit: (retryAfter, options, octokit, retryCount) => {
        octokit.log.warn(
          // Types are busted from octokit
          //@ts-expect-error
          `Request quota exhausted for request ${options.method} ${options.url}`
        );

        octokit.log.info(`Retrying after ${retryAfter} seconds!`);
        return true;
      },
      onSecondaryRateLimit: (_retryAfter, options, octokit) => {
        octokit.log.warn(
          // Types are busted from octokit
          //@ts-expect-error
          `SecondaryRateLimit detected for request ${options.method} ${options.url}`
        );
        return true;
      },
    },
  });

  return octokit;
}

async function getDeployments({
  octokit,
  owner,
  repo,
  environment,
  ref,
  paginateOptions,
}: {
  octokit: ReturnType<typeof getOctokit>;
  owner: string;
  repo: string;
  environment: string;
  ref?: string;
  paginateOptions: {
    numOfPages: number | "all";
    startingPage?: number;
  };
}) {
  const listDeploymentsSharedArgs: Parameters<
    typeof octokit.repos.listDeployments
  >[0] = {
    owner,
    repo,
    environment,
    ref,
    per_page: 100,
    request: {
      retries: 20,
    },
  };

  if (paginateOptions.numOfPages === "all") {
    info(`Fetching all deployments`);
    const deployments = await octokit.paginate(octokit.repos.listDeployments, {
      ...listDeploymentsSharedArgs,
    });

    return deployments;
  } else {
    info(
      `Fetching ${
        paginateOptions.numOfPages * listDeploymentsSharedArgs.per_page!
      } deployments`
    );
    const deployments: Awaited<
      ReturnType<typeof octokit.repos.listDeployments>
    >["data"] = [];

    const offset = paginateOptions.startingPage || 0;
    for (let i = offset; i < paginateOptions.numOfPages + offset; i++) {
      const deploymentPage = await octokit.repos.listDeployments({
        ...listDeploymentsSharedArgs,
        page: i,
      });

      deployments.push(...deploymentPage.data);
    }

    return deployments;
  }
}
