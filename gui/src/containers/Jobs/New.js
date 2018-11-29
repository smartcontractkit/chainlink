import React from 'react'
import { connect } from 'react-redux'
import { Link } from 'react-static'
import Grid from '@material-ui/core/Grid'
import Breadcrumb from 'components/Breadcrumb'
import BreadcrumbItem from 'components/BreadcrumbItem'
import Title from 'components/Title'
import PaddedCard from 'components/PaddedCard'
import Form from 'components/Jobs/Form'
import ErrorMessage from 'components/Notifications/DefaultError'
import Content from 'components/Content'
import { createJobSpec } from 'actions'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'

const SuccessNotification = ({data}) => (
  <React.Fragment>
    Successfully created job <Link id='created-job' to={`/jobs/${data.id}`}>{data.id}</Link>
  </React.Fragment>
)

const New = props => (
  <Content>
    <Breadcrumb>
      <BreadcrumbItem href='/'>Dashboard</BreadcrumbItem>
      <BreadcrumbItem>></BreadcrumbItem>
      <BreadcrumbItem href='/jobs'>Jobs</BreadcrumbItem>
      <BreadcrumbItem>></BreadcrumbItem>
      <BreadcrumbItem>New</BreadcrumbItem>
    </Breadcrumb>
    <Title>New Job</Title>

    <Grid container spacing={40}>
      <Grid item xs={12}>
        <PaddedCard>
          <Form
            actionText='Create Job'
            onSubmit={props.createJobSpec}
            onSuccess={SuccessNotification}
            onError={ErrorMessage}
            {...(props.location && props.location.state)}
          />
        </PaddedCard>
      </Grid>
    </Grid>
  </Content>
)

export const ConnectedNew = connect(
  null,
  matchRouteAndMapDispatchToProps({createJobSpec})
)(New)

export default ConnectedNew
