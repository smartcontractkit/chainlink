import * as React from 'react'

import { render, screen } from 'support/test-utils'

import { EVMAccountRow } from './EVMAccountRow'
import { buildETHKey } from 'support/factories/gql/fetchETHKeys'
import { shortenHex } from 'src/utils/shortenHex'

const { queryByText } = screen

describe('EVMAccountRow', () => {
  function renderComponent(ethKey: EthKeysPayload_ResultsFields) {
    render(
      <table>
        <tbody>
          <EVMAccountRow ethKey={ethKey} />
        </tbody>
      </table>,
    )
  }

  it('renders a row', () => {
    const ethKey = buildETHKey()

    renderComponent(ethKey)

    expect(
      queryByText(shortenHex(ethKey.address, { start: 6, end: 6 })),
    ).toBeInTheDocument()
    expect(queryByText(ethKey.chain.id)).toBeInTheDocument()
    expect(queryByText('Regular')).toBeInTheDocument()
    expect(queryByText('1.00000000')).toBeInTheDocument()
    expect(queryByText('0.100000000000000000')).toBeInTheDocument()
    expect(queryByText('1 minute ago')).toBeInTheDocument()
  })
})
