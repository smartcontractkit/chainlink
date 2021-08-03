import React from 'react'
import RegionalNav from './RegionalNav'
import mountWithTheme from 'test-helpers/mountWithTheme'

describe('pages/Jobs/Runs/RegionalNav', () => {
  it('displays an overview & json tab by default', () => {
    const component = mountWithTheme(<RegionalNav jobRunId="1" jobId="1" />)

    expect(component.text()).toContain('Overview')
    expect(component.text()).toContain('JSON')
    expect(component.text()).not.toContain('Error Log')
  })
})
