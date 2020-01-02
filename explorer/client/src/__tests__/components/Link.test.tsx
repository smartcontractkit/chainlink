import React from 'react'
import { mount } from 'enzyme'
import Link from '../../components/Link'

describe('components/Link', () => {
  it('renders internal links', () => {
    const wrapper = mount(<Link to="/foo" />)
    expect(wrapper.find('a').prop('href')).toContain('/foo')
  })

  it('renders external http links', () => {
    const wrapper = mount(<Link to="http://foo" />)
    expect(wrapper.find('a').prop('href')).toContain('http://foo')
  })

  it('renders external https links', () => {
    const wrapper = mount(<Link to="https://foo" />)
    expect(wrapper.find('a').prop('href')).toContain('https://foo')
  })
})
