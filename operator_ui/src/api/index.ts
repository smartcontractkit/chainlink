import {
  AttributesObject,
  ErrorsObject,
  JsonApiResponse,
  ResourceObject
} from 'json-api-normalizer'

export function TODO(..._: any[]): any {
  throw Error('not implemented')
}

export type Response<T extends AttributesObject> = Promise<
  JsonApiResponse<
    T extends Array<infer U> ? ResourceObject<U>[] : ResourceObject<T>
  >
>

export type PaginatedResponse<T extends AttributesObject> = Promise<
  JsonApiResponse<
    T extends Array<infer U> ? ResourceObject<U>[] : ResourceObject<T>,
    ErrorsObject[],
    never,
    {
      count: number
    },
    Partial<{
      prev: string
      next: string
    }>
  >
>
