import React from 'react'
import { act } from 'react-dom/test-utils'
import { JobsDefinition } from 'pages/Jobs/Definition'
import jsonApiJobSpecFactory from 'factories/jsonApiJobSpec'
import { Route } from 'react-router-dom'
import { mountWithProviders } from 'test-helpers/mountWithTheme'
import syncFetch from 'test-helpers/syncFetch'
import globPath from 'test-helpers/globPath'

const JOB_SPEC_ID = 'c60b9927eeae43168ddbe92584937b1b'

describe('pages/Jobs/Definition', () => {
  it('renders the job definition component', async () => {
    global.fetch.getOnce(
      globPath(`/v2/specs/${JOB_SPEC_ID}`),
      jsonApiJobSpecFactory({
        id: JOB_SPEC_ID,
      }),
    )

    const wrapper = mountWithProviders(
      <Route path="/jobs/:jobSpecId/definition" component={JobsDefinition} />,
      {
        initialEntries: [`/jobs/${JOB_SPEC_ID}/definition`],
      },
    )

    await act(async () => {
      await syncFetch(wrapper)
      wrapper.update()
    })

    expect(wrapper.text()).toContain(
      'Definition{  "initiators": [    {      "id": 1,      "type": "web",      "jobSpecId": "c60b9927eeae43168ddbe92584937b1b"    }  ],  "tasks": [    {      "confirmations": 0,      "type": "httpget",      "params": {        "get": "https://bitstamp.net/api/ticker/"      }    }  ],  "startAt": "2020-09-22T11:49:50.410Z",  "endAt": "2020-09-22T11:59:50.410Z"}',
    )
  })
})
