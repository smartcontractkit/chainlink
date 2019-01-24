import React from 'react'
import { connect } from 'react-redux'
import { Link } from 'react-router-dom'
import Grid from '@material-ui/core/Grid'
import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import Divider from '@material-ui/core/Divider'
import Typography from '@material-ui/core/Typography'
import Form from 'components/Jobs/Form'
import ErrorMessage from 'components/Notifications/DefaultError'
import Content from 'components/Content'
import { createJobSpec } from 'actions'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'

const SuccessNotification = ({ data }) => (
  <React.Fragment>
    Successfully created job <Link id='created-job' to={`/jobs/${data.id}`}>{data.id}</Link>
  </React.Fragment>
)

const New = props => (
  <Content>
    <Grid container spacing={40}>
      <Grid item xs={12} md={11} lg={9}>
        <Card>
          <CardContent>
            <Typography variant='h5' color='secondary'>
              New Job
            </Typography>
          </CardContent>

          <Divider />

          <CardContent>
            <Form
              actionText='Create Job'
              onSubmit={props.createJobSpec}
              onSuccess={SuccessNotification}
              onError={ErrorMessage}
              {...(props.location && props.location.state)}
            />
          </CardContent>
        </Card>
      </Grid>
    </Grid>
  </Content>
)

export const ConnectedNew = connect(
  null,
  matchRouteAndMapDispatchToProps({ createJobSpec })
)(New)

export default ConnectedNew
