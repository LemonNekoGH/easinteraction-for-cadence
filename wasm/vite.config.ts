import { defineConfig } from 'vite'
import dts from 'vite-plugin-dts'

export default defineConfig({
  plugins: [
    // 用于生成类型文件
    dts(),
  ],
  // 在 lib 模式下的特有配置
  build: {
    lib: {
      entry: './index.ts',
      name: 'EasiGen',
      fileName: (format) => `easi-gen.${format}.js`
    }
  }
})
