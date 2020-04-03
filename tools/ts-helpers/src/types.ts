export type DispatchBinding<T extends (...args: any[]) => any> = (
  ...args: Parameters<T>
) => void
