import * as React from 'react'

import { render, screen } from 'support/test-utils'

import {
  EVMAccountsCard,
  Props as EVMAccountsCardProps,
} from './EVMAccountsCard'
import { buildETHKeys } from 'support/factories/gql/fetchETHKeys'
import { shortenHex } from 'src/utils/shortenHex'

const { queryByRole, queryByText } = screen

function renderComponent(cardProps: EVMAccountsCardProps) {
  render(<EVMAccountsCard {...cardProps} />)
}

describe('ETHKeysCard', () => {
  it('renders the keys', () => {
    const ethKeys = buildETHKeys()

    renderComponent({
      loading: false,
      data: {
        ethKeys: {
          results: ethKeys,
        },
      },
    })

    expect(queryByText('Address')).toBeInTheDocument()
    expect(queryByText('Chain ID')).toBeInTheDocument()
    expect(queryByText('Type')).toBeInTheDocument()
    expect(queryByText('LINK Balance')).toBeInTheDocument()
    expect(queryByText('ETH Balance')).toBeInTheDocument()
    expect(queryByText('Created')).toBeInTheDocument()

    expect(
      queryByText(shortenHex(ethKeys[0].address, { start: 6, end: 6 })),
    ).toBeInTheDocument()
    expect(
      queryByText(shortenHex(ethKeys[1].address, { start: 6, end: 6 })),
    ).toBeInTheDocument()
  })

  it('renders no content', () => {
    renderComponent({
      loading: false,
      data: {
        ethKeys: {
          results: [],
        },
      },
    })

    expect(queryByText('No entries to show')).toBeInTheDocument()
  })

  it('renders a loading spinner', () => {
    renderComponent({
      loading: true,
    })

    expect(queryByRole('progressbar')).toBeInTheDocument()
  })

  it('renders an error message', () => {
    renderComponent({
      loading: false,
      errorMsg: 'error message',
    })

    expect(queryByText('error message')).toBeInTheDocument()
  })
})
