import React from 'react'
import render from '../../support/test-helpers/renderWithinRouter'
import BaseLink from '../../src/components/BaseLink'

describe('components/BaseLink', () => {
  it('renders an anchor', () => {
    const component = render(<BaseLink href="/foo">My Link</BaseLink>)
    expect(component.text()).toContain('My Link')
    expect(component.prop('href')).toEqual('/foo')
  })

  it('can render an id', () => {
    const component = render(
      <BaseLink id="my-id" href="/foo">
        My Link
      </BaseLink>,
    )
    expect(component.prop('id')).toEqual('my-id')
  })

  it('can render a css class', () => {
    const component = render(
      <BaseLink className="my-css-class" href="/foo">
        My Link
      </BaseLink>,
    )
    expect(component.prop('class')).toEqual('my-css-class')
  })
})
