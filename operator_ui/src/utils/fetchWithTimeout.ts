import 'isomorphic-unfetch'

export default function (url: string, options: any, timeout: number = 20000): Promise<any> {
  return Promise.race([
    fetch(url, options),
    new Promise((_, reject) =>
      setTimeout(() => reject(new Error('timeout')), timeout)
    )
  ])
}
