import { BookOpen, Download, Info, Printer, RotateCcw } from 'lucide-react'

import { AnalysisTabs } from '@/components/analysis-tabs'
import { LontarLens, type LensView } from '@/components/lontar-lens'
import { Pipeline } from '@/components/pipeline'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import type { AnalysisResponse, Confidence } from '@/lib/analysis'

const confidenceLabel: Record<Confidence, string> = {
  low: 'Keyakinan rendah',
  medium: 'Keyakinan sedang',
  high: 'Keyakinan tinggi',
}

const confidenceVariant: Record<
  Confidence,
  'success' | 'brass' | 'destructive'
> = {
  low: 'destructive',
  medium: 'brass',
  high: 'success',
}

interface ResultWorkspaceProps {
  result: AnalysisResponse
  originalUrl: string
  enhancedUrl?: string | null
  view: LensView
  onViewChange: (view: LensView) => void
  onReset: () => void
  onExport: () => void
  onPrint: () => void
  isMock: boolean
}

export function ResultWorkspace({
  result,
  originalUrl,
  enhancedUrl,
  view,
  onViewChange,
  onReset,
  onExport,
  onPrint,
  isMock,
}: ResultWorkspaceProps) {
  const source = result.source_reference
  const plant = result.ethnomedicinal_entity
  const confidence = result.paleography.confidence

  return (
    <section className="mx-auto max-w-[1440px] px-5 py-8 md:px-8 md:py-10">
      <div className="flex flex-wrap items-end justify-between gap-4">
        <div>
          <p className="text-sm font-medium text-muted-foreground">
            Tahap 3 · Hasil
          </p>
          <h1 className="mt-1 font-serif text-3xl tracking-tight md:text-4xl">
            Hasil analisis
          </h1>
          <p className="mt-1.5 text-sm text-muted-foreground">
            {source.lontar_group} · Lempir{' '}
            {source.lempir_number ?? 'tidak tercatat'} · {source.language}
            <span className="sm:hidden">
              {' '}
              · {confidenceLabel[confidence]}
            </span>
          </p>
        </div>

        <div className="flex flex-wrap items-center gap-2">
          {isMock ? <Badge variant="brass">Data contoh</Badge> : null}
          <Badge
            variant={confidenceVariant[confidence]}
            className="hidden sm:inline-flex"
          >
            {confidenceLabel[confidence]}
          </Badge>
          <div className="flex flex-wrap items-center gap-2 print:hidden">
            <Button
              type="button"
              variant="outline"
              size="sm"
              onClick={onExport}
              aria-label="Unduh hasil sebagai JSON"
              className="size-11 px-0 sm:h-9 sm:w-auto sm:px-3"
            >
              <Download aria-hidden="true" />
              <span className="hidden sm:inline">JSON</span>
            </Button>
            <Button
              type="button"
              variant="outline"
              size="sm"
              onClick={onPrint}
              aria-label="Cetak hasil analisis"
              className="size-11 px-0 sm:h-9 sm:w-auto sm:px-3"
            >
              <Printer aria-hidden="true" />
              <span className="hidden sm:inline">Cetak</span>
            </Button>
            <Button
              type="button"
              variant="ghost"
              size="sm"
              onClick={onReset}
              aria-label="Mulai analisis naskah baru"
              className="size-11 px-0 sm:h-9 sm:w-auto sm:px-3"
            >
              <RotateCcw aria-hidden="true" />
              <span className="hidden sm:inline">Naskah baru</span>
            </Button>
          </div>
        </div>
      </div>

      <div className="mt-5 grid items-start gap-5 lg:grid-cols-[1.08fr_0.92fr]">
        <div className="order-2 lg:order-1">
          <LontarLens
            key={originalUrl}
            originalUrl={originalUrl}
            enhancedUrl={enhancedUrl}
            view={view}
            onViewChange={onViewChange}
            alt={`Lempir lontar ${source.lontar_group}`}
            zoomable
          />
        </div>

        <div className="order-1 lg:order-2">
          <Card className="flex flex-col overflow-hidden">
            <div className="flex items-center justify-between gap-3 border-b border-border px-4 py-3.5 sm:px-5">
              <div className="flex items-center gap-2">
                <BookOpen className="size-4 text-primary" aria-hidden="true" />
                <h2 className="text-sm font-semibold">Bacaan naskah</h2>
              </div>
              <span className="truncate text-xs text-muted-foreground">
                {plant.plant_local_name}
              </span>
            </div>

            <div className="p-4 sm:p-5">
              <AnalysisTabs result={result} />
            </div>

            <div className="mt-auto border-t border-border bg-brass-soft/55 px-4 py-3.5 sm:px-5">
              <div className="flex items-start gap-2">
                <Info
                  className="mt-0.5 size-3.5 shrink-0 text-brass-strong"
                  aria-hidden="true"
                />
                <p className="text-xs leading-relaxed text-brass-strong">
                  {result.scientific_validation.safety_disclaimer}
                </p>
              </div>
            </div>
          </Card>
          <Pipeline current={4} />
        </div>
      </div>
    </section>
  )
}
