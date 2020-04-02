declare module 'docker-registry-client' {
  export function createClientV2({ name: string }): RegistryClientV2

  class RegistryClientV2 {
    listTags(cb: (err: Error, tags: Tags) => void)
    close(): void
  }

  export interface Tags {
    name: string
    tags: string[]
  }
}
