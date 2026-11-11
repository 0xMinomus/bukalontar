import type { AnalysisResponse, Recipe } from './analysis'

const NO_PLANT = /^(tidak disebutkan|tidak ada|tidak terdeteksi|bukan naskah usada|tidak ditemukan)/i
const EMPTY_VALUE =
  /^(tidak disebutkan|tidak ada|tidak terdeteksi|belum tersedia|belum ada)\.?$/i

export function resultRecipes(result: AnalysisResponse): Recipe[] {
  if (result.recipes && result.recipes.length > 0) {
    return result.recipes
  }
  const entity = result.ethnomedicinal_entity
  const validation = result.scientific_validation
  return [
    {
      plant_local_name: entity.plant_local_name,
      plant_scientific_name: entity.plant_scientific_name,
      botanical_family: entity.botanical_family,
      kb_match: entity.kb_match,
      plant_part_used: entity.plant_part_used,
      energy_character: entity.energy_character,
      formulation_type: entity.formulation_type,
      adjuvants: entity.adjuvants,
      target_ailment: entity.target_ailment,
      preparation_instructions: entity.preparation_instructions,
      phytochemical_compounds: validation.phytochemical_compounds,
      modern_pharmacological_action: validation.modern_pharmacological_action,
      reference_studies: validation.reference_studies,
      reference_links: validation.reference_links,
    },
  ]
}

export function hasBalinese(value: string | null | undefined): boolean {
  return !!value && /[\u1B00-\u1B7F]/.test(value)
}

export function safeReferenceLinks(recipe: Recipe) {
  if (!recipe.kb_match) return []
  return (recipe.reference_links ?? []).filter((link) => {
    try {
      const url = new URL(link.url)
      return (url.protocol === 'https:' || url.protocol === 'http:') && !!url.hostname
    } catch {
      return false
    }
  })
}

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

export function analysisSummary(result: AnalysisResponse): string | null {
  const entity = result.ethnomedicinal_entity
  const source = usefulText(entity.preparation_instructions)
    ? entity.preparation_instructions
    : result.paleography.indonesian_translation
  return usefulText(source) ? formatTranslation(source) : null
}
