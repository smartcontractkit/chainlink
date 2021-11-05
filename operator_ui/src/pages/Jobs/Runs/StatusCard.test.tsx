import React from 'react'
import StatusCard from './StatusCard'
import { cleanup, render, screen } from 'support/test-utils'

const { queryByText } = screen

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
    render(<StatusCard title={'pending_incoming_confirmations'} />)

    expect(queryByText('Pending Incoming Confirmations')).toBeInTheDocument()
  })

  it('can display the elapsed time for finished jobruns', () => {
    render(<StatusCard title="errored" {...erroredRun} />)
    expect(queryByText('1m0s')).toBeInTheDocument()

    cleanup()

    render(<StatusCard title="completed" {...completedRun} />)
    expect(queryByText('1m0s')).toBeInTheDocument()
  })

  it('displays a live elapsed time for pending job runs', () => {
    const now2m = '2020-01-03T22:47:00.166261Z'

    jest
      .spyOn(Date, 'now')
      .mockImplementationOnce(() => new Date(now2m).valueOf())

    render(<StatusCard title="pending" {...pendingRun} />)

    expect(queryByText('2m0s')).toBeInTheDocument()
  })

  it('can display link earned for completed jobs', () => {
    render(<StatusCard title="completed" {...completedRun} />)
    expect(queryByText('+2 Link')).toBeInTheDocument()
  })

  it('will not display link earned for errored or pending jobs', () => {
    render(<StatusCard title="errored" {...erroredRun} />)
    expect(queryByText('Link')).not.toBeInTheDocument()

    cleanup()

    render(
      <StatusCard title="pending_incoming_confirmations" {...pendingRun} />,
    )
    expect(queryByText('Link')).not.toBeInTheDocument()
  })
})
