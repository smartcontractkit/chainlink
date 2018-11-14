import React from 'react'
import { connect } from 'react-redux'
import { Link } from 'react-static'
import Grid from '@material-ui/core/Grid'
import PaddedCard from 'components/PaddedCard'
import Title from 'components/Title'
import Breadcrumb from 'components/Breadcrumb'
import BreadcrumbItem from 'components/BreadcrumbItem'
import Form from 'components/Bridges/Form'
import ErrorMessage from 'components/Notifications/DefaultError'
import Content from 'components/Content'
import { createBridge } from 'actions'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'

const SuccessNotification = ({name}) => (<React.Fragment>
  Successfully created bridge <Link to={`/bridges/${name}`}>{name}</Link>
</React.Fragment>)

const New = props => (
  <Content>
    <Breadcrumb>
      <BreadcrumbItem href='/'>Dashboard</BreadcrumbItem>
      <BreadcrumbItem>></BreadcrumbItem>
      <BreadcrumbItem href='/bridges'>Bridges</BreadcrumbItem>
      <BreadcrumbItem>></BreadcrumbItem>
      <BreadcrumbItem>New</BreadcrumbItem>
    </Breadcrumb>
    <Title>New Bridge</Title>

    <Grid container spacing={40}>
      <Grid item xs={12}>
        <PaddedCard>
          <Form
            actionText='Create Bridge'
            onSubmit={props.createBridge}
            onSuccess={SuccessNotification}
            onError={ErrorMessage}
          />
        </PaddedCard>
      </Grid>
    </Grid>
  </Content>
)

export const ConnectedNew = connect(
  null,
  matchRouteAndMapDispatchToProps({createBridge})
)(New)

export default ConnectedNew
