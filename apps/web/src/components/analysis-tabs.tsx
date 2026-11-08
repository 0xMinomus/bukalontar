import { Check, Copy, Leaf, ShieldCheck } from 'lucide-react'
import * as React from 'react'

import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from '@/components/ui/tabs'
import type { AnalysisResponse, Confidence } from '@/lib/analysis'

const confidenceLabel: Record<Confidence, string> = {
  low: 'Keyakinan rendah',
  medium: 'Keyakinan sedang',
  high: 'Keyakinan tinggi',
}

function CopyButton({ text, label = 'Salin' }: { text: string; label?: string }) {
  const [state, setState] = React.useState<'idle' | 'copied' | 'error'>('idle')

  async function copy() {
    try {
      await navigator.clipboard.writeText(text)
      setState('copied')
      window.setTimeout(() => setState('idle'), 1600)
    } catch {
      setState('error')
    }
  }

  return (
    <Button
      type="button"
      variant="ghost"
      size="sm"
      onClick={copy}
      aria-label={`${label} teks`}
      className="h-9 shrink-0 px-2 text-xs print:hidden"
    >
      {state === 'copied' ? (
        <Check aria-hidden="true" />
      ) : (
        <Copy aria-hidden="true" />
      )}
      {state === 'copied' ? 'Tersalin' : state === 'error' ? 'Gagal' : label}
    </Button>
  )
}

export function AnalysisTabs({ result }: { result: AnalysisResponse }) {
  const entity = result.ethnomedicinal_entity
  const validation = result.scientific_validation

  return (
    <Tabs defaultValue="transliterasi" className="w-full">
      <TabsList aria-label="Bagian hasil analisis">
        <TabsTrigger value="transliterasi">Transliterasi</TabsTrigger>
        <TabsTrigger value="terjemahan">Terjemahan</TabsTrigger>
        <TabsTrigger value="kartu">Kartu usada</TabsTrigger>
        <TabsTrigger value="validasi">Validasi</TabsTrigger>
      </TabsList>

      <TabsContent value="transliterasi">
        <div className="mt-4 max-h-[470px] overflow-y-auto pr-1 [scrollbar-width:thin]">
          <div className="flex items-center justify-between gap-3">
            <Badge variant="outline">Aksara Bali terdeteksi</Badge>
            <CopyButton text={result.paleography.transliteration} />
          </div>
          <p className="prose-lontar mt-4">
            {result.paleography.transliteration}
          </p>
          <div className="mt-4 flex flex-wrap items-center gap-x-3 gap-y-1 text-xs text-muted-foreground">
            <span>{confidenceLabel[result.paleography.confidence]}</span>
            <span aria-hidden="true">·</span>
            <span>{result.source_reference.language}</span>
            <span aria-hidden="true">·</span>
            <span>
              {result.source_reference.lontar_group}, lempir{' '}
              {result.source_reference.lempir_number ?? 'tidak tercatat'}
            </span>
          </div>
        </div>
      </TabsContent>

      <TabsContent value="terjemahan">
        <div className="mt-4 max-h-[470px] overflow-y-auto pr-1 [scrollbar-width:thin]">
          <section>
            <div className="flex items-center justify-between gap-3">
              <h3 className="text-sm font-semibold">Indonesia</h3>
              <CopyButton text={result.paleography.indonesian_translation} />
            </div>
            <p className="mt-2 text-sm leading-relaxed text-foreground">
              {result.paleography.indonesian_translation}
            </p>
          </section>

          <Separator className="my-5" />

          <section>
            <div className="flex items-center justify-between gap-3">
              <h3 className="text-sm font-semibold">English</h3>
              <CopyButton text={result.paleography.english_translation} />
            </div>
            <p className="mt-2 text-sm leading-relaxed text-muted-foreground">
              {result.paleography.english_translation}
            </p>
          </section>
        </div>
      </TabsContent>

      <TabsContent value="kartu">
        <div className="mt-4 max-h-[470px] overflow-y-auto pr-1 [scrollbar-width:thin]">
          <div className="flex items-start justify-between gap-3 border-b border-border pb-4">
            <div>
              <h3 className="font-heading text-2xl leading-tight tracking-tight">
                {entity.plant_local_name}
              </h3>
              <p className="mt-1 text-sm italic text-muted-foreground">
                {entity.plant_scientific_name} · {entity.botanical_family}
              </p>
            </div>
            <Badge
              variant={entity.kb_match ? 'success' : 'outline'}
              className="shrink-0"
            >
              <Leaf aria-hidden="true" className="size-3" />
              {entity.kb_match ? 'Cocok KB' : 'Tidak cocok'}
            </Badge>
          </div>

          <dl className="mt-4 grid gap-4 sm:grid-cols-2">
            <div>
              <dt className="label-xs">Bagian dipakai</dt>
              <dd className="mt-1 text-sm font-medium">
                {entity.plant_part_used}
              </dd>
            </div>
            <div>
              <dt className="label-xs">Sifat tradisional</dt>
              <dd className="mt-1 text-sm font-medium">
                {entity.energy_character}
              </dd>
            </div>
            <div>
              <dt className="label-xs">Bentuk racikan</dt>
              <dd className="mt-1 text-sm font-medium">
                {entity.formulation_type}
              </dd>
            </div>
            <div>
              <dt className="label-xs">Indikasi</dt>
              <dd className="mt-1 text-sm font-medium">
                {entity.target_ailment}
              </dd>
            </div>
          </dl>

          <div className="mt-5 border-t border-border pt-4">
            <p className="label-xs">Cara racik</p>
            <p className="mt-1.5 text-sm leading-relaxed">
              {entity.preparation_instructions}
            </p>
          </div>

          <div className="mt-4">
            <p className="label-xs">Bahan pelengkap</p>
            <p className="mt-1.5 text-sm">
              {entity.adjuvants.length > 0
                ? entity.adjuvants.join(' · ')
                : 'Tidak ada'}
            </p>
          </div>
        </div>
      </TabsContent>

      <TabsContent value="validasi">
        <div className="mt-4 max-h-[470px] overflow-y-auto pr-1 [scrollbar-width:thin]">
          <div className="flex items-start gap-3 rounded-lg border border-border bg-secondary/50 p-3">
            <ShieldCheck
              className="mt-0.5 size-4 shrink-0 text-primary"
              aria-hidden="true"
            />
            <p className="text-xs leading-relaxed text-muted-foreground">
              Senyawa dan aksi di bawah berasal dari basis data Taru Pramana,
              bukan tebakan model.
            </p>
          </div>

          <section className="mt-4">
            <p className="label-xs">Senyawa aktif</p>
            <p className="mt-1.5 text-sm font-medium">
              {validation.phytochemical_compounds.join(' · ')}
            </p>
          </section>

          <Separator className="my-4" />

          <section>
            <p className="label-xs">Aksi farmakologis</p>
            <p className="mt-1.5 text-sm leading-relaxed">
              {validation.modern_pharmacological_action}
            </p>
          </section>

          <Separator className="my-4" />

          <section>
            <p className="label-xs">Rujukan</p>
            <ul className="mt-1.5 space-y-1.5">
              {validation.reference_studies.map((reference) => (
                <li
                  key={reference}
                  className="flex gap-2 text-sm leading-relaxed text-muted-foreground"
                >
                  <span
                    className="mt-2 size-1 shrink-0 rounded-full bg-brass"
                    aria-hidden="true"
                  />
                  {reference}
                </li>
              ))}
            </ul>
          </section>
        </div>
      </TabsContent>
    </Tabs>
  )
}
