import * as React from 'react'

import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'support/test-utils'
import userEvent from '@testing-library/user-event'

import { TransactionRow } from './TransactionRow'
import { buildEthTx } from 'support/factories/gql/fetchEthTransactions'
import { shortenHex } from 'src/utils/shortenHex'

const { findByText, getByRole, queryByText } = screen

function renderComponent(tx: EthTransactionsPayload_ResultsFields) {
  renderWithRouter(
    <>
      <Route exact path="/">
        <table>
          <tbody>
            <TransactionRow tx={tx} />
          </tbody>
        </table>
      </Route>
      <Route exact path="/transactions/:id">
        Link Success
      </Route>
    </>,
  )
}

describe('TransactionRow', () => {
  it('renders the row', () => {
    const tx = buildEthTx()

    renderComponent(tx)

    expect(
      queryByText(shortenHex(tx.hash, { start: 8, end: 8 })),
    ).toBeInTheDocument()
    expect(
      queryByText(shortenHex(tx.from, { start: 8, end: 8 })),
    ).toBeInTheDocument()
    expect(
      queryByText(shortenHex(tx.to, { start: 8, end: 8 })),
    ).toBeInTheDocument()
    expect(queryByText(tx.chain.id)).toBeInTheDocument()
    expect(queryByText(tx.nonce as string)).toBeInTheDocument()
    expect(queryByText(tx.sentAt as string)).toBeInTheDocument()
  })

  it('links to the job details', async () => {
    const tx = buildEthTx()

    renderComponent(tx)

    const link = getByRole('link', {
      name: shortenHex(tx.hash, { start: 8, end: 8 }),
    })
    expect(link).toHaveAttribute('href', `/transactions/${tx.hash}`)

    userEvent.click(link)

    expect(await findByText('Link Success')).toBeInTheDocument()
  })
})
