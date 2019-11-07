import { mount } from 'enzyme'
import React from 'react'
import Details from '../../../components/JobRuns/Details'
import { JobRun, ChainlinkNode } from 'explorer/models'
import { mockPartial } from '../../support/mocks'

describe('components/JobRuns/Details', () => {
  it('hides error when not present', () => {
    const chainlinkNode = mockPartial<ChainlinkNode>({})
    const jobRun = mockPartial<JobRun>({ chainlinkNode })

    const wrapper = mount(<Details jobRun={jobRun} etherscanHost="" />)

    expect(wrapper.text()).not.toContain('Error')
  })

  it('displays error when present', () => {
    const chainlinkNode = mockPartial<ChainlinkNode>({})
    const jobRun = mockPartial<JobRun>({ error: 'Failure!', chainlinkNode })

    const wrapper = mount(<Details jobRun={jobRun} etherscanHost="" />)

    expect(wrapper.text()).toContain('Error')
    expect(wrapper.text()).toContain('Failure!')
  })
})
