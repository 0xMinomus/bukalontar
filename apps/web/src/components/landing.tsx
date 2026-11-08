import { Info } from 'lucide-react'

import { Dropzone } from '@/components/dropzone'
import { Ornament } from '@/components/ornament'
import { Reveal } from '@/components/reveal'
import { SamplePicker } from '@/components/sample-picker'
import type { Sample } from '@/lib/samples'

interface LandingProps {
  onFile: (file: File) => void
  onSample: (sample: Sample) => void
  busySampleId?: string | null
}

export function Landing({ onFile, onSample, busySampleId }: LandingProps) {
  return (
    <section
      id="workspace"
      className="relative overflow-hidden border-b border-border/70"
    >
      <div className="relative mx-auto max-w-5xl px-5 pb-14 pt-12 text-center md:px-8 md:pb-20 md:pt-20">
        <Reveal>
          <p className="small-caps text-xs font-semibold text-brass-strong">
            Bali DigiFest 2026 · Arsip Usada
          </p>

          <h1 className="mt-6 font-heading text-[clamp(3.3rem,9vw,7.4rem)] font-medium leading-[0.9] tracking-[-0.035em] text-foreground">
            Baca lontar.
            <br />
            <span className="italic text-primary">Temukan usadha.</span>
          </h1>

          <Ornament className="mx-auto mt-7 w-44 text-brass-strong/75" />

          <p className="mx-auto mt-5 max-w-[58ch] text-base leading-relaxed text-muted-foreground sm:text-lg">
            Foto lempir lontar, lalu baca transliterasi aksara Bali, terjemahan
            bahasa Kawi, dan kartu tanaman obat dalam satu alur.
          </p>
        </Reveal>

        <Reveal delay={130} className="mx-auto mt-9 max-w-4xl text-left">
          <Dropzone onFile={onFile} />
          <SamplePicker
            onSelect={onSample}
            busyId={busySampleId}
            disabled={Boolean(busySampleId)}
          />

          <div className="mt-5 flex items-start justify-center gap-2 text-xs leading-relaxed text-muted-foreground">
            <Info className="mt-0.5 size-3.5 shrink-0" aria-hidden="true" />
            <p>
              Pratinjau antarmuka. Hasil AI masih memakai data contoh sampai
              backend Go terhubung.
            </p>
          </div>
        </Reveal>

        <Reveal
          delay={220}
          className="mx-auto mt-8 flex max-w-4xl flex-wrap items-center justify-center gap-x-5 gap-y-2 border-t border-border/70 pt-5 text-xs text-muted-foreground"
        >
          <span>
            <span className="small-caps font-semibold text-brass-strong">
              Aksara
            </span>{' '}
            Bali
          </span>
          <span aria-hidden="true" className="text-border">
            ·
          </span>
          <span>
            <span className="small-caps font-semibold text-brass-strong">
              Bahasa
            </span>{' '}
            Kawi
          </span>
          <span aria-hidden="true" className="text-border">
            ·
          </span>
          <span>
            <span className="small-caps font-semibold text-brass-strong">
              Rujukan
            </span>{' '}
            Taru Pramana
          </span>
        </Reveal>
      </div>
    </section>
  )
}
