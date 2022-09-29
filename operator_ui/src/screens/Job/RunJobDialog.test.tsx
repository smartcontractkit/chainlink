import * as React from 'react'

import { render, screen } from 'support/test-utils'
import userEvent from '@testing-library/user-event'

import { RunJobDialog } from './RunJobDialog'

const { getByRole, queryByText } = screen

describe('RunJobDialog', () => {
  let handleOnClose: jest.Mock
  let handleOnRun: jest.Mock

  function renderComponent() {
    render(
      <RunJobDialog open={true} onClose={handleOnClose} onRun={handleOnRun} />,
    )
  }

  beforeEach(() => {
    handleOnClose = jest.fn()
    handleOnRun = jest.fn()
  })

  it('renders the dialog', () => {
    renderComponent()

    expect(queryByText('Pipeline Input'))
  })

  it('can close the dialog', () => {
    renderComponent()

    userEvent.click(getByRole('button', { name: /close/i }))

    expect(handleOnClose).toHaveBeenCalled()
  })

  it('submits the form', () => {
    renderComponent()

    userEvent.type(getByRole('textbox'), '{someinput}')
    userEvent.click(getByRole('button', { name: /run job/i }))

    expect(handleOnRun).toHaveBeenCalled()
    expect(handleOnClose).toHaveBeenCalled()
  })
})
