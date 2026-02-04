import React, { Component } from 'react';
import { AlertTriangle, RefreshCw, Home } from 'lucide-react';

/**
 * ErrorBoundary catches JavaScript errors anywhere in its child component tree,
 * logs those errors, and displays a fallback UI.
 */
class ErrorBoundary extends Component {
    constructor(props) {
        super(props);
        this.state = { hasError: false, error: null, errorInfo: null };
    }

    static getDerivedStateFromError(error) {
        // Update state so the next render shows the fallback UI
        return { hasError: true, error };
    }

    componentDidCatch(error, errorInfo) {
        // Log error to console (could send to monitoring service)
        console.error('ErrorBoundary caught an error:', error, errorInfo);
        this.setState({ errorInfo });
    }

    handleRetry = () => {
        this.setState({ hasError: false, error: null, errorInfo: null });
    };

    handleGoHome = () => {
        window.location.href = '/dashboard';
    };

    render() {
        if (this.state.hasError) {
            // Fallback UI
            return (
                <div className="error-boundary">
                    <div className="error-boundary-content">
                        <div className="error-boundary-icon">
                            <AlertTriangle size={64} />
                        </div>
                        <h1>Oops! Something went wrong</h1>
                        <p className="error-boundary-message">
                            We're sorry, but something unexpected happened. Please try again.
                        </p>
                        {process.env.NODE_ENV === 'development' && this.state.error && (
                            <details className="error-boundary-details">
                                <summary>Error Details</summary>
                                <pre>{this.state.error.toString()}</pre>
                                {this.state.errorInfo && (
                                    <pre>{this.state.errorInfo.componentStack}</pre>
                                )}
                            </details>
                        )}
                        <div className="error-boundary-actions">
                            <button 
                                className="error-boundary-btn error-boundary-btn-primary"
                                onClick={this.handleRetry}
                            >
                                <RefreshCw size={18} />
                                Try Again
                            </button>
                            <button 
                                className="error-boundary-btn error-boundary-btn-secondary"
                                onClick={this.handleGoHome}
                            >
                                <Home size={18} />
                                Go to Dashboard
                            </button>
                        </div>
                    </div>
                    <style>{`
                        .error-boundary {
                            display: flex;
                            align-items: center;
                            justify-content: center;
                            min-height: 100vh;
                            background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
                            padding: 2rem;
                        }
                        .error-boundary-content {
                            text-align: center;
                            max-width: 500px;
                            background: rgba(255, 255, 255, 0.05);
                            backdrop-filter: blur(10px);
                            border-radius: 20px;
                            padding: 3rem;
                            border: 1px solid rgba(255, 255, 255, 0.1);
                        }
                        .error-boundary-icon {
                            color: #f59e0b;
                            margin-bottom: 1.5rem;
                        }
                        .error-boundary h1 {
                            color: #fff;
                            font-size: 1.75rem;
                            margin-bottom: 1rem;
                            font-weight: 600;
                        }
                        .error-boundary-message {
                            color: rgba(255, 255, 255, 0.7);
                            margin-bottom: 2rem;
                            line-height: 1.6;
                        }
                        .error-boundary-details {
                            text-align: left;
                            background: rgba(0, 0, 0, 0.3);
                            border-radius: 8px;
                            padding: 1rem;
                            margin-bottom: 2rem;
                            color: #ef4444;
                            font-size: 0.875rem;
                        }
                        .error-boundary-details summary {
                            cursor: pointer;
                            color: rgba(255, 255, 255, 0.7);
                            margin-bottom: 0.5rem;
                        }
                        .error-boundary-details pre {
                            overflow-x: auto;
                            white-space: pre-wrap;
                            word-break: break-word;
                        }
                        .error-boundary-actions {
                            display: flex;
                            gap: 1rem;
                            justify-content: center;
                            flex-wrap: wrap;
                        }
                        .error-boundary-btn {
                            display: inline-flex;
                            align-items: center;
                            gap: 0.5rem;
                            padding: 0.75rem 1.5rem;
                            border-radius: 10px;
                            font-weight: 500;
                            cursor: pointer;
                            transition: all 0.2s ease;
                            border: none;
                        }
                        .error-boundary-btn-primary {
                            background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
                            color: white;
                        }
                        .error-boundary-btn-primary:hover {
                            transform: translateY(-2px);
                            box-shadow: 0 10px 20px rgba(99, 102, 241, 0.3);
                        }
                        .error-boundary-btn-secondary {
                            background: rgba(255, 255, 255, 0.1);
                            color: white;
                            border: 1px solid rgba(255, 255, 255, 0.2);
                        }
                        .error-boundary-btn-secondary:hover {
                            background: rgba(255, 255, 255, 0.15);
                        }
                    `}</style>
                </div>
            );
        }

        return this.props.children;
    }
}

export default ErrorBoundary;
