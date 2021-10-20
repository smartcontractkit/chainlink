import React from 'react'
import { jsonApiOcrJobRun } from 'factories/jsonApiOcrJobRun'
import { mountWithProviders } from 'test-helpers/mountWithTheme'
import { transformPipelineJobRun } from '../../transformJobRuns'
import { PipelineJobRunOverview } from './PipelineJobRunOverview'

describe('PipelineJobRunOverview', () => {
  it('displays an overview & json tab by default', () => {
    const error = 'something something error'
    const apiResponse = jsonApiOcrJobRun({
      id: '1',
      errors: [error],
      dotDagSource: `
    fetch    [type=http method=POST url="http://localhost:8001" requestData="{\\"hi\\": \\"hello\\"}"];
    parse    [type=jsonparse path="data,result"];
    multiply [type=multiply times=100];
    fetch -> parse -> multiply;
`,
    })

    const jobRun = transformPipelineJobRun('1')(apiResponse.data)

    const component = mountWithProviders(
      <PipelineJobRunOverview jobRun={jobRun} />,
    )

    // There should be 2 dividers for 3 tasks
    expect(component.find('Divider').length).toEqual(2)

    // Should contain attributes of the failed task
    // eslint-disable-next-line no-useless-escape
    expect(component.text()).toContain(`requestData: {\"hi\": \"hello\"}`)
    expect(component.text()).toContain('url: http://localhost:8001')
    expect(component.text()).toContain('method: POST')

    // Should not contain attributes of tasks that were not run
    expect(component.text()).not.toContain('path: data,result')
    expect(component.text()).not.toContain('times: 100')
  })
})
