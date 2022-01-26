import React from 'react'
import { partialAsFull } from 'support/test-helpers/partialAsFull'
import Activity from '../../../src/components/Dashboards/Activity'
import { JobRunV2 } from 'core/store/models'
import { cleanup, renderWithRouter, screen } from 'support/test-utils'

const { queryByText } = screen

const CREATED_AT = '2019-06-11T14:37:42.077995-07:00'

describe('components/Dashboards/Activity', () => {
  it('displays the given runs', () => {
    const runs = [
      {
        id: 'runA',
        type: 'RUN',
        attributes: partialAsFull<JobRunV2>({
          createdAt: CREATED_AT,
          errors: [],
          pipelineSpec: {
            jobID: '1',
            ID: 1,
            dotDagSource: 'dot',
            CreatedAt: '1',
          },
        }),
      },
    ]

    renderWithRouter(<Activity runs={runs} pageSize={1} />)

    expect(queryByText('Run: runA')).toBeInTheDocument()
  })

  it('displays a "View More" link when there is more than 1 page of runs', () => {
    const runs = [
      {
        id: 'runB',
        type: 'RUN',
        attributes: partialAsFull<JobRunV2>({
          createdAt: CREATED_AT,
          errors: [],
          pipelineSpec: {
            jobID: '1',
            ID: 1,
            dotDagSource: 'dot',
            CreatedAt: '1',
          },
        }),
      },
      {
        id: 'runC',
        type: 'RUN',
        attributes: partialAsFull<JobRunV2>({
          createdAt: CREATED_AT,
          errors: [],
          pipelineSpec: {
            jobID: '1',
            ID: 1,
            dotDagSource: 'dot',
            CreatedAt: '1',
          },
        }),
      },
    ]

    renderWithRouter(<Activity runs={runs} pageSize={1} count={2} />)
    expect(queryByText('View More')).toBeInTheDocument()

    cleanup()

    renderWithRouter(<Activity runs={runs} pageSize={2} count={2} />)
    expect(queryByText('View More')).toBeNull()
  })

  it('can show a loading message', () => {
    renderWithRouter(<Activity pageSize={1} />)
    expect(queryByText('Loading ...')).toBeInTheDocument()
  })

  it('can show a no activity message', () => {
    renderWithRouter(<Activity runs={[]} pageSize={1} />)
    expect(queryByText('No recent activity')).toBeInTheDocument()
  })
})
