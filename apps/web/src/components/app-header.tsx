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
  document.getElementById(target)?.scrollIntoView({ behavior: 'smooth', block: 'start' })
}

export function AppHeader({ showNav = true, onHome }: AppHeaderProps) {
  return (
    <header className="app-header no-print">
      <div className="header-inner">
        <Link to="/" onClick={onHome} className="brand-link">
          <BrandMark className="size-11" />
          <span className="flex flex-col leading-none">
            <span className="font-heading text-[1.8rem] font-semibold tracking-tight">BukaLontar</span>
            <span className="mt-1 text-[10px] text-muted-foreground">Aksara, bahasa, pengetahuan Bali</span>
          </span>
        </Link>
        {showNav ? (
          <nav aria-label="Navigasi utama" className="header-nav">
            {sections.map(([label, target]) => (
              <button key={target} type="button" onClick={() => scrollToSection(target)}>
                {label}
              </button>
            ))}
          </nav>
        ) : (
          <span className="hidden text-xs text-muted-foreground sm:block">Ruang baca naskah Bali</span>
        )}
      </div>
    </header>
  )
}
