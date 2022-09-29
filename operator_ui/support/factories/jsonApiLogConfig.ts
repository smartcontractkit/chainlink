import { ApiResponse } from 'utils/json-api-client'
import { LogConfig } from 'core/store/models'

export function logConfigFactory(config: Partial<LogConfig>) {
  return {
    data: {
      id: 'log',
      type: 'logs',
      attributes: {
        serviceName: config.serviceName || ['Global', 'IsSqlEnabled'],
        logLevel: config.logLevel || ['info', 'true'],
        defaultLogLevel: config.defaultLogLevel || 'info',
      },
    },
  } as ApiResponse<LogConfig>
}
