import React from 'react'

import { MockedProvider, MockedResponse } from '@apollo/client/testing'
import { GraphQLError } from 'graphql'
import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'test-utils'
import userEvent from '@testing-library/user-event'

import { BridgesScreen, BRIDGES_QUERY } from './BridgesScreen'
import { buildBridge, buildBridges } from 'support/factories/gql/fetchBridges'
import { waitForLoading } from 'support/test-helpers/wait'

const { findAllByRole, findByText, getByRole, queryByText } = screen

function renderComponent(mocks: MockedResponse[]) {
  renderWithRouter(
    <>
      <Route exact path="/bridges">
        <MockedProvider mocks={mocks} addTypename={false}>
          <BridgesScreen />
        </MockedProvider>
      </Route>

      <Route path="/feeds_manager/new">Redirect Success</Route>
    </>,
    { initialEntries: ['/bridges?per=2'] },
  )
}

describe('BridgesScreen', () => {
  it('renders the list of bridges', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: BRIDGES_QUERY,
          variables: { offset: 0, limit: 2 },
        },
        result: {
          data: {
            bridges: {
              results: buildBridges(),
              metadata: {
                total: 2,
              },
            },
          },
        },
      },
    ]

    renderComponent(mocks)

    const rows = await findAllByRole('row')

    expect(rows).toHaveLength(3)
  })

  it('can page through the list of bridges', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: BRIDGES_QUERY,
          variables: { offset: 0, limit: 2 },
        },
        result: {
          data: {
            bridges: {
              results: buildBridges(),
              metadata: {
                total: 3,
              },
            },
          },
        },
      },
      {
        request: {
          query: BRIDGES_QUERY,
          variables: { offset: 2, limit: 2 },
        },
        result: {
          data: {
            bridges: {
              results: [
                buildBridge({
                  name: 'bridge-api3',
                }),
              ],
              metadata: {
                total: 3,
              },
            },
          },
        },
      },
      {
        request: {
          query: BRIDGES_QUERY,
          variables: { offset: 0, limit: 2 },
        },
        result: {
          data: {
            bridges: {
              results: buildBridges(),
              metadata: {
                total: 3,
              },
            },
          },
        },
      },
    ]

    renderComponent(mocks)

    // Page 1
    await waitForLoading()

    expect(queryByText('bridge-api1')).toBeInTheDocument()
    expect(queryByText('bridge-api2')).toBeInTheDocument()
    expect(queryByText('1-2 of 3')).toBeInTheDocument()
    expect(getByRole('button', { name: /prev-page/i })).toBeDisabled()

    // Page 2
    userEvent.click(getByRole('button', { name: /next-page/i }))

    await waitForLoading()

    expect(queryByText('bridge-api3')).toBeInTheDocument()
    expect(queryByText('3-3 of 3')).toBeInTheDocument()
    expect(getByRole('button', { name: /next-page/i })).toBeDisabled()

    // Page 1
    userEvent.click(getByRole('button', { name: /prev-page/i }))

    await waitForLoading()

    expect(queryByText('bridge-api1')).toBeInTheDocument()
    expect(queryByText('bridge-api2')).toBeInTheDocument()
    expect(queryByText('1-2 of 3')).toBeInTheDocument()
  })

  it('renders GQL errors', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: BRIDGES_QUERY,
          variables: { offset: 0, limit: 2 },
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
