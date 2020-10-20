import React from 'react'
import { act } from 'react-dom/test-utils'
import { JobsDefinition } from 'pages/Jobs/Definition'
import jsonApiJobSpecFactory from 'factories/jsonApiJobSpec'
import { MemoryRouter, Route } from 'react-router-dom'
import mountWithTheme from 'test-helpers/mountWithTheme'
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

    const wrapper = mountWithTheme(
      <MemoryRouter initialEntries={[`/jobs/${JOB_SPEC_ID}/definition`]}>
        <Route path="/jobs/:jobSpecId/definition" component={JobsDefinition} />
      </MemoryRouter>,
    )

    await act(async () => {
      await syncFetch(wrapper)
      wrapper.update()
    })

    expect(wrapper.text()).toContain(
      '{  "initiators": [    {      "type": "web"    }  ],  "tasks": [    {      "confirmations": 0,      "type": "httpget",      "url": "https://bitstamp.net/api/ticker/"    }  ],  "startAt": undefined,  "endAt": undefined}',
    )
  })
})
