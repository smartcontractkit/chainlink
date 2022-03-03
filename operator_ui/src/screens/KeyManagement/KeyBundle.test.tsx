import React from 'react'

import { render, screen } from 'support/test-utils'

import { KeyBundle } from './KeyBundle'

const { getByText } = screen

describe('pages/Keys/KeyBundle', () => {
  it('renders key bundle cell', async () => {
    const expectedPrimary = 'Primary information'
    const expectedSecondary = [
      'Secondary information 1',
      'Secondary information 2',
    ]
    render(
      <KeyBundle primary={expectedPrimary} secondary={expectedSecondary} />,
    )

    expect(getByText(expectedPrimary)).toBeInTheDocument()
    expect(getByText(expectedSecondary[0])).toBeInTheDocument()
    expect(getByText(expectedSecondary[1])).toBeInTheDocument()
  })
})
