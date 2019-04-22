import React from 'react'
import { mount } from 'enzyme'
import Table from '../../components/Table'

const HEADERS = ['First Name', 'Last Name']

describe('components/Table', () => {
  it('renders table headers', () => {
    const wrapper = mount(<Table headers={HEADERS} />)

    expect(wrapper.text()).toContain('First Name')
    expect(wrapper.text()).toContain('Last Name')
  })

  it('renders the columns from each row', () => {
    const rows = [
      [{ type: 'text', text: 'Michael' }, { type: 'text', text: 'Jordan' }],
      [{ type: 'text', text: 'Charles' }, { type: 'text', text: 'Barkley' }]
    ]
    const wrapper = mount(<Table headers={HEADERS} rows={rows} />)

    expect(wrapper.text()).toContain('Michael')
    expect(wrapper.text()).toContain('Jordan')
    expect(wrapper.text()).toContain('Charles')
    expect(wrapper.text()).toContain('Barkley')
  })

  it('renders a default loading message when rows are undefined', () => {
    const wrapper = mount(<Table headers={HEADERS} rows={undefined} />)

    expect(wrapper.text()).toContain('Loading...')
  })

  it('can override the default loading message', () => {
    const wrapper = mount(
      <Table headers={HEADERS} rows={undefined} loadingMsg="LOADING" />
    )

    expect(wrapper.text()).not.toContain('Loading...')
    expect(wrapper.text()).toContain('LOADING')
  })

  it('renders a default empty message when there are no rows', () => {
    const wrapper = mount(<Table headers={HEADERS} rows={[]} />)

    expect(wrapper.text()).toContain('No results')
  })

  it('can override the default empty message', () => {
    const wrapper = mount(
      <Table headers={HEADERS} rows={[]} emptyMsg="EMPTY" />
    )

    expect(wrapper.text()).not.toContain('No results')
    expect(wrapper.text()).toContain('EMPTY')
  })
})
