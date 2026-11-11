export type Confidence = 'low' | 'medium' | 'high'

export type LontarGroup =
  | 'Usadha Taru Pramana'
  | 'Usadha Tiwang'
  | 'Usadha Dalem'
  | 'Usadha Yeh'
  | 'Tidak diketahui'
  | 'Non-Usadha'

export type PlantPart =
  | 'Don (Daun)'
  | 'Babakan (Kulit Kayu)'
  | 'Akah (Akar)'
  | 'Buah'
  | 'Sekar (Bunga)'
  | 'Campuran'
  | 'Tidak disebutkan.'

export type EnergyCharacter = 'Panes (Hangat)' | 'Tis (Dingin)' | 'Dumelada (Sedang)' | 'Tidak disebutkan.'

export type FormulationType =
  | 'Loloh (Minum)'
  | 'Boreh (Balur)'
  | 'Tutuh (Tetes)'
  | 'Sasap (Kompres)'
  | 'Tidak disebutkan.'

export interface ReferenceLink {
  title: string
  url: string
  kind: 'literature_search' | 'literature'
}

export interface Recipe {
  plant_local_name: string
  plant_scientific_name: string
  botanical_family: string
  kb_match: boolean
  plant_part_used: PlantPart
  energy_character: EnergyCharacter
  formulation_type: FormulationType
  adjuvants: string[]
  target_ailment: string
  preparation_instructions: string
  phytochemical_compounds: string[]
  modern_pharmacological_action: string
  reference_studies: string[]
  reference_links?: ReferenceLink[]
}

export interface AnalysisResponse {
  manuscript_id: string
  source_reference: {
    lontar_group: LontarGroup
    lempir_number?: number
    language: 'Kawi / Jawa Kuno' | 'Bali Alus' | 'Tidak diketahui'
  }
  paleography: {
    aksara_bali_detected: boolean
    confidence: Confidence
    balinese_transcription?: string
    transliteration: string
    indonesian_translation: string
    english_translation: string
  }
  ethnomedicinal_entity: {
    plant_local_name: string
    plant_scientific_name: string
    botanical_family: string
    kb_match: boolean
    plant_part_used: PlantPart
    energy_character: EnergyCharacter
    formulation_type: FormulationType
    adjuvants: string[]
    target_ailment: string
    preparation_instructions: string
  }
  scientific_validation: {
    phytochemical_compounds: string[]
    modern_pharmacological_action: string
    reference_studies: string[]
    reference_links?: ReferenceLink[]
    safety_disclaimer: string
  }
  recipes?: Recipe[]
  reading_assistance?: {
    mode: 'none' | 'examples' | 'ocr' | 'both'
    ocr_model?: string
    ocr_ms: number
    focus_regions: number
    example_count: number
    example_sources: string[]
    warnings: string[]
  }
}

export const SAFETY_DISCLAIMER =
  'Informasi ini untuk preservasi warisan budaya dan edukasi fitofarmaka. Untuk penanganan medis, konsultasikan dengan dokter atau tenaga kesehatan profesional.'

const taruSotong: AnalysisResponse = {
  manuscript_id: 'taru-pramana-1384-10-lp-03',
  source_reference: {
    lontar_group: 'Usadha Taru Pramana',
    lempir_number: 3,
    language: 'Kawi / Jawa Kuno',
  },
  paleography: {
    aksara_bali_detected: true,
    confidence: 'high',
    transliteration:
      'Tityang taru sotong, daging anget, rasa sepet, dados anggen tamba pangemped mising, sa, muncuk tityange ulig, anggen papuser ring pungsed, ra, katumbar bolong, 3, lunak.',
    indonesian_translation:
      'Hamba Pohon Jambu Biji (Sotong), sifat daging buah hamba hangat, rasanya sepet, dapat dipergunakan sebagai obat penahan mencret (diare). Sarana: ambil pucuk daun hamba kemudian ditumbuk atau diulek, lalu ditempelkan pada pusar. Ramuan pelengkap: campurkan tiga butir ketumbar bolong dan asam jawa.',
    english_translation:
      'I am the Guava Tree (Sotong). My fruit flesh is warm and astringent, and can be used as a remedy for diarrhea. Take my young leaves, crush them, and apply them to the navel. Add three coriander seeds and tamarind.',
  },
  ethnomedicinal_entity: {
    plant_local_name: 'Taru Sotong',
    plant_scientific_name: 'Psidium guajava',
    botanical_family: 'Myrtaceae',
    kb_match: true,
    plant_part_used: 'Don (Daun)',
    energy_character: 'Panes (Hangat)',
    formulation_type: 'Sasap (Kompres)',
    adjuvants: ['Ketumbar bolong', 'Asam jawa'],
    target_ailment: 'Mising (diare)',
    preparation_instructions:
      'Pucuk daun ditumbuk halus, dicampur tiga butir ketumbar bolong dan asam jawa, lalu ditempelkan pada pusar.',
  },
  scientific_validation: {
    phytochemical_compounds: ['Tannin', 'Quercetin', 'Flavonoid'],
    modern_pharmacological_action:
      'Antimikroba terhadap Escherichia coli dan antispasmodik untuk meredakan kontraksi usus.',
    reference_studies: [
      'Tinjauan fitokimia Psidium guajava',
      'Uji antimikroba ekstrak daun jambu biji',
    ],
    safety_disclaimer: SAFETY_DISCLAIMER,
  },
}

const taruPule: AnalysisResponse = {
  manuscript_id: 'taru-pramana-1384-10-lp-07',
  source_reference: {
    lontar_group: 'Usadha Tiwang',
    lempir_number: 12,
    language: 'Kawi / Jawa Kuno',
  },
  paleography: {
    aksara_bali_detected: true,
    confidence: 'high',
    transliteration:
      'Tityang taru pule, daging tis, kulit getep panes, dados anggen tamba cika, sa, babakan tityang kaurab adas, temutitis, ulig anggen boreh.',
    indonesian_translation:
      'Hamba Pohon Pule, daging buah hamba sejuk, kulit batang hamba hangat, dapat digunakan sebagai obat kram perut mendadak (cika). Sarana: kulit kayu hamba dicampur adas dan temutitis, diulek lalu dibalurkan sebagai boreh.',
    english_translation:
      'I am the Pule Tree. My fruit flesh is cool and my bark is warm. I can be used for sudden stomach cramps (cika). Mix my bark with fennel and temutitis, crush the mixture, and apply it as a balm.',
  },
  ethnomedicinal_entity: {
    plant_local_name: 'Taru Pule',
    plant_scientific_name: 'Alstonia scholaris',
    botanical_family: 'Apocynaceae',
    kb_match: true,
    plant_part_used: 'Babakan (Kulit Kayu)',
    energy_character: 'Panes (Hangat)',
    formulation_type: 'Boreh (Balur)',
    adjuvants: ['Adas', 'Temutitis'],
    target_ailment: 'Cika (kram perut)',
    preparation_instructions:
      'Kulit batang diulek bersama adas dan temutitis, lalu dibalurkan pada perut sebagai boreh.',
  },
  scientific_validation: {
    phytochemical_compounds: ['Echitamine', 'Alstonine', 'Alkaloid'],
    modern_pharmacological_action:
      'Analgesik dan antipiretik; alkaloid kulit batang pule diteliti untuk efek antinyeri dan penurun demam.',
    reference_studies: [
      'Tinjauan alkaloid Alstonia scholaris',
      'Kajian etnofarmakologi kulit batang pule',
    ],
    safety_disclaimer: SAFETY_DISCLAIMER,
  },
}

export function mockAnalysis(sampleId?: string): AnalysisResponse {
  const data = sampleId === 'sukawati' ? taruPule : taruSotong
  return structuredClone(data)
}
