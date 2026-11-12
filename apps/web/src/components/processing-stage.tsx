import { Check, Circle, Loader2, X } from 'lucide-react'

import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Progress } from '@/components/ui/progress'
import { cn } from '@/lib/utils'

const steps = [
  { title: 'Menyiapkan citra', detail: 'Foto asli tanpa perubahan aksara' },
  { title: 'Membaca aksara', detail: 'Aksara Bali dan pasang aksara' },
  { title: 'Menyusun terjemahan', detail: 'Indonesia dan Inggris' },
  { title: 'Memeriksa konteks', detail: 'Sumber naskah dan resep Usadha bila ada' },
]

interface ProcessingStageProps {
  imageUrl: string
  step: number
  onCancel: () => void
}

export function ProcessingStage({
  imageUrl,
  step,
  onCancel,
}: ProcessingStageProps) {
  const progress = Math.min(100, Math.round((step / steps.length) * 100))

  return (
    <section className="mx-auto max-w-3xl px-5 py-10 md:px-8 md:py-16">
      <Card className="overflow-hidden">
        <div className="grid gap-0 md:grid-cols-[0.42fr_0.58fr]">
          <div className="relative hidden min-h-[280px] overflow-hidden bg-mat md:block">
            <img
              src={imageUrl}
              alt=""
              className="absolute inset-0 size-full object-cover opacity-70"
            />
            <div className="absolute inset-0 bg-gradient-to-t from-mat via-transparent to-transparent" />
            <p className="absolute inset-x-4 bottom-4 text-[11px] font-medium text-white">
              Lempir sedang dibaca
            </p>
          </div>

          <div className="p-5 sm:p-7">
            <p className="small-caps text-xs font-semibold text-brass-strong">
              Membaca lontar
            </p>
            <h1 className="mt-1 font-heading text-4xl tracking-tight">
              Menyusun hasil
            </h1>
            <p className="mt-2 text-sm leading-relaxed text-muted-foreground">
              Aksara sedang dibaca dan diterjemahkan. Jika ada resep Usadha,
              hasilnya dilengkapi rujukan tanaman yang tersedia.
            </p>

            <ol className="mt-6 space-y-0" aria-live="polite">
              {steps.map((item, index) => {
                const number = index + 1
                const done = step > number
                const active = step === number
                return (
                  <li key={item.title} className="flex gap-3">
                    <div className="flex flex-col items-center">
                      <span
                        className={cn(
                          'flex size-6 shrink-0 items-center justify-center rounded-full border text-[11px] font-semibold',
                          done && 'border-primary bg-primary text-primary-foreground',
                          active &&
                            'border-brass bg-brass-soft text-brass-strong',
                          !done &&
                            !active &&
                            'border-border bg-card text-muted-foreground',
                        )}
                      >
                        {done ? (
                          <Check className="size-3.5" aria-hidden="true" />
                        ) : active ? (
                          <Loader2
                            className="size-3.5 animate-spin"
                            aria-hidden="true"
                          />
                        ) : (
                          <Circle className="size-2" aria-hidden="true" />
                        )}
                      </span>
                      {index < steps.length - 1 ? (
                        <span
                          className={cn(
                            'my-1 w-px flex-1',
                            done ? 'bg-primary/40' : 'bg-border',
                          )}
                        />
                      ) : null}
                    </div>
                    <div className="pb-4">
                      <p
                        className={cn(
                          'text-sm font-semibold',
                          !done && !active && 'text-muted-foreground',
                        )}
                      >
                        {item.title}
                      </p>
                      <p className="mt-0.5 text-xs text-muted-foreground">
                        {item.detail}
                      </p>
                    </div>
                  </li>
                )
              })}
            </ol>

            <Progress
              value={progress}
              className="mt-1"
              aria-label="Kemajuan analisis"
            />

            <div className="mt-4 flex items-center justify-between gap-4">
              <p className="num text-xs text-muted-foreground">
                {progress}% selesai
              </p>
              <Button type="button" variant="ghost" size="sm" onClick={onCancel}>
                <X aria-hidden="true" />
                Batalkan
              </Button>
            </div>
          </div>
        </div>
      </Card>
    </section>
  )
}
