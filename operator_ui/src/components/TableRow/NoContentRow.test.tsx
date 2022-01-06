import React from 'react'
import { render, screen } from '@testing-library/react'

import { NoContentRow } from './NoContentRow'

const { getAllByRole, queryAllByRole, queryByText } = screen

describe('ErrorRow', () => {
  function renderComponent(visible: boolean, children?: React.ReactNode) {
    render(
      <table>
        <tbody>
          <NoContentRow visible={visible}>{children}</NoContentRow>
        </tbody>
      </table>,
    )
  }

  it('renders default text', () => {
    renderComponent(true)

    expect(getAllByRole('row')).toHaveLength(1)
    expect(queryByText('No entries to show')).toBeInTheDocument()
  })

  it('renders custom child', () => {
    renderComponent(true, 'custom message')

    expect(getAllByRole('row')).toHaveLength(1)
    expect(queryByText('custom message')).toBeInTheDocument()
  })

  it('renders nothing', () => {
    renderComponent(false)

    expect(queryAllByRole('row')).toHaveLength(0)
    expect(queryByText('error message')).toBeNull()
  })
})
