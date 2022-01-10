import React from 'react'

import { render, screen } from 'test-utils'

import { buildEthTx } from 'support/factories/gql/fetchEthTransaction'
import { TransactionView } from './TransactionView'

const { getByRole, getByText } = screen

describe('TransactionView', () => {
  function renderComponent(tx: EthTransactionPayloadFields) {
    render(<TransactionView tx={tx} />)
  }

  it('renders the view', async () => {
    const tx = buildEthTx()
    renderComponent(tx)

    expect(
      getByRole('heading', { name: 'Transaction Details' }),
    ).toBeInTheDocument()
    expect(getByText(tx.hash)).toBeInTheDocument()
  })
})
