import React from 'react'
import RegionalNav from 'components/JobRuns/RegionalNav'
import mountWithTheme from 'test-helpers/mountWithTheme'

const CREATED_AT = '2019-06-11T14:37:42.077995-07:00'

describe('components/JobRuns/RegionalNav', () => {
  it('displays an overview & json tab', () => {
    const component = mountWithTheme(<RegionalNav />)

    expect(component.text()).toContain('Overview')
    expect(component.text()).toContain('JSON')
    expect(component.text()).not.toContain('Error Log')
  })

  it('displays an error log tab when the status is "errored"', () => {
    const jobRun = { status: 'errored', createdAt: CREATED_AT }
    const component = mountWithTheme(<RegionalNav jobRun={jobRun} />)

    expect(component.text()).toContain('Error Log')
  })
})
