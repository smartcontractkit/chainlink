import React from 'react'
import { shallow } from 'enzyme'
import { useLoadingPlaceholder } from './useLoadingPlaceholder'

describe('useLoadingPlaceholder', () => {
  it('renders "Loading..." text while loading', () => {
    const { LoadingPlaceholder } = useLoadingPlaceholder(true)
    const wrapper = shallow(<LoadingPlaceholder />)

    expect(wrapper.text()).toContain('<WithStyles(Typography) />')
  })

  it('defaults to false and renders an empty component', () => {
    const { LoadingPlaceholder } = useLoadingPlaceholder()
    const wrapper = shallow(<LoadingPlaceholder />)

    expect(wrapper.text()).toBe('')
  })

  it('exposes "isLoading" variable', () => {
    const { isLoading } = useLoadingPlaceholder(true)

    expect(isLoading).toBe(true)
  })
})
