import React from 'react'
import { syncFetch } from 'test-helpers/syncFetch'
import { mountWithProviders } from 'test-helpers/mountWithTheme'
import { Route } from 'react-router-dom'
import { Show } from 'pages/Bridges/Show'
import globPath from 'test-helpers/globPath'

describe('pages/Bridges/Show', () => {
  it('renders the details of the bridge spec', async () => {
    expect.assertions(5)
    const response = {
      data: {
        id: 'tallbridge',
        type: 'bridges',
        attributes: {
          name: 'Tall Bridge',
          url: 'https://localhost.com:712/endpoint',
          confirmations: 9,
          outgoingToken: 'outgoingToken',
        },
      },
    }

    global.fetch.getOnce(globPath(`/v2/bridge_types/tallbridge`), response)

    const wrapper = mountWithProviders(
      <Route path="/bridges/:bridgeId" component={Show} />,
      {
        initialEntries: [`/bridges/tallbridge`],
      },
    )

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('Tall Bridge')
    expect(wrapper.text()).toContain('Confirmations')
    expect(wrapper.text()).toContain('https://localhost.com:712/endpoint')
    expect(wrapper.text()).toContain('outgoingToken')
    expect(wrapper.text()).toContain('9')
  })

  it('deletes a bridge', async () => {
    const response = {
      data: {
        id: 'tallbridge',
        type: 'bridges',
        attributes: {
          name: 'Tall Bridge',
          url: 'https://localhost.com:712/endpoint',
          confirmations: 9,
          outgoingToken: 'outgoingToken',
        },
      },
    }

    global.fetch.getOnce(globPath(`/v2/bridge_types/tallbridge`), response)

    const wrapper = mountWithProviders(
      <Route path="/bridges/:bridgeId" component={Show} />,
      {
        initialEntries: [`/bridges/tallbridge`],
      },
    )

    await syncFetch(wrapper)

    wrapper
      .find('Button')
      .find({ children: 'Delete' })
      .first()
      .simulate('click')

    await syncFetch(wrapper)

    global.fetch.deleteOnce(globPath(`/v2/bridge_types/tallbridge`), {})

    wrapper
      .find('Button')
      .find({ children: 'Confirm' })
      .first()
      .simulate('click')

    await syncFetch(wrapper)

    const routerComponentProps: any = wrapper.find('Router').props()
    expect(routerComponentProps?.history?.location?.pathname).toEqual(
      '/bridges',
    )
  })

  it('fails to delete a bridge', async () => {
    const response = {
      data: {
        id: 'tallbridge',
        type: 'bridges',
        attributes: {
          name: 'Tall Bridge',
          url: 'https://localhost.com:712/endpoint',
          confirmations: 9,
          outgoingToken: 'outgoingToken',
        },
      },
    }

    global.fetch.getOnce(globPath(`/v2/bridge_types/tallbridge`), response)

    const wrapper = mountWithProviders(
      <Route path="/bridges/:bridgeId" component={Show} />,
      {
        initialEntries: [`/bridges/tallbridge`],
      },
    )

    await syncFetch(wrapper)

    wrapper
      .find('Button')
      .find({ children: 'Delete' })
      .first()
      .simulate('click')

    await syncFetch(wrapper)

    global.fetch.deleteOnce(globPath(`/v2/bridge_types/tallbridge`), {
      body: {
        errors: [
          {
            detail: "can't remove the bridge",
          },
        ],
      },
      status: 409,
    })

    wrapper
      .find('Button')
      .find({ children: 'Confirm' })
      .first()
      .simulate('click')

    await syncFetch(wrapper)

    const routerComponentProps: any = wrapper.find('Router').props()
    expect(routerComponentProps?.history?.location?.pathname).toEqual(
      '/bridges/tallbridge',
    )
  })
})
