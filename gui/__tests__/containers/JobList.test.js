/* eslint-env jest */
import React from 'react'
import { mount } from 'enzyme'
import { JobList } from 'containers/JobList.js'
import jobSpecFactory from 'factories/jobSpec'
import syncFetch from 'test-helpers/syncFetch'

describe('containers/JobList', () => {
  it('renders the list of jobs', async () => {
    expect.assertions(3)

    const response = jobSpecFactory([{
      id: 'c60b9927eeae43168ddbe92584937b1b',
      initiators: [{'type': 'web'}],
      createdAt: '2018-05-10T00:41:54.531043837Z'
    }])
    global.fetch.getOnce('/v2/specs', response)

    const wrapper = mount(<JobList />)

    await syncFetch(wrapper).then(() => {
      expect(wrapper.text()).toContain('c60b9927eeae43168ddbe92584937b1b')
      expect(wrapper.text()).toContain('web')
      expect(wrapper.text()).toContain('2018-05-10T00:41:54.531043837Z')
    })
  })

  it('displays an error message when the jobs network request fails', async () => {
    expect.assertions(1)

    global.fetch.catch(() => { throw new TypeError('Failed to fetch') })

    const wrapper = mount(<JobList />)

    await syncFetch(wrapper).then(() => {
      expect(wrapper.text()).toContain(
        'There was an error fetching the jobs. Please reload the page.'
      )
    })
  })
})
