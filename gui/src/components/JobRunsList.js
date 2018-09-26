import React from 'react'
import PropTypes from 'prop-types'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import JobRunsHead from 'components/JobRunsHead'
import JobRunsRow from 'components/JobRunsRow'
import Card from '@material-ui/core/Card'
import { withStyles } from '@material-ui/core/styles'

const styles = theme => ({
  jobRunsCard: {
    overflow: 'auto'
  }
})

const JobRunsList = ({jobSpecId, runs, classes}) => (
  <Card className={classes.jobRunsCard}>
    <Table>
      <JobRunsHead />
      <TableBody>
        {runs.map(r => <JobRunsRow key={r.id} jobSpecId={jobSpecId} {...r} />)}
      </TableBody>
    </Table>
  </Card>
)

JobRunsList.propTypes = {
  jobSpecId: PropTypes.string.isRequired,
  runs: PropTypes.array.isRequired
}

export default withStyles(styles)(JobRunsList)
