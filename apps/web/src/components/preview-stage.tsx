import { ArrowLeft, Check, Loader2, ScanText } from 'lucide-react'

import { LontarLens, type LensView } from '@/components/lontar-lens'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Separator } from '@/components/ui/separator'
import { formatBytes } from '@/lib/utils'

interface PreviewStageProps {
  fileName: string
  fileSize: number
  dimensions?: { width: number; height: number }
  originalUrl: string
  enhancedUrl?: string | null
  enhancing: boolean
  enhanceError?: string | null
  enhanceElapsedMs?: number
  view: LensView
  onViewChange: (view: LensView) => void
  onAnalyze: () => void
  onReset: () => void
}

const tips = [
  'Seluruh lempir masuk ke dalam frame.',
  'Cahaya merata, tanpa bayangan tangan.',
  'Goresan aksara terlihat tajam saat diperbesar.',
]

export function PreviewStage({
  fileName,
  fileSize,
  dimensions,
  originalUrl,
  enhancedUrl,
  enhancing,
  enhanceError,
  enhanceElapsedMs,
  view,
  onViewChange,
  onAnalyze,
  onReset,
}: PreviewStageProps) {
  return (
    <section className="mx-auto max-w-[1440px] px-5 py-8 md:px-8 md:py-10">
      <div className="mb-5 flex flex-wrap items-end justify-between gap-4">
        <div>
          <p className="text-sm font-medium text-muted-foreground">
            Tahap 1 · Periksa citra
          </p>
          <h1 className="mt-1 font-serif text-3xl tracking-tight md:text-4xl">
            Pastikan aksara terbaca
          </h1>
        </div>
        <Button type="button" variant="ghost" onClick={onReset}>
          <ArrowLeft aria-hidden="true" />
          Ganti foto
        </Button>
      </div>

      <div className="grid gap-5 lg:grid-cols-[1.35fr_0.65fr]">
        <LontarLens
          key={originalUrl}
          originalUrl={originalUrl}
          enhancedUrl={enhancedUrl}
          view={view}
          onViewChange={onViewChange}
          alt={`Pratinjau ${fileName}`}
        />

        <Card className="h-fit p-5">
          <div className="flex items-start justify-between gap-3">
            <div className="min-w-0">
              <p className="label-xs">Berkas</p>
              <h2 className="mt-1 break-all text-sm font-semibold">
                {fileName}
              </h2>
            </div>
            <Badge variant="outline" className="shrink-0">
              {dimensions
                ? `${dimensions.width} × ${dimensions.height}`
                : 'Menyiapkan'}
            </Badge>
          </div>

          <Separator className="my-4" />

          <dl className="space-y-2.5 text-sm">
            <div className="flex items-center justify-between gap-4">
              <dt className="text-muted-foreground">Ukuran</dt>
              <dd className="num font-medium">{formatBytes(fileSize)}</dd>
            </div>
            <div className="flex items-center justify-between gap-4">
              <dt className="text-muted-foreground">Kontras adaptif</dt>
              <dd className="font-medium">
                {enhancing
                  ? 'Diproses'
                  : enhancedUrl
                    ? `${enhanceElapsedMs ?? 0} ms`
                    : 'Tidak tersedia'}
              </dd>
            </div>
            <div className="flex items-center justify-between gap-4">
              <dt className="text-muted-foreground">Batas kirim</dt>
              <dd className="num font-medium">8 MB</dd>
            </div>
          </dl>

          {enhancing ? (
            <div className="mt-4 flex items-center gap-2 text-xs text-muted-foreground">
              <Loader2 className="size-3.5 animate-spin" aria-hidden="true" />
              Meningkatkan kontras di browser
            </div>
          ) : null}

          {enhanceError ? (
            <p className="mt-4 rounded-md border border-destructive/30 bg-destructive/10 p-2.5 text-xs leading-relaxed text-destructive">
              {enhanceError}
            </p>
          ) : null}

          <Button
            type="button"
            size="lg"
            className="mt-6 w-full"
            onClick={onAnalyze}
            disabled={enhancing}
          >
            <ScanText aria-hidden="true" />
            Analisis naskah
          </Button>

          <p className="mt-2 text-center text-xs leading-relaxed text-muted-foreground">
            AI akan membaca aksara, menerjemahkan Kawi, lalu memetakan tanaman.
          </p>

          <Separator className="my-5" />

          <div>
            <p className="label-xs">Tips foto</p>
            <ul className="mt-2 space-y-2">
              {tips.map((tip) => (
                <li
                  key={tip}
                  className="flex items-start gap-2 text-sm leading-relaxed text-muted-foreground"
                >
                  <Check
                    className="mt-0.5 size-3.5 shrink-0 text-primary"
                    aria-hidden="true"
                  />
                  {tip}
                </li>
              ))}
            </ul>
          </div>
        </Card>
      </div>
    </section>
  )
}
