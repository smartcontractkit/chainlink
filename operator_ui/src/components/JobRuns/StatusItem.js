import React from 'react'
import PropTypes from 'prop-types'
import ExpansionPanel from '@material-ui/core/ExpansionPanel'
import ExpansionPanelSummary from '@material-ui/core/ExpansionPanelSummary'
import ExpansionPanelDetails from '@material-ui/core/ExpansionPanelDetails'
import Typography from '@material-ui/core/Typography'
import ExpandMoreIcon from '@material-ui/icons/ExpandMore'
import StatusIcon from 'components/JobRuns/StatusIcon'
import classNames from 'classnames'
import { makeStyles } from '@material-ui/styles'

const useStyles = makeStyles(theme => {
  return {
    borderTop: {
      borderTop: 'solid 1px',
      borderTopColor: theme.palette.divider
    },
    item: {
      position: 'relative',
      paddingLeft: 50
    },
    status: {
      position: 'absolute',
      top: 0,
      left: 0,
      paddingTop: 25,
      paddingLeft: 30,
      borderRight: 'solid 1px',
      borderRightColor: theme.palette.divider,
      width: 50,
      height: '100%'
    },
    summary: {
      minHeight: '0 !important'
    },
    content: {
      margin: '12px 0 !important'
    },
    details: {
      padding: theme.spacing(2)
    },
    expansionPanel: {
      boxShadow: 'none'
    }
  }
})

const render = (summary, children) => {
  const classes = useStyles()
  if (children) {
    return (
      <ExpansionPanel className={classes.expansionPanel}>
        <ExpansionPanelSummary
          className={classes.summary}
          classes={{ content: classes.content }}
          expandIcon={<ExpandMoreIcon />}
        >
          <Typography variant="h5">{summary}</Typography>
        </ExpansionPanelSummary>
        <ExpansionPanelDetails>{children}</ExpansionPanelDetails>
      </ExpansionPanel>
    )
  }

  return <Typography>{summary}</Typography>
}

const StatusItem = ({ status, summary, borderTop, children }) => {
  const classes = useStyles()
  return (
    <div
      className={classNames(classes.item, { [classes.borderTop]: borderTop })}
    >
      <div className={classes.status}>
        <StatusIcon width={38} height={38}>
          {status}
        </StatusIcon>
      </div>
      <div className={classes.details}>
        {render(summary, children, classes)}
      </div>
    </div>
  )
}

StatusItem.defaultProps = {
  borderTop: true
}

StatusItem.propTypes = {
  status: PropTypes.string.isRequired,
  borderTop: PropTypes.bool.isRequired
}

export default StatusItem
