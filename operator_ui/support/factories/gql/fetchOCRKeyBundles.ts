// buildOCRKeyBundle builds a ocr key bundle for the FetchOCRKeyBundles query.
export function buildOCRKeyBundle(
  overrides?: Partial<OcrKeyBundlesPayload_ResultsFields>,
): OcrKeyBundlesPayload_ResultsFields {
  return {
    __typename: 'OCRKeyBundle',
    id: '0cb6a1cb6c83bf0a94ba4b18ddc3385c7a52823795fc4a2424a42109258b7641',
    configPublicKey:
      'ocrcfg_2d19a30110a75886162ca1b0fc2d812377f917c386f2bdd746f89100ba79b364',
    onChainSigningAddress: 'ocrsad_0x04Ab7a825E09ab3bcf7697e8aAaB4A565Be9D2b5',
    offChainPublicKey:
      'ocroff_d94746863f04e638e6e0328433cf8e874fa1d26cf1662e7d14b723e983e0465e',
    ...overrides,
  }
}

// buildOCRKeyBundles builds a list of ocr key bundles.
export function buildOCRKeyBundles(): ReadonlyArray<OcrKeyBundlesPayload_ResultsFields> {
  return [
    buildOCRKeyBundle(),
    buildOCRKeyBundle({
      id: '4fadc92ce0b3deff6b2e2ef49cfc26cc39f8500818ad6591fb68f6c6ad0bb0dc',
      configPublicKey:
        'ocrcfg_e990a63c766c2f9d3ca33500b52a356afd16ca0cf9f3eeefe7c88026daf11d78',
      onChainSigningAddress:
        'ocrsad_0x2eB9410b954cbB18A83653B73cEaBa123dB19E9D',
      offChainPublicKey:
        'ocroff_f8014d2c2cc7730285d3fc7f2f0e5d8b3189675dfbd5dff68eb4a057f50afaf8',
    }),
  ]
}
