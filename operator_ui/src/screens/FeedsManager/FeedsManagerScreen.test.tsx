import * as React from 'react'

import { GraphQLError } from 'graphql'
import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'support/test-utils'
import { MockedProvider, MockedResponse } from '@apollo/client/testing'

import { FeedsManagerScreen } from './FeedsManagerScreen'
import { buildFeedsManagerResultFields } from 'support/factories/gql/fetchFeedsManagersWithProposals'
import { FEEDS_MANAGERS_WITH_PROPOSALS_QUERY } from 'src/hooks/queries/useFeedsManagersWithProposalsQuery'

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
              results: [buildFeedsManagerResultFields()],
            },
          },
        },
      },
    ]

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
