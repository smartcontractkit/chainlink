import * as React from 'react'

import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'support/test-utils'

import { buildEthTxs } from 'support/factories/gql/fetchEthTransactions'
import {
  TransactionsView,
  Props as TransactionsViewProps,
} from './TransactionsView'

const { getAllByRole, queryByText } = screen

function renderComponent(viewProps: TransactionsViewProps) {
  renderWithRouter(
    <>
      <Route exact path="/transactions">
        <TransactionsView {...viewProps} />
      </Route>
    </>,
    { initialEntries: ['/transactions'] },
  )
}

describe('TransactionsView', () => {
  it('renders the transactions table', () => {
    const txs = buildEthTxs()

    renderComponent({
      data: {
        ethTransactions: {
          results: txs,
          metadata: { total: txs.length },
        },
      },
      loading: false,
      page: 1,
      pageSize: 10,
    })

    expect(getAllByRole('row')).toHaveLength(3)

    expect(queryByText('Txn Hash')).toBeInTheDocument()
    expect(queryByText('Chain ID')).toBeInTheDocument()
    expect(queryByText('From')).toBeInTheDocument()
    expect(queryByText('To')).toBeInTheDocument()
    expect(queryByText('Nonce')).toBeInTheDocument()
    expect(queryByText('Block')).toBeInTheDocument()

    expect(queryByText('2-2 of 2'))
  })
})
