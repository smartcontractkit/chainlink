import React from 'react'
import { gql } from '@apollo/client'
import { useHistory } from 'react-router-dom'

import Card from '@material-ui/core/Card'
import Grid from '@material-ui/core/Grid'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableHead from '@material-ui/core/TableHead'
import TablePagination from '@material-ui/core/TablePagination'
import TableRow from '@material-ui/core/TableRow'

import BaseLink from 'components/BaseLink'
import { BridgeRow } from './BridgeRow'
import Button from 'components/Button'
import Content from 'components/Content'
import { Heading1 } from 'src/components/Heading/Heading1'

export const BRIDGES_PAYLOAD__RESULTS_FIELDS = gql`
  fragment BridgesPayload_ResultsFields on Bridge {
    id
    name
    url
    confirmations
    minimumContractPayment
  }
`

export interface Props {
  bridges: ReadonlyArray<BridgesPayload_ResultsFields>
  page: number
  pageSize: number
  total: number
}

export const BridgesView: React.FC<Props> = ({
  bridges,
  page,
  pageSize,
  total,
}) => {
  const history = useHistory()

  return (
    <Content>
      <Grid container spacing={32}>
        <Grid item xs={9}>
          <Heading1>Bridges</Heading1>
        </Grid>
        <Grid item xs={3}>
          <Grid container justify="flex-end">
            <Grid item>
              <Button
                variant="secondary"
                component={BaseLink}
                href={'/bridges/new'}
              >
                New Bridge
              </Button>
            </Grid>
          </Grid>
        </Grid>

        <Grid item xs={12}>
          <Card>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Name</TableCell>
                  <TableCell>URL</TableCell>
                  <TableCell>Default Confirmations</TableCell>
                  <TableCell>Minimum Contract Payment</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {bridges.map((bridge) => (
                  <BridgeRow bridge={bridge} key={bridge.name} />
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
                history.push(`/bridges?page=${p + 1}&per=${pageSize}`)
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
