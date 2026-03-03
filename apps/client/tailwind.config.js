/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ['IBM Plex Sans Thai', 'Inter', 'system-ui', '-apple-system', 'sans-serif'],
      },
      colors: {
        messenger: {
          50: '#f0f4ff',
          100: '#e0e9ff',
          200: '#c7d7fe',
          300: '#a5bbfc',
          400: '#8196f8',
          500: '#6366f1',
          600: '#4f46e5',
          700: '#4338ca',
          800: '#3730a3',
          900: '#312e81',
        },
        purple: {
          50: '#faf5ff',
          100: '#f3e8ff',
          200: '#e9d5ff',
          300: '#d8b4fe',
          400: '#c084fc',
          500: '#a855f7',
          600: '#9333ea',
          700: '#7e22ce',
          800: '#6b21a8',
          900: '#581c87',
        }
      },
      backgroundImage: {
        'messenger-gradient': 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
        'messenger-gradient-alt': 'linear-gradient(135deg, #6366f1 0%, #a855f7 100%)',
        'messenger-card': 'linear-gradient(135deg, #f8faff 0%, #f3f4ff 100%)',
      }
    },
  },
  plugins: [],
}
