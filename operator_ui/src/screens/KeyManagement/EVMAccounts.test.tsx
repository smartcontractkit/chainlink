import * as React from 'react'

import { GraphQLError } from 'graphql'
import { renderWithRouter, screen } from 'support/test-utils'
import { MockedProvider, MockedResponse } from '@apollo/client/testing'

import { EVMAccounts } from './EVMAccounts'
import { buildETHKeys } from 'support/factories/gql/fetchETHKeys'
import { ETH_KEYS_QUERY } from 'hooks/queries/useEVMAccountsQuery'
import { waitForLoading } from 'support/test-helpers/wait'
import { shortenHex } from 'src/utils/shortenHex'

const { findByText } = screen

function renderComponent(mocks: MockedResponse[]) {
  renderWithRouter(
    <>
      <MockedProvider mocks={mocks} addTypename={false}>
        <EVMAccounts />
      </MockedProvider>
    </>,
  )
}

function fetchEthKeysQuery(
  ethKeys: ReadonlyArray<EthKeysPayload_ResultsFields>,
) {
  return {
    request: {
      query: ETH_KEYS_QUERY,
    },
    result: {
      data: {
        ethKeys: {
          results: ethKeys,
        },
      },
    },
  }
}

describe('EVMAccount', () => {
  it('renders the page', async () => {
    const payload = buildETHKeys()
    const mocks: MockedResponse[] = [fetchEthKeysQuery(payload)]

    renderComponent(mocks)

    await waitForLoading()

    expect(
      await findByText(shortenHex(payload[0].address, { start: 6, end: 6 })),
    ).toBeInTheDocument()
  })

  it('renders GQL query errors', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: ETH_KEYS_QUERY,
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
