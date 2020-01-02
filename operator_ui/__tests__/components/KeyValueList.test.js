import { KeyValueList } from '@chainlink/styleguide'
import { render } from 'enzyme'
import React from 'react'

describe('components/KeyValueList', () => {
  it('can display a title', () => {
    const withTitle = render(<KeyValueList entries={[]} title="My Title" />)
    expect(withTitle.text()).toContain('My Title')
  })

  it('can display header columns', () => {
    const withHead = render(<KeyValueList entries={[]} showHead />)
    expect(withHead.text()).toContain('Key')
    expect(withHead.text()).toContain('Value')

    const withoutHead = render(<KeyValueList entries={[]} />)
    expect(withoutHead.text()).not.toContain('Key')
    expect(withoutHead.text()).not.toContain('Value')
  })

  it('renders entry pairs', () => {
    const entries = [
      ['CHAINLINK_DEV', 'true'],
      ['DATABASE_TIMEOUT', 1000],
    ]

    const wrapper = render(<KeyValueList entries={entries} />)
    expect(wrapper.text()).toContain('CHAINLINK_DEV')
    expect(wrapper.text()).toContain('true')
    expect(wrapper.text()).toContain('DATABASE_TIMEOUT')
    expect(wrapper.text()).toContain('1000')
  })

  it('can titleize keys', () => {
    const entries = [
      ['gasLimit', 50000],
      ['hash', 'abc123'],
    ]

    const withoutTitleize = render(<KeyValueList entries={entries} />)
    expect(withoutTitleize.text()).toContain('gasLimit')
    expect(withoutTitleize.text()).toContain('hash')
    expect(withoutTitleize.text()).toContain('50000')
    expect(withoutTitleize.text()).toContain('abc123')

    const withTitleize = render(<KeyValueList entries={entries} titleize />)
    expect(withTitleize.text()).toContain('Gas Limit')
    expect(withTitleize.text()).toContain('Hash')
    expect(withTitleize.text()).toContain('50000')
    expect(withTitleize.text()).toContain('abc123')
  })
})
