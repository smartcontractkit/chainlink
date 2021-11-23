// buildJobProposal builds a job proposal for the FetchJobProposal query.
export function buildJobProposal(
  overrides?: Partial<JobProposalPayloadFields>,
): JobProposalPayloadFields {
  return {
    __typename: 'JobProposal',
    id: '1',
    spec: 'name="test"',
    status: 'PENDING',
    externalJobID: null,
    ...overrides,
  }
}
