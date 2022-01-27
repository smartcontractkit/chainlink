import * as React from 'react'

import { renderWithRouter, screen } from 'support/test-utils'

import {
  AccountBalanceCard,
  Props as AccountBalanceProps,
} from './AccountBalanceCard'
import {
  buildETHKey,
  buildETHKeys,
} from 'support/factories/gql/fetchAccountBalances'
import { fromJuels } from 'src/utils/tokens/link'

const { getAllByText, queryByText, queryByRole } = screen

function renderComponent(cardProps: AccountBalanceProps) {
  renderWithRouter(<AccountBalanceCard {...cardProps} />)
}

describe('AccountBalanceCard', () => {
  it('renders the account balance', () => {
    const ethKey = buildETHKey()

    renderComponent({
      loading: false,
      data: {
        ethKeys: {
          results: [ethKey],
        },
      },
    })

    expect(queryByText(ethKey.address)).toBeInTheDocument()
    expect(
      queryByText(fromJuels(ethKey.linkBalance as string)),
    ).toBeInTheDocument()
    expect(queryByText(ethKey.ethBalance as string)).toBeInTheDocument()

    // Does not appear if there is only one account
    expect(queryByRole('link', { name: /view more accounts/i })).toBeNull()
  })

  it('renders the empty balances for an account', () => {
    const ethKey = buildETHKey({
      linkBalance: undefined,
      ethBalance: undefined,
    })

    renderComponent({
      loading: false,
      data: {
        ethKeys: {
          results: [ethKey],
        },
      },
    })

    expect(getAllByText('--')).toHaveLength(2)
  })

  it('shows the view more button', () => {
    renderComponent({
      loading: false,
      data: {
        ethKeys: {
          results: buildETHKeys(),
        },
      },
    })

    expect(
      queryByRole('link', { name: /view more accounts/i }),
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

    expect(queryByText('No account available')).toBeInTheDocument()
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
