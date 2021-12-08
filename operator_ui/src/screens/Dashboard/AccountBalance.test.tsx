import * as React from 'react'

import { GraphQLError } from 'graphql'
import { renderWithRouter, screen } from 'support/test-utils'
import { MockedProvider, MockedResponse } from '@apollo/client/testing'

import { AccountBalance, ACCOUNT_BALANCES_QUERY } from './AccountBalance'
import { buildETHKeys } from 'support/factories/gql/fetchAccountBalances'
import { waitForLoading } from 'support/test-helpers/wait'

const { findByText, queryByText } = screen

function renderComponent(mocks: MockedResponse[]) {
  renderWithRouter(
    <MockedProvider mocks={mocks} addTypename={false}>
      <AccountBalance />
    </MockedProvider>,
  )
}

function fetchAccountBalancesQuery(
  accounts: ReadonlyArray<AccountBalancesPayload_ResultsFields>,
) {
  return {
    request: {
      query: ACCOUNT_BALANCES_QUERY,
    },
    result: {
      data: {
        ethKeys: {
          results: accounts,
        },
      },
    },
  }
}

describe('Activity', () => {
  it('renders the first account balance', async () => {
    const payload = buildETHKeys()
    const mocks: MockedResponse[] = [fetchAccountBalancesQuery(payload)]

    renderComponent(mocks)

    await waitForLoading()

    expect(await findByText(payload[0].address)).toBeInTheDocument()
    expect(queryByText(payload[1].address)).toBeNull()
  })

  it('renders GQL query errors', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: ACCOUNT_BALANCES_QUERY,
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
