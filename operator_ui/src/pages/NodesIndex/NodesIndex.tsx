import React from 'react'

import { v2 } from 'api'
import Content from 'components/Content'
import NodesList from './NodesList'
import * as models from 'core/store/models'
import { Title } from 'components/Title'
import { useErrorHandler } from 'hooks/useErrorHandler'
import { useLoadingPlaceholder } from 'hooks/useLoadingPlaceholder'

import Grid from '@material-ui/core/Grid'
import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import SearchIcon from '@material-ui/icons/Search'
import {
  createStyles,
  withStyles,
  WithStyles,
  Theme,
} from '@material-ui/core/styles'
import TextField from '@material-ui/core/TextField'

export type NodeSpecV2 = models.Resource<models.Node>

async function getNodes() {
  return Promise.all([v2.nodes.getNodes()]).then(([v2Nodes]) => {
    const nodesByDate = v2Nodes.data.sort((a: NodeSpecV2, b: NodeSpecV2) => {
      const nodeA = new Date(a.attributes.createdAt).getTime()
      const nodeB = new Date(b.attributes.createdAt).getTime()
      return nodeA > nodeB ? -1 : 1
    })

    return nodesByDate
  })
}

const searchIncludes = (searchParam: string) => {
  const lowerCaseSearchParam = searchParam.toLowerCase()

  return (stringToSearch: string) => {
    return stringToSearch.toLowerCase().includes(lowerCaseSearchParam)
  }
}

export const simpleNodeFilter = (search: string) => (node: NodeSpecV2) => {
  if (search === '') {
    return true
  }

  return matchSimple(node, search)
}

// matchSimple does a simple match on the id, name, and EVM chain ID
function matchSimple(node: NodeSpecV2, term: string) {
  const match = searchIncludes(term)

  const dataset: string[] = [
    node.id,
    node.attributes.name,
    node.attributes.evmChainID,
  ]

  return dataset.some(match)
}

const styles = (theme: Theme) =>
  createStyles({
    card: {
      padding: theme.spacing.unit,
      marginBottom: theme.spacing.unit * 3,
    },
    search: {
      marginBottom: theme.spacing.unit,
    },
  })

export const NodesIndex = ({
  classes,
}: {
  classes: WithStyles<typeof styles>['classes']
}) => {
  const [search, setSearch] = React.useState('')
  const [nodes, setNodes] = React.useState<NodeSpecV2[]>()
  const { error, ErrorComponent, setError } = useErrorHandler()
  const { LoadingPlaceholder } = useLoadingPlaceholder(!error && !nodes)

  React.useEffect(() => {
    document.title = 'Nodes'
  }, [])

  React.useEffect(() => {
    getNodes().then(setNodes).catch(setError)
  }, [setError])

  const nodeFilter = React.useMemo(() => simpleNodeFilter(search.trim()), [
    search,
  ])

  return (
    <Content>
      <Grid container>
        <Grid item xs={9}>
          <Title>Nodes</Title>
        </Grid>

        <Grid item xs={12}>
          <ErrorComponent />
          <LoadingPlaceholder />
          {!error && nodes && (
            <Card className={classes.card}>
              <CardContent>
                <Grid
                  container
                  spacing={8}
                  alignItems="flex-end"
                  className={classes.search}
                >
                  <Grid item>
                    <SearchIcon />
                  </Grid>
                  <Grid item>
                    <TextField
                      label="Search"
                      value={search}
                      name="search"
                      onChange={(event) => setSearch(event.target.value)}
                    />
                  </Grid>
                </Grid>

                <NodesList nodes={nodes} nodeFilter={nodeFilter} />
              </CardContent>
            </Card>
          )}
        </Grid>
      </Grid>
    </Content>
  )
}

export default withStyles(styles)(NodesIndex)
