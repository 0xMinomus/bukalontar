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
  isMock: boolean
}

export function Landing({ onFile, onSample, busySampleId, isMock }: LandingProps) {
  return (
    <section className="manuscript-landing">
      <div className="landing-folio mx-auto max-w-5xl px-5 text-center md:px-8">
        <div className="hero-copy">
          <h1 className="font-heading">
            Baca lontar,<br />
            <em>pahami isinya.</em>
          </h1>
          <p className="mx-auto max-w-[46ch] text-sm leading-relaxed text-muted-foreground sm:text-base">
            Dari foto lempir ke aksara Latin, terjemahan, dan konteks naskah Bali.
          </p>
        </div>

        <div className="upload-scene">
          <span className="scene-leaf scene-leaf-left" aria-hidden="true" />
          <span className="scene-leaf scene-leaf-right" aria-hidden="true" />
          <Dropzone onFile={onFile} />
        </div>
        <p className="upload-note text-xs leading-relaxed text-muted-foreground">
          JPG, PNG, WEBP · Maks. 8 MB · Bisa tempel dengan Ctrl V
        </p>

        <Reveal className="sample-folio text-left">
          <SamplePicker onSelect={onSample} busyId={busySampleId} disabled={Boolean(busySampleId)} />
          {isMock ? (
            <div className="mt-4 flex items-start gap-2 text-xs leading-relaxed text-muted-foreground">
              <Info className="mt-0.5 size-3.5 shrink-0" aria-hidden="true" />
              <p>Mode contoh. Hasil memakai data demonstrasi, bukan pembacaan foto Anda.</p>
            </div>
          ) : null}
        </Reveal>

        <div className="landing-colophon">
          <Ornament className="mx-auto w-32 text-brass-strong" />
          <p className="mt-3 text-xs leading-relaxed text-muted-foreground">
            Kakawin, babad, geguritan, tutur, dan usadha.<br className="sm:hidden" /> Setiap lontar dibaca sesuai isinya.
          </p>
        </div>
      </div>
    </section>
  )
}
