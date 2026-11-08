import { ArrowUp } from 'lucide-react'

import { BrandMark } from '@/components/brand-mark'

export function SiteFooter() {
  return (
    <footer className="border-t border-border">
      <div className="mx-auto max-w-[1440px] px-5 py-8 md:px-8">
        <div className="flex flex-wrap items-start justify-between gap-6">
          <div className="flex items-center gap-3">
            <span className="flex size-8 items-center justify-center rounded-lg bg-primary text-primary-foreground">
              <BrandMark className="size-4" />
            </span>
            <div>
              <p className="font-serif text-lg leading-none tracking-tight">
                UsadhaAI
              </p>
              <p className="mt-1.5 text-xs text-muted-foreground">
                Bali DigiFest 2026 · AI Web Innovation Challenge
              </p>
            </div>
          </div>

          <div className="flex flex-wrap items-center gap-x-5 gap-y-2 text-sm font-medium text-muted-foreground">
            <a
              className="transition-colors hover:text-foreground"
              href="https://github.com/0xMinomus/digifest"
              target="_blank"
              rel="noreferrer"
            >
              Repositori
            </a>
            <a
              className="transition-colors hover:text-foreground"
              href="https://github.com/0xMinomus/digifest/tree/main/docs"
              target="_blank"
              rel="noreferrer"
            >
              Dokumentasi
            </a>
            <a
              className="inline-flex items-center gap-1 transition-colors hover:text-foreground"
              href="#workspace"
            >
              Kembali ke atas
              <ArrowUp className="size-3.5" aria-hidden="true" />
            </a>
          </div>
        </div>

        <div className="mt-6 flex flex-wrap items-center justify-between gap-3 border-t border-border pt-4 text-[11px] leading-relaxed text-muted-foreground">
          <p>
            Foto contoh oleh Eka343, lisensi{' '}
            <a
              href="https://creativecommons.org/licenses/by-sa/4.0/"
              target="_blank"
              rel="noreferrer"
              className="underline underline-offset-2 hover:text-foreground"
            >
              CC BY-SA 4.0
            </a>
            . Dipakai untuk menguji antarmuka.
          </p>
          <p>UsadhaAI, 2026.</p>
        </div>
      </div>
    </footer>
  )
}
