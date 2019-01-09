import React from 'react'
import { render } from 'enzyme'
import KeyValueList from 'components/KeyValueList'

const renderComponent = props => (
  render(
    <KeyValueList {...props} />
  )
)

describe('components/KeyValueList', () => {
  it('can display a header', () => {
    let withHeader = renderComponent({entries: [], showHead: true})
    expect(withHeader.text()).toContain('Key')
    expect(withHeader.text()).toContain('Value')

    let withoutHeader = renderComponent({entries: []})
    expect(withoutHeader.text()).not.toContain('Key')
    expect(withoutHeader.text()).not.toContain('Value')
  })
})
