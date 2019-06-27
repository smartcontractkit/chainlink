import React from 'react'
import { withStyles } from '@material-ui/core/styles'
import { useHooks, useState } from 'use-react-hooks'
import { connect } from 'react-redux'
import { bindActionCreators } from 'redux'
import PaddedCard from '@chainlink/styleguide/components/PaddedCard'
import Divider from '@material-ui/core/Divider'
import Typography from '@material-ui/core/Typography'
import Button from 'components/Button'
import Grid from '@material-ui/core/Grid'
import Slide from '@material-ui/core/Slide'
import { deleteCompletedJobRuns, deleteErroredJobRuns } from 'actions'

const styles = theme => {
  return {
    deleteRunsDivider: {
      marginTop: theme.spacing.unit * 3,
      marginBottom: theme.spacing.unit * 2
    }
  }
}

const WEEK_MS = 1000 * 60 * 60 * 24 * 7

const DeleteJobRuns = useHooks(
  ({ classes, deleteCompletedJobRuns, deleteErroredJobRuns }) => {
    const updatedBefore = new Date(Date.now() - WEEK_MS).toISOString()
    const [showCompletedConfirm, setCompletedConfirm] = useState(false)
    const [showErroredConfirm, setErroredConfirm] = useState(false)

    const confirmCompletedDelete = () => {
      deleteCompletedJobRuns(updatedBefore)
      setCompletedConfirm(false)
    }
    const confirmErroredDeleted = () => {
      deleteErroredJobRuns(updatedBefore)
      setErroredConfirm(false)
    }

    return (
      <PaddedCard>
        <Typography variant="h5" color="secondary">
          Delete Runs
        </Typography>

        <Typography variant="subtitle1" color="textSecondary">
          Reduce your database size by deleting completed and or errored job
          runs. This action will delete runs a week ago and before.
        </Typography>

        <Divider className={classes.deleteRunsDivider} />

        <Grid container spacing={16}>
          <Grid item xs={12}>
            {!showCompletedConfirm && (
              <Slide
                direction="right"
                in={!showCompletedConfirm}
                mountOnEnter
                unmountOnExit
              >
                <Button onClick={() => setCompletedConfirm(true)}>
                  Delete Completed Jobs
                </Button>
              </Slide>
            )}
            {showCompletedConfirm && (
              <Slide
                direction="right"
                in={showCompletedConfirm}
                mountOnEnter
                unmountOnExit
              >
                <Button
                  variant="danger"
                  onClick={() => confirmCompletedDelete()}
                >
                  Confirm delete all completed job runs up to {updatedBefore}
                </Button>
              </Slide>
            )}
          </Grid>
          <Grid item xs={12}>
            {!showErroredConfirm && (
              <Slide
                direction="right"
                in={!showErroredConfirm}
                mountOnEnter
                unmountOnExit
              >
                <Button onClick={() => setErroredConfirm(true)}>
                  Delete Errored Jobs
                </Button>
              </Slide>
            )}
            {showErroredConfirm && (
              <Slide
                direction="right"
                in={showErroredConfirm}
                mountOnEnter
                unmountOnExit
              >
                <Button
                  variant="danger"
                  onClick={() => confirmErroredDeleted()}
                >
                  Confirm delete all errored job runs up to {updatedBefore}
                </Button>
              </Slide>
            )}
          </Grid>
        </Grid>
      </PaddedCard>
    )
  }
)

const mapStateToProps = state => ({})

const mapDispatchToProps = dispatch =>
  bindActionCreators(
    {
      deleteCompletedJobRuns,
      deleteErroredJobRuns
    },
    dispatch
  )

export const ConnectedDeleteJobRuns = connect(
  mapStateToProps,
  mapDispatchToProps
)(DeleteJobRuns)

export default withStyles(styles)(ConnectedDeleteJobRuns)
