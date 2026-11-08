import type { AnalysisResponse } from './analysis'

export type PersistedView = 'original' | 'enhanced'

export interface PersistedSession {
  result: AnalysisResponse
  original: Blob
  enhanced: Blob | null
  enhancedMeta?: {
    width: number
    height: number
    elapsedMs: number
  }
  fileName: string
  sourceId?: string
  view: PersistedView
  savedAt: number
}

const DB_NAME = 'usadha-ai'
const DB_VERSION = 1
const STORE_NAME = 'session'
const ACTIVE_KEY = 'active'

function openDatabase() {
  const { promise, resolve, reject } = Promise.withResolvers<IDBDatabase>()
  const request = indexedDB.open(DB_NAME, DB_VERSION)

  request.onupgradeneeded = () => {
    if (!request.result.objectStoreNames.contains(STORE_NAME)) {
      request.result.createObjectStore(STORE_NAME)
    }
  }
  request.onsuccess = () => resolve(request.result)
  request.onerror = () => reject(request.error)

  return promise
}

export async function saveAnalysisSession(session: PersistedSession) {
  const database = await openDatabase()
  try {
    const { promise, resolve, reject } = Promise.withResolvers<void>()
    const transaction = database.transaction(STORE_NAME, 'readwrite')
    transaction.objectStore(STORE_NAME).put(session, ACTIVE_KEY)
    transaction.oncomplete = () => resolve()
    transaction.onerror = () => reject(transaction.error)
    transaction.onabort = () => reject(transaction.error)
    await promise
  } finally {
    database.close()
  }
}

export async function loadAnalysisSession(): Promise<PersistedSession | null> {
  const database = await openDatabase()
  try {
    const { promise, resolve, reject } =
      Promise.withResolvers<PersistedSession | null>()
    const transaction = database.transaction(STORE_NAME, 'readonly')
    const request = transaction.objectStore(STORE_NAME).get(ACTIVE_KEY)
    request.onsuccess = () =>
      resolve((request.result as PersistedSession | undefined) ?? null)
    request.onerror = () => reject(request.error)
    return await promise
  } finally {
    database.close()
  }
}

export async function clearAnalysisSession() {
  const database = await openDatabase()
  try {
    const { promise, resolve, reject } = Promise.withResolvers<void>()
    const transaction = database.transaction(STORE_NAME, 'readwrite')
    transaction.objectStore(STORE_NAME).delete(ACTIVE_KEY)
    transaction.oncomplete = () => resolve()
    transaction.onerror = () => reject(transaction.error)
    transaction.onabort = () => reject(transaction.error)
    await promise
  } finally {
    database.close()
  }
}
