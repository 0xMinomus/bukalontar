import { Maximize2, Minus, Plus, RotateCcw } from 'lucide-react'
import * as React from 'react'

import { Button } from '@/components/ui/button'
import {
  ToggleGroup,
  ToggleGroupItem,
} from '@/components/ui/toggle-group'
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip'
import { cn } from '@/lib/utils'

export type LensView = 'original' | 'enhanced'

interface LontarLensProps {
  originalUrl: string
  enhancedUrl?: string | null
  view: LensView
  onViewChange: (view: LensView) => void
  alt: string
  zoomable?: boolean
  className?: string
}

interface DragOrigin {
  clientX: number
  clientY: number
  offsetX: number
  offsetY: number
}

function clamp(value: number, min: number, max: number) {
  return Math.min(max, Math.max(min, value))
}

function maxOffset(size: number, scale: number) {
  return (size * (scale - 1)) / 2
}

export function LontarLens({
  originalUrl,
  enhancedUrl,
  view,
  onViewChange,
  alt,
  zoomable = false,
  className,
}: LontarLensProps) {
  const [zoom, setZoom] = React.useState(1)
  const [offset, setOffset] = React.useState({ x: 0, y: 0 })
  const [dragging, setDragging] = React.useState(false)
  const viewportRef = React.useRef<HTMLDivElement>(null)
  const dragRef = React.useRef<DragOrigin | null>(null)

  const hasEnhanced = Boolean(enhancedUrl)
  const activeUrl =
    view === 'enhanced' && enhancedUrl ? enhancedUrl : originalUrl

  function changeZoom(delta: number) {
    const next = clamp(zoom + delta, 1, 3)
    setZoom(next)
    if (next === 1) {
      setOffset({ x: 0, y: 0 })
      return
    }
    const viewport = viewportRef.current
    if (!viewport) return
    const limitX = maxOffset(viewport.clientWidth, next)
    const limitY = maxOffset(viewport.clientHeight, next)
    setOffset((current) => ({
      x: clamp(current.x, -limitX, limitX),
      y: clamp(current.y, -limitY, limitY),
    }))
  }

  function handlePointerDown(event: React.PointerEvent<HTMLDivElement>) {
    if (!zoomable || zoom <= 1) return
    event.currentTarget.setPointerCapture(event.pointerId)
    dragRef.current = {
      clientX: event.clientX,
      clientY: event.clientY,
      offsetX: offset.x,
      offsetY: offset.y,
    }
    setDragging(true)
  }

  function handlePointerMove(event: React.PointerEvent<HTMLDivElement>) {
    const origin = dragRef.current
    if (!origin) return
    const viewport = viewportRef.current
    if (!viewport) return
    const limitX = maxOffset(viewport.clientWidth, zoom)
    const limitY = maxOffset(viewport.clientHeight, zoom)
    setOffset({
      x: clamp(
        origin.offsetX + event.clientX - origin.clientX,
        -limitX,
        limitX,
      ),
      y: clamp(
        origin.offsetY + event.clientY - origin.clientY,
        -limitY,
        limitY,
      ),
    })
  }

  function handlePointerUp(event: React.PointerEvent<HTMLDivElement>) {
    dragRef.current = null
    setDragging(false)
    if (event.currentTarget.hasPointerCapture(event.pointerId)) {
      event.currentTarget.releasePointerCapture(event.pointerId)
    }
  }

  function handleKeyDown(event: React.KeyboardEvent<HTMLDivElement>) {
    if (!zoomable || zoom <= 1) return
    const step = 28
    const moves: Record<string, [number, number]> = {
      ArrowLeft: [step, 0],
      ArrowRight: [-step, 0],
      ArrowUp: [0, step],
      ArrowDown: [0, -step],
    }
    const move = moves[event.key]
    if (!move) return
    event.preventDefault()
    const viewport = viewportRef.current
    if (!viewport) return
    const limitX = maxOffset(viewport.clientWidth, zoom)
    const limitY = maxOffset(viewport.clientHeight, zoom)
    setOffset((current) => ({
      x: clamp(current.x + move[0], -limitX, limitX),
      y: clamp(current.y + move[1], -limitY, limitY),
    }))
  }

  return (
    <div
      className={cn(
        'overflow-hidden rounded-xl border border-mat bg-mat',
        className,
      )}
    >
      <div className="flex flex-wrap items-center justify-between gap-2 border-b border-white/10 px-2.5 py-2 print:hidden">
        <div className="flex min-w-0 flex-wrap items-center gap-2">
          <span className="hidden px-1 text-[11px] font-semibold text-white/60 sm:inline">
            Lensa lontar
          </span>
          <ToggleGroup
            type="single"
            value={view}
            onValueChange={(value) => {
              if (value === 'original' || value === 'enhanced') {
                onViewChange(value)
              }
            }}
            aria-label="Pilih tampilan citra"
            className="border-white/15 bg-black/30"
          >
            <ToggleGroupItem
              value="original"
              className="h-7 px-2.5 text-white/65 data-[state=on]:bg-white data-[state=on]:text-mat"
            >
              Asli
            </ToggleGroupItem>
            <ToggleGroupItem
              value="enhanced"
              disabled={!hasEnhanced}
              className="h-7 px-2.5 text-white/65 data-[state=on]:bg-white data-[state=on]:text-mat"
            >
              Diperjelas
            </ToggleGroupItem>
          </ToggleGroup>
        </div>

        {zoomable ? (
          <div className="flex items-center gap-0.5">
            <Tooltip>
              <TooltipTrigger asChild>
                <Button
                  type="button"
                  variant="ghost"
                  size="icon"
                  className="size-9 text-white/75 hover:bg-white/10 hover:text-white"
                  onClick={() => changeZoom(-0.5)}
                  disabled={zoom <= 1}
                  aria-label="Perkecil citra"
                >
                  <Minus />
                </Button>
              </TooltipTrigger>
              <TooltipContent>Perkecil</TooltipContent>
            </Tooltip>

            <span className="num w-9 text-center text-[11px] font-semibold text-white/60">
              {zoom.toFixed(1)}x
            </span>

            <Tooltip>
              <TooltipTrigger asChild>
                <Button
                  type="button"
                  variant="ghost"
                  size="icon"
                  className="size-9 text-white/75 hover:bg-white/10 hover:text-white"
                  onClick={() => changeZoom(0.5)}
                  disabled={zoom >= 3}
                  aria-label="Perbesar citra"
                >
                  <Plus />
                </Button>
              </TooltipTrigger>
              <TooltipContent>Perbesar</TooltipContent>
            </Tooltip>

            <Tooltip>
              <TooltipTrigger asChild>
                <Button
                  type="button"
                  variant="ghost"
                  size="icon"
                  className="size-9 text-white/75 hover:bg-white/10 hover:text-white"
                  onClick={() => {
                    setZoom(1)
                    setOffset({ x: 0, y: 0 })
                  }}
                  aria-label="Kembalikan ukuran citra"
                >
                  <RotateCcw />
                </Button>
              </TooltipTrigger>
              <TooltipContent>Reset</TooltipContent>
            </Tooltip>
          </div>
        ) : (
          <span className="hidden items-center gap-1.5 px-1 text-[11px] font-medium text-white/60 sm:inline-flex">
            <Maximize2 className="size-3.5" aria-hidden="true" />
            Pratinjau
          </span>
        )}
      </div>

      <div
        ref={viewportRef}
        className={cn(
          'lens-viewport flex min-h-[220px] select-none items-center justify-center overflow-hidden p-2 sm:min-h-[320px]',
          zoomable &&
            'focus-visible:outline-white/70',
          zoomable &&
            zoom > 1 &&
            'cursor-grab touch-none active:cursor-grabbing',
        )}
        onPointerDown={handlePointerDown}
        onPointerMove={handlePointerMove}
        onPointerUp={handlePointerUp}
        onPointerCancel={handlePointerUp}
        onKeyDown={handleKeyDown}
        tabIndex={zoomable ? 0 : undefined}
        role={zoomable ? 'group' : undefined}
        aria-label={
          zoomable
            ? 'Lensa lontar. Setelah diperbesar, geser dengan tombol panah.'
            : undefined
        }
      >
        <img
          src={activeUrl}
          alt={alt}
          draggable={false}
          className={cn(
            'lens-image max-h-[300px] w-full object-contain sm:max-h-[440px]',
            !dragging && 'transition-transform duration-200 ease-out',
          )}
          style={{
            transform: `translate(${offset.x}px, ${offset.y}px) scale(${zoom})`,
            transformOrigin: 'center',
          }}
        />
      </div>
    </div>
  )
}
