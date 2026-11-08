import { useCallback, useEffect, useRef, useState } from 'react'

import { AppHeader } from '@/components/app-header'
import { ErrorState } from '@/components/error-state'
import { HowItWorks } from '@/components/how-it-works'
import { Landing } from '@/components/landing'
import { type LensView } from '@/components/lontar-lens'
import { PreviewStage } from '@/components/preview-stage'
import { ProcessingStage } from '@/components/processing-stage'
import { ResultWorkspace } from '@/components/result-workspace'
import { SiteFooter } from '@/components/site-footer'
import { TooltipProvider } from '@/components/ui/tooltip'
import { analyzeManuscript } from '@/lib/api'
import type { AnalysisResponse } from '@/lib/analysis'
import { enhanceLontar, type EnhancedImage } from '@/lib/image'
import type { Sample } from '@/lib/samples'

type Stage = 'idle' | 'preview' | 'processing' | 'result' | 'error'

const MAX_BYTES = 8 * 1024 * 1024
const ACCEPTED_TYPES = ['image/jpeg', 'image/png', 'image/webp']
const IS_MOCK = import.meta.env.VITE_USE_MOCK !== 'false'

export default function App() {
  const [stage, setStage] = useState<Stage>('idle')
  const [file, setFile] = useState<File | null>(null)
  const [sourceId, setSourceId] = useState<string | undefined>()
  const [originalUrl, setOriginalUrl] = useState<string | null>(null)
  const [enhanced, setEnhanced] = useState<EnhancedImage | null>(null)
  const [enhancing, setEnhancing] = useState(false)
  const [enhanceError, setEnhanceError] = useState<string | null>(null)
  const [view, setView] = useState<LensView>('original')
  const [step, setStep] = useState(1)
  const [result, setResult] = useState<AnalysisResponse | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [busySampleId, setBusySampleId] = useState<string | null>(null)

  const loadIdRef = useRef(0)
  const abortRef = useRef<AbortController | null>(null)

  useEffect(() => {
    return () => {
      if (originalUrl) URL.revokeObjectURL(originalUrl)
    }
  }, [originalUrl])

  useEffect(() => {
    return () => {
      if (enhanced?.url) URL.revokeObjectURL(enhanced.url)
    }
  }, [enhanced?.url])

  useEffect(() => {
    if (stage !== 'idle') window.scrollTo({ top: 0 })
  }, [stage])

  const reset = useCallback(() => {
    loadIdRef.current += 1
    abortRef.current?.abort()
    abortRef.current = null
    setFile(null)
    setSourceId(undefined)
    setOriginalUrl(null)
    setEnhanced(null)
    setEnhancing(false)
    setEnhanceError(null)
    setView('original')
    setStep(1)
    setResult(null)
    setError(null)
    setBusySampleId(null)
    setStage('idle')
  }, [])

  const loadFile = useCallback(async (nextFile: File, sampleId?: string) => {
    if (!ACCEPTED_TYPES.includes(nextFile.type)) {
      setError('Format foto harus JPG, PNG, atau WEBP.')
      setStage('error')
      return
    }
    if (nextFile.size > MAX_BYTES) {
      setError('Ukuran foto melebihi 8 MB. Kompres atau pilih foto lain.')
      setStage('error')
      return
    }

    const loadId = loadIdRef.current + 1
    loadIdRef.current = loadId
    abortRef.current?.abort()
    abortRef.current = null

    setFile(nextFile)
    setSourceId(sampleId)
    setOriginalUrl(URL.createObjectURL(nextFile))
    setEnhanced(null)
    setEnhanceError(null)
    setEnhancing(true)
    setResult(null)
    setError(null)
    setStep(1)
    setView('original')
    setStage('preview')

    try {
      const nextEnhanced = await enhanceLontar(nextFile)
      if (loadIdRef.current !== loadId) {
        URL.revokeObjectURL(nextEnhanced.url)
        return
      }
      setEnhanced(nextEnhanced)
      setView('enhanced')
    } catch (cause) {
      if (loadIdRef.current === loadId) {
        setEnhanceError(
          cause instanceof Error
            ? cause.message
            : 'Peningkatan kontras gagal.',
        )
      }
    } finally {
      if (loadIdRef.current === loadId) setEnhancing(false)
    }
  }, [])

  async function handleSample(sample: Sample) {
    if (busySampleId) return
    setBusySampleId(sample.id)
    try {
      const response = await fetch(sample.src)
      if (!response.ok) throw new Error('Contoh arsip tidak bisa dimuat.')
      const blob = await response.blob()
      const nextFile = new File([blob], `${sample.id}.jpg`, {
        type: blob.type || 'image/jpeg',
      })
      await loadFile(nextFile, sample.id)
    } catch (cause) {
      setError(
        cause instanceof Error ? cause.message : 'Contoh arsip tidak bisa dimuat.',
      )
      setStage('error')
    } finally {
      setBusySampleId(null)
    }
  }

  async function startAnalysis() {
    if (!file) return
    abortRef.current?.abort()
    const controller = new AbortController()
    abortRef.current = controller

    setStage('processing')
    setStep(1)
    setError(null)

    try {
      const nextResult = await analyzeManuscript({
        file,
        enhanced: enhanced?.blob ?? null,
        sampleId: sourceId,
        signal: controller.signal,
        onStage: (nextStep) => {
          if (!controller.signal.aborted) setStep(nextStep)
        },
      })
      if (controller.signal.aborted) return
      setResult(nextResult)
      setStage('result')
    } catch (cause) {
      if (
        controller.signal.aborted ||
        (cause instanceof DOMException && cause.name === 'AbortError')
      ) {
        setStage('preview')
        return
      }
      setError(
        cause instanceof Error ? cause.message : 'Analisis gagal. Coba lagi.',
      )
      setStage('error')
    } finally {
      if (abortRef.current === controller) abortRef.current = null
    }
  }

  function cancelAnalysis() {
    abortRef.current?.abort()
    abortRef.current = null
    setStage('preview')
    setStep(1)
  }

  function downloadJson() {
    if (!result) return
    const blob = new Blob([JSON.stringify(result, null, 2)], {
      type: 'application/json',
    })
    const url = URL.createObjectURL(blob)
    const anchor = document.createElement('a')
    anchor.href = url
    anchor.download = `${result.manuscript_id}.json`
    anchor.click()
    URL.revokeObjectURL(url)
  }

  useEffect(() => {
    function onPaste(event: ClipboardEvent) {
      if (stage !== 'idle' && stage !== 'error') return
      const items = event.clipboardData?.items
      if (!items) return
      const imageItem = Array.from(items).find((item) =>
        item.type.startsWith('image/'),
      )
      const pasted = imageItem?.getAsFile()
      if (pasted) {
        event.preventDefault()
        void loadFile(pasted)
      }
    }
    window.addEventListener('paste', onPaste)
    return () => window.removeEventListener('paste', onPaste)
  }, [loadFile, stage])

  return (
    <TooltipProvider>
      <div className="flex min-h-[100dvh] flex-col">
        <div className="no-print">
          <AppHeader showNav={stage === 'idle'} />
        </div>

        <main className="flex-1">
          <div id="workspace">
            {stage === 'idle' ? (
              <Landing
                onFile={(nextFile) => void loadFile(nextFile)}
                onSample={handleSample}
                busySampleId={busySampleId}
              />
            ) : null}

            {stage === 'preview' && file && originalUrl ? (
              <PreviewStage
                fileName={file.name}
                fileSize={file.size}
                dimensions={
                  enhanced
                    ? { width: enhanced.width, height: enhanced.height }
                    : undefined
                }
                originalUrl={originalUrl}
                enhancedUrl={enhanced?.url}
                enhancing={enhancing}
                enhanceError={enhanceError}
                enhanceElapsedMs={enhanced?.elapsedMs}
                view={view}
                onViewChange={setView}
                onAnalyze={() => void startAnalysis()}
                onReset={reset}
              />
            ) : null}

            {stage === 'processing' && (originalUrl || enhanced?.url) ? (
              <ProcessingStage
                imageUrl={enhanced?.url ?? originalUrl ?? ''}
                step={step}
                onCancel={cancelAnalysis}
              />
            ) : null}

            {stage === 'result' && result && originalUrl ? (
              <ResultWorkspace
                result={result}
                originalUrl={originalUrl}
                enhancedUrl={enhanced?.url}
                view={view}
                onViewChange={setView}
                onReset={reset}
                onExport={downloadJson}
                onPrint={() => window.print()}
                isMock={IS_MOCK}
              />
            ) : null}

            {stage === 'error' ? (
              <ErrorState
                message={error ?? 'Terjadi kesalahan yang belum dikenali.'}
                canRetry={Boolean(originalUrl)}
                onRetry={() => setStage('preview')}
                onReset={reset}
              />
            ) : null}
          </div>

          {stage === 'idle' ? (
            <div className="no-print">
              <HowItWorks />
            </div>
          ) : null}
        </main>

        <div className="no-print">
          <SiteFooter />
        </div>
      </div>
    </TooltipProvider>
  )
}
