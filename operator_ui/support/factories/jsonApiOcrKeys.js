import uuid from 'uuid/v4'

export const jsonApiOcrKeys = (keys) => {
  const k = keys || []

  return {
    data: k.map((c) => {
      const config = c || {}
      const id = config.id || uuid().replace(/-/g, '')
      const ConfigPublicKey = config.ConfigPublicKey || uuid().replace(/-/g, '')
      const OffChainPublicKey =
        config.OffChainPublicKey || uuid().replace(/-/g, '')
      const OnChainSigningAddress =
        config.OnChainSigningAddress || uuid().replace(/-/g, '')

      return {
        id,
        type: 'encryptedKeyBundles',
        attributes: {
          ConfigPublicKey,
          OffChainPublicKey,
          OnChainSigningAddress,
          CreatedAt: new Date().toISOString(),
          UpdatedAt: new Date().toISOString(),
        },
      }
    }),
  }
}
