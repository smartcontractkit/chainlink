import { mount } from 'enzyme'
import React from 'react'
import { partialAsFull } from '@chainlink/ts-test-helpers'
import { JobRun, ChainlinkNode } from 'explorer/models'
import Details from '../../../components/JobRuns/Details'

describe('components/JobRuns/Details', () => {
  it('hides error when not present', () => {
    const chainlinkNode = partialAsFull<ChainlinkNode>({})
    const jobRun = partialAsFull<JobRun>({ chainlinkNode })

    const wrapper = mount(<Details jobRun={jobRun} etherscanHost="" />)

    expect(wrapper.text()).not.toContain('Error')
  })

  it('displays error when present', () => {
    const chainlinkNode = partialAsFull<ChainlinkNode>({})
    const jobRun = partialAsFull<JobRun>({ error: 'Failure!', chainlinkNode })

    const wrapper = mount(<Details jobRun={jobRun} etherscanHost="" />)

    expect(wrapper.text()).toContain('Error')
    expect(wrapper.text()).toContain('Failure!')
  })
})
