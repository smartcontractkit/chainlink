import React from 'react'
import { partialAsFull } from '@chainlink/ts-helpers'
import { JobRun } from 'operator_ui'
import mountWithTheme from '../../../support/test-helpers/mountWithTheme'
import Activity from '../../../src/components/Dashboards/Activity'

const CREATED_AT = '2019-06-11T14:37:42.077995-07:00'

describe('components/Dashboards/Activity', () => {
  it('displays the given runs', () => {
    const runs = [partialAsFull<JobRun>({ id: 'runA', createdAt: CREATED_AT })]
    const component = mountWithTheme(<Activity runs={runs} pageSize={1} />)
    expect(component.text()).toContain('Run: runA')
  })

  it('can show a loading message', () => {
    const component = mountWithTheme(<Activity pageSize={1} />)
    expect(component.text()).toContain('Loading ...')
  })

  it('can show a no activity message', () => {
    const component = mountWithTheme(<Activity runs={[]} pageSize={1} />)
    expect(component.text()).toContain('No recent activity')
  })
})
