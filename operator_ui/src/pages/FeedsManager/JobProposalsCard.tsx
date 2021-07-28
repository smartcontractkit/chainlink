import React from 'react'

import Card from '@material-ui/core/Card'
import CardHeader from '@material-ui/core/CardHeader'
import { createStyles, WithStyles, withStyles } from '@material-ui/core/styles'
import Tab from '@material-ui/core/Tab'
import Tabs from '@material-ui/core/Tabs'
import Table from '@material-ui/core/Table'
import TableCell from '@material-ui/core/TableCell'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'

const styles = () => {
  return createStyles({
    tabsRoot: {
      borderBottom: '1px solid #e8e8e8',
    },
  })
}

interface Props extends WithStyles<typeof styles> {
  proposals: any // TODO - Implement job proposals
}

// Placeholder
export const JobProposalsCard = withStyles(styles)(({ classes }: Props) => {
  const [tabValue, setTabValue] = React.useState(0)

  return (
    <Card>
      <CardHeader title="Job Proposals" />

      <Tabs
        value={tabValue}
        className={classes.tabsRoot}
        indicatorColor="primary"
        onChange={(_, value) => {
          setTabValue(value)
        }}
      >
        <Tab label="Pending" />
        <Tab label="Approved" />
        <Tab label="Rejected" />
      </Tabs>

      <Table>
        <TableHead>
          <TableRow>
            <TableCell>ID</TableCell>
            <TableCell>External Job ID</TableCell>
            <TableCell>Status</TableCell>
            <TableCell>Proposed At</TableCell>
          </TableRow>
        </TableHead>
      </Table>
    </Card>
  )
})
