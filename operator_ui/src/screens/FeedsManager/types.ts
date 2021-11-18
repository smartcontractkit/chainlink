// FetchResults types the results of fetching all the feed managers
export type FetchResults =
  FetchFeedManagersWithProposals['feedsManagers']['results']
// FeedsManager types a feedsManager result
export type FeedsManager = FetchResults[number]
// JobProposals types the job proposals field of a feeds manager
export type JobProposals = FeedsManager['jobProposals']
// JobProposals types a job proposal
export type JobProposal = JobProposals[number]
