# UsadhaAI API

Go stdlib backend for the UsadhaAI web app. It accepts one lontar image, sends it to Gemini through 9router in OpenAI-compatible JSON mode, validates the result, grounds botanical fields against the embedded Taru Pramana KB, and returns the `AnalysisResponse` contract.

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
AI_API_KEY=...
AI_TIMEOUT_SECONDS=120
RATE_LIMIT_PER_MINUTE=60
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
- Each request runs two independent flash-low analysis passes. If they disagree or neither finds a plant, one focused recipe pass runs before returning. Every request calls the model; no cached response is served.
- Botanical scientific name, family, active compounds, references, and energy character are overwritten from the KB on a match. On no match, taxonomy and citations are cleared rather than hallucinated.
- The canonical medical disclaimer is always written by the server.
- 9router/upstream calls retry on 429, 5xx, timeout, and transport errors; the optional fallback provider is tried after the primary fails.
