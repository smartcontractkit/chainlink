import * as React from 'react'

import { render, screen } from 'support/test-utils'
import userEvent from '@testing-library/user-event'

import { buildP2PKey } from 'support/factories/gql/fetchP2PKeys'
import { P2PKeyRow } from './P2PKeyRow'

const { getByRole, queryByText } = screen

describe('P2PKeyRow', () => {
  let handleDelete: jest.Mock

  beforeEach(() => {
    handleDelete = jest.fn()
  })

  function renderComponent(bundle: P2PKeysPayload_ResultsFields) {
    render(
      <table>
        <tbody>
          <P2PKeyRow p2pKey={bundle} onDelete={handleDelete} />
        </tbody>
      </table>,
    )
  }

  it('renders a row', () => {
    const p2pKey = buildP2PKey()

    renderComponent(p2pKey)

    expect(queryByText(`Peer ID: ${p2pKey.peerID}`)).toBeInTheDocument()
    expect(queryByText(`Public Key: ${p2pKey.publicKey}`)).toBeInTheDocument()
  })

  it('calls delete', () => {
    const p2pKey = buildP2PKey()

    renderComponent(p2pKey)

    userEvent.click(getByRole('button', { name: /delete/i }))

    expect(handleDelete).toHaveBeenCalled()
  })
})
