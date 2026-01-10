// @ts-check
import withNuxt from './.nuxt/eslint.config.mjs'
import eslintPluginPrettierRecommended from 'eslint-plugin-prettier/recommended'

export default withNuxt(
  // Global rules for all files.
  {
    rules: {
      // Warn on console.log in production code.
      'no-console': ['warn', { allow: ['warn', 'error'] }],
      // Disable multi-word component names rule for flexibility.
      'vue/multi-word-component-names': 'off'
    }
  },
  // Test files can use console freely.
  {
    files: ['**/*.test.ts', '**/*.spec.ts'],
    rules: {
      'no-console': 'off'
    }
  },
  // Prettier integration (must be last to override other formatting rules).
  eslintPluginPrettierRecommended
)
