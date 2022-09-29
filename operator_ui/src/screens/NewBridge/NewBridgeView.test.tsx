import * as React from 'react'

import { render, screen } from 'support/test-utils'

import { NewBridgeView } from './NewBridgeView'

const { getByTestId, getByText } = screen

describe('NewBridgeView', () => {
  it('renders with initial values', () => {
    const handleSubmit = jest.fn()

    render(<NewBridgeView onSubmit={handleSubmit} />)

    expect(getByText('New Bridge')).toBeInTheDocument()
    expect(getByTestId('bridge-form')).toHaveFormValues({
      name: '',
      url: '',
      minimumContractPayment: '0',
      confirmations: 0,
    })
  })
})
