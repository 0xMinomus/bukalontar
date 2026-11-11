# BukaLontar API

Go stdlib backend for the BukaLontar web app. It accepts one Balinese lontar image, sends it to a multimodal model through 9router in OpenAI-compatible JSON mode, validates the result, grounds botanical fields against the embedded Taru Pramana KB, and returns the `AnalysisResponse` contract. Non-Usadha lontar images return a reading and translation; Usadha leaves may also return recipe cards.

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
curl -sS -F image=@apps/web/public/samples/lontar-wanagiri-usada-3.jpg \
  -F source_id=wanagiri \
  http://127.0.0.1:8080/api/analyze-manuscript | jq .
```

## Design notes

- No third-party Go dependencies; deployment is one binary plus `.env`.
- Image bytes are read into memory; no temp files.
- Upload is capped at 8 MB before preprocessing and again before the AI call.
- AI output is parsed even when wrapped in a markdown fence, then validated. One repair attempt is made with the validation error.
- Each request first classifies only the target image, with no references or plant catalog. Only `balinese/high` reaches one reading pass, with up to six preparations. Both calls use Gemini 3.8 Flash Low under the request deadline; rejection skips the reader. Results are not unioned across stochastic readings, and distinct preparations using the same plant are not collapsed by plant name. Every request calls the model; no cached response is served.
- Botanical scientific name, family, compounds, references, reference links, pharmacological action, and energy character are overwritten from the KB on a match. On no match, taxonomy, citations, and model-supplied efficacy claims are cleared. Unknown preparation details remain unspecified, not drinking-medicine defaults.
- The canonical medical disclaimer is always written by the server.
- 9router/upstream calls retry on 429, 5xx, timeout, and transport errors, without switching models.

`AI_READING_ASSISTANCE` accepts `none`, `examples` (default), `chart`, `leaves`, `ocr`, `both`, or `all`. `examples` adds three embedded expert-labelled glyph references; `chart` adds one labelled contact sheet of all 55 published glyph classes; `leaves` adds full-leaf enlargements of the same photograph; `ocr` adds learned text-focus crops. None replaces the original image. Measured Balinese-script error did **not** improve with more images; `ocr`, `leaves`, and `all` remain experimental. Set `none` to reproduce the pre-assistance pipeline for evaluation.

Only clearly Balinese-script targets are in scope. A separate classification call (256 output tokens maximum) excludes reference glyphs and plant vocabulary. Non-Balinese, mixed, or uncertain decisions skip the reader; malformed decisions fail closed with an error. The server clears all readings, translations, provenance, and recipes on reader rejection or a conflicting script in its transcription/transliteration. Both stages still depend on Gemini, not an independently trained image classifier. Kawi language written in Balinese script remains in scope. Experimental Unicode transcription is retained in the API, not shown in the UI/PDF because measured reading accuracy is inadequate.

## Model evaluation

`cmd/evaluate` runs the same service against real images with repeated calls, exact model IDs, image SHA-256, latency, and structured output recorded as JSONL. See [the evaluation report](../../docs/AI_EVALUATION.md) for reproducible commands and the Sol rejection evidence. `AI_REASONING_EFFORT` affects `cx/` models only; it does not fine-tune model weights.
