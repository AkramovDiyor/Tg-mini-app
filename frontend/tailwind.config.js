/** @type {import('tailwindcss').Config} */
export default {
    content: ['./index.html', './src/**/*.{js,jsx}'],
    theme: {
      extend: {
        fontFamily: {
          sans: ['Inter', 'system-ui', '-apple-system', 'Segoe UI', 'Roboto', 'sans-serif'],
        },
        boxShadow: {
          card: '0 1px 2px rgba(16,24,40,.05), 0 4px 16px -4px rgba(16,24,40,.07)',
        },
        keyframes: {
          'fade-in': { from: { opacity: '0' }, to: { opacity: '1' } },
          'fade-up': { from: { opacity: '0', transform: 'translateY(10px)' }, to: { opacity: '1', transform: 'translateY(0)' } },
          'slide-up': { from: { transform: 'translateY(100%)' }, to: { transform: 'translateY(0)' } },
          'pop-in': { '0%': { transform: 'scale(.6)', opacity: '0' }, '100%': { transform: 'scale(1)', opacity: '1' } },
          'toast-in': { from: { opacity: '0', transform: 'translateY(14px) scale(.96)' }, to: { opacity: '1', transform: 'translateY(0) scale(1)' } },
        },
        animation: {
          'fade-in': 'fade-in .25s ease-out both',
          'fade-up': 'fade-up .3s ease-out both',
          'slide-up': 'slide-up .35s cubic-bezier(.16,1,.3,1) both',
          'pop-in': 'pop-in .35s cubic-bezier(.34,1.56,.64,1) both',
          'toast-in': 'toast-in .3s ease-out both',
        },
      },
    },
    plugins: [],
  }