import { generateUuid } from '../test-helpers/generateUuid'
import * as models from 'core/store/models'

export interface OcrKeyBundle extends models.OcrKey {
  id?: string
}

export const jsonApiOcrKeys = (keys: OcrKeyBundle[]) => {
  const k = keys || []

  return {
    data: k.map((c) => {
      const config = c || {}
      const id = config.id || generateUuid()
      const configPublicKey = config.configPublicKey || generateUuid()
      const offChainPublicKey = config.offChainPublicKey || generateUuid()
      const onChainSigningAddress =
        config.onChainSigningAddress || generateUuid()

      return {
        id,
        type: 'encryptedKeyBundles',
        attributes: {
          configPublicKey,
          offChainPublicKey,
          onChainSigningAddress,
          createdAt: new Date().toISOString(),
          UpdatedAt: new Date().toISOString(),
        },
      }
    }),
  }
}
