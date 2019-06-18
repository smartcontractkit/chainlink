import React from 'react'
import StatusCard from 'components/JobRuns/StatusCard'
import mountWithTheme from 'test-helpers/mountWithTheme'

describe('components/JobRuns/StatusCard', () => {
  it('converts the given title to title case', () => {
    let component = mountWithTheme(
      <StatusCard title={'pending_confirmations'} />
    )
    expect(component.text()).toContain('Pending Confirmations')
  })

  it('can display children', () => {
    let withChildren = mountWithTheme(
      <StatusCard title={'pending_confirmations'}>I am a child</StatusCard>
    )
    expect(withChildren.text()).toContain('I am a child')
  })
})
