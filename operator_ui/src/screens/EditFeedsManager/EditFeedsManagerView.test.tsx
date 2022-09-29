import * as React from 'react'
import { render, screen } from 'support/test-utils'

import { EditFeedsManagerView } from './EditFeedsManagerView'
import { buildFeedsManager } from 'support/factories/gql/fetchFeedsManagers'

const { getByTestId, getByText } = screen

describe('EditFeedsManagerView', () => {
  it('renders with initial values', () => {
    const handleSubmit = jest.fn()
    const manager = buildFeedsManager()

    render(<EditFeedsManagerView data={manager} onSubmit={handleSubmit} />)

    expect(getByText('Edit Feeds Manager')).toBeInTheDocument()
    expect(getByTestId('feeds-manager-form')).toHaveFormValues({
      name: 'Chainlink Feeds Manager',
      uri: manager.uri,
      publicKey: manager.publicKey,
    })
  })
})
