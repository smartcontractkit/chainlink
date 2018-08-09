import React from 'react'
import { mount } from 'enzyme'
import TokenBalance from 'components/TokenBalance.js'

describe('components/TokenBalance', () => {
  it('renders the title and a loading indicator when it is fetching', () => {
    const wrapper = mount(<TokenBalance title='Ethereum' value={null} />)

    expect(wrapper.text()).toContain('Ethereum...')
  })

  it('renders the title and the error message', () => {
    const wrapper = mount(<TokenBalance title='Ethereum' error='An Error' />)

    expect(wrapper.text()).toContain('EthereumAn Error')
  })

  it('renders the title and the formatted balance', () => {
    const wrapper = mount(<TokenBalance title='Ethereum' value='10120000000000000000000' />)

    expect(wrapper.text()).toContain('Ethereum10.12k')
  })
})
