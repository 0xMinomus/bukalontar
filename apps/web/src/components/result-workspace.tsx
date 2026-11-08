import {
  ArrowUpRight,
  BookOpen,
  Download,
  Info,
  Printer,
  RotateCcw,
} from 'lucide-react'

import { AnalysisTabs } from '@/components/analysis-tabs'
import { LontarLens, type LensView } from '@/components/lontar-lens'
import { Pipeline } from '@/components/pipeline'
import { PrintReport } from '@/components/print-report'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import type { AnalysisResponse, Confidence } from '@/lib/analysis'
import { hasMedicinalPlant } from '@/lib/reading'

const confidenceLabel: Record<Confidence, string> = {
  low: 'Keterbacaan rendah',
  medium: 'Keterbacaan sedang',
  high: 'Keterbacaan tinggi',
}

const confidenceShort: Record<Confidence, string> = {
  low: 'Rendah',
  medium: 'Sedang',
  high: 'Tinggi',
}

const confidenceVariant: Record<
  Confidence,
  'success' | 'brass' | 'destructive'
> = {
  low: 'destructive',
  medium: 'brass',
  high: 'success',
}

const glossary = [
  { term: 'mising', meaning: 'diare' },
  { term: 'cika', meaning: 'kram perut' },
  { term: 'taru', meaning: 'pohon' },
  { term: 'panes', meaning: 'hangat' },
  { term: 'tis', meaning: 'dingin' },
  { term: 'dumelada', meaning: 'sedang' },
]

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
  const hasPlant = hasMedicinalPlant(plant.plant_local_name)

  return (
    <section className="result-page mx-auto max-w-[1440px] animate-[fade-rise_500ms_ease-out_both] px-5 py-8 md:px-8 md:py-10">
      <div className="screen-only flex flex-wrap items-end justify-between gap-4">
        <div>
          <p className="small-caps text-xs font-semibold text-brass-strong">
            Hasil pembacaan
          </p>
          <h1 className="mt-1 font-heading text-4xl tracking-tight md:text-5xl">
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
              aria-label="Unduh hasil analisis"
              className="size-11 px-0 sm:h-9 sm:w-auto sm:px-3"
            >
              <Download aria-hidden="true" />
              <span className="hidden sm:inline">Unduh</span>
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

      <dl className="screen-only mt-5 grid grid-cols-2 gap-px overflow-hidden rounded-lg border border-border bg-border/70 sm:grid-cols-3">
        <div className="bg-card px-4 py-3">
          <dt className="label-xs">Aksara Bali</dt>
          <dd className="mt-1 text-sm font-medium">
            {result.paleography.aksara_bali_detected
              ? 'Terbaca'
              : 'Tidak terbaca'}
          </dd>
        </div>
        <div className="bg-card px-4 py-3">
          <dt className="label-xs">Keterbacaan</dt>
          <dd className="mt-1 text-sm font-medium">
            {confidenceShort[confidence]}
          </dd>
        </div>
        <div className="bg-card px-4 py-3">
          <dt className="label-xs">Rujukan tanaman</dt>
          <dd className="mt-1 text-sm font-medium">
            {!hasPlant
              ? 'Tidak ada'
              : plant.kb_match
                ? 'Tercatat'
                : 'Belum tercatat'}
          </dd>
        </div>
      </dl>

      <div className="screen-only mt-5 grid items-start gap-5 lg:grid-cols-[1.08fr_0.92fr]">
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
          <Card className="paper-panel flex flex-col overflow-hidden">
            <div className="flex items-center justify-between gap-3 border-b border-border px-4 py-3.5 sm:px-5">
              <div className="flex items-center gap-2">
                <BookOpen className="size-4 text-primary" aria-hidden="true" />
                <h2 className="font-heading text-lg leading-none">Bacaan naskah</h2>
              </div>
              <span className="truncate text-xs text-muted-foreground">
                {hasPlant ? plant.plant_local_name : 'Tidak ada tanaman obat'}
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

      <section className="screen-only important-info mt-8 grid gap-7 border-t border-border/70 pt-7 md:grid-cols-[1.1fr_0.95fr_0.95fr]">
          <div>
            <p className="small-caps text-xs font-semibold text-brass-strong">
              Cara membaca
            </p>
            <p className="mt-2 text-sm leading-relaxed text-muted-foreground">
              Transliterasi mengikuti bunyi aksara Bali. Terjemahan Kawi
              bersifat semantik, bukan kata per kata. Periksa lempir asli dan
              rujukan filologi untuk keputusan penting.
            </p>
          </div>

          <div>
            <p className="small-caps text-xs font-semibold text-brass-strong">
              Istilah kunci
            </p>
            <dl className="mt-2 grid grid-cols-2 gap-x-4 gap-y-1.5 text-sm">
              {glossary.map((item) => (
                <div key={item.term} className="flex gap-1.5">
                  <dt className="font-serif italic text-foreground">
                    {item.term}
                  </dt>
                  <dd className="text-muted-foreground">{item.meaning}</dd>
                </div>
              ))}
            </dl>
          </div>

          <div>
            <p className="small-caps text-xs font-semibold text-brass-strong">
              Catatan keterbacaan
            </p>
            <p className="mt-2 text-sm leading-relaxed text-muted-foreground">
              Goresan aus, tumpukan pasang aksara, dan istilah arkais bisa
              salah dibaca. Periksa lempir asli dan rujukan filologi. Ini bukan
              pengganti saran medis.
            </p>
            <a
              href="https://github.com/0xMinomus/digifest/tree/main/docs"
              target="_blank"
              rel="noreferrer"
              className="small-caps mt-3 inline-flex items-center gap-1.5 border-b border-brass/50 pb-1 text-[11px] font-semibold text-primary transition-colors hover:text-brass-strong print:hidden"
            >
              Panduan membaca
              <ArrowUpRight className="size-3.5" aria-hidden="true" />
            </a>
          </div>
      </section>

      <PrintReport result={result} originalUrl={originalUrl} />
    </section>
  )
}
