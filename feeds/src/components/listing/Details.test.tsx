import { partialAsFull } from '@chainlink/ts-helpers'
import React from 'react'
import { MemoryRouter } from 'react-router-dom'
import '@testing-library/jest-dom/extend-expect'
import { render } from '@testing-library/react'
import { FeedConfig } from 'config'
import { Provider as ReduxProvider } from 'react-redux'
import createStore from '../../state/createStore'
import { DetailsContent } from './Details'
import { humanizeUnixTimestamp } from '../../utils'

const AllTheProviders: React.FC = ({ children }) => {
  const { store } = createStore()

  return (
    <ReduxProvider store={store}>
      <MemoryRouter>{children}</MemoryRouter>
    </ReduxProvider>
  )
}

const feed = partialAsFull<FeedConfig>({
  name: 'pair name',
  path: '/link',
  valuePrefix: '$',
  sponsored: ['sponsor 1', 'sponsor 2'],
})

describe('components/listing/Details', () => {
  it('renders popover with details', () => {
    const { container } = render(
      <AllTheProviders>
        <DetailsContent
          feed={feed}
          answer={'10.1'}
          healthCheckPrice={10.2}
          healthCheckStatus={{ result: 'OK', errors: [] }}
          answerTimestamp={1591005300}
          healthClasses={'ok'}
        />
      </AllTheProviders>,
    )
    expect(container).toHaveTextContent('$ 10.1')
    expect(container).toHaveTextContent('$ 10.2')
    expect(container).toHaveTextContent(
      humanizeUnixTimestamp(1591005300, 'LLL'),
    )
  })
})
