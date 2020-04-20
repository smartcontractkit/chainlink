declare module 'feeds' {
  interface FeedConfig {
    contractAddress: string
    listing: boolean
    contractVersion?: number
    contractType: string
    name: string
    valuePrefix: string
    pair: string[]
    heartbeat?: number
    path: string
    networkId: number
    history: boolean
    decimalPlaces?: number
    multiply?: string
    sponsored?: string[]
    threshold: number
    compareOffchain?: string
  }
}
