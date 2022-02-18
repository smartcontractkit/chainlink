import React from 'react'
import { render, screen } from '@testing-library/react'

import { ErrorRow } from './ErrorRow'

const { getAllByRole, queryAllByRole, queryByText } = screen

describe('ErrorRow', () => {
  function renderComponent(msg?: string) {
    render(
      <table>
        <tbody>
          <ErrorRow msg={msg} />
        </tbody>
      </table>,
    )
  }

  it('renders an error', () => {
    renderComponent('error message')

    expect(getAllByRole('row')).toHaveLength(1)
    expect(queryByText('error message')).toBeInTheDocument()
  })

  it('renders nothing', () => {
    renderComponent()

    expect(queryAllByRole('row')).toHaveLength(0)
    expect(queryByText('error message')).toBeNull()
  })
})
