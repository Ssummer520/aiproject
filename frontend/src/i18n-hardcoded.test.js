import { describe, expect, it } from 'vitest'
import { readFileSync, readdirSync, statSync } from 'node:fs'
import { join } from 'node:path'
import messages from './i18n.js'

const SRC_DIR = new URL('./', import.meta.url).pathname
const allowedLiteralPatterns = [
  /<span>ChinaTravel<\/span>/,
  /<option value="Hangzhou">Hangzhou<\/option>/,
  /<option value="Shanghai">Shanghai<\/option>/,
  /<option value="Beijing">Beijing<\/option>/,
  /<option value="Xi'an">Xi'an<\/option>/,
  /<option value="Chengdu">Chengdu<\/option>/,
]

function walk(dir) {
  return readdirSync(dir).flatMap(name => {
    const path = join(dir, name)
    if (name === 'dist' || name === 'node_modules') return []
    if (statSync(path).isDirectory()) return walk(path)
    return path.endsWith('.vue') ? [path] : []
  })
}

function stripAllowed(line) {
  return allowedLiteralPatterns.reduce((current, pattern) => current.replace(pattern, ''), line)
}

function hasMessageKey(localeMessages, key) {
  return key.split('.').every(part => (localeMessages = localeMessages?.[part]) !== undefined)
}

describe('i18n static copy coverage', () => {
  it('keeps Vue templates free of inline locale ternaries and hardcoded UI attributes', () => {
    const violations = []
    for (const file of walk(SRC_DIR)) {
      const lines = readFileSync(file, 'utf8').split('\n')
      for (const [index, rawLine] of lines.entries()) {
        const line = stripAllowed(rawLine)
        if (line.includes("locale === 'zh'") || line.includes("locale.value === 'zh'")) {
          violations.push(`${file}:${index + 1}: inline locale ternary`)
        }
        if (/(?<!:)\b(?:placeholder|title)="[A-Za-z\u4e00-\u9fff]/.test(line)) {
          violations.push(`${file}:${index + 1}: hardcoded translated attribute`)
        }
        if (/>\s*[\u4e00-\u9fff][^<{]*</.test(line)) {
          violations.push(`${file}:${index + 1}: hardcoded Chinese text`)
        }
      }
    }
    expect(violations).toEqual([])
  })

  it('renders backend demo text through localization helpers', () => {
    const forbidden = [
      /{{\s*product\.(?:name|subtitle|description|policy|usage|meeting_point)\s*}}/,
      /{{\s*d\.(?:name|city)\s*}}/,
      /{{\s*deal\.(?:title|description)\s*}}/,
      /{{\s*item\.(?:title|note|product_name|package_name|city)\s*}}/,
      /{{\s*merchant\.(?:name|city|status)\s*}}/,
      /{{\s*article\.(?:title|city|category|status)\s*}}/,
      /{{\s*review\.(?:content|merchant_reply)\s*}}/,
    ]
    const violations = []
    for (const file of walk(SRC_DIR)) {
      const lines = readFileSync(file, 'utf8').split('\n')
      for (const [index, line] of lines.entries()) {
        if (forbidden.some(pattern => pattern.test(line))) {
          violations.push(`${file}:${index + 1}: dynamic text bypasses localization helper`)
        }
      }
    }
    expect(violations).toEqual([])
  })

  it('defines every static translation key in both locales', () => {
    const keyPatterns = [
      /\$t\(['"]([^'"]+)['"]/g,
      /\bt\(['"]([^'"]+)['"]/g,
      /labelKey:\s*['"]([^'"]+)['"]/g,
      /nameKey:\s*['"]([^'"]+)['"]/g,
    ]
    const keys = new Set()

    for (const file of walk(SRC_DIR)) {
      const text = readFileSync(file, 'utf8')
      for (const pattern of keyPatterns) {
        let match
        while ((match = pattern.exec(text))) keys.add(match[1])
      }
    }

    const violations = []
    for (const key of keys) {
      for (const lang of ['en', 'zh']) {
        if (!hasMessageKey(messages.global.messages.value[lang], key)) {
          violations.push(`${lang}: ${key}`)
        }
      }
    }

    expect(violations).toEqual([])
  })

  it('does not expose missing auto translation keys in templates', () => {
    const violations = []
    for (const file of walk(SRC_DIR)) {
      const lines = readFileSync(file, 'utf8').split('\n')
      for (const [index, line] of lines.entries()) {
        if (/auto\.auto_[0-9a-f]{8}/.test(line) && !/\$t\(['"]auto\.auto_[0-9a-f]{8}['"]\)|\bt\(['"]auto\.auto_[0-9a-f]{8}['"]\)|\w+Key:\s*['"]auto\.auto_[0-9a-f]{8}['"]|\bt\([^)]*\|\|\s*['"]auto\.auto_[0-9a-f]{8}['"]\)/.test(line)) {
          violations.push(`${file}:${index + 1}: possible raw auto key`)
        }
      }
    }
    expect(violations).toEqual([])
  })

})
