import { Bar, Doughnut } from 'react-chartjs-2'
import {
    Chart as ChartJS,
    CategoryScale,
    LinearScale,
    BarElement,
    ArcElement,
    Title,
    Tooltip,
    Legend
} from 'chart.js'

ChartJS.register(
    CategoryScale,
    LinearScale,
    BarElement,
    ArcElement,
    Title,
    Tooltip,
    Legend
)

function Dashboard({ data }) {
    const { comparison, impact, pre_launch_reviews, post_launch_reviews } = data

    // Chart options
    const chartOptions = {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
            legend: {
                labels: {
                    color: '#a1a1aa',
                    font: { size: 12 }
                }
            }
        },
        scales: {
            x: {
                ticks: { color: '#71717a' },
                grid: { color: 'rgba(255,255,255,0.05)' }
            },
            y: {
                ticks: { color: '#71717a' },
                grid: { color: 'rgba(255,255,255,0.05)' }
            }
        }
    }

    const doughnutOptions = {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
            legend: {
                position: 'bottom',
                labels: {
                    color: '#a1a1aa',
                    font: { size: 12 },
                    padding: 20
                }
            }
        }
    }

    // Sentiment comparison chart data
    const sentimentChartData = {
        labels: ['Positive', 'Neutral', 'Negative'],
        datasets: [
            {
                label: 'Pre-Launch',
                data: [
                    comparison.pre_launch_sentiment.positive,
                    comparison.pre_launch_sentiment.neutral,
                    comparison.pre_launch_sentiment.negative
                ],
                backgroundColor: 'rgba(99, 102, 241, 0.7)',
                borderRadius: 6
            },
            {
                label: 'Post-Launch',
                data: [
                    comparison.post_launch_sentiment.positive,
                    comparison.post_launch_sentiment.neutral,
                    comparison.post_launch_sentiment.negative
                ],
                backgroundColor: 'rgba(16, 185, 129, 0.7)',
                borderRadius: 6
            }
        ]
    }

    // Post-launch sentiment distribution
    const postSentimentData = {
        labels: ['Positive', 'Neutral', 'Negative'],
        datasets: [{
            data: [
                comparison.post_launch_sentiment.positive,
                comparison.post_launch_sentiment.neutral,
                comparison.post_launch_sentiment.negative
            ],
            backgroundColor: [
                'rgba(16, 185, 129, 0.8)',
                'rgba(99, 102, 241, 0.8)',
                'rgba(239, 68, 68, 0.8)'
            ],
            borderWidth: 0
        }]
    }

    const sentimentShift = comparison.sentiment_shift
    const avgRatingChange = comparison.post_launch_sentiment.average_rating - comparison.pre_launch_sentiment.average_rating

    return (
        <div className="results-section">
            {/* Executive Summary */}
            <section className="executive-summary">
                <h1 className="section-title">Launch Impact Analysis</h1>
                <p className="section-subtitle">AI-powered analysis of your feature launch based on customer feedback</p>

                <div className="glass-card">
                    <div className="summary-header">
                        <div className="score-display">
                            <span className="score-value">{impact.success_score.toFixed(0)}</span>
                            <span className="score-label">/ 100</span>
                        </div>
                        <div className={`success-badge ${impact.overall_success ? 'success' : 'failure'}`}>
                            {impact.overall_success ? '✓ Launch Successful' : '✗ Needs Improvement'}
                        </div>
                    </div>
                    <p className="executive-text">{impact.executive_summary}</p>
                </div>
            </section>

            {/* Stats Grid */}
            <div className="stats-grid">
                <div className="glass-card stat-card">
                    <div className="stat-icon">📊</div>
                    <div className={`stat-value ${sentimentShift >= 0 ? 'positive' : 'negative'}`}>
                        {sentimentShift >= 0 ? '+' : ''}{sentimentShift.toFixed(1)}%
                    </div>
                    <div className="stat-label">Sentiment Shift</div>
                </div>
                <div className="glass-card stat-card">
                    <div className="stat-icon">⭐</div>
                    <div className={`stat-value ${avgRatingChange >= 0 ? 'positive' : 'negative'}`}>
                        {avgRatingChange >= 0 ? '+' : ''}{avgRatingChange.toFixed(1)}
                    </div>
                    <div className="stat-label">Avg Rating Change</div>
                </div>
                <div className="glass-card stat-card">
                    <div className="stat-icon">📋</div>
                    <div className="stat-value">{pre_launch_reviews.count}</div>
                    <div className="stat-label">Pre-Launch Reviews</div>
                </div>
                <div className="glass-card stat-card">
                    <div className="stat-icon">🚀</div>
                    <div className="stat-value">{post_launch_reviews.count}</div>
                    <div className="stat-label">Post-Launch Reviews</div>
                </div>
            </div>

            {/* Charts */}
            <div className="charts-grid">
                <div className="glass-card chart-card">
                    <h3 className="chart-title">📈 Sentiment Comparison</h3>
                    <div className="chart-container">
                        <Bar data={sentimentChartData} options={chartOptions} />
                    </div>
                </div>
                <div className="glass-card chart-card">
                    <h3 className="chart-title">🎯 Post-Launch Sentiment</h3>
                    <div className="chart-container">
                        <Doughnut data={postSentimentData} options={doughnutOptions} />
                    </div>
                </div>
            </div>

            {/* Themes */}
            <section className="glass-card" style={{ marginBottom: '48px' }}>
                <h3 className="chart-title">🏷️ Key Themes Analysis</h3>
                <div className="themes-grid">
                    {comparison.themes.map((theme, index) => (
                        <div key={index} className="theme-card">
                            <div className="theme-header">
                                <span className="theme-name">{theme.theme}</span>
                                <span className={`theme-change ${theme.change_rate >= 0 ? 'positive' : 'negative'}`}>
                                    {theme.change_rate >= 0 ? '+' : ''}{theme.change_rate.toFixed(0)}%
                                </span>
                            </div>
                            <div className="theme-stats">
                                <span>Pre: {theme.pre_count} mentions</span>
                                <span>Post: {theme.post_count} mentions</span>
                                <span>Sentiment: {theme.sentiment}</span>
                            </div>
                        </div>
                    ))}
                </div>
            </section>

            {/* Lists */}
            <div className="lists-grid">
                <div className="glass-card list-card improvements">
                    <h3>✅ Key Improvements</h3>
                    <ul>
                        {impact.key_improvements.map((item, i) => (
                            <li key={i}>{item}</li>
                        ))}
                    </ul>
                </div>
                <div className="glass-card list-card issues">
                    <h3>⚠️ Critical Issues</h3>
                    <ul>
                        {impact.critical_issues.length > 0 ? (
                            impact.critical_issues.map((item, i) => (
                                <li key={i}>{item}</li>
                            ))
                        ) : (
                            <li>No critical issues identified</li>
                        )}
                    </ul>
                </div>
                <div className="glass-card list-card">
                    <h3>💡 Recommendations</h3>
                    <ul>
                        {impact.recommendations.map((item, i) => (
                            <li key={i}>{item}</li>
                        ))}
                    </ul>
                </div>
            </div>
        </div>
    )
}

export default Dashboard
