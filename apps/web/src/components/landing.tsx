import { Info } from 'lucide-react'

import { Dropzone } from '@/components/dropzone'
import { SamplePicker } from '@/components/sample-picker'
import type { Sample } from '@/lib/samples'

interface LandingProps {
  onFile: (file: File) => void
  onSample: (sample: Sample) => void
  busySampleId?: string | null
}

export function Landing({ onFile, onSample, busySampleId }: LandingProps) {
  return (
    <section className="mx-auto max-w-[1440px] px-5 pb-12 pt-10 md:px-8 md:pb-16 md:pt-16">
      <div className="grid items-start gap-10 lg:grid-cols-[1.05fr_0.95fr] lg:gap-16">
        <div className="lg:pt-2">
          <p className="text-sm font-medium text-muted-foreground">
            Bali DigiFest 2026 · Arsip usada
          </p>

          <h1 className="mt-4 font-serif text-[clamp(2.5rem,6vw,4.5rem)] leading-[1.02] tracking-tight text-foreground">
            Baca lontar.
            <br />
            <span className="italic text-primary">Temukan usadha.</span>
          </h1>

          <p className="mt-5 max-w-[52ch] text-base leading-relaxed text-muted-foreground sm:text-lg">
            Foto lempir lontar, lalu baca transliterasi aksara Bali, terjemahan
            bahasa Kawi, dan kartu tanaman obat dalam satu alur.
          </p>

          <blockquote className="mt-7 max-w-[54ch] border-l-2 border-brass pl-4">
            <p className="font-serif text-[1.05rem] italic leading-relaxed text-foreground">
              “Tityang taru sotong, daging anget, rasa sepet, dados anggen tamba
              pangemped mising.”
            </p>
            <cite className="mt-1.5 block text-xs not-italic text-muted-foreground">
              Lontar Taru Pramana, Va. 1384/10
            </cite>
          </blockquote>

          <dl className="mt-8 grid max-w-[520px] grid-cols-3 gap-4 border-t border-border pt-5">
            <div>
              <dt className="text-xs font-semibold text-muted-foreground">
                Aksara
              </dt>
              <dd className="mt-1 font-serif text-base leading-tight sm:text-lg">
                Bali
              </dd>
            </div>
            <div>
              <dt className="text-xs font-semibold text-muted-foreground">
                Bahasa
              </dt>
              <dd className="mt-1 font-serif text-base leading-tight sm:text-lg">
                Kawi
              </dd>
            </div>
            <div>
              <dt className="text-xs font-semibold text-muted-foreground">
                Rujukan
              </dt>
              <dd className="mt-1 font-serif text-base leading-tight sm:text-lg">
                Taru Pramana
              </dd>
            </div>
          </dl>
        </div>

        <div>
          <Dropzone onFile={onFile} />
          <SamplePicker
            onSelect={onSample}
            busyId={busySampleId}
            disabled={Boolean(busySampleId)}
          />

          <div className="mt-4 flex items-start gap-2 text-xs leading-relaxed text-muted-foreground">
            <Info className="mt-0.5 size-3.5 shrink-0" aria-hidden="true" />
            <p>
              Pratinjau antarmuka. Hasil AI masih memakai data contoh sampai
              backend Go terhubung.
            </p>
          </div>
        </div>
      </div>
    </section>
  )
}
