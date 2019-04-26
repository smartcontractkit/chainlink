import React from 'react'
import { mount } from 'enzyme'
import Details from '../../../components/JobRuns/Details'

describe('components/JobRuns/Details', () => {
  it('hides error when not present', () => {
    const chainlinkNode = {}
    const jobRun = { chainlinkNode } as IJobRun
    const wrapper = mount(<Details jobRun={jobRun} />)

    expect(wrapper.text()).not.toContain('Error')
  })

  it('displays error when present', () => {
    const chainlinkNode = {}
    const jobRun = { error: 'Failure!', chainlinkNode } as IJobRun
    const wrapper = mount(<Details jobRun={jobRun} />)

    expect(wrapper.text()).toContain('Error')
    expect(wrapper.text()).toContain('Failure!')
  })
})
