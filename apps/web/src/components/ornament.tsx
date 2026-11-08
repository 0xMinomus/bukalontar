import { cn } from '@/lib/utils'

export function Ornament({ className }: { className?: string }) {
  return (
    <svg
      viewBox="0 0 180 16"
      fill="none"
      aria-hidden="true"
      className={cn('h-4 w-40 text-brass', className)}
    >
      <path d="M0 8h62" stroke="currentColor" strokeWidth="1" opacity=".55" />
      <path d="M118 8h62" stroke="currentColor" strokeWidth="1" opacity=".55" />
      <path d="M86 3l5 5-5 5-5-5 5-5Z" fill="currentColor" />
      <circle cx="72" cy="8" r="1.6" fill="currentColor" />
      <circle cx="108" cy="8" r="1.6" fill="currentColor" />
      <path
        d="M78 8c3-3 6-3 9 0M102 8c-3 3-6 3-9 0"
        stroke="currentColor"
        strokeWidth="1"
        opacity=".7"
      />
    </svg>
  )
}
