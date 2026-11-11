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
AI_READING_ASSISTANCE=both
AI_API_KEY=
AI_TIMEOUT_SECONDS=120
RATE_LIMIT_PER_MINUTE=10
```

For the fallback route, set `AI_FALLBACK_ENABLED=true` and the `AI_FALLBACK_*` values. The fallback uses the same OpenAI-compatible client.

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
- Each request uses one analysis pass, with up to six preparations. Results are not unioned across stochastic readings, and distinct preparations using the same plant are not collapsed by plant name. Every request calls the model; no cached response is served.
- Botanical scientific name, family, compounds, references, reference links, pharmacological action, and energy character are overwritten from the KB on a match. On no match, taxonomy, citations, and model-supplied efficacy claims are cleared. Unknown preparation details remain unspecified, not drinking-medicine defaults.
- The canonical medical disclaimer is always written by the server.
- 9router/upstream calls retry on 429, 5xx, timeout, and transport errors; the optional fallback provider is tried after the primary fails.

`AI_READING_ASSISTANCE` accepts `none`, `examples`, `ocr`, or `both` (default `examples`). `examples` adds three embedded expert-labelled glyph references; `ocr` adds up to two learned text-focus crops of the same photograph. Neither replaces the original image: the model still reads the manuscript photo. `ocr` and `both` remain experimental: measured Balinese-script error did not improve and the focus crops once induced Javanese-script output, which the server now strips. Set `none` to reproduce the pre-assistance pipeline for evaluation.

## Model evaluation

`cmd/evaluate` runs the same service against real images with repeated calls, exact model IDs, image SHA-256, latency, and structured output recorded as JSONL. See [the evaluation report](../../docs/AI_EVALUATION.md) for reproducible commands and the Sol rejection evidence. `AI_REASONING_EFFORT` affects `cx/` models only; it does not fine-tune model weights.
