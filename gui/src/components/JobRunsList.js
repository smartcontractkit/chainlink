import React from 'react'
import PropTypes from 'prop-types'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import JobRunsHead from 'components/JobRunsHead'
import JobRunsRow from 'components/JobRunsRow'

const JobRunsList = ({jobSpecId, runs}) => (
  <Table>
    <JobRunsHead />
    <TableBody>
      {runs.map(r => <JobRunsRow key={r.id} jobSpecId={jobSpecId} {...r} />)}
    </TableBody>
  </Table>
)

JobRunsList.propTypes = {
  jobSpecId: PropTypes.string.isRequired,
  runs: PropTypes.array.isRequired
}

export default JobRunsList
