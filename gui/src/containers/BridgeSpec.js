import React, { Component } from "react";
import PropTypes from "prop-types";
import Typography from "@material-ui/core/Typography";
import Grid from "@material-ui/core/Grid";
import PaddedCard from "components/PaddedCard";
import Breadcrumb from "components/Breadcrumb";
import BreadcrumbItem from "components/BreadcrumbItem";
import Card from "@material-ui/core/Card";
import { withStyles } from "@material-ui/core/styles";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import { fetchBridgeSpec } from "actions";

const styles = theme => ({
  title: {
    marginTop: theme.spacing.unit * 5,
    marginBottom: theme.spacing.unit * 5
  },
  definitionTitle: {
    marginBottom: theme.spacing.unit * 3
  },
  breadcrumb: {
    marginTop: theme.spacing.unit * 5,
    marginBottom: theme.spacing.unit * 5
  },
  lastRun: {
    marginTop: theme.spacing.unit * 5,
    marginBottom: theme.spacing.unit * 5
  },
  showMore: {
    marginTop: theme.spacing.unit * 3,
    marginLeft: theme.spacing.unit * 3,
    display: "block"
  }
});

const renderBridgeSpec = ({ classes, name, url, confirmations }) => (
  <Grid container spacing={40}>
    <Grid item xs={8}>
      <PaddedCard>
        <Typography variant="title" className={classes.definitionTitle}>
          Bridge Info:
        </Typography>
        Name: {name} <br/>
        URL: {url} <br/>
        Confirmations: {confirmations} <br/>
      </PaddedCard>
    </Grid>
  </Grid>
);

const renderBridgeAssociations = ({ classes }) => (
  <React.Fragment>
    <Typography variant="title" className={classes.lastRun}>
      Associations with:
    </Typography>
    <Card>WIP</Card>
  </React.Fragment>
);

const renderFetching = () => <div>Fetching...</div>;

const renderDetails = props => {
  if (!props.fetching) {
    return (
      <React.Fragment>
        {renderBridgeSpec(props)}
        {renderBridgeAssociations(props)}
      </React.Fragment>
    );
  } else {
    return renderFetching();
  }
};

export class BridgeSpec extends Component {
  componentDidMount() {
    this.props.fetchBridgeSpec(this.props.match.params.bridgeName)
  }

  render() {
    const { classes, name } = this.props;
    return (
      <div>
        <Breadcrumb className={classes.breadcrumb}>
          <BreadcrumbItem href="/">Dashboard</BreadcrumbItem>
          <BreadcrumbItem>></BreadcrumbItem>
          <BreadcrumbItem href="/bridges">Bridges</BreadcrumbItem>
          <BreadcrumbItem>></BreadcrumbItem>
          <BreadcrumbItem>{name}</BreadcrumbItem>
        </Breadcrumb>
        <Typography variant="display2" color="inherit" className={classes.title}>
          Bridge Spec Details:
        </Typography>
        {renderDetails(this.props)}
      </div>
    );
  }
}

BridgeSpec.propTypes = {
  classes: PropTypes.object.isRequired,
};

const mapStateToProps = (state) => {
  let BridgeSpecError
  if (state.networkError) {
    BridgeSpecError = 'error fetching balance'
  }
  return {
    name: state.bridgeSpec.name,
    url: state.bridgeSpec.url,
    confirmations: state.bridgeSpec.confirmations,
    BridgeSpecError: BridgeSpecError
  }
};

const mapDispatchToProps = dispatch => {
  return bindActionCreators({ fetchBridgeSpec }, dispatch);
};

export const ConnectedBridgeSpec = connect(mapStateToProps, mapDispatchToProps)(BridgeSpec);

export default withStyles(styles)(ConnectedBridgeSpec);
