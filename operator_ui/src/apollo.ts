import { ApolloClient, InMemoryCache, HttpLink } from '@apollo/client'
import generatedIntrospection from 'src/types/generated/possibleTypes'

const baseURL = process.env.CHAINLINK_BASEURL ?? location.origin

const httpLink = new HttpLink({
  uri: `${baseURL}/query`,
  credentials: 'include',
})

export const client = new ApolloClient({
  cache: new InMemoryCache({
    possibleTypes: generatedIntrospection.possibleTypes,
  }),
  link: httpLink,
})
