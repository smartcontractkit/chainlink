import React from "react";
import Button from '@material-ui/core/Button';
import Typography from '@material-ui/core/Typography'
import { withStyles } from '@material-ui/core/styles';
import { CopyToClipboard } from "react-copy-to-clipboard";

const styles = theme => ({
  button: {
    margin: theme.spacing.unit,
  },
  inform: {
	display: "inline-block",
    opacity: "1",
  }
});
    

class Copy extends React.Component {
  state = {
    copied: false
  };

  render() {
    return (
      <div>
        <CopyToClipboard text={this.props.JobSpec} onCopy={() => this.setState({ copied: true })}>
      	<Button variant="contained" className={this.props.classes.button}>
			Copy JobSpec
	  	</Button>
        </CopyToClipboard>
        {
		this.state.copied
				&& 
		<Typography className={this.props.classes.inform} color="primary"> 
			Copied 
		</Typography>
		}
	  </div> 
    );
  }
}

export default withStyles(styles)(Copy);
