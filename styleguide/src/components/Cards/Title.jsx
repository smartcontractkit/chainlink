import React from 'react';
import CardContent from '@material-ui/core/CardContent';
import Divider from '@material-ui/core/Divider';
import Typography from '@material-ui/core/Typography';
const Title = ({ children, divider = false }) => {
    return (<React.Fragment>
      <CardContent>
        <Typography variant="h5" color="secondary">
          {children}
        </Typography>
      </CardContent>

      {divider && <Divider />}
    </React.Fragment>);
};
export default Title;
//# sourceMappingURL=Title.jsx.map