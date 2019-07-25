import React from 'react';
import { Theme, WithStyles } from '@material-ui/core/styles';
declare const styles: (theme: Theme) => Record<"animate" | "@keyframes spin", import("@material-ui/core/styles/withStyles").CSSProperties>;
interface IProps extends WithStyles<typeof styles> {
    src: string;
    width?: number;
    height?: number;
    spin?: boolean;
    alt?: string;
}
declare const _default: React.ComponentType<Pick<IProps, "height" | "width" | "src" | "alt" | "spin"> & import("@material-ui/core/styles").StyledComponentProps<"animate" | "@keyframes spin">>;
export default _default;
