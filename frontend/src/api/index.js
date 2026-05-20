export const BASE = import.meta.env.VITE_API_BASE_URL ?? ''

async function request(method, path) {
  const res = await fetch(`${BASE}${path}`, { method })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: `HTTP ${res.status}` }))
    throw new Error(err.error ?? `HTTP ${res.status}`)
  }
  return res.json()
}

const get = (path) => request('GET', path)

export const api = {
  getPaintings: ()   => get('/api/paintings'),
  getPainting:  (id) => get(`/api/paintings/${id}`),
  getAuthors:   ()   => get('/api/authors'),
  getAuthor:    (id) => get(`/api/authors/${id}`),
  getMaterials: ()   => get('/api/materials'),
}
