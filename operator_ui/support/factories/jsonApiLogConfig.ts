import { ApiResponse } from 'utils/json-api-client'
import { LogConfig } from 'core/store/models'

export function logConfigFactory(config: Partial<LogConfig>) {
  return {
    data: {
      id: 'log',
      type: 'logs',
      attributes: {
        level: config.level || 'debug',
        sqlEnabled: config.sqlEnabled || false,
      },
    },
  } as ApiResponse<LogConfig>
}
