# UsadhaAI Web

Frontend prototype for UsadhaAI: React 19, TypeScript, Vite, Tailwind CSS v4, and shadcn-style primitives (Radix UI + CVA).

## Status

- States implemented: landing, preview, processing, result, error.
- Client preprocessing: downscale, CLAHE, and unsharp mask on Canvas.
- AI responses are mocked until the Go backend is connected (`VITE_USE_MOCK=true`).
- Sample lontar photos are real archive images. Credits and licenses live in `public/samples/CREDITS.md`.

## Run

From the repository root:

```bash
npm install
npm run dev
```

The app runs at `http://localhost:5173`. The Vite dev server proxies `/api` to `http://localhost:8080` for the future Go backend.

Production build:

```bash
npm run build
```

## Environment

Copy `.env.example` to `.env`:

```env
VITE_USE_MOCK=true
VITE_API_BASE_URL=
```

Set `VITE_USE_MOCK=false` only after `/api/analyze-manuscript` exists.

## Design direction

Balinese archival reading room, not a generic museum skin: warm parchment, sepia ink, deep lontar-leaf green, and aged brass. Cormorant Garamond carries display type; Lora carries body and interface text. The paper grain, brass rules, Roman numerals, and real lontar photographs are deliberate. Motion is limited to scroll reveals, tab changes, hover lifts, and the draggable/zoomable lens. About sections appear only on the landing state so the result view stays focused.
