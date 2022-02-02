import React from 'react'
import { gql } from '@apollo/client'
import { Switch, Route, useParams, useRouteMatch } from 'react-router-dom'
import { v2 } from 'api'
import { NodeResource } from './ChainNodes'
import RegionalNav from './RegionalNav'
import { Resource, Chain } from 'core/store/models'
import { ChainNodes } from './ChainNodes'
import { ChainConfig } from './ChainConfig'
import NewChainNode from './NewChainNode'
import UpdateChain from './UpdateChain'
import { ChainPayload } from 'src/types/generated/graphql'

export const CHAIN_PAYLOAD__NODES_FIELDS = gql`
  fragment ChainPayload_NodesFields on Node {
    id
    name
    httpURL
    wsURL
    createdAt
  }
`

export const CHAIN_PAYLOAD__FIELDS = gql`
  ${CHAIN_PAYLOAD__NODES_FIELDS}
  fragment ChainPayload_Fields on Chain {
    id
    enabled
    createdAt
    nodes {
      ...ChainPayload_NodesFields
    }
  }
`

export interface Props {
  chain: ChainPayload_Fields
  onDelete: () => void
}

interface RouteParams {
  id: string
}

export const ChainView: React.FC<Props> = ({ chain, onDelete }) => {
  const { id } = useParams<RouteParams>()
  const { path } = useRouteMatch()

  React.useEffect(() => {
    document.title = `Chain ${id}`
  }, [id])

  return (
    <>
      <RegionalNav chainId={id} chain={chain} />
      <Switch>
        <Route path={`${path}/nodes/new`}>
          {<NewChainNode chain={chain} />}
        </Route>
        <Route path={`${path}/edit`}>{<UpdateChain chain={chain} />}</Route>
        <Route path={`${path}/config-overrides`}>
          {<ChainConfig chain={chain} />}
        </Route>
        <Route path={`${path}`}>
          {<ChainNodes nodes={chain.nodes} chain={chain} />}
        </Route>
      </Switch>
    </>
  )
}
