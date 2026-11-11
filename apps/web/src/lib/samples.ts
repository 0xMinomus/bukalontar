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
    id: 'wanagiri',
    name: 'Lontar Wanagiri',
    detail: 'Usada 3 · beraksara Bali',
    src: '/samples/lontar-wanagiri-usada-3.jpg',
    credit: 'Eka343, CC BY-SA 4.0',
    sourceUrl:
      'https://commons.wikimedia.org/wiki/File:Bali-lontar-Wanagiri-Usada-3.jpg',
    licenseUrl: 'https://creativecommons.org/licenses/by-sa/4.0/',
  },
  {
    id: 'sukawati',
    name: 'Lontar Sukawati',
    detail: 'Usada 07 · beraksara Bali',
    src: '/samples/lontar-sukawati-usada-07.jpg',
    credit: 'Eka343, CC BY-SA 4.0',
    sourceUrl:
      'https://commons.wikimedia.org/wiki/File:Bali-lontar-Sukawati-Usada-07.jpg',
    licenseUrl: 'https://creativecommons.org/licenses/by-sa/4.0/',
  },
]
