import * as React from 'react'

import { GraphQLError } from 'graphql'
import { renderWithRouter, screen } from 'support/test-utils'
import { MockedProvider, MockedResponse } from '@apollo/client/testing'

import { Activity, RECENT_JOB_RUNS_QUERY } from './Activity'
import { buildRuns } from 'support/factories/gql/fetchRecentJobRuns'
import { waitForLoading } from 'support/test-helpers/wait'

const { findByText } = screen

function renderComponent(mocks: MockedResponse[]) {
  renderWithRouter(
    <MockedProvider mocks={mocks} addTypename={false}>
      <Activity />
    </MockedProvider>,
  )
}

function fetchRecentJobRunsQuery(
  runs: ReadonlyArray<RecentJobRunsPayload_ResultsFields>,
) {
  return {
    request: {
      query: RECENT_JOB_RUNS_QUERY,
      variables: { offset: 0, limit: 5 },
    },
    result: {
      data: {
        jobRuns: {
          results: runs,
          metadata: {
            total: runs.length,
          },
        },
      },
    },
  }
}

describe('Activity', () => {
  it('renders the recent jobs', async () => {
    const payload = buildRuns()
    const mocks: MockedResponse[] = [fetchRecentJobRunsQuery(payload)]

    renderComponent(mocks)

    await waitForLoading()

    expect(await findByText(`Job: ${payload[0].job.id}`)).toBeInTheDocument()
    expect(await findByText(`Job: ${payload[1].job.id}`)).toBeInTheDocument()
  })

  it('renders GQL query errors', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: RECENT_JOB_RUNS_QUERY,
          variables: { offset: 0, limit: 5 },
        },
        result: {
          errors: [new GraphQLError('Error!')],
        },
      },
    ]

    renderComponent(mocks)

    expect(await findByText('Error!')).toBeInTheDocument()
  })
})
