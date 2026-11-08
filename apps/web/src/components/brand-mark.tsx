import { cn } from '@/lib/utils'

export function BrandMark({ className }: { className?: string }) {
  return (
    <img
      src="/bukalontar-icon.svg"
      alt=""
      aria-hidden="true"
      width={40}
      height={40}
      className={cn('size-8 object-contain', className)}
    />
  )
}
