---
name: BukaLontar
description: A Balinese lontar reading desk following the supplied parchment reference.
colors:
  paper: "#efe6cf"
  ink: "#302b1c"
  surface: "#f6eedb"
  primary: "#454c2c"
  primary-foreground: "#fbf6ea"
  muted-foreground: "#6a5b49"
  brass-strong: "#705020"
typography:
  display:
    fontFamily: "Cormorant Garamond, Georgia, serif"
    fontSize: "clamp(3.5rem, 6.5vw, 6rem)"
    fontWeight: 500
    lineHeight: 0.98
    letterSpacing: "-0.035em"
  body:
    fontFamily: "Lora Variable, Georgia, serif"
rounded:
  control: "6px"
  header: "0px"
spacing:
  small: "12px"
  medium: "24px"
components:
  button-primary:
    backgroundColor: "{colors.primary}"
    textColor: "{colors.primary-foreground}"
    rounded: "{rounded.control}"
---

## Overview

The supplied antique reference is interpreted as a Balinese reading desk, not a generic museum landing page. The actual workflow remains a web application: choose a manuscript, inspect it, analyze it, read the result. ENERGY 3 / RHYTHM 2 / MOTION 3.

## Colors

Paper and honey-gold reflect dried lontar material. Olive distinguishes actions from the manuscript surface. Sepia keeps explanatory text readable without competing with the reading itself. The fixed light theme follows the user's parchment reference.

## Typography

Cormorant Garamond brings the fine, engraved character of the reference to headings. Lora carries reading text and controls. Preload critical local WOFF2 faces; keep their source licenses alongside them. Use optional font display so late downloads cannot swap already-painted text. A slow connection may retain the initial fallback for that visit.

## Layout

Center the concise upload instruction on the manuscript stack. Preserve existing navigation targets and routes. The sticky header spans the full viewport width, flush with the top, without floating margins, rounded corners, or a shadow. Its inner content is capped at 1200px. Height is 65px on desktop and 84px on the mobile landing page; navigation becomes a compact second row below 640px. Result tabs retain horizontal scrolling when necessary, but hide native scrollbar chrome and prevent vertical overflow.

## Elevation & Depth

Only the header uses backdrop blur. Soft shadows separate overlapping leaves from paper; the header is separated by a thin bottom rule rather than elevation. Reduced transparency uses a solid header.

## Shapes

Controls retain the existing small radius. The upload leaves use long rounded manuscript edges and binding holes. Supporting folios have slightly irregular corner radii to connect them to the reference without clipping text.

## Components

Use the supplied logo for header, footer, and favicon. The complete lontar surface activates the native image picker, with visible keyboard focus and Enter/Space support. Do not place separate upload or camera buttons inside it. Mobile camera options belong to the native picker. Keep leaves in their settled pose on load; hover, drag-over, and pointer tilt respond through eased CSS transforms without a perpetual loop.

Keep the initial viewport stable. Folios below the initial viewport receive one-time intersection transitions. Native CSS view timelines fan the upload leaves as they leave the viewport, settle archive mounts and the source photograph, and draw the existing engraved ornaments. These gestures embody handling a bound manuscript rather than adding a new interaction. Archive corner mounts and binding slots use the existing brass and paper palette; translucent paper highlights express material, not a second accent. Mouse hover briefly catches light across the upload surface and enlarges archive photographs; nothing loops at rest. Browsers without view timelines retain the settled composition and existing reveals. Never use delayed fade-out/fade-in or a timeout to reanimate visible content. Preserve native smooth anchor scrolling and native wheel/touch behavior. Reduced motion shows every folio immediately and disables spatial animations, scroll timelines, hover movement, and tilt. Existing lens gestures, tabs, copy, print, paste, validation, and sessions remain intact.

## Do's and Don'ts

- Describe general lontar reading first. Describe Usadha as conditional enrichment.
- Keep file restrictions outside the short centered upload instruction.
- Retain source credits and responsible-AI warnings.
- Do not add a GitHub button to the header.
- Do not add animation libraries, continuous particle effects, fabricated statistics, or decorative claims.

## Evidence in the final submission

The `v1.1-final-submission` snapshot retains this visual direction. The two selectable examples are actual Balinese-script palm-leaf photographs, not typeset script samples. Original-image analysis and preview enhancement are distinct. Preserve the AI warning and the distinction between detected script and reading quality; a Kawi language label does not mean Kawi script. Do not expose experimental Unicode as a verified transcription or style an AI source quote as expert approval. Results remain reading aids, including when KB taxonomy and literature links are available.

After that immutable tag, main has four archive choices: two Usadha, Carcan, and Kakawin. Preserve the separate photo/reference/license links. Exact-byte archive analysis must visibly say “Berbantuan transkripsi rujukan” in both result and print; never present its source Latin as new OCR or use it as an independent accuracy demonstration. Other photographs remain experimental.
