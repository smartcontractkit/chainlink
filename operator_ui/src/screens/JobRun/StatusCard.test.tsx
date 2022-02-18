import React from 'react'

import { StatusCard } from './StatusCard'
import { render, screen } from 'support/test-utils'

const { queryByText } = screen

describe('components/StatusCard', () => {
  const start = '2020-01-03T22:45:00.166261Z'
  const end1m = '2020-01-03T22:46:00.166261Z'

  it('converts the given title to title case', () => {
    render(<StatusCard status="COMPLETED" startedAt={start} />)

    expect(queryByText('Completed')).toBeInTheDocument()
  })

  it('displays the elapsed time for completed job runs', () => {
    render(
      <StatusCard status="COMPLETED" startedAt={start} finishedAt={end1m} />,
    )
    expect(queryByText('1m0s')).toBeInTheDocument()
  })

  it('displays the elapsed time for errored job runs', () => {
    render(<StatusCard status="ERRORED" startedAt={start} finishedAt={end1m} />)

    expect(queryByText('1m0s')).toBeInTheDocument()
  })

  it('displays a live elapsed time for running job runs', () => {
    const now2m = '2020-01-03T22:47:00.166261Z'

    const spy = jest
      .spyOn(Date, 'now')
      .mockImplementation(() => new Date(now2m).valueOf())

    render(<StatusCard status="RUNNING" startedAt={start} />)

    expect(queryByText('2m0s')).toBeInTheDocument()

    spy.mockReset()
  })
})
