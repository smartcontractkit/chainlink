declare module 'feeds' {
  interface FeedConfig {
    contractAddress: string
    name: string
    valuePrefix: string
    pair: string[]
    counter?: number
    path: string
    contractVersion?: number
    networkId: number
    history: boolean
    decimalPlaces?: number
    multiply?: string
    sponsored?: string[]
    threshold: number
  }
}
