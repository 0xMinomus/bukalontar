import { clsx, type ClassValue } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function formatBytes(bytes: number) {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(0)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}

export function sleep(ms: number, signal?: AbortSignal) {
  const { promise, resolve, reject } = Promise.withResolvers<void>()
  if (signal?.aborted) {
    reject(new DOMException('Aborted', 'AbortError'))
    return promise
  }
  const timer = window.setTimeout(resolve, ms)
  signal?.addEventListener(
    'abort',
    () => {
      window.clearTimeout(timer)
      reject(new DOMException('Aborted', 'AbortError'))
    },
    { once: true },
  )
  return promise
}
