import React, { Component } from 'react'
import { withStyles } from '@material-ui/core/styles'
import Grid from '@material-ui/core/Grid'
import Typography from '@material-ui/core/Typography'
import classNames from 'classnames'
import { withRouteData } from 'react-static'

const styles = theme => ({
  style: {
    backgroundColor: "#F8F8F8",
    borderTop: "1px solid #E7E7E7",
    textAlign: "center",
    padding: "20px",
    position: "fixed",
    left: "0",
    bottom: "0",
    height: "60px",
    width: "100%",
  },
  wrapper: {
    display: 'block',
    padding: '20px',
    height: '60px',
    width: '100%',
  }
})

const Footnote = ({classes, version, sha}) => {
  return (
      <div className={classes.wrapper}>
      <div className={classes.style}>
            <Typography
              className={classNames(classes.footerText, classes.footerSections)}
            >
            Chainlink Node {version} at commit {sha}
            </Typography>
      </div>
      </div>
  )
}

export default withRouteData(withStyles(styles)(Footnote))