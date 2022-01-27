import * as React from 'react'

import { GraphQLError } from 'graphql'
import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'support/test-utils'
import { MockedProvider, MockedResponse } from '@apollo/client/testing'

import { TransactionScreen, ETH_TRANSACTION_QUERY } from './TransactionScreen'
import { waitForLoading } from 'support/test-helpers/wait'
import { buildEthTx } from 'support/factories/gql/fetchEthTransaction'

const { findByTestId, findByText } = screen

function renderComponent(mocks: MockedResponse[]) {
  renderWithRouter(
    <>
      <Route exact path="/transactions/:id">
        <MockedProvider mocks={mocks} addTypename={false}>
          <TransactionScreen />
        </MockedProvider>
      </Route>
    </>,
    { initialEntries: ['/transactions/0x000'] },
  )
}

function fetchTransactionQuery(tx: EthTransactionPayloadFields) {
  return {
    request: {
      query: ETH_TRANSACTION_QUERY,
      variables: { hash: '0x000' },
    },
    result: {
      data: {
        ethTransaction: tx,
      },
    },
  }
}

describe('TransactionScreen', () => {
  it('renders the page', async () => {
    const payload = buildEthTx()
    const mocks: MockedResponse[] = [fetchTransactionQuery(payload)]

    renderComponent(mocks)

    await waitForLoading()

    expect(await findByText(payload.hash)).toBeInTheDocument()
  })

  it('renders the not found page', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: ETH_TRANSACTION_QUERY,
          variables: { hash: '0x000' },
        },
        result: {
          data: {
            ethTransaction: {
              __typename: 'NotFoundError',
              message: 'tx not found',
            },
          },
        },
      },
    ]

    renderComponent(mocks)

    await waitForLoading()

    expect(await findByTestId('not-found-page')).toBeInTheDocument()
  })

  it('renders GQL query errors', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: ETH_TRANSACTION_QUERY,
          variables: { hash: '0x000' },
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
