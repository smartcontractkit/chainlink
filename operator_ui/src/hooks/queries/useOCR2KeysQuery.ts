import { gql, QueryHookOptions, useQuery } from '@apollo/client'

export const OCR2_KEY_BUNDLES_PAYLOAD__RESULTS_FIELDS = gql`
  fragment OCR2KeyBundlesPayload_ResultsFields on OCR2KeyBundle {
    id
    chainType
    configPublicKey
    onChainPublicKey
    offChainPublicKey
  }
`

export const OCR2_KEY_BUNDLES_QUERY = gql`
  ${OCR2_KEY_BUNDLES_PAYLOAD__RESULTS_FIELDS}
  query FetchOCR2KeyBundles {
    ocr2KeyBundles {
      results {
        ...OCR2KeyBundlesPayload_ResultsFields
      }
    }
  }
`

// useOCRKeysQuery fetches the chains
export const useOCR2KeysQuery = (opts: QueryHookOptions = {}) => {
  return useQuery<FetchOcr2KeyBundles>(OCR2_KEY_BUNDLES_QUERY, opts)
}
