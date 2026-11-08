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
    id: 'kemenuh',
    name: 'Lontar Kemenuh',
    detail: 'Kawruhan Desti ring Raga · catatan kebatinan',
    src: '/samples/lontar-kemenuh.jpg',
    credit: 'Eka343, CC BY-SA 4.0',
    sourceUrl:
      'https://commons.wikimedia.org/wiki/File:Bali-lontar-Kemenuh-Kawruhan_Desti_ring_Raga-70.jpg',
    licenseUrl: 'https://creativecommons.org/licenses/by-sa/4.0/',
  },
  {
    id: 'sukawati',
    name: 'Lontar Sukawati',
    detail: 'Pamandisuara · resep usada',
    src: '/samples/lontar-sukawati.jpg',
    credit: 'Eka343, CC BY-SA 4.0',
    sourceUrl:
      'https://commons.wikimedia.org/wiki/File:Bali-lontar-Sukawati-Pamandisuara-113.jpg',
    licenseUrl: 'https://creativecommons.org/licenses/by-sa/4.0/',
  },
]
