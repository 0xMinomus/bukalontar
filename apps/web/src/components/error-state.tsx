import { AlertTriangle, ArrowLeft, RotateCcw } from 'lucide-react'

import { Button } from '@/components/ui/button'

interface ErrorStateProps {
  message: string
  canRetry: boolean
  onRetry: () => void
  onReset: () => void
}

export function ErrorState({
  message,
  canRetry,
  onRetry,
  onReset,
}: ErrorStateProps) {
  return (
    <section className="mx-auto max-w-xl px-5 py-16 md:py-24">
      <div className="flex size-12 items-center justify-center rounded-lg border border-destructive/30 bg-destructive/10 text-destructive">
        <AlertTriangle className="size-5" aria-hidden="true" />
      </div>
      <h1 className="mt-5 font-serif text-3xl tracking-tight">
        Naskah belum bisa dibaca
      </h1>
      <p className="mt-2 text-sm leading-relaxed text-muted-foreground">
        {message}
      </p>
      <div className="mt-6 flex flex-wrap gap-2">
        {canRetry ? (
          <Button type="button" onClick={onRetry}>
            <RotateCcw aria-hidden="true" />
            Coba lagi
          </Button>
        ) : null}
        <Button
          type="button"
          variant={canRetry ? 'outline' : 'default'}
          onClick={onReset}
        >
          <ArrowLeft aria-hidden="true" />
          Kembali ke awal
        </Button>
      </div>
    </section>
  )
}
