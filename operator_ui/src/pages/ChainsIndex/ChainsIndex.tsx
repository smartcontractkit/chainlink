import React from 'react'

import { v2 } from 'api'
import BaseLink from 'components/BaseLink'
import Button from 'components/Button'
import Content from 'components/Content'
import { ChainRow } from './ChainRow'
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
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'
import Typography from '@material-ui/core/Typography'
import TextField from '@material-ui/core/TextField'

export type ChainSpecV2 = models.Resource<models.Chain>

async function getChains() {
  return Promise.all([v2.chains.getChains()]).then(([v2Chains]) => {
    const chainsByDate = v2Chains.data.sort(
      (a: ChainSpecV2, b: ChainSpecV2) => {
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

export const simpleChainFilter = (search: string) => (chain: ChainSpecV2) => {
  if (search === '') {
    return true
  }

  return matchSimple(chain, search)
}

// matchSimple does a simple match on the id
function matchSimple(chain: ChainSpecV2, term: string) {
  const match = searchIncludes(term)

  const dataset: string[] = [chain.id]

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

export const ChainsIndex = ({
  classes,
}: {
  classes: WithStyles<typeof styles>['classes']
}) => {
  const [search, setSearch] = React.useState('')
  const [chains, setChains] = React.useState<ChainSpecV2[]>()
  const { error, ErrorComponent, setError } = useErrorHandler()
  const { LoadingPlaceholder } = useLoadingPlaceholder(!error && !chains)

  React.useEffect(() => {
    document.title = 'Chains'
  }, [])

  React.useEffect(() => {
    getChains().then(setChains).catch(setError)
  }, [setError])

  const chainFilter = React.useMemo(() => simpleChainFilter(search.trim()), [
    search,
  ])

  return (
    <Content>
      <Grid container>
        <Grid item xs={9}>
          <Title>Chains</Title>
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
