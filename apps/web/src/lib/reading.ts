const NO_PLANT = /^(tidak disebutkan|tidak ada|tidak terdeteksi|bukan naskah usada|tidak ditemukan)/i
const EMPTY_VALUE =
  /^(tidak disebutkan|tidak ada|tidak terdeteksi|belum tersedia|belum ada)\.?$/i

export function hasMedicinalPlant(name: string): boolean {
  const value = name.trim()
  return value.length > 0 && !NO_PLANT.test(value)
}

export function usefulText(value: string | null | undefined): value is string {
  if (!value) return false
  const trimmed = value.trim()
  return trimmed.length > 0 && !EMPTY_VALUE.test(trimmed)
}

export function formatTransliteration(value: string): string {
  return value
    .replace(/\r\n?/g, '\n')
    .replace(/[ \t]*Line\s+\d+\s*:\s*/gi, '\n')
    .replace(/[ \t]*(\[Lembar\s+\d+\])/gi, '\n\n$1\n')
    .replace(/\n{3,}/g, '\n\n')
    .trim()
}

export function formatTranslation(value: string): string {
  const normalized = value.replace(/\r\n?/g, '\n').trim()
  if (!normalized) return ''
  if (/\n\s*\n/.test(normalized)) {
    return normalized.replace(/\n{3,}/g, '\n\n')
  }

  const sentences = normalized.match(/[^.!?]+[.!?]+|[^.!?]+$/g)
  if (!sentences || sentences.length < 4) return normalized

  const paragraphs: string[] = []
  for (let index = 0; index < sentences.length; index += 2) {
    paragraphs.push(sentences.slice(index, index + 2).join(' ').trim())
  }
  return paragraphs.join('\n\n')
}
