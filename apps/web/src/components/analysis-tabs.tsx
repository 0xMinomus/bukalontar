import { BookOpen, Check, Copy, Leaf, ShieldCheck } from 'lucide-react'
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
import {
  analysisSummary,
  formatTransliteration,
  formatTranslation,
  hasMedicinalPlant,
  resultRecipes,
  safeReferenceLinks,
  usefulText,
} from '@/lib/reading'

const confidenceLabel: Record<Confidence, string> = {
  low: 'Keterbacaan rendah',
  medium: 'Keterbacaan sedang',
  high: 'Keterbacaan tinggi',
}

function unique(values: string[]): string[] {
  return Array.from(
    new Set(values.map((value) => value.trim()).filter(Boolean)),
  )
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
  const recipes = resultRecipes(result)
  const plantRecipes = recipes.filter((recipe) =>
    hasMedicinalPlant(recipe.plant_local_name),
  )
  const noPlant = plantRecipes.length === 0
  const hasAilment =
    usefulText(result.ethnomedicinal_entity.target_ailment) ||
    recipes.some((recipe) => usefulText(recipe.target_ailment))
  const transliteration =
    formatTransliteration(result.paleography.transliteration) ||
    'Teks aksara belum terbaca.'
  const indonesian = usefulText(result.paleography.indonesian_translation)
    ? formatTranslation(result.paleography.indonesian_translation)
    : 'Terjemahan belum tersedia untuk lembar ini.'
  const english = usefulText(result.paleography.english_translation)
    ? formatTranslation(result.paleography.english_translation)
    : 'Translation is not available for this leaf.'
  const summary = analysisSummary(result)
  const compounds = unique(
    plantRecipes.flatMap((recipe) => recipe.phytochemical_compounds),
  )
  const references = unique(
    plantRecipes.flatMap((recipe) => recipe.reference_studies),
  )
  const actions = plantRecipes.filter((recipe) =>
    usefulText(recipe.modern_pharmacological_action),
  )

  return (
    <Tabs defaultValue="transliterasi" className="w-full">
      <TabsList aria-label="Bagian hasil analisis">
        <TabsTrigger value="transliterasi">Transliterasi</TabsTrigger>
        <TabsTrigger value="terjemahan">Terjemahan</TabsTrigger>
        <TabsTrigger value="kartu">
          {noPlant ? 'Catatan naskah' : 'Kartu usada'}
        </TabsTrigger>
        <TabsTrigger value="validasi">Rujukan</TabsTrigger>
      </TabsList>

      <TabsContent value="transliterasi">
        <div className="mt-4 max-h-[470px] overflow-y-auto pr-1 [scrollbar-width:thin]">
          <div className="flex items-center justify-between gap-3">
            <Badge variant="outline">
              {result.paleography.aksara_bali_detected
                ? 'Aksara Bali terdeteksi'
                : 'Aksara Bali tidak terdeteksi'}
            </Badge>
            <CopyButton text={transliteration} />
          </div>
          <p className="prose-lontar mt-4 whitespace-pre-line">
            {transliteration}
          </p>
          <div className="mt-4 flex flex-wrap items-center gap-x-3 gap-y-1 text-xs text-muted-foreground">
            <span>{confidenceLabel[result.paleography.confidence]}</span>
            <span aria-hidden="true">·</span>
            <span>Bahasa teks: {result.source_reference.language}</span>
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
              <CopyButton text={indonesian} />
            </div>
            <p className="mt-2 whitespace-pre-line text-[15px] leading-7 text-foreground">
              {indonesian}
            </p>
          </section>

          <Separator className="my-5" />

          <section>
            <div className="flex items-center justify-between gap-3">
              <h3 className="text-sm font-semibold">English</h3>
              <CopyButton text={english} />
            </div>
            <p className="mt-2 whitespace-pre-line text-[15px] leading-7 text-muted-foreground">
              {english}
            </p>
          </section>
        </div>
      </TabsContent>

      <TabsContent value="kartu">
        <div className="mt-4 max-h-[470px] overflow-y-auto pr-1 [scrollbar-width:thin]">
          {noPlant ? (
            <div>
              <div className="flex items-start gap-3 rounded-lg border border-border bg-secondary/40 p-4">
                <BookOpen
                  className="mt-0.5 size-4 shrink-0 text-primary"
                  aria-hidden="true"
                />
                <div>
                  <h3 className="font-heading text-2xl leading-tight tracking-tight">
                    {result.source_reference.lontar_group === 'Non-Usadha'
                      ? 'Naskah non-Usadha'
                      : 'Resep belum teridentifikasi'}
                  </h3>
                  <p className="mt-1 text-sm leading-relaxed text-muted-foreground">
                    {hasAilment
                      ? `Hasil pembacaan menyebut ${result.ethnomedicinal_entity.target_ailment}, tetapi belum ada resep tanaman yang berhasil dipetakan.`
                      : 'Belum ada resep yang berhasil diidentifikasi. Bagian yang tidak terbaca tidak membuktikan bahwa naskah tidak memuat resep.'}
                  </p>
                </div>
              </div>
              <p className="mt-4 text-sm leading-relaxed">
                Periksa bacaan dan foto asli sebelum menyimpulkan isi naskah.
              </p>
              {summary ? (
                <div className="mt-4 border-t border-border pt-4">
                  <p className="label-xs">Ringkasan isi</p>
                  <p className="mt-1.5 whitespace-pre-line text-sm leading-7">
                    {summary}
                  </p>
                </div>
              ) : null}
            </div>
          ) : (
            <div>
              {plantRecipes.map((recipe, index) => {
                const taxonomy = [
                  recipe.plant_scientific_name,
                  recipe.botanical_family,
                ].filter(usefulText)
                return (
                  <div
                    key={`${recipe.plant_local_name}-${index}`}
                    className={
                      index > 0 ? 'mt-6 border-t border-border pt-5' : ''
                    }
                  >
                    {plantRecipes.length > 1 ? (
                      <p className="label-xs mb-2">Resep {index + 1}</p>
                    ) : null}
                    <div className="flex items-start justify-between gap-3 border-b border-border pb-4">
                      <div>
                        <h3 className="font-heading text-2xl leading-tight tracking-tight">
                          {recipe.plant_local_name}
                        </h3>
                        <p className="mt-1 text-sm italic text-muted-foreground">
                          {taxonomy.length > 0
                            ? taxonomy.join(' · ')
                            : 'Belum tercatat di Taru Pramana'}
                        </p>
                      </div>
                      <Badge
                        variant={recipe.kb_match ? 'success' : 'outline'}
                        className="shrink-0"
                      >
                        <Leaf aria-hidden="true" className="size-3" />
                        {recipe.kb_match ? 'Tercatat' : 'Belum tercatat'}
                      </Badge>
                    </div>

                    <dl className="mt-4 grid gap-4 sm:grid-cols-2">
                      <div>
                        <dt className="label-xs">Bagian dipakai</dt>
                        <dd className="mt-1 text-sm font-medium">
                          {recipe.plant_part_used}
                        </dd>
                      </div>
                      <div>
                        <dt className="label-xs">Sifat tradisional</dt>
                        <dd className="mt-1 text-sm font-medium">
                          {recipe.energy_character}
                        </dd>
                      </div>
                      <div>
                        <dt className="label-xs">Bentuk racikan</dt>
                        <dd className="mt-1 text-sm font-medium">
                          {recipe.formulation_type}
                        </dd>
                      </div>
                      <div>
                        <dt className="label-xs">Indikasi</dt>
                        <dd className="mt-1 text-sm font-medium">
                          {recipe.target_ailment}
                        </dd>
                      </div>
                    </dl>

                    <div className="mt-5 border-t border-border pt-4">
                      <p className="label-xs">Cara racik</p>
                      {usefulText(recipe.preparation_instructions) ? (
                        <p className="mt-1.5 whitespace-pre-line text-sm leading-7">
                          {formatTranslation(recipe.preparation_instructions)}
                        </p>
                      ) : (
                        <p className="mt-1.5 text-sm text-muted-foreground">
                          Tidak ada cara racik yang tercatat.
                        </p>
                      )}
                    </div>

                    <div className="mt-4">
                      <p className="label-xs">Bahan pelengkap</p>
                      <p className="mt-1.5 text-sm">
                        {recipe.adjuvants.length > 0
                          ? recipe.adjuvants.join(' · ')
                          : 'Tidak ada'}
                      </p>
                    </div>
                    {safeReferenceLinks(recipe).length > 0 ? (
                      <div className="mt-4">
                        <p className="label-xs">Tautan dari knowledge base</p>
                        <ul className="mt-1.5 space-y-1 text-sm">
                          {safeReferenceLinks(recipe).map((link) => (
                            <li key={link.url}>
                              <a className="underline underline-offset-2" href={link.url} target="_blank" rel="noopener noreferrer">{link.title}</a>
                            </li>
                          ))}
                        </ul>
                        <p className="mt-1 text-xs text-muted-foreground">Tautan ini bukan sumber naskah dan tidak membuktikan khasiat resep.</p>
                      </div>
                    ) : null}
                  </div>
                )
              })}
            </div>
          )}
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
              {noPlant
                ? 'Lembar ini tidak memuat tanaman obat, jadi belum ada rujukan tanaman yang bisa ditampilkan.'
                : `${plantRecipes.length} resep ditemukan. Nama dan cara racik berasal dari pembacaan naskah; informasi tambahan berasal dari knowledge base bila tanaman tercatat.`}
            </p>
          </div>

          {!noPlant ? (
            <>
              <section className="mt-4">
                <p className="label-xs">Kandungan aktif</p>
                <p className="mt-1.5 text-sm font-medium">
                  {compounds.length > 0
                    ? compounds.join(' · ')
                    : 'Belum tersedia.'}
                </p>
              </section>

              <Separator className="my-4" />

              <section>
                <p className="label-xs">Khasiat modern</p>
                {actions.length > 0 ? (
                  <ul className="mt-1.5 space-y-2">
                    {actions.map((recipe) => (
                      <li
                        key={recipe.plant_local_name}
                        className="whitespace-pre-line text-sm leading-7"
                      >
                        <span className="font-semibold">
                          {recipe.plant_local_name}:
                        </span>{' '}
                        {formatTranslation(recipe.modern_pharmacological_action)}
                      </li>
                    ))}
                  </ul>
                ) : (
                  <p className="mt-1.5 text-sm text-muted-foreground">
                    Belum tersedia.
                  </p>
                )}
              </section>

              <Separator className="my-4" />

              <section>
                <p className="label-xs">Deskripsi rujukan di knowledge base</p>
                {references.length > 0 ? (
                  <ul className="mt-1.5 space-y-1.5">
                    {references.map((reference) => (
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
                ) : (
                  <p className="mt-1.5 text-sm text-muted-foreground">
                    Belum ada rujukan.
                  </p>
                )}
                {plantRecipes.map((recipe) => {
                  const links = safeReferenceLinks(recipe)
                  return links.length > 0 ? (
                    <div className="mt-4" key={recipe.plant_local_name}>
                      <p className="text-sm font-semibold">{recipe.plant_local_name}: tautan dari knowledge base</p>
                      <ul className="mt-1.5 space-y-1.5">
                        {links.map((link) => (
                          <li key={link.url} className="text-sm">
                            <a className="underline underline-offset-2" href={link.url} target="_blank" rel="noopener noreferrer">{link.title}</a>
                          </li>
                        ))}
                      </ul>
                    </div>
                  ) : null
                })}
                <p className="mt-3 text-xs text-muted-foreground">Tautan literatur bukan sumber naskah; hasil pencarian tidak membuktikan khasiat resep.</p>
              </section>
            </>
          ) : null}
        </div>
      </TabsContent>
    </Tabs>
  )
}
