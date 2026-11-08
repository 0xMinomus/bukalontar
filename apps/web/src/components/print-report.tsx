import type { AnalysisResponse, Confidence } from '@/lib/analysis'
import {
  analysisSummary,
  formatTransliteration,
  formatTranslation,
  hasMedicinalPlant,
  resultRecipes,
  usefulText,
} from '@/lib/reading'

const confidenceLabel: Record<Confidence, string> = {
  low: 'Rendah',
  medium: 'Sedang',
  high: 'Tinggi',
}

interface PrintReportProps {
  result: AnalysisResponse
  originalUrl: string
}

function plain(value: string, fallback = 'Tidak tercatat.') {
  return usefulText(value) ? value : fallback
}

function unique(values: string[]): string[] {
  return Array.from(
    new Set(values.map((value) => value.trim()).filter(Boolean)),
  )
}

export function PrintReport({ result, originalUrl }: PrintReportProps) {
  const source = result.source_reference
  const paleography = result.paleography
  const entity = result.ethnomedicinal_entity
  const validation = result.scientific_validation
  const recipes = resultRecipes(result)
  const plantRecipes = recipes.filter((recipe) =>
    hasMedicinalPlant(recipe.plant_local_name),
  )
  const hasPlant = plantRecipes.length > 0
  const hasAilment =
    usefulText(entity.target_ailment) ||
    recipes.some((recipe) => usefulText(recipe.target_ailment))
  const transliteration =
    formatTransliteration(paleography.transliteration) ||
    'Teks aksara belum terbaca.'
  const indonesian = usefulText(paleography.indonesian_translation)
    ? formatTranslation(paleography.indonesian_translation)
    : 'Terjemahan belum tersedia untuk lembar ini.'
  const english = usefulText(paleography.english_translation)
    ? formatTranslation(paleography.english_translation)
    : 'Translation is not available for this leaf.'
  const summary = hasPlant
    ? plantRecipes.length > 1
      ? `Ditemukan ${plantRecipes.length} resep pada lempir ini. Rincian setiap resep ada di bagian Kartu Usada.`
      : [
          `${plantRecipes[0].plant_local_name} digunakan untuk ${plain(plantRecipes[0].target_ailment)}`,
          plain(
            plantRecipes[0].preparation_instructions,
            'Cara racik belum tercatat.',
          ),
        ].join('\n\n')
    : analysisSummary(result)
  const taxonomy = [
    plantRecipes[0]?.plant_scientific_name,
    plantRecipes[0]?.botanical_family,
  ].filter(usefulText)
  const compounds = unique(
    plantRecipes.flatMap((recipe) => recipe.phytochemical_compounds),
  )
  const references = unique(
    plantRecipes.flatMap((recipe) => recipe.reference_studies),
  )
  const actions = plantRecipes.filter((recipe) =>
    usefulText(recipe.modern_pharmacological_action),
  )
  const hasValidation =
    hasPlant &&
    (compounds.length > 0 || actions.length > 0 || references.length > 0)

  return (
    <article className="print-report">
      <header className="report-header">
        <p className="report-kicker">BukaLontar · Laporan Pembacaan</p>
        <h1 className="report-title">Hasil Pembacaan Lontar</h1>
        <p className="report-subtitle">
          {source.lontar_group} · Lempir{' '}
          {source.lempir_number ?? 'tidak tercatat'} · {source.language}
        </p>
      </header>

      <section className="report-section report-section-overview">
        <div className="report-grid">
          <figure className="report-figure">
            <img
              src={originalUrl}
              alt="Foto lontar yang dianalisis"
              className="report-image"
            />
            <figcaption className="report-caption">
              Foto lontar yang dianalisis.
            </figcaption>
          </figure>

          <dl className="report-meta">
            <div className="report-field">
              <dt>Kelompok naskah</dt>
              <dd>{plain(source.lontar_group)}</dd>
            </div>
            <div className="report-field">
              <dt>Lempir</dt>
              <dd>{source.lempir_number ?? 'Tidak tercatat'}</dd>
            </div>
            <div className="report-field">
              <dt>Bahasa</dt>
              <dd>{plain(source.language)}</dd>
            </div>
            <div className="report-field">
              <dt>Keterbacaan</dt>
              <dd>{confidenceLabel[paleography.confidence]}</dd>
            </div>
            <div className="report-field">
              <dt>Aksara Bali</dt>
              <dd>{paleography.aksara_bali_detected ? 'Terbaca' : 'Tidak terbaca'}</dd>
            </div>
            <div className="report-field">
              <dt>Tanaman</dt>
              <dd>
                {hasPlant
                  ? plantRecipes.length > 1
                    ? `${plantRecipes.length} resep: ${plantRecipes
                        .map((recipe) => recipe.plant_local_name)
                        .join(', ')}`
                    : plantRecipes[0].plant_local_name
                  : 'Tidak ada tanaman obat'}
              </dd>
            </div>
            <div className="report-field">
              <dt>Rujukan tanaman</dt>
              <dd>
                {!hasPlant
                  ? 'Tidak ada'
                  : plantRecipes.length > 1
                    ? `${plantRecipes.filter((recipe) => recipe.kb_match).length}/${plantRecipes.length} tercatat di Taru Pramana`
                    : plantRecipes[0].kb_match
                      ? 'Tercatat di Taru Pramana'
                      : 'Belum tercatat di Taru Pramana'}
              </dd>
            </div>
            <div className="report-field">
              <dt>Nama ilmiah</dt>
              <dd>
                {hasPlant
                  ? plantRecipes.length > 1
                    ? 'Lihat rincian resep'
                    : taxonomy.length > 0
                      ? taxonomy.join(' · ')
                      : 'Belum terverifikasi'
                  : 'Tidak ada'}
              </dd>
            </div>
          </dl>
        </div>
      </section>

      <section className="report-section">
        <h2 className="report-section-title">Ringkasan</h2>
        <p className="report-prose">
          {summary ??
            (hasPlant
              ? 'Ringkasan isi belum tersedia.'
              : 'Lembar ini tidak memuat resep ramuan obat. Tidak ada tanaman obat atau cara racik yang bisa dipetakan dari teks ini.')}
        </p>
      </section>

      <section className="report-section">
        <h2 className="report-section-title">Transliterasi</h2>
        <p className="report-prose">{transliteration}</p>
      </section>

      <section className="report-section">
        <h2 className="report-section-title">Terjemahan Indonesia</h2>
        <p className="report-prose">{indonesian}</p>
      </section>

      <section className="report-section">
        <h2 className="report-section-title">English Translation</h2>
        <p className="report-prose">{english}</p>
      </section>

      <section
        className={
          plantRecipes.length <= 2
            ? 'report-section report-section-card'
            : 'report-section'
        }
      >
        <h2 className="report-section-title">
          {hasPlant
            ? 'Kartu Usada'
            : hasAilment
              ? 'Catatan Usada'
              : 'Catatan Naskah'}
        </h2>
        {hasPlant ? (
          plantRecipes.map((recipe, index) => (
            <div
              className="report-recipe"
              key={`${recipe.plant_local_name}-${index}`}
            >
              <h3 className="report-subheading">
                {plantRecipes.length > 1
                  ? `Resep ${index + 1} · ${recipe.plant_local_name}`
                  : recipe.plant_local_name}
              </h3>
              <dl className="report-kv">
                <div className="report-field">
                  <dt>Tanaman</dt>
                  <dd>{recipe.plant_local_name}</dd>
                </div>
                <div className="report-field">
                  <dt>Bagian dipakai</dt>
                  <dd>{plain(recipe.plant_part_used)}</dd>
                </div>
                <div className="report-field">
                  <dt>Sifat tradisional</dt>
                  <dd>{plain(recipe.energy_character)}</dd>
                </div>
                <div className="report-field">
                  <dt>Bentuk racikan</dt>
                  <dd>{plain(recipe.formulation_type)}</dd>
                </div>
                <div className="report-field">
                  <dt>Indikasi</dt>
                  <dd>{plain(recipe.target_ailment)}</dd>
                </div>
                <div className="report-field">
                  <dt>Bahan pelengkap</dt>
                  <dd>
                    {recipe.adjuvants.length > 0
                      ? recipe.adjuvants.join(' · ')
                      : 'Tidak ada'}
                  </dd>
                </div>
              </dl>
              <h3 className="report-subheading">Cara racik</h3>
              <p className="report-prose">
                {usefulText(recipe.preparation_instructions)
                  ? formatTranslation(recipe.preparation_instructions)
                  : 'Cara racik belum tercatat.'}
              </p>
            </div>
          ))
        ) : (
          <p className="report-prose">
            {hasAilment
              ? `Lembar ini membahas ${entity.target_ailment}, tetapi tidak menyebut ramuan tanaman yang bisa dipetakan.`
              : 'Tidak ada tanaman obat atau cara racik pada lembar ini.'}
          </p>
        )}
      </section>

      {hasValidation ? (
        <section className="report-section report-section-validation">
          <h2 className="report-section-title">Kandungan dan Rujukan</h2>
          <dl className="report-kv">
            <div className="report-field">
              <dt>Kandungan aktif</dt>
              <dd>
                {compounds.length > 0
                  ? compounds.join(' · ')
                  : 'Belum tersedia'}
              </dd>
            </div>
          </dl>
          <h3 className="report-subheading">Khasiat modern</h3>
          {actions.length > 0 ? (
            <ul className="report-list">
              {actions.map((recipe) => (
                <li key={recipe.plant_local_name}>
                  <strong>{recipe.plant_local_name}:</strong>{' '}
                  {formatTranslation(recipe.modern_pharmacological_action)}
                </li>
              ))}
            </ul>
          ) : (
            <p className="report-prose">Belum tersedia.</p>
          )}
          <h3 className="report-subheading">Sumber rujukan</h3>
          {references.length > 0 ? (
            <ul className="report-list">
              {references.map((reference) => (
                <li key={reference}>{reference}</li>
              ))}
            </ul>
          ) : (
            <p className="report-prose">Belum ada rujukan.</p>
          )}
        </section>
      ) : null}

      <section className="report-section report-section-notes">
        <h2 className="report-section-title">Catatan</h2>
        <p className="report-prose">
          Tingkat keterbacaan: {confidenceLabel[paleography.confidence]}.
          Goresan aus, tumpukan pasang aksara, dan istilah arkais bisa salah
          dibaca. Periksa lempir asli dan rujukan filologi sebelum mengambil
          keputusan penting.
        </p>
        <p className="report-disclaimer">{validation.safety_disclaimer}</p>
      </section>

      <footer className="report-footer">
        <span>BukaLontar · Laporan Pembacaan Lontar</span>
        <span>{source.lontar_group}</span>
      </footer>
    </article>
  )
}
