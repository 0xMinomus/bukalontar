import { Check } from 'lucide-react'

import { cn } from '@/lib/utils'

const steps = ['Unggah', 'Tingkatkan', 'Baca', 'Konteks']

export function Pipeline({ current = steps.length }: { current?: number }) {
  return (
    <div className="mt-4 rounded-xl border border-border bg-card px-4 py-3">
      <ol
        className="grid grid-cols-2 gap-x-4 gap-y-2 sm:flex sm:items-center sm:gap-0"
        aria-label="Alur analisis"
      >
        {steps.map((step, index) => {
          const done = index < current
          return (
            <li
              key={step}
              className="flex min-w-0 items-center gap-2 sm:flex-1 sm:last:flex-none"
            >
              <span
                className={cn(
                  'flex size-5 shrink-0 items-center justify-center rounded-full border text-[10px] font-semibold',
                  done
                    ? 'border-primary bg-primary text-primary-foreground'
                    : 'border-border bg-card text-muted-foreground',
                )}
              >
                {done ? (
                  <Check className="size-3" aria-hidden="true" />
                ) : (
                  index + 1
                )}
              </span>
              <span
                className={cn(
                  'whitespace-nowrap text-xs font-medium',
                  done ? 'text-foreground' : 'text-muted-foreground',
                )}
              >
                {step}
              </span>
              {index < steps.length - 1 ? (
                <span
                  className={cn(
                    'mx-2 hidden h-px min-w-4 flex-1 sm:block',
                    done ? 'bg-primary/40' : 'bg-border',
                  )}
                  aria-hidden="true"
                />
              ) : null}
            </li>
          )
        })}
      </ol>
    </div>
  )
}
