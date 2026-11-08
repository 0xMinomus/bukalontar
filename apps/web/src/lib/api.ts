import { mockAnalysis, type AnalysisResponse } from './analysis'
import { sleep } from './utils'

export interface AnalyzeInput {
  file: File
  enhanced: Blob | null
  sampleId?: string
  signal?: AbortSignal
  onStage?: (stage: number) => void
}

export async function analyzeManuscript({
  file,
  enhanced,
  sampleId,
  signal,
  onStage,
}: AnalyzeInput): Promise<AnalysisResponse> {
  const useMock = import.meta.env.VITE_USE_MOCK === 'true'
  const baseUrl = import.meta.env.VITE_API_BASE_URL ?? ''

  if (useMock) {
    onStage?.(1)
    await sleep(650, signal)
    onStage?.(2)
    await sleep(900, signal)
    onStage?.(3)
    await sleep(850, signal)
    onStage?.(4)
    await sleep(350, signal)
    return mockAnalysis(sampleId)
  }

  const form = new FormData()
  form.append('image', enhanced ?? file, 'lontar.jpg')
  form.append('source_id', sampleId ?? 'upload')

  onStage?.(2)
  const stageTimer = setTimeout(() => onStage?.(3), 2200)
  try {
    const response = await fetch(`${baseUrl}/api/analyze-manuscript`, {
      method: 'POST',
      body: form,
      signal,
    })

    if (!response.ok) {
      const detail = await response.text().catch(() => '')
      throw new Error(
        detail || `Analisis gagal dengan status ${response.status}.`,
      )
    }

    onStage?.(4)
    return (await response.json()) as AnalysisResponse
  } finally {
    clearTimeout(stageTimer)
  }
}
