import React from 'react'
import { mountWithProviders } from 'test-helpers/mountWithTheme'
import { KeyBundle } from './KeyBundle'

describe('pages/Keys/KeyBundle', () => {
  it('renders key bundle cell', async () => {
    const expectedPrimary = 'Primary information'
    const expectedSecondary = [
      'Secondary information 1',
      'Secondary information 2',
    ]
    const wrapper = mountWithProviders(
      <KeyBundle primary={expectedPrimary} secondary={expectedSecondary} />,
    )

    expect(wrapper.text()).toContain(expectedPrimary)
    expect(wrapper.text()).toContain(expectedSecondary[0])
    expect(wrapper.text()).toContain(expectedSecondary[1])
  })
})
