import { partialAsFull } from '@chainlink/ts-helpers'
import React from 'react'
import { MemoryRouter } from 'react-router-dom'
import '@testing-library/jest-dom/extend-expect'
import { render } from '@testing-library/react'
import { FeedConfig } from 'config'
import { Provider as ReduxProvider } from 'react-redux'
import createStore from '../../state/createStore'
import { GridItem } from './GridItem'

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

describe('components/listing/GridItem', () => {
  it('renders answer value with prefix', () => {
    const { container } = render(
      <AllTheProviders>
        <GridItem
          feed={feed}
          answer={'10.1'}
          fetchLatestData={jest.fn()}
          fetchHealthStatus={jest.fn()}
          enableDetails={false}
        />
      </AllTheProviders>,
    )

    expect(container).toHaveTextContent('10.1')
    expect(container).toHaveTextContent('$')
  })
})
