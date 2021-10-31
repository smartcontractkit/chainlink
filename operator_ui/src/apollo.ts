import { ApolloClient, InMemoryCache, HttpLink } from '@apollo/client'
import generatedIntrospection from 'src/types/generated/possibleTypes'

const httpLink = new HttpLink({
  uri: `${process.env.CHAINLINK_BASEURL}/query`,
  credentials: 'same-origin',
})

export const client = new ApolloClient({
  cache: new InMemoryCache({
    possibleTypes: generatedIntrospection.possibleTypes,
  }),
  link: httpLink,
})
