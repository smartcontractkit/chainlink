import React from 'react';
import TimeAgoNoTooltip from 'react-time-ago/no-tooltip';
import Tooltip from './Tooltip';
import localizedTimestamp from '../utils/localizedTimestamp';
const TimeAgo = ({ children, tooltip = false }) => {
    const date = Date.parse(children);
    const ago = <TimeAgoNoTooltip tooltip={false}>{date}</TimeAgoNoTooltip>;
    if (tooltip) {
        return (<Tooltip title={localizedTimestamp(children)}>
        <span>{ago}</span>
      </Tooltip>);
    }
    return ago;
};
export default TimeAgo;
//# sourceMappingURL=TimeAgo.jsx.map