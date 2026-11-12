# Product

<!-- impeccable:product-schema 1 -->

## Platform

web

## Product Purpose

BukaLontar is an experimental aid for reading photographs of Balinese-script lontar through proposed Latin transliteration, Indonesian and English translations, manuscript context, and a printable report. It does not replace a philologist. Source: ../../docs/PROJECT_CONCEPT.md. Final submission snapshot: `v1.1-final-submission`.

## Users

Students, manuscript researchers, curators, and readers of Balinese manuscripts. Usadha researchers can inspect extracted plant recipes and their available references.

## Capabilities and Constraints

Preserve the existing routes, upload, paste, sample selection, preview enhancement, lens controls, analysis, session restoration, result tabs, copying, and printing. UI-only changes must not silently change the API contract or AI workflow. Live analysis sends the original image, uses only Gemini 3.8 Flash Low, and displays uncertainty. Usadha recipes remain automatic and conditional on backend quotation checks, not the purpose of every manuscript. Kawi language written in Balinese script remains in scope; non-Balinese writing systems are outside the intended scope, but rejection is not guaranteed.

## Brand Commitments

Use the supplied ../../logos/bukalontar-icon.svg logo. The user requests a Balinese interpretation of the supplied antique parchment reference, with a central interactive lontar upload surface, concise centered upload instructions, lightweight motion, and a sticky translucent header without a GitHub button.

## Evidence on Hand

Product scope: ../../docs/PROJECT_CONCEPT.md and ../../README.md. Supplied visual reference: ../../_.jpeg. Existing archival sample photographs retain their Wikimedia attribution and CC BY-SA license. AI output remains a reading aid, not a verified edition or medical advice.

The final code upgrade improves recipe-evidence boundaries, not proven OCR accuracy. Full-page sequence OCR and a second visual verifier were not integrated. Use ../../docs/AI_EVALUATION.md for measured outcomes; do not present model confidence, a matching AI quote, or an HTTP 200 response as a verified reading.
