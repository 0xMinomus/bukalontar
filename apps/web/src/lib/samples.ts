export interface Sample {
  id: string
  name: string
  detail: string
  src: string
  credit: string
  sourceUrl: string
  licenseUrl: string
}

export const samples: Sample[] = [
  {
    id: 'usada-lare',
    name: 'Usada Lare',
    detail: 'Aksara Bali · sisi depan',
    src: '/samples/lontar-usada-lare-depan.jpg',
    credit: 'Eka343, CC BY-SA 4.0',
    sourceUrl:
      'https://commons.wikimedia.org/wiki/File:Bali-lontar-Balai_Bahasa_Bali-Usada_Lare.jpg',
    licenseUrl: 'https://creativecommons.org/licenses/by-sa/4.0/',
  },
  {
    id: 'usada-rare-154',
    name: 'Usada Rare 154',
    detail: 'Aksara Bali · sisi depan',
    src: '/samples/lontar-usada-rare-154-depan.jpg',
    credit: 'Eka343, CC BY-SA 4.0',
    sourceUrl:
      'https://commons.wikimedia.org/wiki/File:Bali-lontar-Balai_Bahasa_Bali-Usada_Rare_154.jpg',
    licenseUrl: 'https://creativecommons.org/licenses/by-sa/4.0/',
  },
]
