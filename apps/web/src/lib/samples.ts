export interface Sample {
  id: string
  name: string
  genre: 'Usadha' | 'Carcan' | 'Kakawin'
  detail: string
  src: string
  credit: string
  sourceUrl: string
  licenseUrl: string
  referenceUrl: string
}

export const samples: Sample[] = [
  {
    id: 'usada-sukawati',
    name: 'Usadha Sukawati',
    genre: 'Usadha',
    detail: 'Racikan tanaman · sisi depan',
    src: '/samples/lontar-usadha-sukawati-depan.jpg',
    credit: 'Eka343 · CC BY-SA 4.0',
    sourceUrl:
      'https://commons.wikimedia.org/wiki/File:Bali-lontar-Sukawati-Usada-07.jpg',
    licenseUrl: 'https://creativecommons.org/licenses/by-sa/4.0/',
    referenceUrl:
      'https://ban.wikisource.org/w/index.php?title=Kaca:Bali-lontar-Sukawati-Usada-07.jpg&oldid=98029',
  },
  {
    id: 'usada-rare-kemenuh',
    name: 'Usadha Rare Kemenuh',
    genre: 'Usadha',
    detail: 'Racikan untuk rare · sisi depan',
    src: '/samples/lontar-usadha-rare-kemenuh-depan.jpg',
    credit: 'Eka343 · CC BY-SA 4.0',
    sourceUrl:
      'https://commons.wikimedia.org/wiki/File:Bali-lontar-Kemenuh-Usadha_Rare-56.jpg',
    licenseUrl: 'https://creativecommons.org/licenses/by-sa/4.0/',
    referenceUrl:
      'https://ban.wikisource.org/w/index.php?title=Kaca:Bali-lontar-Kemenuh-Usadha_Rare-56.jpg&oldid=107871',
  },
  {
    id: 'carcan-kucing',
    name: 'Carcan Kucing',
    genre: 'Carcan',
    detail: 'Ciri dan pertanda kucing · lempir 1B',
    src: '/samples/lontar-carcan-kucing-1b.jpg',
    credit: 'Internet Archive · PD-ID-exempt',
    sourceUrl:
      'https://commons.wikimedia.org/wiki/File:Bali-lontar-carcan-kucing-250ppi.pdf',
    licenseUrl:
      'https://commons.wikimedia.org/wiki/File:Bali-lontar-carcan-kucing-250ppi.pdf#Licensing',
    referenceUrl:
      'https://ban.wikisource.org/w/index.php?title=Carcan_Kucing&oldid=107176',
  },
  {
    id: 'kakawin-bharatayuda',
    name: 'Kakawin Bharatayuda',
    genre: 'Kakawin',
    detail: 'Susastra · sisi depan',
    src: '/samples/lontar-kakawin-bharatayuda-depan.jpg',
    credit: 'Eka343 · CC BY-SA 4.0',
    sourceUrl:
      'https://commons.wikimedia.org/wiki/File:Bali-lontar-Balai_Bahasa_Bali-Kekawin_Bharatayuda_Maarti_01.jpg',
    licenseUrl: 'https://creativecommons.org/licenses/by-sa/4.0/',
    referenceUrl:
      'https://ban.wikisource.org/w/index.php?title=Kekawin_Bharatayuda_Maarti_(Balai_Bahasa_Prov._Bali)&oldid=134900',
  },
]
