import { mount } from 'enzyme'
import React from 'react'
import TaskRuns from '../../../components/JobRuns/TaskRuns'

const etherscanHost = 'ropsten.etherscan.io'

describe('components/JobRuns/TaskRuns', () => {
  it('hides incoming pending confirmations with NO minimumConfirmations', () => {
    const taskRuns = [
      {
        id: 1,
        status: 'completed',
        type: 'httpget'
      } as ITaskRun
    ]

    const wrapper = mount(
      <TaskRuns taskRuns={taskRuns} etherscanHost={etherscanHost} />
    )
    expect(wrapper.text()).not.toContain('pending confirmations')
  })

  it('hides incoming pending confirmations with 0 minimumConfirmations', () => {
    const taskRuns = [
      {
        confirmations: 0,
        id: 1,
        minimumConfirmations: 0,
        status: 'completed',
        type: 'httpget'
      } as ITaskRun
    ]

    const wrapper = mount(
      <TaskRuns taskRuns={taskRuns} etherscanHost={etherscanHost} />
    )
    expect(wrapper.text()).not.toContain('pending confirmations')
  })

  it('displays incoming pending confirmations with 0/3 pending confirmations', () => {
    const taskRuns = [
      {
        confirmations: 0,
        id: 1,
        minimumConfirmations: 3,
        status: 'completed',
        type: 'httpget'
      } as ITaskRun
    ]

    const wrapper = mount(
      <TaskRuns taskRuns={taskRuns} etherscanHost={etherscanHost} />
    )
    expect(wrapper.text()).toContain('pending confirmations')
    expect(wrapper.text()).toContain('0')
    expect(wrapper.text()).toContain('3')
  })

  it('displays incoming pending confirmations with 1/3 pending confirmations', () => {
    const taskRuns = [
      {
        confirmations: 1,
        id: 1,
        minimumConfirmations: 3,
        status: 'completed',
        type: 'httpget'
      } as ITaskRun
    ]

    const wrapper = mount(
      <TaskRuns taskRuns={taskRuns} etherscanHost={etherscanHost} />
    )
    expect(wrapper.text()).toContain('pending confirmations')
    expect(wrapper.text()).toContain('1')
    expect(wrapper.text()).toContain('3')
  })
})
