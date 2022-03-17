import React from 'react'

import { gql } from '@apollo/client'
import { useHistory } from 'react-router-dom'

import Card from '@material-ui/core/Card'
import Grid from '@material-ui/core/Grid'
import Table from '@material-ui/core/Table'
import TableHead from '@material-ui/core/TableHead'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TablePagination from '@material-ui/core/TablePagination'
import TableRow from '@material-ui/core/TableRow'

import Content from 'src/components/Content'
import { Heading1 } from 'src/components/Heading/Heading1'
import { NodeRow } from './NodeRow'
import { SearchTextField } from 'src/components/Search/SearchTextField'

export const NODES_PAYLOAD__RESULTS_FIELDS = gql`
  fragment NodesPayload_ResultsFields on Node {
    id
    chain {
      id
    }
    name
    createdAt
    state
  }
`

const searchIncludes = (searchParam: string) => {
  const lowerCaseSearchParam = searchParam.toLowerCase()

  return (stringToSearch: string) => {
    return stringToSearch.toLowerCase().includes(lowerCaseSearchParam)
  }
}

export const simpleNodeFilter =
  (search: string) => (node: NodesPayload_ResultsFields) => {
    if (search === '') {
      return true
    }

    return matchSimple(node, search)
  }

// matchSimple does a simple match on the id, name, and EVM chain ID
function matchSimple(node: NodesPayload_ResultsFields, term: string) {
  const match = searchIncludes(term)

  const dataset: string[] = [node.id, node.name, node.chain.id]

  return dataset.some(match)
}

export interface Props {
  nodes: ReadonlyArray<NodesPayload_ResultsFields>
  page: number
  pageSize: number
  total: number
}

export const NodesView: React.FC<Props> = ({
  nodes,
  page,
  pageSize,
  total,
}) => {
  const history = useHistory()
  const [search, setSearch] = React.useState('')

  const filteredNodes = React.useMemo(
    () => nodes.filter(simpleNodeFilter(search.trim())),
    [search, nodes],
  )

  return (
    <Content>
      <Grid container>
        <Grid item xs={9}>
          <Heading1>Nodes</Heading1>
        </Grid>

        <Grid item xs={12}>
          <SearchTextField
            value={search}
            onChange={setSearch}
            placeholder="Search nodes"
          />

          <Card>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>ID</TableCell>
                  <TableCell>Name</TableCell>
                  <TableCell>EVM Chain ID</TableCell>
                  <TableCell>Created</TableCell>
                  <TableCell>State</TableCell>
                </TableRow>
              </TableHead>

              <TableBody>
                {filteredNodes.length === 0 && (
                  <TableRow>
                    <TableCell component="th" scope="row" colSpan={3}>
                      No nodes found
                    </TableCell>
                  </TableRow>
                )}

                {filteredNodes.map((node) => (
                  <NodeRow key={node.id} node={node} />
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
                history.push(`/nodes?page=${p + 1}&per=${pageSize}`)
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
