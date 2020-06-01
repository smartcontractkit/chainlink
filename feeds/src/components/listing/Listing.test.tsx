import { partialAsFull } from '@chainlink/ts-helpers'
import React from 'react'
import { MemoryRouter } from 'react-router-dom'
import '@testing-library/jest-dom/extend-expect'
import { render } from '@testing-library/react'
import { FeedConfig } from 'config'
import { Provider as ReduxProvider } from 'react-redux'
import { ListingGroup } from 'state/ducks/listing/selectors'
import createStore from '../../state/createStore'
import { Listing } from './Listing'
import { Contract } from 'ethers'
import * as utils from '../../contracts/utils'

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
  beforeAll(() => {
    jest.spyOn(utils, 'formatAnswer').mockImplementation(answer => answer)
    jest.spyOn(utils, 'createContract').mockImplementation(() => {
      return partialAsFull<Contract>({
        latestAnswer: () => 'latestAnswer',
        currentAnswer: () => 'currentAnswer',
        latestTimestamp: () => 1590703158,
      })
    })
  })

  it('renders a loading message', () => {
    const { container } = render(
      <AllTheProviders>
        <Listing
          loadingFeeds={true}
          feedGroups={[]}
          fetchFeeds={jest.fn()}
          enableDetails={false}
        />
      </AllTheProviders>,
    )

    expect(container).toHaveTextContent('Loading Feeds...')
  })

  it('renders the name from a list of groups', () => {
    const { container } = render(
      <AllTheProviders>
        <Listing
          loadingFeeds={false}
          feedGroups={listingGroups}
          fetchFeeds={jest.fn()}
          enableDetails={false}
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
          loadingFeeds={false}
          feedGroups={listingGroups}
          fetchFeeds={jest.fn()}
          enableDetails={false}
        />
      </AllTheProviders>,
    )

    expect(container).toHaveTextContent('pair name 1')
    expect(container).toHaveTextContent('pair name 2')
    expect(container).toHaveTextContent('pair name 3')
    expect(container).toHaveTextContent('pair name 4')
  })
})
