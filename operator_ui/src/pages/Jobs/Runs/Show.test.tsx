import React from 'react'
import { Route } from 'react-router-dom'
import { jobRunAPIResponse } from 'factories/jsonApiOcrJobRun'
import { Show } from './Show'
import isoDate, { MINUTE_MS } from 'test-helpers/isoDate'
import globPath from 'test-helpers/globPath'
import { renderWithRouter, screen } from 'support/test-utils'

const { findByText } = screen

describe('pages/JobRuns/Show/Overview', () => {
  const JOB_ID = '200'
  const RUN_ID = '100'

  it('renders the details of the job spec and its latest runs', async () => {
    const minuteAgo = isoDate(Date.now() - MINUTE_MS)

    // Mock the job runs fetch
    const run = {
      jobID: JOB_ID,
      createdAt: minuteAgo,
      outputs: [null],
      errors: ['task inputs: too many errors'],
      inputs: {
        ds: {},
        ds_parse: {},
        jobRun: {
          meta: null,
          requestBody: '',
        },
        jobSpec: {
          databaseID: 38,
          externalJobID: '0eec7e1d-d0d2-476c-a1a8-72dfb6633f46',
          name: '',
        },
      },
      taskRuns: [
        {
          type: 'http',
          createdAt: '2021-07-27T19:26:37.312358+08:00',
          finishedAt: '2021-07-27T19:26:38.116848+08:00',
          output: null,
          error:
            'got error from https://chain.link/ETH-USD: (status code 404) ',
          dotId: 'ds',
        },
        {
          type: 'jsonparse',
          createdAt: '2021-07-27T19:26:38.11691+08:00',
          finishedAt: '2021-07-27T19:26:38.117104+08:00',
          output: null,
          error: 'task inputs: too many errors',
          dotId: 'ds_parse',
        },
      ],
      pipelineSpec: {
        ID: 40,
        dotDagSource:
          '\t\t\t\tds          [type=http method=GET url="https://chain.link/ETH-USD"];\n\t\t\t\tds_parse    [type=jsonparse path="data,price"];\n\t\t\t\tds -\u003e ds_parse;\n\t\t\t',
        CreatedAt: '2021-07-27T19:26:38.117104+08:00',
        jobID: JOB_ID,
      },
    }

    global.fetch.getOnce(
      globPath(`/v2/jobs/${JOB_ID}/runs/${RUN_ID}`),
      jobRunAPIResponse(run),
    )

    renderWithRouter(
      <Route path={`/jobs/:jobId/runs/:jobRunId`} component={Show} />,
      {
        initialEntries: [`/jobs/${JOB_ID}/runs/${RUN_ID}`],
      },
    )

    expect(await findByText('Errored')).toBeInTheDocument()
    expect(await findByText('Task list')).toBeInTheDocument()
    expect(await findByText('Errors')).toBeInTheDocument()
    expect(await findByText('task inputs: too many errors')).toBeInTheDocument()
  })
})
