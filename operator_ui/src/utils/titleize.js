// insight from - https://github.com/sindresorhus/titleize
export default (input) => {
  if (typeof input !== 'string') {
    return input
  }

  return input
    .toLowerCase()
    .replace(/_/g, ' ')
    .replace(/(?:^|\s|-)\S/g, (x) => x.toUpperCase())
}
