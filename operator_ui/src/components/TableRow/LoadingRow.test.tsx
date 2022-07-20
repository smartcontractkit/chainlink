import React from 'react'
import { render, screen } from '@testing-library/react'

import { LoadingRow } from './LoadingRow'

const { getAllByRole, queryAllByRole, queryByTestId } = screen

describe('LoadingRow', () => {
  function renderComponent(visible: boolean) {
    render(
      <table>
        <tbody>
          <LoadingRow visible={visible} />
        </tbody>
      </table>,
    )
  }

  it('renders an loading spinner', () => {
    renderComponent(true)

    expect(getAllByRole('row')).toHaveLength(1)
    expect(queryByTestId('loading-spinner')).toBeInTheDocument()
  })

  it('renders nothing', () => {
    renderComponent(false)

    expect(queryAllByRole('row')).toHaveLength(0)
    expect(queryByTestId('loading-spinner')).toBeNull()
  })
})
