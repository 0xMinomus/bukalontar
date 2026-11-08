export interface EnhancedImage {
  url: string
  blob: Blob
  width: number
  height: number
  elapsedMs: number
}

const MAX_DIMENSION = 2000
const MAX_TILES = 8
const MIN_TILE_PX = 48
const CLIP_LIMIT = 2.2
const UNSHARP_AMOUNT = 0.55

export async function enhanceLontar(source: Blob): Promise<EnhancedImage> {
  const started = performance.now()
  const bitmap = await createImageBitmap(source)
  const scale = Math.min(
    1,
    MAX_DIMENSION / Math.max(bitmap.width, bitmap.height),
  )
  const width = Math.max(1, Math.round(bitmap.width * scale))
  const height = Math.max(1, Math.round(bitmap.height * scale))

  const canvas = document.createElement('canvas')
  canvas.width = width
  canvas.height = height
  const ctx = canvas.getContext('2d', { willReadFrequently: true })
  if (!ctx) {
    bitmap.close()
    throw new Error('Canvas tidak tersedia di browser ini.')
  }

  ctx.drawImage(bitmap, 0, 0, width, height)
  bitmap.close()

  const image = ctx.getImageData(0, 0, width, height)
  const luminance = new Uint8ClampedArray(width * height)

  for (let i = 0, p = 0; i < image.data.length; i += 4, p += 1) {
    luminance[p] = Math.round(
      0.2126 * image.data[i] +
        0.7152 * image.data[i + 1] +
        0.0722 * image.data[i + 2],
    )
  }

  const equalized = clahe(luminance, width, height)
  const sharpened = unsharpMask(equalized, width, height, UNSHARP_AMOUNT)

  for (let i = 0, p = 0; i < image.data.length; i += 4, p += 1) {
    const value = sharpened[p]
    // A slight warm bias keeps the enhanced view tied to palm leaf instead of
    // looking like a clinical scan.
    image.data[i] = value
    image.data[i + 1] = value
    image.data[i + 2] = Math.round(value * 0.96)
    image.data[i + 3] = 255
  }

  ctx.putImageData(image, 0, 0)

  const { promise, resolve, reject } = Promise.withResolvers<Blob>()
  canvas.toBlob(
    (blob) => {
      if (blob) resolve(blob)
      else reject(new Error('Gagal membuat pratinjau yang ditingkatkan.'))
    },
    'image/jpeg',
    0.9,
  )
  const blob = await promise

  return {
    url: URL.createObjectURL(blob),
    blob,
    width,
    height,
    elapsedMs: Math.round(performance.now() - started),
  }
}

function clahe(
  source: Uint8ClampedArray,
  width: number,
  height: number,
): Uint8ClampedArray {
  const tilesX = Math.max(1, Math.min(MAX_TILES, Math.floor(width / MIN_TILE_PX)))
  const tilesY = Math.max(
    1,
    Math.min(MAX_TILES, Math.floor(height / MIN_TILE_PX)),
  )
  const tileWidth = Math.ceil(width / tilesX)
  const tileHeight = Math.ceil(height / tilesY)
  const luts: Uint8ClampedArray[] = []

  for (let ty = 0; ty < tilesY; ty += 1) {
    for (let tx = 0; tx < tilesX; tx += 1) {
      const x0 = tx * tileWidth
      const y0 = ty * tileHeight
      const x1 = Math.min(width, x0 + tileWidth)
      const y1 = Math.min(height, y0 + tileHeight)
      const histogram = new Uint32Array(256)
      let count = 0

      for (let y = y0; y < y1; y += 1) {
        for (let x = x0; x < x1; x += 1) {
          histogram[source[y * width + x]] += 1
          count += 1
        }
      }

      const clip = Math.max(1, Math.floor((CLIP_LIMIT * count) / 256))
      let excess = 0
      for (let i = 0; i < 256; i += 1) {
        if (histogram[i] > clip) {
          excess += histogram[i] - clip
          histogram[i] = clip
        }
      }

      const share = Math.floor(excess / 256)
      let remainder = excess % 256
      for (let i = 0; i < 256; i += 1) {
        histogram[i] += share
        if (remainder > 0) {
          histogram[i] += 1
          remainder -= 1
        }
      }

      const lut = new Uint8ClampedArray(256)
      let cumulative = 0
      for (let i = 0; i < 256; i += 1) {
        cumulative += histogram[i]
        lut[i] = Math.round((cumulative / count) * 255)
      }
      luts.push(lut)
    }
  }

  const output = new Uint8ClampedArray(width * height)

  for (let y = 0; y < height; y += 1) {
    const gridY = y / tileHeight - 0.5
    const ty0 = clamp(Math.floor(gridY), 0, tilesY - 1)
    const ty1 = clamp(ty0 + 1, 0, tilesY - 1)
    const fy = clamp(gridY - ty0, 0, 1)

    for (let x = 0; x < width; x += 1) {
      const gridX = x / tileWidth - 0.5
      const tx0 = clamp(Math.floor(gridX), 0, tilesX - 1)
      const tx1 = clamp(tx0 + 1, 0, tilesX - 1)
      const fx = clamp(gridX - tx0, 0, 1)
      const value = source[y * width + x]

      const topLeft = luts[ty0 * tilesX + tx0][value]
      const topRight = luts[ty0 * tilesX + tx1][value]
      const bottomLeft = luts[ty1 * tilesX + tx0][value]
      const bottomRight = luts[ty1 * tilesX + tx1][value]

      const top = topLeft + (topRight - topLeft) * fx
      const bottom = bottomLeft + (bottomRight - bottomLeft) * fx
      output[y * width + x] = Math.round(top + (bottom - top) * fy)
    }
  }

  return output
}

function unsharpMask(
  source: Uint8ClampedArray,
  width: number,
  height: number,
  amount: number,
): Uint8ClampedArray {
  const blurred = new Uint8ClampedArray(source.length)
  const output = new Uint8ClampedArray(source.length)

  for (let y = 0; y < height; y += 1) {
    for (let x = 0; x < width; x += 1) {
      let sum = 0
      let neighbours = 0
      for (let dy = -1; dy <= 1; dy += 1) {
        const yy = y + dy
        if (yy < 0 || yy >= height) continue
        for (let dx = -1; dx <= 1; dx += 1) {
          const xx = x + dx
          if (xx < 0 || xx >= width) continue
          sum += source[yy * width + xx]
          neighbours += 1
        }
      }
      blurred[y * width + x] = sum / neighbours
    }
  }

  for (let i = 0; i < source.length; i += 1) {
    output[i] = clamp(
      Math.round(source[i] + amount * (source[i] - blurred[i])),
      0,
      255,
    )
  }

  return output
}

function clamp(value: number, min: number, max: number) {
  return Math.min(max, Math.max(min, value))
}
