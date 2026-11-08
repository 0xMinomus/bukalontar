import { Code2 } from 'lucide-react'
import { Link } from 'react-router-dom'

import { BrandMark } from '@/components/brand-mark'

interface AppHeaderProps {
  showNav?: boolean
  onHome?: () => void
}

const sections = [
  ['Ruang kerja', 'top'],
  ['Alur', 'alur'],
  ['Tentang', 'tentang'],
] as const

function scrollToSection(target: string) {
  if (target === 'top') {
    window.scrollTo({ top: 0, behavior: 'smooth' })
    return
  }
  document
    .getElementById(target)
    ?.scrollIntoView({ behavior: 'smooth', block: 'start' })
}

export function AppHeader({ showNav = true, onHome }: AppHeaderProps) {
  return (
    <header className="border-b border-border/80 bg-background">
      <div className="relative mx-auto flex h-[4.5rem] max-w-[1440px] items-center justify-between gap-4 px-5 md:px-8">
        <Link
          to="/"
          onClick={onHome}
          className="flex items-center gap-3 rounded-lg focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background"
        >
          <span className="flex size-10 items-center justify-center rounded-full border border-brass/45 bg-primary text-primary-foreground shadow-[inset_0_0_0_1px_rgb(255_255_255_/_0.15)]">
            <BrandMark className="size-5" />
          </span>
          <span className="flex flex-col leading-none">
            <span className="font-heading text-2xl tracking-tight">
              BukaLontar
            </span>
            <span className="small-caps mt-1 text-[10px] font-semibold text-muted-foreground">
              Arsip Usada Bali
            </span>
          </span>
        </Link>

        {showNav ? (
          <nav
            aria-label="Navigasi utama"
            className="absolute left-1/2 hidden -translate-x-1/2 items-center gap-6 md:flex"
          >
            {sections.map(([label, target]) => (
              <button
                key={target}
                type="button"
                onClick={() => scrollToSection(target)}
                className="group relative cursor-pointer border-0 bg-transparent py-2 text-sm font-medium text-muted-foreground transition-colors hover:text-foreground"
              >
                {label}
                <span
                  className="absolute inset-x-0 -bottom-0.5 h-px origin-left scale-x-0 bg-brass transition-transform duration-300 group-hover:scale-x-100"
                  aria-hidden="true"
                />
              </button>
            ))}
          </nav>
        ) : null}

        <div>
          <a
            href="https://github.com/0xMinomus/bukalontar"
            target="_blank"
            rel="noreferrer"
            className="inline-flex h-10 items-center gap-2 rounded-lg border border-border/70 px-3 text-xs font-semibold text-muted-foreground transition-colors hover:border-brass/50 hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background"
            aria-label="Buka repositori GitHub BukaLontar"
          >
            <Code2 className="size-4" aria-hidden="true" />
            <span className="small-caps hidden sm:inline">GitHub</span>
          </a>
        </div>
      </div>
    </header>
  )
}
