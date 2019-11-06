/**
 * The type of any function that is deployable
 */
type Deployable = {
  deploy: (...deployArgs: any[]) => Promise<any>
}

/**
 * Get the return type of a function, and unbox any promises
 */
export type Instance<T extends Deployable> = T extends {
  deploy: (...deployArgs: any[]) => Promise<infer U>
}
  ? U
  : never
