/**
 * Parse the current query string out of the browser location
 *
 * @param location The location value to use, not hardcoded so we can inject
 * mock values for testing
 */
export function searchQuery(location: Location = window.location): string {
  const searchParams = new URLSearchParams(location.search)
  const search = searchParams.get('search')

  return search ?? ''
}
