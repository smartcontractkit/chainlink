// insight from - https://github.com/sindresorhus/titleize
export default (input) => {
  const normalized = input || ''
  return normalized
    .toLowerCase()
    .replace(/_/g, ' ')
    .replace(/(?:^|\s|-)\S/g, (x) => x.toUpperCase())
}
