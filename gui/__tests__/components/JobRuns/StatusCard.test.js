import React from 'react'
import StatusCard from 'components/JobRuns/StatusCard'
import mountWithTheme from 'test-helpers/mountWithTheme'

describe('components/JobRuns/StatusCard', () => {
  it('can display a title', () => {
    let withTitle = mountWithTheme(
      <StatusCard>pending_confirmations</StatusCard>
    )
    expect(withTitle.text()).toContain('Pending Confirmations')
  })
})
