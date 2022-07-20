import React, { useMemo } from 'react'

import { useHistory } from 'react-router-dom'

import Grid from '@material-ui/core/Grid'
import Card from '@material-ui/core/Card'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableHead from '@material-ui/core/TableHead'
import TablePagination from '@material-ui/core/TablePagination'
import TableRow from '@material-ui/core/TableRow'

import BaseLink from 'src/components/BaseLink'
import { BetaAlert } from './BetaAlert'
import Button from 'src/components/Button'
import { ChainRow } from './ChainRow'
import Content from 'src/components/Content'
import { Heading1 } from 'src/components/Heading/Heading1'
import { SearchTextField } from 'src/components/Search/SearchTextField'

const searchIncludes = (searchParam: string) => {
  const lowerCaseSearchParam = searchParam.toLowerCase()

  return (stringToSearch: string) => {
    return stringToSearch.toLowerCase().includes(lowerCaseSearchParam)
  }
}

export const simpleChainFilter =
  (search: string) => (chain: ChainsPayload_ResultsFields) => {
    if (search === '') {
      return true
    }

    return matchSimple(chain, search)
  }

// matchSimple does a simple match on the id
function matchSimple(chain: ChainsPayload_ResultsFields, term: string) {
  const match = searchIncludes(term)

  const dataset: string[] = [chain.id]

  return dataset.some(match)
}

export interface Props {
  chains: ReadonlyArray<ChainsPayload_ResultsFields>
  page: number
  pageSize: number
  total: number
}

export const ChainsView: React.FC<Props> = ({
  chains,
  page,
  pageSize,
  total,
}) => {
  const history = useHistory()
  const [search, setSearch] = React.useState('')

  const filteredChains = useMemo(
    () => chains.filter(simpleChainFilter(search.trim())),
    [search, chains],
  )

  return (
    <Content>
      <Grid container>
        <Grid item xs={12}>
          <BetaAlert />
        </Grid>

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
          <SearchTextField
            value={search}
            onChange={setSearch}
            placeholder="Search chains"
          />

          <Card>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Chain ID</TableCell>
                  <TableCell>Enabled</TableCell>
                  <TableCell>Created</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {filteredChains.length === 0 && (
                  <TableRow>
                    <TableCell component="th" scope="row" colSpan={3}>
                      No chains found
                    </TableCell>
                  </TableRow>
                )}

                {filteredChains.map((chain) => (
                  <ChainRow key={chain.id} chain={chain} />
                ))}
              </TableBody>
            </Table>
            <TablePagination
              component="div"
              count={total}
              rowsPerPage={pageSize}
              rowsPerPageOptions={[pageSize]}
              page={page - 1}
              onChangePage={(_, p) => {
                history.push(`/chains?page=${p + 1}&per=${pageSize}`)
              }}
              onChangeRowsPerPage={() => {}} /* handler required by component, so make it a no-op */
              backIconButtonProps={{ 'aria-label': 'prev-page' }}
              nextIconButtonProps={{ 'aria-label': 'next-page' }}
            />
          </Card>
        </Grid>
      </Grid>
    </Content>
  )
}
