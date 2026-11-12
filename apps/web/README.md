# BukaLontar Web

Frontend for BukaLontar's experimental palm-leaf reading aid: React 19, TypeScript, Vite, Tailwind CSS v4, and shadcn-style primitives (Radix UI + CVA). Final submission snapshot: [`v1.1-final-submission`](https://github.com/0xMinomus/bukalontar/tree/v1.1-final-submission).

## Status

- States implemented: landing, preview, processing, result, error.
- Canvas downscale, CLAHE, and unsharp mask affect the preview only. Live analysis uploads the original image.
- Live AI: Go API at `/api/analyze-manuscript`, using only `ag/gemini-3.8-flash-low` through 9router.
- Mock mode remains available with `VITE_USE_MOCK=true`. Tanpa env tersebut, frontend memanggil API live.
- Non-usada leaves show a reading summary instead of empty ethnomedicinal fields; KB-grounded plants keep the usada card.
- Print view renders the full result as a report, not only the active tab.
- Four selectable examples: Usadha Sukawati, Usadha Rare Kemenuh, Carcan Kucing, and Kakawin Bharatayuda. Each has separate photo, reference, and license links. Credits and crop hashes are in `public/samples/CREDITS.md`.
- Balinese-script detection and reading quality are separate. Kawi may identify a language written in Balinese script, not the detected writing system. Experimental Unicode is not displayed in the UI/PDF.
- Exact archive bytes use cited community transcription supplied by the backend, with “Berbantuan transkripsi rujukan” shown in result and PDF. These are not fresh OCR readings or cached model answers. Other uploads remain experimental. Quotes must match the canonical source reading; unknown spans, ellipses, and leaf boundaries cannot be bridged.

## Routes

- `/` landing and about sections
- `/analisis` upload preview, processing, and error states
- `/hasil` analysis result

Routing uses `BrowserRouter`, so URLs are path based, never hash based. The header logo returns to `/` and clears the active manuscript. The footer is shared by every route. The active result, original photo, enhanced preview, file name, and view state are persisted in IndexedDB, so `/hasil` survives a refresh; the logo or `Naskah baru` clears the session. A direct `/hasil` or `/analisis` load without a stored session redirects to `/`. Clearing this browser data does not prove deletion at the AI gateway.

Vercel Services must rewrite `/analisis` and `/hasil` to the web service's `/index.html`; forwarding those paths unchanged returns a platform 404 before React can restore the session. The root `vercel.json` defines these rewrites before the general web rule, leaving API, health, and static assets unchanged. After deployment, verify both direct URLs and refresh an existing result. Rollback consists of reverting this routing rule.

## Run

From the repository root:

Use Node.js 22.12+ (or Node.js 20.19 on the 20.x line), npm, and Go 1.26.

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

Main now includes reference-assisted archive examples after the immutable `v1.1-final-submission` snapshot. Live analysis still calls Gemini; `VITE_USE_MOCK` is not evidence. Photo-only accuracy and archive-assisted behavior must be reported separately. Production health does not expose a deployment commit. See [AI evaluation](../../docs/AI_EVALUATION.md).

## Link previews

`index.html` includes static Open Graph and Twitter Card metadata so link crawlers do not need JavaScript. The public image is `/og-bukalontar-v1.jpg`: a 1200 × 630 JPEG using the existing logo, local Cormorant Garamond and Lora fonts, parchment palette, and a real lontar photograph. Large title text leads the card; the photograph shows the manuscript rather than a simulated reading result. Source attribution is printed on the image and documented in `public/samples/CREDITS.md`.

Metadata uses absolute URLs on `https://www.bukalontar.my.id`. Deploy the HTML and public image together. Verify that the image URL returns `200` with `Content-Type: image/jpeg`, not the SPA HTML fallback. When changing the artwork, increment its filename and update both image URLs in the metadata; messaging apps may retain old previews. Existing WhatsApp messages may keep their cached preview after deployment.


## Design direction

Balinese lontar reading desk, following the supplied antique parchment reference. Warm paper and sepia preserve the manuscript material; olive marks primary actions, and honey-gold leaves echo the supplied logo. Cormorant Garamond carries display text; Lora carries body text. Four critical WOFF2 assets are preloaded from `src/assets/fonts` in development and hashed by Vite for production. Their OFL licenses accompany the source assets. Font faces use `font-display: optional`: a slow connection keeps the initial fallback for that visit rather than swapping fonts after paint. Latin Extended subsets remain available for transliteration.

The entire lontar surface opens the native image picker by click, Enter, or Space; there are no separate upload or camera buttons. Mobile camera availability follows the device's native image picker. Drag-and-drop and paste retain the existing upload workflow. Leaves keep their settled pose on refresh and respond to hover, drag-over, and pointer movement using transform transitions, without an animation loop.

Initially visible content does not run entrance animations. Below-fold folios transition once on intersection, without a timeout that later hides or moves already-visible content. Navigation keeps native smooth scrolling; wheel and touch scrolling are not intercepted. Reduced motion disables tilt and reveals. The compact sticky header is full-width and flush with the top, without floating margins, rounded corners, or a shadow. It retains the only backdrop blur and a solid reduced-transparency fallback. Result tabs hide native scrollbar arrows while preserving horizontal scrolling and keyboard selection. Analysis, preprocessing, saved sessions, result tabs, copying, lens controls, and full-report printing remain unchanged. Copy covers all Balinese lontar; Usadha remains conditional enrichment.
