export function parseQRCode(text: string): [smdp: string, matchingId: string] | undefined {
  const segments = text.split('$')
  if (segments.length != 3) return
  if (segments[0].trim() !== 'LPA:1') return
  return [segments[1].trim(), segments[2].trim()]
}
