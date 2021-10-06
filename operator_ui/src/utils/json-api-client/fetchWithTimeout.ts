import 'isomorphic-unfetch'

export function fetchWithTimeout(
  url: string,
  options: Parameters<typeof fetch>[1],
  timeout = 20000,
): Promise<Response> {
  return Promise.race([
    fetch(url, options),
    new Promise((_, reject) =>
      setTimeout(() => reject(new Error('timeout')), timeout),
    ) as any as Response,
  ])
}
