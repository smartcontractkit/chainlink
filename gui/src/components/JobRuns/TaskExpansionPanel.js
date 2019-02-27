import React from 'react'
import PropTypes from 'prop-types'
import Grid from '@material-ui/core/Grid'
import StatusItem from 'components/JobRuns/StatusItem'
import PrettyJson from 'components/PrettyJson'
import capitalize from 'lodash/capitalize'

const renderConfirmations = (confirmations, minimumConfirmations) => {
  if (minimumConfirmations) {
    return (
      <div>
        {confirmations}/{minimumConfirmations} confirmations
      </div>
    )
  }
}

const TaskExpansionPanel = ({ children }) => {
  const initiator = children.initiator
  const taskRuns = children.taskRuns

  return (
    <Grid container spacing={0}>
      <Grid item xs={12}>
        <StatusItem
          summary={capitalize(initiator.type)}
          status={children.status}
          borderTop={false}
        >
          <PrettyJson object={initiator.params} />
        </StatusItem>
      </Grid>
      {taskRuns.map(taskRun => (
        <Grid item xs={12} key={taskRun.id}>
          <StatusItem
            summary={capitalize(taskRun.task.type)}
            status={taskRun.status}
          >
            {renderConfirmations(
              taskRun.task.confirmations,
              taskRun.minimumConfirmations
            )}
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

export default TaskExpansionPanel
