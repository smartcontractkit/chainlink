import '@testing-library/jest-dom'

import * as React from 'react'
import { render, screen } from '@testing-library/react'

import { EditFeedsManagerView } from './EditFeedsManagerView'
import { buildFeedsManager } from 'support/test-helpers/factories/feedsManager'

test('renders with initial values', () => {
  const handleSubmit = jest.fn()
  const manager = buildFeedsManager()

  render(<EditFeedsManagerView data={manager} onSubmit={handleSubmit} />)

  expect(screen.queryByTestId('feeds-manager-form')).toHaveFormValues({
    name: 'Chainlink Feeds Manager',
    uri: manager.uri,
    publicKey: manager.publicKey,
    jobTypes: manager.jobTypes,
    isBootstrapPeer: manager.isBootstrapPeer,
  })
})
