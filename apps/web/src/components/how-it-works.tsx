import { ArrowUpRight } from 'lucide-react'

const layers = [
  {
    number: '01',
    title: 'Aksara',
    body: 'Membaca goresan aksara Bali pada lempir, termasuk gantungan, gempelan, dan pangangge yang bertumpuk.',
  },
  {
    number: '02',
    title: 'Bahasa',
    body: 'Menerjemahkan bahasa Kawi klasik ke Indonesia dan Inggris, dengan konteks istilah medis dan etnobotani.',
  },
  {
    number: '03',
    title: 'Botani',
    body: 'Memetakan nama tanaman lokal ke nama ilmiah, famili, dan senyawa aktif dari basis data Taru Pramana.',
  },
]

export function HowItWorks() {
  return (
    <>
      <section id="alur" className="border-t border-border">
        <div className="mx-auto max-w-[1440px] px-5 py-12 md:px-8 md:py-16">
          <div className="max-w-[62ch]">
            <h2 className="font-serif text-3xl tracking-tight md:text-4xl">
              Tiga lapis yang dibaca
            </h2>
            <p className="mt-3 text-sm leading-relaxed text-muted-foreground sm:text-base">
              Lontar usada tidak hanya sulit karena aksaranya. Bahasa dan nama
              tanamannya juga sudah jauh dari percakapan sehari-hari. UsadhaAI
              membaca ketiganya dalam satu alur.
            </p>
          </div>

          <div className="mt-8 grid gap-8 md:grid-cols-3 md:gap-0 md:divide-x md:divide-border">
            {layers.map((layer, index) => (
              <article
                key={layer.number}
                className={
                  index === 0
                    ? 'md:pr-8'
                    : index === 1
                      ? 'md:px-8'
                      : 'md:pl-8'
                }
              >
                <p className="num text-sm font-semibold text-brass-strong">
                  {layer.number}
                </p>
                <h3 className="mt-2 font-serif text-2xl tracking-tight">
                  {layer.title}
                </h3>
                <p className="mt-2 text-sm leading-relaxed text-muted-foreground">
                  {layer.body}
                </p>
              </article>
            ))}
          </div>

          <a
            href="https://github.com/0xMinomus/digifest/tree/main/docs"
            target="_blank"
            rel="noreferrer"
            className="mt-8 inline-flex items-center gap-1.5 text-sm font-semibold text-primary underline-offset-4 hover:underline"
          >
            Baca catatan teknis
            <ArrowUpRight className="size-4" aria-hidden="true" />
          </a>
        </div>
      </section>

      <section id="tentang" className="border-t border-border">
        <div className="mx-auto grid max-w-[1440px] gap-8 px-5 py-12 md:grid-cols-[1.1fr_0.9fr] md:px-8 md:py-16">
          <div className="max-w-[62ch]">
            <h2 className="font-serif text-3xl tracking-tight md:text-4xl">
              Tentang UsadhaAI
            </h2>
            <p className="mt-3 text-sm leading-relaxed text-muted-foreground sm:text-base">
              UsadhaAI adalah purwarupa untuk membuka kembali lontar usada Bali.
              Fokusnya bukan sekadar membaca huruf, tetapi menyusun pengetahuan
              obat tradisional yang bisa diperiksa: tanaman apa, bagian mana,
              bagaimana diracik, dan apa kata fitofarmaka modern.
            </p>
            <p className="mt-3 text-sm leading-relaxed text-muted-foreground sm:text-base">
              Antarmuka ini masih memakai data contoh. Backend Go dan pemanggilan
              Gemini 3.8 Flash akan dihubungkan pada tahap berikutnya.
            </p>
          </div>

          <div className="border-t border-border pt-6 md:border-l md:border-t-0 md:pl-8 md:pt-0">
            <p className="label-xs">Rujukan naskah</p>
            <ul className="mt-2 space-y-2 text-sm leading-relaxed text-muted-foreground">
              <li>Lontar Taru Pramana, Gedong Kirtya Va. 1384/10</li>
              <li>Lontar Usadha Rare, Gedong Kirtya Va. 1391/11</li>
              <li>Katalog usada Bali dan taksonomi botani modern</li>
            </ul>

            <p className="label-xs mt-6">Repositori</p>
            <a
              href="https://github.com/0xMinomus/digifest"
              target="_blank"
              rel="noreferrer"
              className="mt-1 inline-flex items-center gap-1.5 text-sm font-semibold text-primary underline-offset-4 hover:underline"
            >
              github.com/0xMinomus/digifest
              <ArrowUpRight className="size-4" aria-hidden="true" />
            </a>
          </div>
        </div>
      </section>
    </>
  )
}
