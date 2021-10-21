import { ApiResponse } from 'utils/json-api-client'
import { FeatureFlag } from 'core/store/models'

export const jsonApiFeatureFlags = () => {
  return {
    data: [
      {
        id: 'csa',
        type: 'features',
        attributes: {
          enabled: false,
        },
      },
    ],
  } as ApiResponse<FeatureFlag[]>
}
