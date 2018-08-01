export const get = key => {
  const localStorageItem = global.localStorage.getItem(`chainlink.${key}`)
  let obj = {}

  if (localStorageItem) {
    try { return JSON.parse(localStorageItem) } catch (e) {}
  }

  return obj
}

export const set = (key, obj) => {
  global.localStorage.setItem(`chainlink.${key}`, JSON.stringify(obj))
}
