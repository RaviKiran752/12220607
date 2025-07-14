import { useState } from 'react'
import './App.css'
import axios from 'axios'

const API_BASE = 'http://localhost:3000'

function App() {
  const [url, setUrl] = useState('')
  const [validity, setValidity] = useState(30)
  const [shortcode, setShortcode] = useState('')
  const [result, setResult] = useState<{shortcode: string, expiry: string} | null>(null)
  const [error, setError] = useState('')
  const [statsCode, setStatsCode] = useState('')
  const [stats, setStats] = useState<any>(null)
  const [statsError, setStatsError] = useState('')
  const [copied, setCopied] = useState(false)
  const [loading, setLoading] = useState(false)
  const [statsLoading, setStatsLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setResult(null)
    setLoading(true)
    try {
      const res = await axios.post(`${API_BASE}/shorturls`, {
        url,
        validity,
        shortcode: shortcode || undefined,
      })
      setResult(res.data)
      setCopied(false)
    } catch (err: any) {
      setError(err.response?.data || 'Error creating short URL')
    } finally {
      setLoading(false)
    }
  }

  const handleStats = async (e: React.FormEvent) => {
    e.preventDefault()
    setStatsError('')
    setStats(null)
    setStatsLoading(true)
    try {
      const res = await axios.get(`${API_BASE}/shorturls/${statsCode}`)
      setStats(res.data)
    } catch (err: any) {
      setStatsError(err.response?.data || 'Error fetching stats')
    } finally {
      setStatsLoading(false)
    }
  }

  const handleCopy = async (text: string) => {
    try {
      await navigator.clipboard.writeText(text)
      setCopied(true)
      setTimeout(() => setCopied(false), 1200)
    } catch {}
  }

  return (
    <div className="app-bg">
      <div className="container">
        <h1>URL Shortener</h1>
        <form onSubmit={handleSubmit}>
          <input
            type="url"
            placeholder="Enter URL"
            value={url}
            onChange={e => setUrl(e.target.value)}
            required
          />
          <input
            type="number"
            placeholder="Validity (seconds)"
            value={validity}
            onChange={e => setValidity(Number(e.target.value))}
            min={1}
          />
          <input
            type="text"
            placeholder="Custom shortcode (optional)"
            value={shortcode}
            onChange={e => setShortcode(e.target.value)}
          />
          <button type="submit" disabled={loading}>{loading ? 'Shortening...' : 'Shorten'}</button>
        </form>
        {error && <div className="error">{error}</div>}
        {result && (
          <div className="result">
            <div className="short-url-row">
              <span>Short URL: <a href={`${API_BASE}/${result.shortcode}`} target="_blank" rel="noopener noreferrer">{`${API_BASE}/${result.shortcode}`}</a></span>
              <button className="copy-btn" onClick={() => handleCopy(`${API_BASE}/${result.shortcode}`)}>{copied ? 'Copied!' : 'Copy'}</button>
            </div>
            <p>Expires at: {result.expiry}</p>
          </div>
        )}

        <h2>Get Short URL Stats</h2>
        <form onSubmit={handleStats}>
          <input
            type="text"
            placeholder="Enter shortcode"
            value={statsCode}
            onChange={e => setStatsCode(e.target.value)}
            required
          />
          <button type="submit" disabled={statsLoading}>{statsLoading ? 'Loading...' : 'Get Stats'}</button>
        </form>
        {statsError && <div className="error">{statsError}</div>}
        {stats && (
          <div className="stats">
            <p>Original URL: <a href={stats.url} target="_blank" rel="noopener noreferrer">{stats.url}</a></p>
            <p>Created at: {stats.created_at}</p>
            <p>Expires at: {stats.expiry}</p>
            <p>Hits: {stats.hits}</p>
            {Array.isArray(stats.clicks) && stats.clicks.length > 0 && (
              <div style={{marginTop: '1rem'}}>
                <strong>Click Details:</strong>
                <table style={{width: '100%', marginTop: '0.5rem', fontSize: '0.97rem', borderCollapse: 'collapse'}}>
                  <thead>
                    <tr>
                      <th style={{textAlign: 'left', borderBottom: '1px solid #ccc'}}>Timestamp</th>
                      <th style={{textAlign: 'left', borderBottom: '1px solid #ccc'}}>Referrer</th>
                      <th style={{textAlign: 'left', borderBottom: '1px solid #ccc'}}>Location</th>
                    </tr>
                  </thead>
                  <tbody>
                    {stats.clicks.map((click: any, idx: number) => (
                      <tr key={idx}>
                        <td>{click.timestamp}</td>
                        <td>{click.referrer || '-'}</td>
                        <td>{click.location}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  )
}

export default App
