import React from 'react'
import { mount } from 'enzyme'
import TokenBalanceCard from 'components/Cards/TokenBalance.js'

describe('components/Cards/TokenBalance', () => {
  it('renders the title and a loading indicator when it is fetching', () => {
    const wrapper = mount(
      <TokenBalanceCard title="Ether Balance" value={null} />
    )

    expect(wrapper.text()).toContain('Ether Balance...')
  })

  it('renders the title and the error message', () => {
    const wrapper = mount(
      <TokenBalanceCard title="Ether Balance" error="An Error" />
    )

    expect(wrapper.text()).toContain('Ether BalanceAn Error')
  })

  it('renders the title and the formatted balance', () => {
    const wrapper = mount(
      <TokenBalanceCard title="Ether Balance" value="7779070000000000000000" />
    )

    expect(wrapper.text()).toContain('Ether Balance7.779070k')
  })
})
