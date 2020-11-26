import React from 'react'
import StatusCard from './StatusCard'
import mountWithTheme from 'test-helpers/mountWithTheme'

describe('components/StatusCard', () => {
  const start = '2020-01-03T22:45:00.166261Z'
  const end1m = '2020-01-03T22:46:00.166261Z'
  const pendingRun = {
    id: 'runA',
    status: 'pending',
    result: {},
    createdAt: start,
    finishedAt: null,
  }
  const completedRun = {
    id: 'runA',
    status: 'completed',
    createdAt: start,
    finishedAt: end1m,
    payment: '2000000000000000000',
  }
  const erroredRun = {
    id: 'runA',
    status: 'errored',
    result: {},
    createdAt: start,
    finishedAt: end1m,
  }
  it('converts the given title to title case', () => {
    const component = mountWithTheme(
      <StatusCard title={'pending_incoming_confirmations'} />,
    )
    expect(component.text()).toContain('Pending Incoming Confirmations')
  })

  it('can display the elapsed time for finished jobruns', () => {
    const erroredStatus = mountWithTheme(
      <StatusCard title="errored" {...erroredRun} />,
    )
    const completedStatus = mountWithTheme(
      <StatusCard title="completed" {...completedRun} />,
    )

    expect(erroredStatus.text()).toContain('1m')
    expect(completedStatus.text()).toContain('1m')
  })

  it('displays a live elapsed time for pending job runs', () => {
    const now2m = '2020-01-03T22:47:00.166261Z'

    jest
      .spyOn(Date, 'now')
      .mockImplementationOnce(() => new Date(now2m).valueOf())

    const pendingStatus = mountWithTheme(
      <StatusCard title="pending" {...pendingRun} />,
    )

    expect(pendingStatus.html()).toContain('2m')
  })

  it('can display link earned for completed jobs', () => {
    const completedStatus = mountWithTheme(
      <StatusCard title="completed" {...completedRun} />,
    )
    expect(completedStatus.text()).toContain('+2 Link')
  })

  it('will not display link earned for errored or pending jobs', () => {
    const erroredStatus = mountWithTheme(
      <StatusCard title="errored" {...erroredRun} />,
    )
    const pendingStatus = mountWithTheme(
      <StatusCard title="pending_incoming_confirmations" {...pendingRun} />,
    )
    expect(erroredStatus.text()).not.toContain('Link')
    expect(pendingStatus.text()).not.toContain('Link')
  })
})
