import { cn } from '@/lib/utils'

export function BrandMark({ className }: { className?: string }) {
  return (
    <svg
      viewBox="0 0 32 32"
      aria-hidden="true"
      className={cn('size-8', className)}
      fill="none"
    >
      <rect
        x="4"
        y="6.5"
        width="24"
        height="8"
        rx="3"
        stroke="currentColor"
        strokeWidth="1.5"
      />
      <circle cx="8.5" cy="10.5" r="1" fill="currentColor" />
      <path
        d="M12 10.5h12"
        stroke="currentColor"
        strokeWidth="1.5"
        strokeLinecap="round"
      />
      <rect
        x="4"
        y="17.5"
        width="19"
        height="8"
        rx="3"
        stroke="currentColor"
        strokeWidth="1.5"
      />
      <circle cx="8.5" cy="21.5" r="1" fill="currentColor" />
      <path
        d="M12 21.5h7"
        stroke="currentColor"
        strokeWidth="1.5"
        strokeLinecap="round"
      />
    </svg>
  )
}
