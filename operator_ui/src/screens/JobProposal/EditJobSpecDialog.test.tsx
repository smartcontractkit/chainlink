import * as React from 'react'

import userEvent from '@testing-library/user-event'
import { render, screen, waitFor } from 'support/test-utils'

import { EditJobSpecDialog } from './EditJobSpecDialog'

const { findByText, getByRole, getByText, queryByRole } = screen

describe('EditJobSpecDialog', () => {
  it('renders the dialog', async () => {
    const handleSubmit = jest.fn()

    render(
      <EditJobSpecDialog
        open={true}
        onClose={() => null}
        initialValues={{ definition: 'name=test', id: '1' }}
        onSubmit={handleSubmit}
      />,
    )

    expect(
      queryByRole('heading', { name: /edit job spec/i }),
    ).toBeInTheDocument()
  })

  it('does not render the dialog when not open', async () => {
    const handleSubmit = jest.fn()

    render(
      <EditJobSpecDialog
        open={false}
        onClose={() => null}
        initialValues={{ definition: 'name=test', id: '1' }}
        onSubmit={handleSubmit}
      />,
    )

    expect(queryByRole('heading', { name: /edit job spec/i })).toBeNull()
  })

  it('validates the form', async () => {
    const handleSubmit = jest.fn()

    render(
      <EditJobSpecDialog
        open={true}
        onClose={() => null}
        initialValues={{ definition: '', id: '1' }}
        onSubmit={handleSubmit}
      />,
    )

    userEvent.click(getByRole('button', { name: /submit/i }))

    expect(await findByText('Required')).toBeInTheDocument()
  })

  it('submits the form', async () => {
    const handleSubmit = jest.fn()

    render(
      <EditJobSpecDialog
        open={true}
        onClose={() => null}
        initialValues={{ definition: 'test', id: '1' }}
        onSubmit={handleSubmit}
      />,
    )

    userEvent.click(getByText(/submit/i))

    await waitFor(() =>
      expect(handleSubmit).toHaveBeenCalledWith(
        { definition: 'test', id: '1' },
        expect.anything(),
      ),
    )
  })

  it('closes the form', async () => {
    const handleSubmit = jest.fn()
    const handleClose = jest.fn()

    render(
      <EditJobSpecDialog
        open={true}
        onClose={handleClose}
        initialValues={{ definition: 'test', id: '1' }}
        onSubmit={handleSubmit}
      />,
    )

    userEvent.click(getByText(/cancel/i))

    await waitFor(() => expect(handleClose).toHaveBeenCalled())
  })
})
