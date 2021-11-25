import React from 'react'

import { v2 } from 'api'
import BaseLink from 'components/BaseLink'
import Button from 'components/Button'
import Content from 'components/Content'
import { ChainRow } from './ChainRow'
import { Resource, Chain } from 'core/store/models'
import { SearchTextField } from 'src/components/SearchTextField'
import { useErrorHandler } from 'hooks/useErrorHandler'
import { useLoadingPlaceholder } from 'hooks/useLoadingPlaceholder'

import Grid from '@material-ui/core/Grid'
import Card from '@material-ui/core/Card'
import CardHeader from '@material-ui/core/CardHeader'
import CardContent from '@material-ui/core/CardContent'
import { createStyles, withStyles, WithStyles } from '@material-ui/core/styles'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'
import Typography from '@material-ui/core/Typography'

import { Heading1 } from 'src/components/Heading/Heading1'

export type ChainResource = Resource<Chain>

async function getChains() {
  return Promise.all([v2.chains.getChains()]).then(([v2Chains]) => {
    const chainsByDate = v2Chains.data.sort(
      (a: ChainResource, b: ChainResource) => {
        const chainA = new Date(a.attributes.createdAt).getTime()
        const chainB = new Date(b.attributes.createdAt).getTime()
        return chainA > chainB ? -1 : 1
      },
    )

    return chainsByDate
  })
}

const searchIncludes = (searchParam: string) => {
  const lowerCaseSearchParam = searchParam.toLowerCase()

  return (stringToSearch: string) => {
    return stringToSearch.toLowerCase().includes(lowerCaseSearchParam)
  }
}

export const simpleChainFilter = (search: string) => (chain: ChainResource) => {
  if (search === '') {
    return true
  }

  return matchSimple(chain, search)
}

// matchSimple does a simple match on the id
function matchSimple(chain: ChainResource, term: string) {
  const match = searchIncludes(term)

  const dataset: string[] = [chain.id]

  return dataset.some(match)
}

const styles = () =>
  createStyles({
    cardHeader: {
      borderBottom: 0,
    },
  })

export const ChainsIndex = ({
  classes,
}: {
  classes: WithStyles<typeof styles>['classes']
}) => {
  const [search, setSearch] = React.useState('')
  const [chains, setChains] = React.useState<ChainResource[]>()
  const { error, ErrorComponent, setError } = useErrorHandler()
  const { LoadingPlaceholder } = useLoadingPlaceholder(!error && !chains)

  React.useEffect(() => {
    document.title = 'Chains'
  }, [])

  React.useEffect(() => {
    getChains().then(setChains).catch(setError)
  }, [setError])

  const chainFilter = React.useMemo(
    () => simpleChainFilter(search.trim()),
    [search],
  )

  return (
    <Content>
      <Grid container>
        <Grid item xs={9}>
          <Heading1>Chains</Heading1>
        </Grid>

        <Grid item xs={3}>
          <Grid container justify="flex-end">
            <Grid item>
              <Button
                variant="secondary"
                component={BaseLink}
                href={'/chains/new'}
              >
                New Chain
              </Button>
            </Grid>
          </Grid>
        </Grid>

        <Grid item xs={12}>
          <ErrorComponent />
          <LoadingPlaceholder />
          {!error && chains && (
            <Card>
              <CardHeader
                title={<SearchTextField value={search} onChange={setSearch} />}
                className={classes.cardHeader}
              />
              <CardContent>
                <Table>
                  <TableHead>
                    <TableRow>
                      <TableCell>
                        <Typography variant="body1" color="textSecondary">
                          Chain ID
                        </Typography>
                      </TableCell>

                      <TableCell>
                        <Typography variant="body1" color="textSecondary">
                          Enabled
                        </Typography>
                      </TableCell>

                      <TableCell>
                        <Typography variant="body1" color="textSecondary">
                          Config overrides
                        </Typography>
                      </TableCell>

                      <TableCell>
                        <Typography variant="body1" color="textSecondary">
                          Created
                        </Typography>
                      </TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {chains.filter(chainFilter).map((chain) => (
                      <ChainRow key={chain.id} chain={chain} />
                    ))}
                  </TableBody>
                </Table>
              </CardContent>
            </Card>
          )}
        </Grid>
      </Grid>
    </Content>
  )
}

export default withStyles(styles)(ChainsIndex)
