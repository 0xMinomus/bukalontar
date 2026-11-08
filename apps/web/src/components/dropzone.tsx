import { Camera, ImagePlus } from 'lucide-react'
import * as React from 'react'

import { BrandMark } from '@/components/brand-mark'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { cn } from '@/lib/utils'

interface DropzoneProps {
  onFile: (file: File) => void
  disabled?: boolean
}

export function Dropzone({ onFile, disabled }: DropzoneProps) {
  const inputRef = React.useRef<HTMLInputElement>(null)
  const cameraRef = React.useRef<HTMLInputElement>(null)
  const [dragging, setDragging] = React.useState(false)

  function handleFiles(files: FileList | null) {
    const file = files?.[0]
    if (file) onFile(file)
  }

  return (
    <Card
      className={cn(
        'paper-panel frame-antique relative overflow-hidden border-dashed p-5 transition-[border-color,background-color,box-shadow] duration-300 sm:p-6',
        dragging && 'border-brass bg-accent/60',
        disabled && 'opacity-60',
      )}
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
      <div className="flex items-start gap-4">
        <span
          className={cn(
            'flex size-12 shrink-0 items-center justify-center rounded-full border border-brass/40 bg-primary text-primary-foreground shadow-[inset_0_0_0_1px_rgb(255_255_255_/_0.14)] transition-transform duration-300',
            dragging && 'scale-105',
          )}
        >
          <BrandMark className="size-6" />
        </span>
        <div className="min-w-0">
          <p className="small-caps text-[11px] font-semibold text-brass-strong">
            Berkas lontar
          </p>
          <h2 className="mt-1 font-heading text-3xl leading-none tracking-tight sm:text-4xl">
            Tarik foto lontar ke sini
          </h2>
          <p className="mt-2 max-w-[42ch] text-sm leading-relaxed text-muted-foreground">
            Satu lempir atau satu halaman. Pastikan aksara terlihat jelas dan
            seluruh bagian naskah masuk ke frame.
          </p>
        </div>
      </div>

      <div className="mt-5 flex flex-wrap items-center gap-2.5">
        <Button
          type="button"
          onClick={() => inputRef.current?.click()}
          disabled={disabled}
        >
          <ImagePlus aria-hidden="true" />
          Pilih foto lontar
        </Button>
        <Button
          type="button"
          variant="outline"
          className="sm:hidden"
          onClick={() => cameraRef.current?.click()}
          disabled={disabled}
        >
          <Camera aria-hidden="true" />
          Kamera
        </Button>
        <span className="inline-flex items-center gap-1.5 text-xs text-muted-foreground">
          atau tempel dengan
          <kbd className="rounded-sm border border-border bg-secondary px-1.5 py-0.5 font-sans text-[11px] font-semibold text-foreground">
            Ctrl V
          </kbd>
        </span>
      </div>

      <p className="mt-5 border-t border-border/70 pt-3 text-xs text-muted-foreground">
        JPG, PNG, atau WEBP. Maksimal 8 MB. Foto diproses tanpa disimpan.
      </p>

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
      <input
        ref={cameraRef}
        type="file"
        accept="image/*"
        capture="environment"
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
