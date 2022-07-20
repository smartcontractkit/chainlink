import * as React from 'react'

import { render, screen } from 'support/test-utils'

import { buildEthTx } from 'support/factories/gql/fetchEthTransaction'
import { TransactionCard } from './TransactionCard'
import titleize from 'src/utils/titleize'

const { queryByText } = screen

describe('TransactionCard', () => {
  function renderComponent(tx: EthTransactionPayloadFields) {
    render(<TransactionCard tx={tx} />)
  }

  it('renders a node', () => {
    const tx = buildEthTx()

    renderComponent(tx)

    expect(queryByText(tx.chain.id)).toBeInTheDocument()
    expect(queryByText(tx.data)).toBeInTheDocument()
    expect(queryByText(tx.from)).toBeInTheDocument()
    expect(queryByText(tx.gasLimit)).toBeInTheDocument()
    expect(queryByText(tx.gasPrice)).toBeInTheDocument()
    expect(queryByText(tx.hash)).toBeInTheDocument()
    expect(queryByText(tx.hex)).toBeInTheDocument()
    expect(queryByText(tx.nonce as string)).toBeInTheDocument()
    expect(queryByText(tx.sentAt as string)).toBeInTheDocument()
    expect(queryByText(titleize(tx.state))).toBeInTheDocument()
    expect(queryByText(tx.to)).toBeInTheDocument()
    expect(queryByText(tx.value)).toBeInTheDocument()
  })
})
