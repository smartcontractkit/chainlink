/* eslint-env jest */
import { JobsIndex as Index } from 'pages/Jobs/Index'
import jsonApiJobSpecsFactory from 'factories/jsonApiJobSpecs'
import React from 'react'
import { Route } from 'react-router-dom'
import { act } from 'react-dom/test-utils'
import syncFetch from 'test-helpers/syncFetch'
import globPath from 'test-helpers/globPath'
import { mountWithProviders } from 'test-helpers/mountWithTheme'

describe('pages/Jobs/Index', () => {
  it('renders the list of jobs', async () => {
    const jobSpecsResponse = jsonApiJobSpecsFactory([
      {
        id: 'c60b9927eeae43168ddbe92584937b1b',
        initiators: [{ type: 'web' }],
        createdAt: new Date().toISOString(),
      },
    ])
    global.fetch.getOnce(globPath('/v2/specs'), jobSpecsResponse)

    const wrapper = mountWithProviders(<Route component={Index} />)

    await act(async () => {
      await syncFetch(wrapper)
    })
    expect(wrapper.text()).toContain('c60b9927eeae43168ddbe92584937b1b')
    expect(wrapper.text()).toContain('web')
    expect(wrapper.text()).toContain('just now')
  })
})
