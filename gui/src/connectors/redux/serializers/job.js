export default json => ({
  id: json.id,
  createdAt: Date.parse(json.attributes.createdAt),
  initiators: json.attributes.initiators
})
