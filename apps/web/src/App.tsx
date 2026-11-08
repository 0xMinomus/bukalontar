import { Loader2 } from 'lucide-react'
import { useCallback, useEffect, useRef, useState } from 'react'
import {
  Navigate,
  Route,
  Routes,
  useLocation,
  useNavigate,
} from 'react-router-dom'

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
import {
  clearAnalysisSession,
  loadAnalysisSession,
  saveAnalysisSession,
} from '@/lib/session-store'

type Stage = 'idle' | 'preview' | 'processing' | 'result' | 'error'

const MAX_BYTES = 8 * 1024 * 1024
const ACCEPTED_TYPES = ['image/jpeg', 'image/png', 'image/webp']
const IS_MOCK = import.meta.env.VITE_USE_MOCK === 'true'

function SessionLoading() {
  return (
    <div className="mx-auto max-w-5xl px-5 py-24 text-center">
      <Loader2
        className="mx-auto size-6 animate-spin text-brass-strong"
        aria-hidden="true"
      />
      <p className="small-caps mt-4 text-xs font-semibold text-brass-strong">
        Memuat sesi
      </p>
      <p className="mt-2 font-heading text-3xl">Menyiapkan hasil</p>
    </div>
  )
}

export default function App() {
  const navigate = useNavigate()
  const location = useLocation()
  const isLanding = location.pathname === '/'

  const [hydrating, setHydrating] = useState(
    () => location.pathname === '/hasil' || location.pathname === '/analisis',
  )
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

  useEffect(() => {
    if (!hydrating) return
    let cancelled = false
    void (async () => {
      try {
        const session = await loadAnalysisSession()
        if (cancelled || !session) return
        const original = session.original
        setFile(
          new File([original], session.fileName, {
            type: original.type || 'image/jpeg',
          }),
        )
        setSourceId(session.sourceId)
        setOriginalUrl(URL.createObjectURL(original))
        if (session.enhanced) {
          setEnhanced({
            url: URL.createObjectURL(session.enhanced),
            blob: session.enhanced,
            width: session.enhancedMeta?.width ?? 0,
            height: session.enhancedMeta?.height ?? 0,
            elapsedMs: session.enhancedMeta?.elapsedMs ?? 0,
          })
          setView(session.view)
        } else {
          setView('original')
        }
        setResult(session.result)
        setStage('result')
      } catch {
        // Corrupt or unavailable store; the route guard sends the user home.
      } finally {
        if (!cancelled) setHydrating(false)
      }
    })()
    return () => {
      cancelled = true
    }
  }, [hydrating])

  const clearState = useCallback(() => {
    void clearAnalysisSession().catch(() => {})
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

  const reset = useCallback(() => {
    clearState()
    navigate('/')
  }, [clearState, navigate])

  const handleLogo = useCallback(() => {
    clearState()
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }, [clearState])

  const loadFile = useCallback(
    async (nextFile: File, sampleId?: string) => {
      if (!ACCEPTED_TYPES.includes(nextFile.type)) {
        setError('Format foto harus JPG, PNG, atau WEBP.')
        setStage('error')
        navigate('/analisis')
        return
      }
      if (nextFile.size > MAX_BYTES) {
        setError('Ukuran foto melebihi 8 MB. Kompres atau pilih foto lain.')
        setStage('error')
        navigate('/analisis')
        return
      }

      const loadId = loadIdRef.current + 1
      loadIdRef.current = loadId
      abortRef.current?.abort()
      abortRef.current = null
      void clearAnalysisSession().catch(() => {})

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
      navigate('/analisis')

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
    },
    [navigate],
  )

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
        cause instanceof Error
          ? cause.message
          : 'Contoh arsip tidak bisa dimuat.',
      )
      setStage('error')
      navigate('/analisis')
    } finally {
      setBusySampleId(null)
    }
  }

  async function startAnalysis() {
    if (!file) return
    const manuscriptFile = file
    const currentEnhanced = enhanced
    const currentView = view
    void clearAnalysisSession().catch(() => {})
    abortRef.current?.abort()
    const controller = new AbortController()
    abortRef.current = controller

    setStage('processing')
    setStep(1)
    setError(null)
    navigate('/analisis')

    try {
      const nextResult = await analyzeManuscript({
        file: manuscriptFile,
        enhanced: currentEnhanced?.blob ?? null,
        sampleId: sourceId,
        signal: controller.signal,
        onStage: (nextStep) => {
          if (!controller.signal.aborted) setStep(nextStep)
        },
      })
      if (controller.signal.aborted) return
      try {
        await saveAnalysisSession({
          result: nextResult,
          original: manuscriptFile,
          enhanced: currentEnhanced?.blob ?? null,
          enhancedMeta: currentEnhanced
            ? {
                width: currentEnhanced.width,
                height: currentEnhanced.height,
                elapsedMs: currentEnhanced.elapsedMs,
              }
            : undefined,
          fileName: manuscriptFile.name,
          sourceId,
          view: currentView,
          savedAt: Date.now(),
        })
      } catch {
        // Storage unavailable; the result remains usable without refresh.
      }
      if (controller.signal.aborted) return
      setResult(nextResult)
      setStage('result')
      navigate('/hasil')
    } catch (cause) {
      if (
        controller.signal.aborted ||
        (cause instanceof DOMException && cause.name === 'AbortError')
      ) {
        setStage('preview')
        navigate('/analisis')
        return
      }
      setError(
        cause instanceof Error ? cause.message : 'Analisis gagal. Coba lagi.',
      )
      setStage('error')
      navigate('/analisis')
    } finally {
      if (abortRef.current === controller) abortRef.current = null
    }
  }

  function cancelAnalysis() {
    abortRef.current?.abort()
    abortRef.current = null
    setStage('preview')
    setStep(1)
    navigate('/analisis')
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
        <div className="no-print app-header-shell">
          <AppHeader showNav={isLanding} onHome={handleLogo} />
        </div>

        <main className="flex-1">
          <Routes>
            <Route
              path="/"
              element={
                <Landing
                  onFile={(nextFile) => void loadFile(nextFile)}
                  onSample={handleSample}
                  busySampleId={busySampleId}
                  isMock={IS_MOCK}
                />
              }
            />

            <Route
              path="/analisis"
              element={
                hydrating ? (
                  <SessionLoading />
                ) : stage === 'error' ? (
                  <ErrorState
                    message={error ?? 'Terjadi kesalahan yang belum dikenali.'}
                    canRetry={Boolean(originalUrl)}
                    onRetry={() => {
                      setStage('preview')
                      navigate('/analisis')
                    }}
                    onReset={reset}
                  />
                ) : !file || !originalUrl ? (
                  <Navigate to="/" replace />
                ) : stage === 'processing' ? (
                  <ProcessingStage
                    imageUrl={enhanced?.url ?? originalUrl}
                    step={step}
                    onCancel={cancelAnalysis}
                  />
                ) : (
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
                )
              }
            />

            <Route
              path="/hasil"
              element={
                hydrating ? (
                  <SessionLoading />
                ) : result && originalUrl ? (
                  <ResultWorkspace
                    result={result}
                    originalUrl={originalUrl}
                    enhancedUrl={enhanced?.url}
                    view={view}
                    onViewChange={setView}
                    onReset={reset}
                    onPrint={() => window.print()}
                    isMock={IS_MOCK}
                  />
                ) : (
                  <Navigate to="/" replace />
                )
              }
            />

            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>

          {isLanding ? (
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
