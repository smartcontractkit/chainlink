import React from 'react'
import { Switch, Route, useParams, useRouteMatch } from 'react-router-dom'
import { v2 } from 'api'
import { NodeSpecV2 } from '../NodesIndex/NodesIndex'
import RegionalNav from './RegionalNav'
import { ChainSpecV2 } from '../ChainsIndex/ChainsIndex'
import { ChainNodes } from './ChainNodes'
import { ChainConfig } from './ChainConfig'
import NewChainNode from './NewChainNode'
import UpdateChain from './UpdateChain'

interface RouteParams {
  chainId: string
}

export const ChainsShow = () => {
  const { chainId } = useParams<RouteParams>()
  const { path } = useRouteMatch()
  const [chain, setChain] = React.useState<ChainSpecV2>()
  const [nodes, setNodes] = React.useState<NodeSpecV2[]>([])

  const getNodes = async () => {
    const nodes = await v2.nodes.getNodes()
    setNodes(nodes.data)
  }

  React.useEffect(() => {
    getNodes()
  }, [])

  React.useEffect(() => {
    document.title = `Chain ${chainId}`
  }, [chainId])

  React.useEffect(() => {
    Promise.all([v2.chains.getChains()])
      .then(([v2Chains]) =>
        v2Chains.data.find((chain: ChainSpecV2) => chain.id === chainId),
      )
      .then(setChain)
  }, [chainId])

  return (
    <>
      <RegionalNav chainId={chainId} chain={chain} />
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
