import React from 'react';
import { createStyles, withStyles } from '@material-ui/core/styles';
import MuiTooltip from '@material-ui/core/Tooltip';
const styles = ({ palette, shadows, typography }) => createStyles({
    lightTooltip: Object.assign({ background: palette.primary.contrastText, color: palette.text.primary, boxShadow: shadows[24] }, typography.h6)
});
const Tooltip = ({ title, children, classes }) => {
    return (<MuiTooltip title={title} classes={{ tooltip: classes.lightTooltip }}>
      {children}
    </MuiTooltip>);
};
export default withStyles(styles)(Tooltip);
//# sourceMappingURL=Tooltip.jsx.map