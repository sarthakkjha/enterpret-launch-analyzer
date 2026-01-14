import { useState } from 'react'
import './App.css'
import FileUpload from './components/FileUpload'
import Dashboard from './components/Dashboard'

function App() {
  const [preLaunchFile, setPreLaunchFile] = useState(null)
  const [postLaunchFile, setPostLaunchFile] = useState(null)
  const [analysisResult, setAnalysisResult] = useState(null)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState(null)

  const API_BASE = import.meta.env.VITE_API_BASE || 'http://localhost:8080/api'

  const handleAnalyze = async () => {
    if (!preLaunchFile || !postLaunchFile) {
      setError('Please upload both pre-launch and post-launch CSV files')
      return
    }

    setIsLoading(true)
    setError(null)

    try {
      // Upload files
      const formData = new FormData()
      formData.append('preLaunch', preLaunchFile)
      formData.append('postLaunch', postLaunchFile)

      const uploadRes = await fetch(`${API_BASE}/upload`, {
        method: 'POST',
        body: formData
      })

      if (!uploadRes.ok) {
        const errData = await uploadRes.json()
        throw new Error(errData.error || 'Upload failed')
      }

      // Run analysis
      const analyzeRes = await fetch(`${API_BASE}/analyze`, {
        method: 'POST'
      })

      if (!analyzeRes.ok) {
        const errData = await analyzeRes.json()
        throw new Error(errData.error || 'Analysis failed')
      }

      const result = await analyzeRes.json()
      setAnalysisResult(result)
    } catch (err) {
      setError(err.message)
    } finally {
      setIsLoading(false)
    }
  }

  const handleReset = () => {
    setPreLaunchFile(null)
    setPostLaunchFile(null)
    setAnalysisResult(null)
    setError(null)
  }

  return (
    <div className="app">
      <header className="header">
        <div className="logo">
          <div className="logo-icon">E</div>
          <div>
            <div className="logo-text">Enterpret</div>
            <div className="logo-subtitle">Launch Impact Analyzer</div>
          </div>
        </div>
        {analysisResult && (
          <button className="btn btn-secondary" onClick={handleReset}>
            ← New Analysis
          </button>
        )}
      </header>

      <main className="main-content">
        {!analysisResult ? (
          <section className="upload-section">
            <h1 className="section-title">Analyze Your Launch Impact</h1>
            <p className="section-subtitle">
              Upload pre-launch and post-launch customer feedback CSVs to discover how your feature launch performed using AI-powered analysis.
            </p>

            <div className="glass-card">
              <div className="upload-grid">
                <FileUpload
                  label="Pre-Launch Reviews"
                  icon="📋"
                  file={preLaunchFile}
                  onFileSelect={setPreLaunchFile}
                  subtitle="Customer feedback before the feature launch"
                />
                <FileUpload
                  label="Post-Launch Reviews"
                  icon="🚀"
                  file={postLaunchFile}
                  onFileSelect={setPostLaunchFile}
                  subtitle="Customer feedback after the feature launch"
                />
              </div>

              {error && (
                <div className="error-message">
                  ⚠️ {error}
                </div>
              )}

              <button
                className="btn btn-primary btn-full"
                onClick={handleAnalyze}
                disabled={!preLaunchFile || !postLaunchFile || isLoading}
              >
                {isLoading ? (
                  <>
                    <span className="spinner"></span>
                    Analyzing with AI...
                  </>
                ) : (
                  <>
                    ✨ Analyze Launch Impact
                  </>
                )}
              </button>
            </div>
          </section>
        ) : (
          <Dashboard data={analysisResult} />
        )}
      </main>
    </div>
  )
}

export default App
