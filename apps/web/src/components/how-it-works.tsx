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
    body: 'Menerjemahkan bahasa Kawi klasik ke Indonesia dan Inggris, dengan konteks istilah medis dan etnobotani.',
  },
  {
    numeral: 'III',
    title: 'Botani',
    body: 'Memetakan nama tanaman lokal ke nama ilmiah, famili, dan senyawa aktif dari basis data Taru Pramana.',
  },
]

export function HowItWorks() {
  return (
    <>
      <section id="alur" className="border-t border-border/70">
        <div className="mx-auto max-w-5xl px-5 py-14 md:px-8 md:py-20">
          <Reveal className="text-center">
            <p className="small-caps text-xs font-semibold text-brass-strong">
              Alur kerja
            </p>
            <h2 className="mt-3 font-heading text-4xl tracking-tight md:text-5xl">
              Tiga lapis yang dibaca
            </h2>
            <Ornament className="mx-auto mt-4 w-36 text-brass-strong/70" />
            <p className="mx-auto mt-4 max-w-[62ch] text-sm leading-relaxed text-muted-foreground sm:text-base">
              Lontar usada tidak hanya sulit karena aksaranya. Bahasa dan nama
              tanamannya juga sudah jauh dari percakapan sehari-hari. BukaLontar
              membaca ketiganya dalam satu alur.
            </p>
          </Reveal>

          <div className="mt-10 grid gap-8 md:grid-cols-3 md:gap-0 md:divide-x md:divide-border/70">
            {layers.map((layer, index) => (
              <Reveal
                key={layer.numeral}
                delay={index * 90}
                className={
                  index === 0
                    ? 'md:pr-8'
                    : index === 1
                      ? 'md:px-8'
                      : 'md:pl-8'
                }
              >
                <article>
                  <p className="font-heading text-3xl leading-none text-brass-strong">
                    {layer.numeral}
                  </p>
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

      <section id="tentang" className="border-t border-border/70">
        <div className="mx-auto max-w-5xl px-5 py-14 md:px-8 md:py-20">
          <Reveal>
            <div className="paper-panel frame-antique p-6 md:p-10">
              <div className="grid gap-8 md:grid-cols-[1.05fr_0.95fr] md:gap-12">
                <div>
                  <p className="small-caps text-xs font-semibold text-brass-strong">
                    Tentang
                  </p>
                  <h2 className="mt-2 font-heading text-4xl tracking-tight md:text-5xl">
                    BukaLontar
                  </h2>
                  <p className="mt-4 text-sm leading-relaxed text-muted-foreground sm:text-base">
                    BukaLontar adalah purwarupa untuk membuka kembali lontar usada
                    Bali. Fokusnya bukan sekadar membaca huruf, tetapi menyusun
                    pengetahuan obat tradisional yang bisa diperiksa: tanaman
                    apa, bagian mana, bagaimana diracik, dan bagaimana
                    khasiatnya menurut riset modern.
                  </p>
                  <blockquote className="mt-6 border-l-2 border-brass pl-4">
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
                    BukaLontar membaca lontar secara langsung. Hasil ditampilkan
                    bersama rujukan tanaman dan pengingat medis.
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
