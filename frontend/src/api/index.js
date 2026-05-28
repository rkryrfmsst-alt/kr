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

async function put(path, body) {
  const res = await fetch(`${BASE}${path}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: `HTTP ${res.status}` }))
    throw new Error(err.error ?? `HTTP ${res.status}`)
  }
  return res.json()
}

async function del(path) {
  const res = await fetch(`${BASE}${path}`, { method: 'DELETE' })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: `HTTP ${res.status}` }))
    throw new Error(err.error ?? `HTTP ${res.status}`)
  }
  return res.json()
}

async function post(path, body) {
  const res = await fetch(`${BASE}${path}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: `HTTP ${res.status}` }))
    throw new Error(err.error ?? `HTTP ${res.status}`)
  }
  return res.json()
}

export const api = {
  getPaintings: ()          => get('/api/paintings'),
  getPainting:  (id)        => get(`/api/paintings/${id}`),
  getAuthors:   ()          => get('/api/authors'),
  getAuthor:    (id)        => get(`/api/authors/${id}`),
  getMaterials: ()          => get('/api/materials'),
  getStyles:    ()          => get('/api/styles'),
  getPlots:     ()          => get('/api/plots'),

  login:    (email, password)           => post('/api/auth/login',    { email, password }),
  register: (email, username, password) => post('/api/auth/register', { email, username, password }),
  logout:   ()                          => post('/api/auth/logout',   {}),
  me:       ()                          => get('/api/auth/me'),

  createAuthor: (data)    => post('/api/admin/authors', data),
  updateAuthor: (id, data) => put(`/api/admin/authors/${id}`, data),
  deleteAuthor: (id)       => del(`/api/admin/authors/${id}`),

  createPainting: (formData) => {
    return fetch(`${BASE}/api/admin/paintings`, { method: 'POST', body: formData })
      .then(async res => {
        if (!res.ok) {
          const err = await res.json().catch(() => ({ error: `HTTP ${res.status}` }))
          throw new Error(err.error ?? `HTTP ${res.status}`)
        }
        return res.json()
      })
  },

  updatePainting: (id, formData) => {
    return fetch(`${BASE}/api/admin/paintings/${id}`, { method: 'PUT', body: formData })
      .then(async res => {
        if (!res.ok) {
          const err = await res.json().catch(() => ({ error: `HTTP ${res.status}` }))
          throw new Error(err.error ?? `HTTP ${res.status}`)
        }
        return res.json()
      })
  },

  deletePainting: (id) => del(`/api/admin/paintings/${id}`),
}
