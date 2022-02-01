import React from 'react'
import { Switch, Route, useParams, useRouteMatch } from 'react-router-dom'
import { v2 } from 'api'
import { NodeResource } from './ChainNodes'
import RegionalNav from './RegionalNav'
import { Resource, Chain } from 'core/store/models'
import { ChainNodes } from './ChainNodes'
import { ChainConfig } from './ChainConfig'
import NewChainNode from './NewChainNode'
import UpdateChain from './UpdateChain'

export type ChainResource = Resource<Chain>

interface RouteParams {
  id: string
}

export const ChainsShow = () => {
  const { id } = useParams<RouteParams>()
  const { path } = useRouteMatch()
  const [chain, setChain] = React.useState<ChainResource>()
  const [nodes, setNodes] = React.useState<NodeResource[]>([])

  const getNodes = async () => {
    const nodes = await v2.nodes.getNodes()
    setNodes(nodes.data)
  }

  React.useEffect(() => {
    getNodes()
  }, [])

  React.useEffect(() => {
    document.title = `Chain ${id}`
  }, [id])

  React.useEffect(() => {
    Promise.all([v2.chains.getChains()])
      .then(([v2Chains]) =>
        v2Chains.data.find((chain: ChainResource) => chain.id === id),
      )
      .then(setChain)
  }, [id])

  return (
    <>
      <RegionalNav chainId={id} chain={chain} />
      <Switch>
        <Route path={`${path}/nodes/new`}>
          {chain && <NewChainNode chain={chain} />}
        </Route>
        <Route path={`${path}/edit`}>
          {chain && <UpdateChain chain={chain} />}
        </Route>
        <Route path={`${path}/config-overrides`}>
          {chain && <ChainConfig chain={chain} />}
        </Route>
        <Route path={`${path}`}>
          {chain && <ChainNodes nodes={nodes} chain={chain} />}
        </Route>
      </Switch>
    </>
  )
}

export default ChainsShow
