import React from 'react'
import { partialAsFull } from 'support/test-helpers/partialAsFull'
import mountWithTheme from '../../../support/test-helpers/mountWithTheme'
import Activity from '../../../src/components/Dashboards/Activity'
import { JobRunV2 } from 'core/store/models'

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
    const component = mountWithTheme(<Activity runs={runs} pageSize={1} />)
    expect(component.text()).toContain('Run: runA')
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

    const componentWithMore = mountWithTheme(
      <Activity runs={runs} pageSize={1} count={2} />,
    )
    expect(componentWithMore.text()).toContain('View More')

    const componentWithoutMore = mountWithTheme(
      <Activity runs={runs} pageSize={2} count={2} />,
    )
    expect(componentWithoutMore.text()).not.toContain('View More')
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
