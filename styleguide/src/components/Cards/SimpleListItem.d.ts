import React from 'react';
import { Theme, WithStyles } from '@material-ui/core/styles';
declare const styles: (theme: Theme) => Record<"cell", import("@material-ui/core/styles/withStyles").CSSProperties>;
interface IProps extends WithStyles<typeof styles> {
    children: React.ReactNode;
}
declare const _default: React.ComponentType<Pick<IProps, "children"> & import("@material-ui/core/styles").StyledComponentProps<"cell">>;
export default _default;
