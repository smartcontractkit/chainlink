import React from 'react'
import Card from '@material-ui/core/Card'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import JobRunsHead from 'components/JobRunsHead'
import JobRunsRow from 'components/JobRunsRow'

const JobRunsList = ({runs}) => (
  <Card>
    <Table>
      <JobRunsHead />
      <TableBody>
        {runs.map(r => <JobRunsRow key={r.id} {...r} />)}
      </TableBody>
    </Table>
  </Card>
)

export default JobRunsList
