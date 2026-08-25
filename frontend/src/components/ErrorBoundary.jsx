import { Component } from 'react'

export class ErrorBoundary extends Component {
  constructor(props) {
    super(props)
    this.state = { hasError: false }
  }

  static getDerivedStateFromError() {
    return { hasError: true }
  }

  componentDidCatch(error, info) {
    console.error('ErrorBoundary:', error, info)
  }

  handleRetry = () => {
    this.setState({ hasError: false })
  }

  render() {
    if (this.state.hasError) {
      return (
        <div className="flex min-h-screen justify-center font-sans">
          <div className="relative flex min-h-screen w-full max-w-[420px] items-center justify-center bg-[#F9FAFB] px-6">
            <div className="text-center">
              <h2 className="mb-2 text-lg font-bold text-slate-800">Что-то сломалось</h2>
              <p className="mb-4 text-sm text-slate-500">
                Перезагрузите экран. Если ошибка повторится — откройте Mini App заново.
              </p>
              <button
                type="button"
                onClick={this.handleRetry}
                className="rounded-xl bg-emerald-600 px-5 py-2.5 text-sm font-bold text-white"
              >
                Попробовать снова
              </button>
            </div>
          </div>
        </div>
      )
    }

    return this.props.children
  }
}
