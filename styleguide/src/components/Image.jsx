import React from 'react';
import { createStyles, withStyles } from '@material-ui/core/styles';
const styles = (theme) => createStyles({
    animate: {
        animation: 'spin 4s linear infinite'
    },
    '@keyframes spin': {
        '100%': {
            transform: 'rotate(360deg)'
        }
    }
});
const Image = ({ src, width, height, alt, classes, spin = false }) => {
    return (<img src={src} className={spin ? classes.animate : ''} alt={alt} width={width} height={height}/>);
};
export default withStyles(styles)(Image);
//# sourceMappingURL=Image.jsx.map