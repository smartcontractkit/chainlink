import * as React from 'react'

import globPath from 'test-helpers/globPath'
import { GraphQLError } from 'graphql'
import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'support/test-utils'
import { MockedProvider, MockedResponse } from '@apollo/client/testing'

import {
  FeedsManagerScreen,
  FEEDS_MANAGERS_WITH_PROPOSALS_QUERY,
} from './FeedsManagerScreen'
import { buildFeedsManagerFields } from 'support/factories/gql/fetchFeedsManagersWithProposals'

const { findByText } = screen

function renderComponent(mocks: MockedResponse[]) {
  renderWithRouter(
    <>
      <Route exact path="/">
        <MockedProvider mocks={mocks} addTypename={false}>
          <FeedsManagerScreen />
        </MockedProvider>
      </Route>

      <Route path="/feeds_manager/new">Redirect Success</Route>
    </>,
  )
}

describe('FeedsManagerScreen', () => {
  it('renders the page', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: FEEDS_MANAGERS_WITH_PROPOSALS_QUERY,
        },
        result: {
          data: {
            feedsManagers: {
              results: [buildFeedsManagerFields()],
            },
          },
        },
      },
    ]

    // Temporary until we switch it out for GQL
    global.fetch.getOnce(globPath('/v2/job_proposals'), { data: [] })

    renderComponent(mocks)

    expect(await findByText('Feeds Manager')).toBeInTheDocument()
    expect(await findByText('Job Proposals')).toBeInTheDocument()
  })

  it('redirects when a manager does not exists', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: FEEDS_MANAGERS_WITH_PROPOSALS_QUERY,
        },
        result: {
          data: {
            feedsManagers: {
              results: [],
            },
          },
        },
      },
    ]

    renderComponent(mocks)

    expect(await findByText('Redirect Success')).toBeInTheDocument()
  })

  it('renders GQL errors', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: FEEDS_MANAGERS_WITH_PROPOSALS_QUERY,
        },
        result: {
          errors: [new GraphQLError('Error!')],
        },
      },
    ]

    renderComponent(mocks)

    expect(await findByText('Error: Error!')).toBeInTheDocument()
  })
})
