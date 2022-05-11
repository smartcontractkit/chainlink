import React from 'react'

import { MockedProvider, MockedResponse } from '@apollo/client/testing'
import { GraphQLError } from 'graphql'
import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'test-utils'
import userEvent from '@testing-library/user-event'

import { ChainsScreen } from './ChainsScreen'
import { buildChain, buildChains } from 'support/factories/gql/fetchChains'
import { CHAINS_QUERY } from 'hooks/queries/useChainsQuery'
import { waitForLoading } from 'support/test-helpers/wait'

const { findAllByRole, findByText, getByRole, queryByText } = screen

function renderComponent(mocks: MockedResponse[]) {
  renderWithRouter(
    <>
      <Route exact path="/chains">
        <MockedProvider mocks={mocks} addTypename={false}>
          <ChainsScreen />
        </MockedProvider>
      </Route>
    </>,
    { initialEntries: ['/chains?per=2'] },
  )
}

describe('ChainsScreen', () => {
  it('renders the list of chains', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: CHAINS_QUERY,
          variables: { offset: 0, limit: 2 },
        },
        result: {
          data: {
            chains: {
              results: buildChains(),
              metadata: {
                total: 2,
              },
            },
          },
        },
      },
    ]

    renderComponent(mocks)

    expect(await findAllByRole('row')).toHaveLength(3)
  })

  it('can page through the list of chains', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: CHAINS_QUERY,
          variables: { offset: 0, limit: 2 },
        },
        result: {
          data: {
            chains: {
              results: buildChains(),
              metadata: {
                total: 3,
              },
            },
          },
        },
      },
      {
        request: {
          query: CHAINS_QUERY,
          variables: { offset: 2, limit: 2 },
        },
        result: {
          data: {
            chains: {
              results: [
                buildChain({
                  id: '4',
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
          query: CHAINS_QUERY,
          variables: { offset: 0, limit: 2 },
        },
        result: {
          data: {
            chains: {
              results: buildChains(),
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

    expect(queryByText('5')).toBeInTheDocument()
    expect(queryByText('42')).toBeInTheDocument()
    expect(queryByText('1-2 of 3')).toBeInTheDocument()
    expect(getByRole('button', { name: /prev-page/i })).toBeDisabled()

    // Page 2
    userEvent.click(getByRole('button', { name: /next-page/i }))

    // screen.logTestingPlaygroundURL()

    await waitForLoading()

    expect(queryByText('4')).toBeInTheDocument()
    expect(queryByText('3-3 of 3')).toBeInTheDocument()
    expect(getByRole('button', { name: /next-page/i })).toBeDisabled()

    // Page 1
    userEvent.click(getByRole('button', { name: /prev-page/i }))

    await waitForLoading()

    expect(queryByText('5')).toBeInTheDocument()
    expect(queryByText('42')).toBeInTheDocument()
    expect(queryByText('1-2 of 3')).toBeInTheDocument()
  })

  it('renders GQL errors', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: CHAINS_QUERY,
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
