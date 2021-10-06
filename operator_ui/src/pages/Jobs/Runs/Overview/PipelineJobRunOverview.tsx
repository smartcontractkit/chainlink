import React from 'react'
import { Card, Divider, Typography } from '@material-ui/core'
import StatusIcon from 'components/StatusIcon'
import { theme } from 'theme'
import { PipelineJobRun } from '../../sharedTypes'
import { augmentOcrTasksList } from '../augmentOcrTasksList'

export const PipelineJobRunOverview = ({
  jobRun,
}: {
  jobRun: PipelineJobRun
}) => (
  <Card>
    {augmentOcrTasksList({ jobRun }).map((node, index) => {
      const { error, status, type, output, ...customAttributes } =
        node.attributes
      return (
        <React.Fragment key={node.id}>
          {index > 0 && <Divider />}
          <div
            style={{
              display: 'flex',
              padding: theme.spacing.unit * 2,
            }}
          >
            <span
              style={{
                paddingRight: theme.spacing.unit * 2,
              }}
            >
              <StatusIcon
                width={theme.spacing.unit * 5}
                height={theme.spacing.unit * 5}
              >
                {status}
              </StatusIcon>
            </span>
            <span>
              <Typography
                style={{
                  lineHeight: `${theme.spacing.unit * 5}px`,
                }}
                variant="h5"
              >
                {node.id}{' '}
                <small
                  style={{
                    color: theme.palette.grey[500],
                  }}
                >
                  {type}
                </small>
              </Typography>

              {['completed', 'errored'].includes(status) && (
                <Typography
                  style={{
                    marginBottom: theme.spacing.unit,
                    marginTop: theme.spacing.unit,
                  }}
                  variant="body1"
                >
                  {status === 'completed' ? output : error}
                </Typography>
              )}

              {status !== 'not_run' &&
                Object.entries(customAttributes).map(([key, value]) => (
                  <Typography key={key} variant="body1">
                    <span
                      style={{
                        fontWeight: theme.typography.fontWeightLight,
                      }}
                    >
                      {key}
                    </span>
                    : {value}
                  </Typography>
                ))}
            </span>
          </div>
        </React.Fragment>
      )
    })}
  </Card>
)
