import React from 'react'
import PropTypes from 'prop-types'
import { withStyles } from '@material-ui/core/styles'
import Grid from '@material-ui/core/Grid'
import StatusItem from 'components/StatusItem'
import PrettyJson from 'components/PrettyJson'
import capitalize from 'lodash/capitalize'

const styles = theme => {
  return {
  }
}

const renderConfirmations = (confirmations, minimumConfirmations) => {
  if (minimumConfirmations) {
    return <div>{confirmations}/{minimumConfirmations} confirmations</div>
  }
}

const TaskExpansionPanel = ({children, classes}) => {
  const initiator = children.initiator
  const taskRuns = children.taskRuns

  return (
    <Grid container>
      <Grid item xs={12}>
        <StatusItem summary={capitalize(initiator.type)} status={children.status}>
          <PrettyJson object={initiator.params} />
        </StatusItem>
      </Grid>
      {taskRuns.map(taskRun => (
        <Grid item xs={12} key={taskRun.id}>
          <StatusItem
            summary={capitalize(taskRun.task.type)}
            status={taskRun.status}
          >
            {renderConfirmations(taskRun.task.confirmations, taskRun.minimumConfirmations)}
            <PrettyJson object={taskRun.task.params} />
          </StatusItem>
        </Grid>
      ))}
    </Grid>
  )
}

TaskExpansionPanel.propTypes = {
  children: PropTypes.object.isRequired
}

export default withStyles(styles)(TaskExpansionPanel)
