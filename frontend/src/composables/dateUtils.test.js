import { describe, expect, it } from 'vitest'
import { formatLocalDate, nextLocalDate } from './dateUtils'

describe('dateUtils', () => {
  it('formats local dates without timezone rollback', () => {
    const date = new Date('2026-05-05T00:00:00')
    expect(formatLocalDate(date)).toBe('2026-05-05')
  })

  it('computes next day in local calendar space', () => {
    expect(nextLocalDate('2026-05-05')).toBe('2026-05-06')
  })

  it('returns empty string for invalid input', () => {
    expect(nextLocalDate('')).toBe('')
    expect(nextLocalDate('not-a-date')).toBe('')
  })
})
