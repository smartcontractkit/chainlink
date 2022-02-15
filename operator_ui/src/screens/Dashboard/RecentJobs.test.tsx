import * as React from 'react'

import { GraphQLError } from 'graphql'
import { renderWithRouter, screen } from 'support/test-utils'
import { MockedProvider, MockedResponse } from '@apollo/client/testing'

import { RecentJobs, RECENT_JOBS_QUERY } from './RecentJobs'
import { buildRecentJobs } from 'support/factories/gql/fetchRecentJobs'
import { waitForLoading } from 'support/test-helpers/wait'

const { findByText } = screen

function renderComponent(mocks: MockedResponse[]) {
  renderWithRouter(
    <MockedProvider mocks={mocks} addTypename={false}>
      <RecentJobs />
    </MockedProvider>,
  )
}

function fetchRecentJobsQuery(
  jobs: ReadonlyArray<RecentJobsPayload_ResultsFields>,
) {
  return {
    request: {
      query: RECENT_JOBS_QUERY,
      variables: { offset: 0, limit: 5 },
    },
    result: {
      data: {
        jobs: {
          results: jobs,
        },
      },
    },
  }
}

describe('RecentJobs', () => {
  it('renders the recent jobs', async () => {
    const payload = buildRecentJobs()
    const mocks: MockedResponse[] = [fetchRecentJobsQuery(payload)]

    renderComponent(mocks)

    await waitForLoading()

    expect(await findByText(payload[0].name)).toBeInTheDocument()
    expect(await findByText(payload[1].name)).toBeInTheDocument()
  })

  it('renders GQL query errors', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: RECENT_JOBS_QUERY,
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
