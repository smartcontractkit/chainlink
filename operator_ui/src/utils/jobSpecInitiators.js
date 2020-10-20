export const isWebInitiator = (initiators) =>
  initiators.find((initiator) => initiator.type === 'web')

export const formatInitiators = (initiators) =>
  (initiators || []).map((i) => i.type).join(', ')
