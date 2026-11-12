# BukaLontar API

Experimental Go backend for BukaLontar. It accepts one palm-leaf photograph, classifies its script, and sends accepted Balinese-script images to Gemini 3.8 Flash Low through 9router. It validates JSON and recipe-evidence boundaries, grounds botanical fields against the embedded Taru Pramana KB, and returns the `AnalysisResponse` contract. A structurally valid response is not a verified reading or medical recommendation.

Final submission snapshot: [`v1.1-final-submission`](https://github.com/0xMinomus/bukalontar/tree/v1.1-final-submission). Requires Go 1.26; the full frontend build also requires Node.js 22.12+ (or Node.js 20.19 on the 20.x line).

## Endpoints

| Method | Path | Purpose |
| --- | --- | --- |
| `POST` | `/api/analyze-manuscript` | `multipart/form-data`: `image` (JPG/PNG/WEBP, max 8 MB), optional `source_id` |
| `GET` | `/health` | Process, KB, and 9router reachability |
| `GET` | `/api/health` | Same health payload |
| `GET` | `/` and SPA routes | Embedded `apps/web/dist` when built with `make build` |

## Configuration

Copy `.env.example` to `.env`, then fill the 9router key. `.env` is gitignored; never commit it.

```env
AI_BASE_URL=http://127.0.0.1:20128/v1
AI_MODEL=ag/gemini-3.8-flash-low
AI_REASONING_EFFORT=low
AI_READING_ASSISTANCE=examples
AI_API_KEY=
AI_TIMEOUT_SECONDS=120
RATE_LIMIT_PER_MINUTE=10
```

The server requires `AI_MODEL=ag/gemini-3.8-flash-low` and refuses to start with another model. There is no fallback model.

## Run

```bash
# from repository root
make dev-api
```

The API listens on `:8080`. Vite proxies `/api` there, so the browser stays on one origin.

## Verification

```bash
cd apps/api
go test ./...
go vet ./...
```

Live smoke test:

```bash
curl -sS http://127.0.0.1:8080/health | jq .
curl -sS -F image=@apps/web/public/samples/lontar-usadha-sukawati-depan.jpg \
  -F source_id=usada-sukawati \
  http://127.0.0.1:8080/api/analyze-manuscript | jq .
```

## Design notes

- No third-party Go dependencies; deployment is one binary plus `.env`.
- Image bytes are read into memory; no temp files.
- Unknown image bytes are sent unchanged to the AI gateway; contrast enhancement is preview-only. Exact registered archive bytes send only their supplied transcript. Upstream retention has not been verified.
- Upload is capped at 8 MB before preprocessing and again before the AI call.
- AI output is parsed even when wrapped in a markdown fence, then validated. One repair attempt is made with the validation error.
- Four exact SHA-256 image matches select cited community transcripts from `internal/ai/manuscript_references.json`. Every request still calls Gemini for translation/extraction; the server preserves source Latin and attaches `reading_reference`. This is disclosed reference-assisted analysis, not fresh OCR or a cached generated answer. Catalogued Carcan/Kakawin cannot acquire medicinal recipes. File names and `source_id` never select a reference; generated provenance is discarded.
- Other images retain the target-only script gate, then one reading pass for `balinese/high`. Both calls use Gemini 3.8 Flash Low under the request deadline. Rejection skips reading; stochastic readings are not unioned.
- Recipe quotes must occur in one uninterrupted passage of the first-leaf transliteration and include the plant name and listed adjuvants. Unknown spans, bracketed annotations, and ellipses cannot be bridged or included in the quote. When leaf labels are present, evidence must follow the first label numbered 1; preambles and other leaves do not supply recipe evidence. A complete passage beside an unreadable span remains eligible. This only checks textual consistency: the model can hallucinate both the reading and the quote, or quote a passage without an actual preparation. Botanical scientific fields and literature are grounded against the KB; traditional energy remains as the model read it, not the KB's default.
- The canonical medical disclaimer is always written by the server.
- 9router/upstream calls retry on 429, 5xx, timeout, and transport errors, without switching models.

`AI_READING_ASSISTANCE` accepts `none`, `examples` (default), `chart`, `leaves`, `ocr`, `both`, or `all`. `examples` adds three embedded expert-labelled glyph references; `chart` adds one labelled contact sheet of all 55 published glyph classes; `leaves` adds full-leaf enlargements of the same photograph; `ocr` adds learned text-focus crops. None replaces the original image. Measured Balinese-script error did **not** improve with more images; `ocr`, `leaves`, and `all` remain experimental. Set `none` to reproduce the pre-assistance pipeline for evaluation.

Only Balinese-script targets are in scope. Unknown images use a separate classification call (256 output tokens maximum) without reference glyphs or plant vocabulary. Non-Balinese, mixed, or uncertain decisions skip the reader; malformed decisions fail closed. Reader rejection or conflicting output clears text and recipes. Exact archive identity comes from the registry, not a new visual classification. Kawi language written in Balinese script remains in scope. Experimental Unicode is not shown in UI/PDF.

## Model evaluation

`cmd/evaluate` runs the real analysis service with repeated calls, image SHA-256, latency, output, and errors recorded as JSONL. See the [archive-assisted demonstration report](../../docs/AI_EVALUATION.md#empat-contoh-arsip-berbantuan-transkripsi), separately from the historical photo-only baseline. The final four-photo batch returned 8/8 responses; this does not measure OCR accuracy because reference text is supplied. AI/model/handler tests and the production build passed; four real landing-button analyses returned HTTP 200, and UI/PDF provenance was checked.
