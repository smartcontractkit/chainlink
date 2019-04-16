export default state => {
  const bridgeIds = state.bridges.currentPage

  return bridgeIds.map(id => state.bridges.items[id]).filter(r => r)
}
