import { Loader2 } from 'lucide-react'

import { type Sample, samples } from '@/lib/samples'

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
      <h2 className="font-heading text-xl tracking-tight">
        Atau mulai dari arsip contoh
      </h2>
      <p className="mt-1 text-xs leading-relaxed text-muted-foreground">
        Dua Usadha, satu Carcan, satu Kakawin. Contoh arsip memakai transkripsi
        rujukan, bukan pembacaan foto murni.
      </p>

      <div className="mt-3 grid gap-2.5 sm:grid-cols-2">
        {samples.map((sample) => {
          const busy = busyId === sample.id
          return (
            <div key={sample.id} className="rounded-lg border border-border bg-card/80">
              <button
                type="button"
                onClick={() => onSelect(sample)}
                disabled={disabled || Boolean(busyId)}
                className="group flex w-full items-center gap-3 rounded-lg p-2.5 text-left transition-colors duration-300 hover:bg-secondary/50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background disabled:pointer-events-none disabled:opacity-60"
              >
                <span className="h-16 w-24 shrink-0 overflow-hidden rounded-sm border border-border/70 bg-[#2a2118]">
                  <img
                    src={sample.src}
                    alt=""
                    width={96}
                    height={64}
                    loading="lazy"
                    decoding="async"
                    className="size-full object-cover sepia-[0.2] contrast-105"
                  />
                </span>
                <span className="min-w-0">
                  <span className="small-caps block text-[10px] text-brass-strong">
                    {sample.genre}
                  </span>
                  <span className="block font-heading text-lg leading-tight">
                    {sample.name}
                  </span>
                  <span className="mt-0.5 block text-xs text-muted-foreground">
                    {sample.detail}
                  </span>
                </span>
                {busy ? (
                  <Loader2
                    className="ml-auto size-4 shrink-0 animate-spin text-muted-foreground"
                    aria-label="Memuat contoh"
                  />
                ) : null}
              </button>
              <div className="border-t border-border/60 px-3 py-2 text-[11px] leading-relaxed text-muted-foreground">
                <p>{sample.credit}</p>
                <div className="mt-1 flex flex-wrap gap-x-3 gap-y-1">
                  <a href={sample.sourceUrl} target="_blank" rel="noreferrer" aria-label={`Sumber foto ${sample.name}`} className="underline underline-offset-2 hover:text-foreground">
                    Sumber foto
                  </a>
                  <a href={sample.referenceUrl} target="_blank" rel="noreferrer" aria-label={`Bacaan rujukan ${sample.name}`} className="underline underline-offset-2 hover:text-foreground">
                    Bacaan rujukan
                  </a>
                  <a href={sample.licenseUrl} target="_blank" rel="noreferrer" aria-label={`Lisensi foto ${sample.name}`} className="underline underline-offset-2 hover:text-foreground">
                    Lisensi foto
                  </a>
                </div>
              </div>
            </div>
          )
        })}
      </div>
      <p className="mt-2.5 text-[11px] leading-relaxed text-muted-foreground">
        Genre mengikuti katalog sumber. Foto unggahan lain tetap memakai
        pembacaan eksperimental; hasil terjemahan dan racikan perlu ditinjau.
      </p>
    </div>
  )
}
