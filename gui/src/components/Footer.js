import React from "react";
import { withStyles } from "@material-ui/core/styles";
import Typography from "@material-ui/core/Typography";
import { withRouteData } from "react-static";
import Card from '@material-ui/core/Card'

const styles = theme => ({
  style: {
    textAlign: "center",
    padding: "20px",
    position: "fixed",
    left: "0",
    bottom: "0",
    height: "60px",
    width: "100%"
  }
});

const Footnote = ({ classes, version, sha }) => {
  return (
      <Card className={classes.style}>
        <Typography>
          Chainlink Node {version} at commit {sha}
        </Typography>
      </Card>
  );
};

export default withRouteData(withStyles(styles)(Footnote));
