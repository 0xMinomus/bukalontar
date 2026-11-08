# BukaLontar Web

Frontend prototype for BukaLontar: React 19, TypeScript, Vite, Tailwind CSS v4, and shadcn-style primitives (Radix UI + CVA).

## Status

- States implemented: landing, preview, processing, result, error.
- Client preprocessing: downscale, CLAHE, and unsharp mask on Canvas.
- Live AI: Go API at `/api/analyze-manuscript` → 9router/Gemini (set `VITE_USE_MOCK=false`).
- Mock mode remains available with `VITE_USE_MOCK=true` or when no env is set.
- Non-usada leaves show a reading summary instead of empty ethnomedicinal fields; KB-grounded plants keep the usada card.
- Print view renders the full result as a report, not only the active tab.
- Sample lontar photos are real archive images. Credits and licenses live in `public/samples/CREDITS.md`.

## Routes

- `/` landing and about sections
- `/analisis` upload preview, processing, and error states
- `/hasil` analysis result

Routing uses `BrowserRouter`, so URLs are path based, never hash based. The header logo returns to `/` and clears the active manuscript. The footer is shared by every route. The active result is persisted in IndexedDB, so `/hasil` survives a refresh; the logo or `Naskah baru` clears it. A direct `/hasil` or `/analisis` load without a stored session redirects to `/`.

## Run

From the repository root:

```bash
npm install
cp apps/api/.env.example apps/api/.env   # isi AI_API_KEY dari 9router
make dev-api                             # terminal 1: Go API :8080
npm run dev                              # terminal 2: Vite :5173
```

Vite mem-proxy `/api` ke `http://localhost:8080`, jadi browser tetap satu origin saat development.

Production build (frontend + API in one binary):

```bash
make build                               # bin/bukalontar
./bin/bukalontar
```

## Environment

Copy `.env.example` to `.env`:

```env
VITE_USE_MOCK=false
VITE_API_BASE_URL=
```

Set `VITE_USE_MOCK=true` only when working on UI without the Go API. The API key never belongs in this file; it lives in `apps/api/.env`, which is ignored.

## Design direction

Balinese archival reading room, not a generic museum skin: warm parchment, sepia ink, deep lontar-leaf green, and aged brass. Cormorant Garamond carries display type; Lora carries body and interface text. The paper grain, brass rules, Roman numerals, and real lontar photographs are deliberate. Motion is limited to scroll reveals, tab changes, hover lifts, and the draggable/zoomable lens. About sections appear only on the landing state so the result view stays focused.
