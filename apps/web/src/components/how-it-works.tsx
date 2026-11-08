import { ArrowUpRight } from 'lucide-react'

import { Ornament } from '@/components/ornament'
import { Reveal } from '@/components/reveal'

const layers = [
  {
    numeral: 'I',
    title: 'Aksara',
    body: 'Membaca goresan aksara Bali pada lempir, termasuk gantungan, gempelan, dan pangangge yang bertumpuk.',
  },
  {
    numeral: 'II',
    title: 'Bahasa',
    body: 'Membaca isi naskah dalam terjemahan Indonesia dan Inggris, sambil tetap menyandingkannya dengan teks asli.',
  },
  {
    numeral: 'III',
    title: 'Konteks',
    body: 'Menampilkan sumber dan keterbacaan. Jika ada resep Usadha, tanaman dicocokkan dengan basis pengetahuan Taru Pramana.',
  },
]

export function HowItWorks() {
  return (
    <>
      <section id="alur" className="reading-section">
        <div className="mx-auto max-w-5xl px-5 py-14 md:px-8 md:py-20">
          <Reveal className="text-center">
            <h2 className="font-heading text-4xl tracking-tight md:text-5xl">
              Dari goresan ke bacaan
            </h2>
            <Ornament className="mx-auto mt-4 w-36 text-brass-strong/70" />
            <p className="mx-auto mt-4 max-w-[62ch] text-sm leading-relaxed text-muted-foreground sm:text-base">
              Unggah satu lempir. Baca teks Latin dan terjemahannya, lalu
              periksa kembali bagian yang kurang jelas pada foto asli.
            </p>
          </Reveal>

          <div className="reading-leaves mt-10">
            {layers.map((layer, index) => (
              <Reveal key={layer.numeral} delay={index * 70} className="reading-leaf">
                <article>
                  <Ornament className="mb-4 w-16 text-brass-strong" />
                  <h3 className="mt-2 font-heading text-2xl tracking-tight">
                    {layer.title}
                  </h3>
                  <p className="mt-2 text-sm leading-relaxed text-muted-foreground">
                    {layer.body}
                  </p>
                </article>
              </Reveal>
            ))}
          </div>

          <Reveal className="mt-9 text-center">
            <a
              href="https://github.com/0xMinomus/bukalontar/tree/main/docs"
              target="_blank"
              rel="noreferrer"
              className="small-caps inline-flex items-center gap-1.5 border-b border-brass/50 pb-1 text-xs font-semibold text-primary transition-colors hover:text-brass-strong"
            >
              Panduan membaca
              <ArrowUpRight className="size-3.5" aria-hidden="true" />
            </a>
          </Reveal>
        </div>
      </section>

      <section id="tentang" className="about-section">
        <div className="mx-auto max-w-5xl px-5 pb-16 pt-8 md:px-8 md:pb-24">
          <Reveal>
            <div className="paper-panel frame-antique about-folio p-6 md:p-10">
              <div className="grid gap-8 md:grid-cols-[1.05fr_0.95fr] md:gap-12">
                <div>
                  <h2 className="font-heading text-4xl tracking-tight md:text-5xl">
                    Lontar menyimpan<br /><em>banyak pengetahuan.</em>
                  </h2>
                  <p className="mt-4 text-sm leading-relaxed text-muted-foreground sm:text-base">
                    Kakawin, babad, geguritan, tutur, hingga usadha. BukaLontar
                    membantu membaca naskah Bali dari foto, tanpa menganggap
                    setiap lempir sebagai catatan pengobatan.
                  </p>
                  <p className="mt-4 text-sm leading-relaxed text-muted-foreground">
                    Bila naskah memuat resep Usadha, kartu tanaman dan rujukan
                    Taru Pramana tetap hadir. Hasil dapat disalin atau dicetak
                    sebagai laporan untuk diperiksa lebih lanjut.
                  </p>
                  <blockquote className="mt-6 border-l border-brass pl-4">
                    <p className="font-serif text-[1.05rem] italic leading-relaxed text-foreground">
                      “Tityang taru sotong, daging anget, rasa sepet, dados
                      anggen tamba pangemped mising.”
                    </p>
                    <cite className="small-caps mt-2 block text-[11px] not-italic text-muted-foreground">
                      Lontar Taru Pramana, Va. 1384/10
                    </cite>
                  </blockquote>
                </div>

                <figure>
                  <div className="overflow-hidden rounded-lg border border-border/80 bg-mat">
                    <img
                      src="/samples/lontar-wanagiri-usada-3.jpg"
                      alt="Dua lempir lontar usadha Bali dengan goresan aksara di atas permukaan gelap"
                      loading="lazy"
                      decoding="async"
                      className="h-48 w-full object-cover sepia-[0.12] md:h-56"
                    />
                  </div>
                  <figcaption className="mt-2 text-[11px] leading-relaxed text-muted-foreground">
                    Lontar Wanagiri, Usada 3. Foto oleh{' '}
                    <a
                      href="https://commons.wikimedia.org/wiki/File:Bali-lontar-Wanagiri-Usada-3.jpg"
                      target="_blank"
                      rel="noreferrer"
                      className="underline underline-offset-2 hover:text-foreground"
                    >
                      Eka343
                    </a>
                    , CC BY-SA 4.0.
                  </figcaption>
                </figure>
              </div>

              <div className="mt-8 grid gap-6 border-t border-border/70 pt-6 md:grid-cols-2">
                <div>
                  <p className="label-xs">Rujukan naskah</p>
                  <ul className="mt-2 space-y-2 text-sm leading-relaxed text-muted-foreground">
                    <li>Lontar Taru Pramana, Gedong Kirtya Va. 1384/10</li>
                    <li>Lontar Usadha Rare, Gedong Kirtya Va. 1391/11</li>
                    <li>Katalog usada Bali dan taksonomi botani modern</li>
                  </ul>
                </div>
                <div>
                  <p className="label-xs">Repositori</p>
                  <a
                    href="https://github.com/0xMinomus/bukalontar"
                    target="_blank"
                    rel="noreferrer"
                    className="mt-1 inline-flex items-center gap-1.5 text-sm font-semibold text-primary underline-offset-4 hover:underline"
                  >
                    github.com/0xMinomus/bukalontar
                    <ArrowUpRight className="size-4" aria-hidden="true" />
                  </a>
                  <p className="mt-3 text-xs leading-relaxed text-muted-foreground">
                    Hasil AI adalah bacaan awal. Periksa lempir asli dan rujukan
                    filologi; informasi Usadha bukan nasihat medis.
                  </p>
                </div>
              </div>
            </div>
          </Reveal>
        </div>
      </section>
    </>
  )
}
