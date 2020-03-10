import React, { useEffect } from 'react'
import { connect, MapDispatchToProps, MapStateToProps } from 'react-redux'
import { RouteComponentProps } from '@reach/router'
import build from 'redux-object'
import {
  createStyles,
  withStyles,
  Theme,
  WithStyles,
} from '@material-ui/core/styles'
import Grid from '@material-ui/core/Grid'
import Title from '../../../components/Title'
import Head from '../../../components/Admin/Heads/Head'
import { fetchAdminHead } from '../../../actions/adminHeads'
import { AppState } from '../../../reducers'
import { HeadShowData } from '../../../reducers/adminHeadsShow'
import { DispatchBinding } from '../../../utils/types'

const styles = ({ breakpoints, spacing }: Theme) =>
  createStyles({
    container: {
      overflow: 'hidden',
      padding: spacing.unit * 2,
      [breakpoints.up('sm')]: {
        padding: spacing.unit * 3,
      },
    },
  })

interface OwnProps extends RouteComponentProps<{ headId: number }> {}

interface StateProps {
  headData: HeadShowData
}

interface DispatchProps {
  fetchAdminHead: DispatchBinding<typeof fetchAdminHead>
}

interface Props
  extends WithStyles<typeof styles>,
    StateProps,
    DispatchProps,
    OwnProps {}

const Show: React.FC<Props> = ({
  classes,
  headData,
  headId,
  fetchAdminHead,
}) => {
  useEffect(() => {
    if (headId) {
      fetchAdminHead(headId)
    }
  }, [fetchAdminHead, headId])

  return (
    <Grid
      container
      spacing={24}
      alignItems="center"
      className={classes.container}
    >
      <Grid item xs={12}>
        <Title>Head Details</Title>
        <Head headData={headData} />
      </Grid>
    </Grid>
  )
}

const adminHeadSelector = (
  headId: number | undefined,
  state: AppState,
): HeadShowData => {
  return build(state, 'adminHeadsShow', headId)
}

const mapStateToProps: MapStateToProps<StateProps, OwnProps, AppState> = (
  state,
  ownProps,
) => {
  return {
    headData: adminHeadSelector(ownProps.headId, state),
  }
}

const mapDispatchToProps: MapDispatchToProps<DispatchProps, OwnProps> = {
  fetchAdminHead,
}

const ConnectedShow = connect(mapStateToProps, mapDispatchToProps)(Show)
const StyledShow = withStyles(styles)(ConnectedShow)

export default StyledShow
