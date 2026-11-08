import { Code2 } from 'lucide-react'

import { BrandMark } from '@/components/brand-mark'

export function AppHeader() {
  return (
    <header className="border-b border-border bg-background">
      <div className="mx-auto flex h-16 max-w-[1440px] items-center gap-8 px-5 md:px-8">
        <a
          href="#workspace"
          className="flex items-center gap-3 rounded-md focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background"
        >
          <span className="flex size-9 items-center justify-center rounded-lg bg-primary text-primary-foreground">
            <BrandMark className="size-5" />
          </span>
          <span className="flex flex-col leading-none">
            <span className="font-serif text-xl tracking-tight">UsadhaAI</span>
            <span className="mt-1 text-[11px] font-medium text-muted-foreground">
              Arsip usada Bali
            </span>
          </span>
        </a>

        <nav
          aria-label="Navigasi utama"
          className="hidden items-center gap-5 text-sm font-medium text-muted-foreground md:flex"
        >
          <a className="transition-colors hover:text-foreground" href="#workspace">
            Ruang kerja
          </a>
          <a className="transition-colors hover:text-foreground" href="#alur">
            Alur
          </a>
          <a className="transition-colors hover:text-foreground" href="#tentang">
            Tentang
          </a>
        </nav>

        <div className="ml-auto flex items-center gap-3">
          <a
            href="https://github.com/0xMinomus/digifest"
            target="_blank"
            rel="noreferrer"
            className="inline-flex h-9 items-center gap-1.5 rounded-lg px-2.5 text-sm font-medium text-muted-foreground transition-colors hover:bg-secondary hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background"
            aria-label="Buka repositori GitHub UsadhaAI"
          >
            <Code2 className="size-4" aria-hidden="true" />
            <span className="hidden sm:inline">GitHub</span>
          </a>
        </div>
      </div>
    </header>
  )
}
