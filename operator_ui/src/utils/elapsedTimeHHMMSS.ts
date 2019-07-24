export default (createdAt: string, finishedAt: string): string => {
    if (!createdAt && !finishedAt) return ''

    const elapsedSecs = +new Date(finishedAt) / 1000 - +new Date(createdAt) / 1000

    const hours = Math.floor(elapsedSecs / 3600)
    const minutes = Math.floor((elapsedSecs % 3600) / 60)
    const seconds = Math.ceil((elapsedSecs % 3600) % 60)
    const concat = `${!!hours ? `${hours}h` : ''}${!!minutes ? `${minutes}m` : ``}${seconds}s`

    return concat
}
