import { CardTitle } from 'components/CardTitle'
import { Card, Grid } from '@material-ui/core'
import Content from 'components/Content'
import React from 'react'
import { ChainSpecV2 } from '../ChainsIndex/ChainsIndex'
import ChainNodesList from '../NodesIndex/NodesList'
import { NodeSpecV2 } from '../NodesIndex/NodesIndex'

interface Props {
  nodes: NodeSpecV2[]
  chain: ChainSpecV2
}

export const ChainNodes = ({ nodes, chain }: Props) => {
  React.useEffect(() => {
    document.title = chain?.id ? `Chain ${chain.id} | Nodes` : 'Chain | Nodes'
  }, [chain])

  return (
    <Content>
      {chain && (
        <Grid container spacing={40}>
          <Grid item xs={12}>
            <Card>
              <CardTitle divider>Nodes</CardTitle>
              <ChainNodesList
                nodes={nodes}
                nodeFilter={(node) => node.attributes.evmChainID === chain.id}
              />
            </Card>
          </Grid>
        </Grid>
      )}
    </Content>
  )
}
