import { ApolloClient, InMemoryCache, HttpLink } from '@apollo/client'

const httpLink = new HttpLink({
  uri: `${process.env.CHAINLINK_BASEURL}/query`,
  credentials: 'same-origin',
})

export const client = new ApolloClient({
  cache: new InMemoryCache(),
  link: httpLink,
})
