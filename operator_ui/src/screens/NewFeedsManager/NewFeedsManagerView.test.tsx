import '@testing-library/jest-dom'

import * as React from 'react'
import { render, screen } from '@testing-library/react'

import { NewFeedsManagerView } from './NewFeedsManagerView'

test('renders with initial values', () => {
  const handleSubmit = jest.fn()

  render(<NewFeedsManagerView onSubmit={handleSubmit} />)

  expect(screen.queryByTestId('feeds-manager-form')).toHaveFormValues({
    name: 'Chainlink Feeds Manager',
    uri: '',
    publicKey: '',
    jobTypes: [],
    isBootstrapPeer: false,
  })
})
