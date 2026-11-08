import { ImagePlus } from 'lucide-react'
import * as React from 'react'

import { Card } from '@/components/ui/card'
import { cn } from '@/lib/utils'

interface DropzoneProps {
  onFile: (file: File) => void
  disabled?: boolean
}

export function Dropzone({ onFile, disabled }: DropzoneProps) {
  const inputRef = React.useRef<HTMLInputElement>(null)
  const [dragging, setDragging] = React.useState(false)
  const stackRef = React.useRef<HTMLDivElement>(null)

  function tiltLeaves(event: React.PointerEvent<HTMLDivElement>) {
    if (
      event.pointerType !== 'mouse' ||
      window.matchMedia('(prefers-reduced-motion: reduce)').matches
    ) return
    const bounds = event.currentTarget.getBoundingClientRect()
    const x = (event.clientX - bounds.left) / bounds.width - 0.5
    const y = (event.clientY - bounds.top) / bounds.height - 0.5
    if (stackRef.current) {
      stackRef.current.style.transform = `perspective(1100px) rotateX(${-y * 3}deg) rotateY(${x * 4}deg)`
    }
  }
  function handleFiles(files: FileList | null) {
    const file = files?.[0]
    if (file) onFile(file)
  }

  return (
    <Card
      className={cn(
        'lontar-dropzone relative',
        dragging && 'is-dragging',
        disabled && 'opacity-60',
      )}
      onPointerMove={tiltLeaves}
      onPointerLeave={() => {
        if (stackRef.current) stackRef.current.style.transform = ''
      }}
      onDragOver={(event) => {
        event.preventDefault()
        if (!disabled) setDragging(true)
      }}
      onDragLeave={(event) => {
        if (!event.currentTarget.contains(event.relatedTarget as Node)) {
          setDragging(false)
        }
      }}
      onDrop={(event) => {
        event.preventDefault()
        setDragging(false)
        if (!disabled) handleFiles(event.dataTransfer.files)
      }}
    >
      <button
        type="button"
        className="lontar-upload-trigger"
        aria-label="Tambahkan foto lontar"
        disabled={disabled}
        onClick={() => inputRef.current?.click()}
      />
      <div ref={stackRef} className="lontar-stack" aria-hidden="true">
        <span className="lontar-lempir lempir-back" />
        <span className="lontar-lempir lempir-middle" />
        <span className="lontar-lempir lempir-front" />
        <span className="lontar-binding binding-left" />
        <span className="lontar-binding binding-right" />
      </div>
      <div className="lontar-upload-content">
        <ImagePlus className="mx-auto size-5 text-brass-strong" aria-hidden="true" />
        <h2 className="mt-3 font-heading text-[clamp(1.75rem,5vw,2.5rem)] leading-tight" aria-live="polite">
          {dragging ? 'Lepaskan foto di sini' : 'Klik atau taruh foto'}
        </h2>
      </div>

      <input
        ref={inputRef}
        type="file"
        accept="image/jpeg,image/png,image/webp"
        tabIndex={-1}
        className="hidden"
        onChange={(event) => {
          handleFiles(event.target.files)
          event.target.value = ''
        }}
        disabled={disabled}
      />
    </Card>
  )
}
