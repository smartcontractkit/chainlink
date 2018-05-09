/* eslint-env jest */
import React from 'react'
import { mount } from 'enzyme'
import { Jobs } from 'containers/Jobs.js'
import jobSpecFactory from 'factories/jobSpec'
import syncFetch from 'test-helpers/syncFetch'

describe('containers/Jobs', () => {
  it('renders the list of jobs', async () => {
    expect.assertions(3)

    const response = jobSpecFactory({
      id: 'c60b9927eeae43168ddbe92584937b1b',
      initiators: [{'type': 'web'}],
      createdAt: '2018-05-10T00:41:54.531043837Z'
    })
    global.fetch.getOnce('/v2/specs', response)

    const classes = {}
    const location = {}
    const wrapper = mount(<Jobs classes={classes} location={location} />)

    await syncFetch(wrapper).then(() => {
      expect(wrapper.text()).toContain('c60b9927eeae43168ddbe92584937b1b')
      expect(wrapper.text()).toContain('web')
      expect(wrapper.text()).toContain('2018-05-10T00:41:54.531043837Z')
    })
  })

  it('displays an error message when the network requst fails', async () => {
    global.fetch.catch(() => { throw new TypeError('Failed to fetch') })

    const classes = {}
    const location = {}
    const wrapper = mount(<Jobs classes={classes} location={location} />)

    await syncFetch(wrapper).then(() => {
      expect(wrapper.text()).toContain('There was an error fetching the jobs. Please reload the page.')
    })
  })
})
