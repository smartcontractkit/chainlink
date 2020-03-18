declare module 'feeds' {
  interface FeedConfig {
    contractAddress: string
    contractVersion?: number
    contractType: string
    name: string
    valuePrefix: string
    pair: string[]
    counter?: number
    path: string
    networkId: number
    history: boolean
    decimalPlaces?: number
    multiply?: string
    sponsored?: string[]
    threshold: number
    compareOffchain?: string
    listing?: boolean
  }
}
