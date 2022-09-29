import * as React from 'react'

import { render, screen } from 'support/test-utils'

import { CSAKeyRow } from './CSAKeyRow'
import { buildCSAKey } from 'support/factories/gql/fetchCSAKeys'

const { queryByText } = screen

describe('CSAKeyRow', () => {
  function renderComponent(csaKey: CsaKeysPayload_ResultsFields) {
    render(
      <table>
        <tbody>
          <CSAKeyRow csaKey={csaKey} />
        </tbody>
      </table>,
    )
  }

  it('renders a row', () => {
    const csaKey = buildCSAKey()

    renderComponent(csaKey)

    expect(queryByText(csaKey.publicKey)).toBeInTheDocument()
  })
})
