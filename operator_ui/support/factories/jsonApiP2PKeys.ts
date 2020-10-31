import { generateUuid } from '../test-helpers/generateUuid'
import * as models from 'core/store/models'

export interface P2PKeyBundle extends models.P2PKey {
  id?: string
}

export const jsonApiP2PKeys = (keys: P2PKeyBundle[]) => {
  const k = keys || []

  return {
    data: k.map((c) => {
      const config = c || {}
      const id = config.id || generateUuid()
      const peerId = config.peerId || generateUuid()
      const publicKey = config.publicKey || generateUuid()

      return {
        id,
        type: 'encryptedKeyBundles',
        attributes: {
          peerId,
          publicKey,
          createdAt: new Date().toISOString(),
          UpdatedAt: new Date().toISOString(),
        },
      }
    }),
  }
}
