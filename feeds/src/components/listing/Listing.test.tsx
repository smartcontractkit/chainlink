import React from 'react'
import { MemoryRouter } from 'react-router-dom'
import '@testing-library/jest-dom/extend-expect'
import { render } from '@testing-library/react'
import { Provider as ReduxProvider } from 'react-redux'
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

const groupMock = [
  {
    name: 'List 1',
    list: [
      {
        answer: 'answer',
        config: {
          name: 'pair name 1',
          path: '/link',
          valuePrefix: 'prefix ',
          sponsored: ['sponsor 1', 'sponsor 2'],
        },
      },
      {
        answer: 'answer2',
        config: {
          name: 'pair name 2',
          path: '/link2',
          valuePrefix: 'prefix2',
          sponsored: ['sponsor 1', 'sponsor 2'],
        },
      },
    ],
  },
  {
    name: 'List 2',
    list: [
      {
        answer: 'answer',
        config: {
          name: 'pair name 3',
          path: '/link',
          valuePrefix: 'prefix',
          sponsored: ['sponsor 1', 'sponsor 2'],
        },
      },
      {
        answer: 'answer2',
        config: {
          name: 'pair name 4',
          path: '/link2',
          valuePrefix: 'prefix2',
          sponsored: ['sponsor 1', 'sponsor 2'],
        },
      },
    ],
  },
]

describe('components/listing/Listing.component', () => {
  it('renders the name from a list of groups', () => {
    const { container } = render(
      <AllTheProviders>
        <Listing
          groups={groupMock}
          fetchAnswers={() => {}}
          fetchHealthStatus={() => {}}
          enableHealth={false}
        />
      </AllTheProviders>,
    )

    expect(container).toHaveTextContent('List 1 Pairs')
    expect(container).toHaveTextContent('List 2 Pairs')
  })

  it('should renders pair name value', () => {
    const { container } = render(
      <AllTheProviders>
        <Listing
          groups={groupMock}
          fetchAnswers={() => {}}
          fetchHealthStatus={() => {}}
          enableHealth={false}
        />
      </AllTheProviders>,
    )

    expect(container).toHaveTextContent('pair name 1')
    expect(container).toHaveTextContent('pair name 2')
    expect(container).toHaveTextContent('pair name 3')
    expect(container).toHaveTextContent('pair name 4')
  })

  it('should renders answer value with prefix', () => {
    const { container } = render(
      <AllTheProviders>
        <Listing
          groups={groupMock}
          fetchAnswers={() => {}}
          fetchHealthStatus={() => {}}
          enableHealth={false}
        />
      </AllTheProviders>,
    )

    expect(container).toHaveTextContent('prefix')
    expect(container).toHaveTextContent('answer')
    expect(container).toHaveTextContent('prefix answer')
    expect(container).toHaveTextContent('prefix2 answer2')
  })

  it('should renders sponsored names', () => {
    const { container } = render(
      <AllTheProviders>
        <Listing
          groups={groupMock}
          fetchAnswers={() => {}}
          fetchHealthStatus={() => {}}
          enableHealth={false}
        />
      </AllTheProviders>,
    )

    expect(container).toHaveTextContent('sponsor 1')
    expect(container).toHaveTextContent('sponsor 2')
  })
})
