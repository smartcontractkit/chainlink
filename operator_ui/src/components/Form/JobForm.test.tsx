import * as React from 'react'

import { render, screen, waitFor } from 'support/test-utils'
import userEvent from '@testing-library/user-event'

import { JobForm, FormValues } from './JobForm'

const { getByRole, findByTestId } = screen

describe('BridgeForm', () => {
  let handleSubmit: jest.Mock
  let handleOnTOMLChange: jest.Mock

  beforeEach(() => {
    handleSubmit = jest.fn()
    handleOnTOMLChange = jest.fn()
  })

  function renderComponent(initialValues: FormValues) {
    render(
      <JobForm
        initialValues={initialValues}
        onSubmit={handleSubmit}
        onTOMLChange={handleOnTOMLChange}
      />,
    )
  }

  it('validates the TOML as required', async () => {
    renderComponent({ toml: '' })

    userEvent.click(getByRole('button', { name: /create job/i }))

    expect(await findByTestId('toml-helper-text')).toHaveTextContent('Required')
  })

  it('validates the TOML as invalid', async () => {
    renderComponent({ toml: 'invalidtoml' })

    userEvent.click(getByRole('button', { name: /create job/i }))

    expect(await findByTestId('toml-helper-text')).toHaveTextContent(
      'Invalid TOML',
    )
  })

  it('submits the form', async () => {
    renderComponent({ toml: '' })

    userEvent.paste(
      getByRole('textbox', { name: /job spec \(toml\) \*/i }),
      'type = "webhook"',
    )

    userEvent.click(getByRole('button', { name: /create job/i }))

    await waitFor(() =>
      expect(handleSubmit).toHaveBeenCalledWith(
        {
          toml: 'type = "webhook"',
        },
        expect.anything(),
      ),
    )
  })

  it('updates onTOML change', async () => {
    renderComponent({ toml: '' })

    userEvent.paste(
      getByRole('textbox', { name: /job spec \(toml\) \*/i }),
      'type = "webhook"',
    )

    await waitFor(() =>
      expect(handleOnTOMLChange).toHaveBeenLastCalledWith('type = "webhook"'),
    )
  })
})
