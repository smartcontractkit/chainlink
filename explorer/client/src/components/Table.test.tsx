import React from 'react'
import { mount } from 'enzyme'
import Table from './Table'

const HEADERS = ['First Name', 'Last Name']

describe('components/Table', () => {
  it('renders table headers', () => {
    const wrapper = mount(
      <Table
        headers={HEADERS}
        loading={false}
        error={false}
        onChangePage={jest.fn()}
      />,
    )

    expect(wrapper.text()).toContain('First Name')
    expect(wrapper.text()).toContain('Last Name')
  })

  it('renders the columns from each row', () => {
    const rows: React.ComponentPropsWithoutRef<typeof Table>['rows'] = [
      [
        { type: 'text', text: 'Michael' },
        { type: 'text', text: 'Jordan' },
      ],
      [
        { type: 'text', text: 'Charles' },
        { type: 'text', text: 'Barkley' },
      ],
    ]
    const wrapper = mount(
      <Table
        rows={rows}
        loading={false}
        error={false}
        headers={HEADERS}
        onChangePage={jest.fn()}
      />,
    )

    expect(wrapper.text()).toContain('Michael')
    expect(wrapper.text()).toContain('Jordan')
    expect(wrapper.text()).toContain('Charles')
    expect(wrapper.text()).toContain('Barkley')
  })

  describe('loading', () => {
    it('renders a default message', () => {
      const wrapper = mount(
        <Table
          loading={true}
          error={false}
          headers={HEADERS}
          onChangePage={jest.fn()}
        />,
      )

      expect(wrapper.text()).toContain('Loading...')
    })

    it('can provide a custom message', () => {
      const wrapper = mount(
        <Table
          loading={true}
          error={false}
          loadingMsg="CUSTOM LOADING..."
          headers={HEADERS}
          onChangePage={jest.fn()}
        />,
      )

      expect(wrapper.text()).not.toContain('Loading...')
      expect(wrapper.text()).toContain('CUSTOM LOADING...')
    })
  })

  describe('empty', () => {
    it('renders a default message', () => {
      const wrapper = mount(
        <Table
          count={0}
          loading={false}
          error={false}
          headers={HEADERS}
          onChangePage={jest.fn()}
        />,
      )

      expect(wrapper.text()).toContain(
        'There are no results added to the Explorer yet',
      )
    })

    it('can provide a custom message', () => {
      const wrapper = mount(
        <Table
          count={0}
          emptyMsg="CUSTOM EMPTY"
          loading={false}
          error={false}
          headers={HEADERS}
          onChangePage={jest.fn()}
        />,
      )

      expect(wrapper.text()).not.toContain(
        'There are no results added to the Explorer yet',
      )
      expect(wrapper.text()).toContain('CUSTOM EMPTY')
    })
  })

  describe('error', () => {
    it('renders a default message', () => {
      const wrapper = mount(
        <Table
          error={true}
          loading={false}
          headers={HEADERS}
          onChangePage={jest.fn()}
        />,
      )

      expect(wrapper.text()).toContain('Error loading resources')
    })

    it('can provide a custom message', () => {
      const wrapper = mount(
        <Table
          error={true}
          loading={false}
          errorMsg="CUSTOM ERROR"
          headers={HEADERS}
          onChangePage={jest.fn()}
        />,
      )

      expect(wrapper.text()).not.toContain('Error loading resources')
      expect(wrapper.text()).toContain('CUSTOM ERROR')
    })
  })
})
