import { Loader2 } from 'lucide-react'

import { type Sample, samples } from '@/lib/samples'
import { cn } from '@/lib/utils'

interface SamplePickerProps {
  onSelect: (sample: Sample) => void
  busyId?: string | null
  disabled?: boolean
}

export function SamplePicker({
  onSelect,
  busyId,
  disabled,
}: SamplePickerProps) {
  return (
    <div className="mt-4">
      <div className="flex items-baseline justify-between gap-4">
        <p className="text-sm font-semibold">Atau mulai dari arsip contoh</p>
        <a
          href="https://commons.wikimedia.org/wiki/Category:Lontar"
          target="_blank"
          rel="noreferrer"
          className="text-xs font-medium text-muted-foreground underline-offset-4 hover:text-foreground hover:underline"
        >
          Sumber foto
        </a>
      </div>

      <div className="mt-2.5 grid gap-2 sm:grid-cols-2">
        {samples.map((sample) => {
          const busy = busyId === sample.id
          return (
            <button
              key={sample.id}
              type="button"
              onClick={() => onSelect(sample)}
              disabled={disabled || Boolean(busyId)}
              className={cn(
                'group flex items-center gap-3 rounded-lg border border-border bg-card p-2 text-left transition-colors hover:border-input hover:bg-secondary/60 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background disabled:pointer-events-none disabled:opacity-60',
              )}
            >
              <span className="h-14 w-20 shrink-0 overflow-hidden rounded-md bg-mat">
                <img
                  src={sample.src}
                  alt=""
                  width={80}
                  height={56}
                  loading="lazy"
                  decoding="async"
                  className="size-full scale-[1.75] object-cover"
                />
              </span>
              <span className="min-w-0">
                <span className="block text-sm font-semibold">
                  {sample.name}
                </span>
                <span className="block truncate text-xs text-muted-foreground">
                  {sample.detail}
                </span>
                <span className="mt-1 block truncate text-[11px] text-muted-foreground">
                  {sample.credit}
                </span>
              </span>
              {busy ? (
                <Loader2
                  className="ml-auto size-4 shrink-0 animate-spin text-muted-foreground"
                  aria-label="Memuat contoh"
                />
              ) : null}
            </button>
          )
        })}
      </div>

      <p className="mt-2 text-[11px] leading-relaxed text-muted-foreground">
        Foto arsip oleh{' '}
        <a
          href={samples[0].sourceUrl}
          target="_blank"
          rel="noreferrer"
          className="underline underline-offset-2 hover:text-foreground"
        >
          Eka343
        </a>
        , lisensi{' '}
        <a
          href={samples[0].licenseUrl}
          target="_blank"
          rel="noreferrer"
          className="underline underline-offset-2 hover:text-foreground"
        >
          CC BY-SA 4.0
        </a>
        . Lempir contoh dipakai untuk menguji antarmuka, bukan klaim koleksi
        resmi.
      </p>
    </div>
  )
}
