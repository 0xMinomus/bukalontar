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

## Link previews

`index.html` includes static Open Graph and Twitter Card metadata so link crawlers do not need JavaScript. The public image is `/og-bukalontar-v1.jpg`: a 1200 × 630 JPEG using the existing logo, local Cormorant Garamond and Lora fonts, parchment palette, and a real lontar photograph. Large title text leads the card; the photograph shows the manuscript rather than a simulated reading result. Source attribution is printed on the image and documented in `public/samples/CREDITS.md`.

Metadata uses absolute URLs on `https://www.bukalontar.my.id`. Deploy the HTML and public image together. Verify that the image URL returns `200` with `Content-Type: image/jpeg`, not the SPA HTML fallback. When changing the artwork, increment its filename and update both image URLs in the metadata; messaging apps may retain old previews. Existing WhatsApp messages may keep their cached preview after deployment.


## Design direction

Balinese lontar reading desk, following the supplied antique parchment reference. Warm paper and sepia preserve the manuscript material; olive marks primary actions, and honey-gold leaves echo the supplied logo. Cormorant Garamond carries display text; Lora carries body text. Four critical WOFF2 assets are preloaded from `src/assets/fonts` in development and hashed by Vite for production. Their OFL licenses accompany the source assets. Font faces use `font-display: optional`: a slow connection keeps the initial fallback for that visit rather than swapping fonts after paint. Latin Extended subsets remain available for transliteration.

The entire lontar surface opens the native image picker by click, Enter, or Space; there are no separate upload or camera buttons. Mobile camera availability follows the device's native image picker. Drag-and-drop and paste retain the existing upload workflow. Leaves keep their settled pose on refresh and respond to hover, drag-over, and pointer movement using transform transitions, without an animation loop.

Initially visible content does not run entrance animations. Below-fold folios transition once on intersection, without a timeout that later hides or moves already-visible content. Navigation keeps native smooth scrolling; wheel and touch scrolling are not intercepted. Reduced motion disables tilt and reveals. The compact sticky header is full-width and flush with the top, without floating margins, rounded corners, or a shadow. It retains the only backdrop blur and a solid reduced-transparency fallback. Result tabs hide native scrollbar arrows while preserving horizontal scrolling and keyboard selection. Analysis, preprocessing, saved sessions, result tabs, copying, lens controls, and full-report printing remain unchanged. Copy covers all Balinese lontar; Usadha remains conditional enrichment.
