import { useLayoutEffect, useRef, type ReactNode } from 'react'

import { cn } from '@/lib/utils'

interface RevealProps {
  children: ReactNode
  className?: string
  delay?: number
}

export function Reveal({ children, className, delay = 0 }: RevealProps) {
  const ref = useRef<HTMLDivElement>(null)

  useLayoutEffect(() => {
    const element = ref.current
    if (
      !element ||
      window.matchMedia('(prefers-reduced-motion: reduce)').matches ||
      element.getBoundingClientRect().top < window.innerHeight ||
      !('IntersectionObserver' in window)
    ) return

    element.dataset.reveal = 'pending'
    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          element.dataset.reveal = 'visible'
          observer.disconnect()
        }
      },
      { rootMargin: '0px 0px -24px 0px' },
    )
    observer.observe(element)
    return () => {
      observer.disconnect()
      delete element.dataset.reveal
    }
  }, [])

  return (
    <div
      ref={ref}
      className={cn('reveal', className)}
      style={{ transitionDelay: `${delay}ms` }}
    >
      {children}
    </div>
  )
}
