import { gql, QueryHookOptions, useQuery } from '@apollo/client'

export const OCR_KEY_BUNDLES_PAYLOAD__RESULTS_FIELDS = gql`
  fragment OCRKeyBundlesPayload_ResultsFields on OCRKeyBundle {
    id
    configPublicKey
    offChainPublicKey
    onChainSigningAddress
  }
`

export const OCR_KEY_BUNDLES_QUERY = gql`
  ${OCR_KEY_BUNDLES_PAYLOAD__RESULTS_FIELDS}
  query FetchOCRKeyBundles {
    ocrKeyBundles {
      results {
        ...OCRKeyBundlesPayload_ResultsFields
      }
    }
  }
`

// useOCRKeysQuery fetches the chains
export const useOCRKeysQuery = (opts: QueryHookOptions = {}) => {
  return useQuery<FetchOcrKeyBundles>(OCR_KEY_BUNDLES_QUERY, opts)
}
