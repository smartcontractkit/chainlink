import React from 'react'
import { jobRunAPIResponse } from 'factories/jsonApiOcrJobRun'
import { transformPipelineJobRun } from '../../transformJobRuns'
import { PipelineJobRunOverview } from './PipelineJobRunOverview'
import { render, screen } from 'support/test-utils'

const { queryByText } = screen

describe('PipelineJobRunOverview', () => {
  it('displays an overview & json tab by default', () => {
    const error = 'something something error'
    const apiResponse = jobRunAPIResponse({
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

    render(<PipelineJobRunOverview jobRun={jobRun} />)

    // Should contain attributes of the failed task
    // eslint-disable-next-line no-useless-escape
    expect(queryByText(/: \{"hi": "hello"\}/i)).toBeInTheDocument()
    expect(queryByText(/: http:\/\/localhost:8001/i)).toBeInTheDocument()

    // Should not contain attributes of tasks that were not run
    expect(queryByText(': data,result')).toBeNull()
    expect(queryByText(': 100')).toBeNull()
  })
})
