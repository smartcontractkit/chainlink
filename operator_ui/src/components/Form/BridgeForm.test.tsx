import * as React from 'react'

import { render, screen, waitFor } from 'support/test-utils'
import userEvent from '@testing-library/user-event'

import { BridgeForm, FormValues } from './BridgeForm'

const { getByRole, findByTestId } = screen

describe('BridgeForm', () => {
  it('validates input', async () => {
    const handleSubmit = jest.fn()
    const initialValues: FormValues = {
      name: '',
      url: '',
      minimumContractPayment: '0',
      confirmations: 0,
    }

    render(
      <BridgeForm
        initialValues={initialValues}
        onSubmit={handleSubmit}
        submitButtonText="Submit"
      />,
    )

    userEvent.click(getByRole('button', { name: /submit/i }))

    expect(await findByTestId('name-helper-text')).toHaveTextContent('Required')
    expect(await findByTestId('url-helper-text')).toHaveTextContent('Required')
  })

  it('disables the name field', async () => {
    const handleSubmit = jest.fn()
    const initialValues: FormValues = {
      name: '',
      url: '',
      minimumContractPayment: '0',
      confirmations: 0,
    }

    render(
      <BridgeForm
        initialValues={initialValues}
        onSubmit={handleSubmit}
        submitButtonText="Submit"
        nameDisabled={true}
      />,
    )

    expect(getByRole('textbox', { name: /name */i })).toBeDisabled()
  })

  it('submits the form', async () => {
    const handleSubmit = jest.fn()
    const initialValues: FormValues = {
      name: '',
      url: '',
      minimumContractPayment: '0',
      confirmations: 0,
    }

    render(
      <BridgeForm
        initialValues={initialValues}
        onSubmit={handleSubmit}
        submitButtonText="Submit"
      />,
    )

    userEvent.type(getByRole('textbox', { name: /name */i }), 'bridge1')
    userEvent.type(
      getByRole('textbox', { name: /bridge url */i }),
      'https://www.test.com',
    )

    const minConfInput = getByRole('textbox', {
      name: /minimum contract payment/i,
    })
    userEvent.clear(minConfInput)
    userEvent.type(minConfInput, '1')

    userEvent.type(getByRole('spinbutton', { name: /confirmations/i }), '2')

    userEvent.click(getByRole('button', { name: /submit/i }))

    await waitFor(() =>
      expect(handleSubmit).toHaveBeenCalledWith(
        {
          name: 'bridge1',
          url: 'https://www.test.com',
          minimumContractPayment: '1',
          confirmations: 2,
        },
        expect.anything(),
      ),
    )
  })
})
