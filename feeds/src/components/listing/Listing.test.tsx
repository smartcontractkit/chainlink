import { partialAsFull } from '@chainlink/ts-helpers'
import '@testing-library/jest-dom/extend-expect'
import { render } from '@testing-library/react'
import { FeedConfig } from 'config'
import React from 'react'
import { Provider as ReduxProvider } from 'react-redux'
import { MemoryRouter } from 'react-router-dom'
import { ListingGroup } from 'state/ducks/listing/selectors'
import createStore from '../../state/createStore'
import { Listing } from './Listing'

const AllTheProviders: React.FC = ({ children }) => {
  const { store } = createStore()

  return (
    <ReduxProvider store={store}>
      <MemoryRouter>{children}</MemoryRouter>
    </ReduxProvider>
  )
}

const listingGroup1 = {
  name: 'List 1',
  feeds: [
    partialAsFull<FeedConfig>({
      name: 'pair name 1',
      path: '/link',
      valuePrefix: '$',
      sponsored: ['sponsor 1', 'sponsor 2'],
    }),
    partialAsFull<FeedConfig>({
      name: 'pair name 2',
      path: '/link2',
      valuePrefix: 'Ξ',
      sponsored: ['sponsor 1', 'sponsor 2'],
    }),
  ],
}
const listingGroup2 = {
  name: 'List 2',
  feeds: [
    partialAsFull<FeedConfig>({
      name: 'pair name 3',
      path: '/link',
      valuePrefix: '$',
      sponsored: ['sponsor 1', 'sponsor 2'],
    }),
    partialAsFull<FeedConfig>({
      name: 'pair name 4',
      path: '/link2',
      valuePrefix: 'Ξ',
      sponsored: ['sponsor 1', 'sponsor 2'],
    }),
  ],
}
const listingGroups: ListingGroup[] = [listingGroup1, listingGroup2]

describe('components/listing/Listing', () => {
  it('renders the name from a list of groups', () => {
    const { container } = render(
      <AllTheProviders>
        <Listing
          groups={listingGroups}
          fetchAnswers={() => {}}
          fetchHealthStatus={() => {}}
          enableHealth={false}
          compareOffchain={false}
        />
      </AllTheProviders>,
    )

    expect(container).toHaveTextContent('List 1 Pairs')
    expect(container).toHaveTextContent('List 2 Pairs')
  })

  it('renders pair name value', () => {
    const { container } = render(
      <AllTheProviders>
        <Listing
          groups={listingGroups}
          fetchAnswers={() => {}}
          fetchHealthStatus={() => {}}
          enableHealth={false}
          compareOffchain={false}
        />
      </AllTheProviders>,
    )

    expect(container).toHaveTextContent('pair name 1')
    expect(container).toHaveTextContent('pair name 2')
    expect(container).toHaveTextContent('pair name 3')
    expect(container).toHaveTextContent('pair name 4')
  })

  it('renders sponsored names', () => {
    const { container } = render(
      <AllTheProviders>
        <Listing
          groups={listingGroups}
          fetchAnswers={() => {}}
          fetchHealthStatus={() => {}}
          enableHealth={false}
          compareOffchain={false}
        />
      </AllTheProviders>,
    )

    expect(container).toHaveTextContent('sponsor 1')
    expect(container).toHaveTextContent('sponsor 2')
  })
})
