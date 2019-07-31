import React from 'react';
import { Theme, WithStyles } from '@material-ui/core/styles';
declare const styles: ({ palette, shadows, typography }: Theme) => Record<"lightTooltip", import("@material-ui/core/styles/withStyles").CSSProperties>;
interface IProps extends WithStyles<typeof styles> {
    children: React.ReactElement<any>;
    title: string;
}
declare const _default: React.ComponentType<Pick<IProps, "children" | "title"> & import("@material-ui/core/styles").StyledComponentProps<"lightTooltip">>;
export default _default;
