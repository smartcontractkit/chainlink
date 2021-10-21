import { generateUuid } from '../test-helpers/generateUuid'
import * as models from 'core/store/models'

export interface CSAKey extends models.CSAKey {
  id?: string
}

export const jsonApiCSAKeys = (keys: CSAKey[]) => {
  const k = keys || []

  return {
    data: k.map((c) => {
      const config = c || {}
      const id = config.id || '1'
      const publicKey = config.publicKey || generateUuid()

      return {
        id,
        type: 'csaKeys',
        attributes: {
          publicKey,
        },
      }
    }),
  }
}
