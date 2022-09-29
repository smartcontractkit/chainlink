import isoDate, { MINUTE_MS } from 'test-helpers/isoDate'

// buildJobProposal builds a job proposal for the FetchJobProposal query.
export function buildJobProposal(
  overrides?: Partial<JobProposalPayloadFields>,
): JobProposalPayloadFields {
  return {
    __typename: 'JobProposal',
    id: '1',
    remoteUUID: '00000000-0000-0000-0000-000000000001',
    status: 'PENDING',
    externalJobID: null,
    specs: [buildJobProposalSpec()],
    ...overrides,
  }
}

// buildJobProposalSpec builds a job proposal spec for the FetchJobProposal query.
export function buildJobProposalSpec(
  overrides?: Partial<JobProposal_SpecsFields>,
): JobProposal_SpecsFields {
  const minuteAgo = isoDate(Date.now() - MINUTE_MS)

  return {
    __typename: 'JobProposalSpec',
    id: '1',
    definition: "name='spec'",
    status: 'PENDING',
    version: 1,
    createdAt: minuteAgo,
    ...overrides,
  }
}
