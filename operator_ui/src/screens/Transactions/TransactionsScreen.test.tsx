import React from 'react'

import { MockedProvider, MockedResponse } from '@apollo/client/testing'
import { GraphQLError } from 'graphql'
import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'test-utils'
import userEvent from '@testing-library/user-event'

import {
  buildEthTx,
  buildEthTxs,
} from 'support/factories/gql/fetchEthTransactions'
import {
  TransactionsScreen,
  ETH_TRANSACTIONS_QUERY,
} from './TransactionsScreen'
import { waitForLoading } from 'support/test-helpers/wait'

const { findAllByRole, findByText, getByRole, queryByText } = screen

function renderComponent(mocks: MockedResponse[]) {
  renderWithRouter(
    <>
      <Route exact path="/transactions">
        <MockedProvider mocks={mocks} addTypename={false}>
          <TransactionsScreen />
        </MockedProvider>
      </Route>
    </>,
    { initialEntries: ['/transactions?per=2'] },
  )
}

describe('TransactionsScreen', () => {
  it('renders the list of transactions', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: ETH_TRANSACTIONS_QUERY,
          variables: { offset: 0, limit: 2 },
        },
        result: {
          data: {
            ethTransactions: {
              results: buildEthTxs(),
              metadata: {
                total: 2,
              },
            },
          },
        },
      },
    ]

    renderComponent(mocks)

    await waitForLoading()

    const rows = await findAllByRole('row')

    expect(rows).toHaveLength(3)
  })

  it('can page through the list of transactions', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: ETH_TRANSACTIONS_QUERY,
          variables: { offset: 0, limit: 2 },
        },
        result: {
          data: {
            ethTransactions: {
              results: buildEthTxs(),
              metadata: {
                total: 3,
              },
            },
          },
        },
      },
      {
        request: {
          query: ETH_TRANSACTIONS_QUERY,
          variables: { offset: 2, limit: 2 },
        },
        result: {
          data: {
            ethTransactions: {
              results: [
                buildEthTx({
                  hash: '0x3333333333333',
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
          query: ETH_TRANSACTIONS_QUERY,
          variables: { offset: 0, limit: 2 },
        },
        result: {
          data: {
            ethTransactions: {
              results: buildEthTxs(),
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

    expect(queryByText('0x111111...11111111')).toBeInTheDocument()
    expect(queryByText('0x222222...22222222')).toBeInTheDocument()
    expect(queryByText('1-2 of 3')).toBeInTheDocument()
    expect(getByRole('button', { name: /prev-page/i })).toBeDisabled()

    // Page 2
    userEvent.click(getByRole('button', { name: /next-page/i }))

    await waitForLoading()

    expect(queryByText('0x333333...33333333')).toBeInTheDocument()
    expect(queryByText('3-3 of 3')).toBeInTheDocument()
    expect(getByRole('button', { name: /next-page/i })).toBeDisabled()

    // Page 1
    userEvent.click(getByRole('button', { name: /prev-page/i }))

    await waitForLoading()

    expect(queryByText('0x111111...11111111')).toBeInTheDocument()
    expect(queryByText('0x222222...22222222')).toBeInTheDocument()
    expect(queryByText('1-2 of 3')).toBeInTheDocument()
  })

  it('renders GQL errors', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: ETH_TRANSACTIONS_QUERY,
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
