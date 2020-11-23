import React from 'react'
import RegionalNav from './RegionalNav'
import mountWithTheme from 'test-helpers/mountWithTheme'

const CREATED_AT = '2019-06-11T14:37:42.077995-07:00'

describe('pages/Jobs/Runs/RegionalNav', () => {
  it('displays an overview & json tab by default', () => {
    const component = mountWithTheme(<RegionalNav jobRunId="1" jobSpecId="1" />)

    expect(component.text()).toContain('Overview')
    expect(component.text()).toContain('JSON')
    expect(component.text()).not.toContain('Error Log')
  })

  it('also displays an error log tab when the status is "errored"', () => {
    const jobRun = { status: 'errored', createdAt: CREATED_AT }
    const component = mountWithTheme(
      <RegionalNav jobRunId="1" jobSpecId="1" jobRun={jobRun} />,
    )

    expect(component.text()).toContain('Error Log')
  })
})
