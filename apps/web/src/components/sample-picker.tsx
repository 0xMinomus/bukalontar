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
    <div className="mt-5">
      <div className="flex items-baseline justify-between gap-4">
        <p className="font-heading text-xl tracking-tight">
          Atau mulai dari arsip contoh
        </p>
        <a
          href="https://commons.wikimedia.org/wiki/Category:Lontar"
          target="_blank"
          rel="noreferrer"
          className="small-caps text-[11px] font-semibold text-muted-foreground underline-offset-4 hover:text-foreground hover:underline"
        >
          Sumber foto
        </a>
      </div>

      <div className="mt-3 grid gap-2.5 sm:grid-cols-2">
        {samples.map((sample) => {
          const busy = busyId === sample.id
          return (
            <button
              key={sample.id}
              type="button"
              onClick={() => onSelect(sample)}
              disabled={disabled || Boolean(busyId)}
              className={cn(
                'group flex items-center gap-3 rounded-lg border border-border bg-card/80 p-2.5 text-left transition-[transform,border-color,box-shadow] duration-300 hover:-translate-y-0.5 hover:border-brass/60 hover:shadow-[0_14px_28px_-22px_rgb(62_42_20_/_0.7)] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background disabled:pointer-events-none disabled:opacity-60',
              )}
            >
              <span className="h-16 w-24 shrink-0 overflow-hidden rounded-sm border border-border/70 bg-[#2a2118]">
                <img
                  src={sample.src}
                  alt=""
                  width={96}
                  height={64}
                  loading="lazy"
                  decoding="async"
                  className="size-full scale-[1.75] object-cover sepia-[0.35] contrast-105 brightness-110 transition-transform duration-500 group-hover:scale-[1.85]"
                />
              </span>
              <span className="min-w-0">
                <span className="block font-heading text-lg leading-tight">
                  {sample.name}
                </span>
                <span className="mt-0.5 block truncate text-xs text-muted-foreground">
                  {sample.detail}
                </span>
                <span className="small-caps mt-1 block truncate text-[11px] text-brass-strong">
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

      <p className="mt-2.5 text-[11px] leading-relaxed text-muted-foreground">
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
