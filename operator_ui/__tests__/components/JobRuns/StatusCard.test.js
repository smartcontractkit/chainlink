import React from 'react'
import StatusCard from 'components/JobRuns/StatusCard'
import mountWithTheme from 'test-helpers/mountWithTheme'
import { MINUTE_MS, TWO_MINUTES_MS } from 'test-helpers/isoDate'

describe('components/JobRuns/StatusCard', () => {
  const pendingRun = {
    id: 'runA',
    status: 'pending',
    result: {},
    createdAt: MINUTE_MS,
    finishedAt: null,
  }
  const completedRun = {
    id: 'runA',
    status: 'completed',
    createdAt: TWO_MINUTES_MS,
    finishedAt: MINUTE_MS,
    payment: 2000000000000000000,
  }
  const erroredRun = {
    id: 'runA',
    status: 'errored',
    result: {},
    createdAt: TWO_MINUTES_MS,
    finishedAt: MINUTE_MS,
  }
  it('converts the given title to title case', () => {
    let component = mountWithTheme(
      <StatusCard title={'pending_confirmations'} />,
    )
    expect(component.text()).toContain('Pending Confirmations')
  })

  it('can display children', () => {
    let withChildren = mountWithTheme(
      <StatusCard title={'pending_confirmations'}>I am a child</StatusCard>,
    )
    expect(withChildren.text()).toContain('I am a child')
  })

  it('can display the elapsed time for jobruns', () => {
    let erroredStatus = mountWithTheme(
      <StatusCard title="errored" jobRun={erroredRun} />,
    )
    let completedStatus = mountWithTheme(
      <StatusCard title="completed" jobRun={completedRun} />,
    )
    let pendingStatus = mountWithTheme(
      <StatusCard title="pending" jobRun={pendingRun} />,
    )

    expect(erroredStatus.text()).toContain('1m')
    expect(completedStatus.text()).toContain('1m')
    expect(pendingStatus.html()).toContain('id="elapsedTime"')
  })

  it('can display link earned for completed jobs', () => {
    let completedStatus = mountWithTheme(
      <StatusCard title="completed" jobRun={completedRun} />,
    )
    expect(completedStatus.text()).toContain('+2 Link')
  })

  it('will not display link earned for errored or pending jobs', () => {
    let erroredStatus = mountWithTheme(
      <StatusCard title="errored" jobRun={erroredRun} />,
    )
    let pendingStatus = mountWithTheme(
      <StatusCard title="pending_confirmations" jobRun={pendingRun} />,
    )
    expect(erroredStatus.text()).not.toContain('Link')
    expect(pendingStatus.text()).not.toContain('Link')
  })
})
