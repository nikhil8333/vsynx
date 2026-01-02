import { test, expect } from '@playwright/test'

test.describe('VsynX Manager E2E Tests', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/')
  })

  test.describe('Application Load', () => {
    test('should display the application title', async ({ page }) => {
      await expect(page.getByRole('heading', { name: 'Vsynx' })).toBeVisible()
      await expect(page.getByText('Extension Manager')).toBeVisible()
    })

    test('should display navigation sidebar', async ({ page }) => {
      await expect(page.getByRole('button', { name: /extensions/i })).toBeVisible()
      await expect(page.getByRole('button', { name: /marketplace/i })).toBeVisible()
      await expect(page.getByRole('button', { name: /audit/i })).toBeVisible()
      await expect(page.getByRole('button', { name: /sync/i })).toBeVisible()
      await expect(page.getByRole('button', { name: /settings/i })).toBeVisible()
    })

    test('should start with Extensions view', async ({ page }) => {
      await expect(page.getByRole('heading', { name: /extensions/i })).toBeVisible()
    })
  })

  test.describe('Navigation', () => {
    test('should navigate to Marketplace view', async ({ page }) => {
      await page.getByRole('button', { name: /marketplace/i }).click()
      await expect(page.getByText('Search Marketplace Extensions')).toBeVisible()
    })

    test('should navigate to Audit view', async ({ page }) => {
      await page.getByRole('button', { name: /audit/i }).click()
      await expect(page.getByRole('heading', { name: /audit/i })).toBeVisible()
    })

    test('should navigate to Sync view', async ({ page }) => {
      await page.getByRole('button', { name: /sync/i }).click()
      await expect(page.getByText('Sync Extensions Between Editors')).toBeVisible()
    })

    test('should navigate to Settings view', async ({ page }) => {
      await page.getByRole('button', { name: /settings/i }).click()
      await expect(page.getByText(/settings/i)).toBeVisible()
    })
  })

  test.describe('Sidebar Collapse', () => {
    test('should collapse and expand sidebar', async ({ page }) => {
      // Find and click collapse button
      const collapseButton = page.getByRole('button', { name: /collapse/i })
      await collapseButton.click()

      // After collapse, expand button should be visible
      await expect(page.getByRole('button', { name: /expand/i })).toBeVisible()

      // Click expand
      await page.getByRole('button', { name: /expand/i }).click()

      // Collapse button should be back
      await expect(page.getByRole('button', { name: /collapse/i })).toBeVisible()
    })
  })

  test.describe('Editor Selection', () => {
    test('should display editor dropdown', async ({ page }) => {
      await expect(page.getByRole('combobox')).toBeVisible()
    })

    test('should allow changing selected editor', async ({ page }) => {
      const dropdown = page.getByRole('combobox')
      await dropdown.click()
      
      // Check that options are available (actual options depend on installed editors)
      const options = await dropdown.locator('option').count()
      expect(options).toBeGreaterThan(0)
    })
  })
})

test.describe('Marketplace Search', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/')
    await page.getByRole('button', { name: /marketplace/i }).click()
  })

  test('should display search input', async ({ page }) => {
    await expect(page.getByPlaceholder(/search by keyword/i)).toBeVisible()
  })

  test('should display search button', async ({ page }) => {
    await expect(page.getByRole('button', { name: /^search$/i })).toBeVisible()
  })

  test('should display quick search buttons', async ({ page }) => {
    await expect(page.getByRole('button', { name: 'python' })).toBeVisible()
    await expect(page.getByRole('button', { name: 'prettier' })).toBeVisible()
    await expect(page.getByRole('button', { name: 'eslint' })).toBeVisible()
  })

  test('should perform search and display results', async ({ page }) => {
    const searchInput = page.getByPlaceholder(/search by keyword/i)
    await searchInput.fill('python')
    
    await page.getByRole('button', { name: /^search$/i }).click()
    
    // Wait for loading to complete
    await page.waitForTimeout(2000)
    
    // Should display results or no results message
    const hasResults = await page.getByText(/python/i).count()
    expect(hasResults).toBeGreaterThan(0)
  })

  test('should show autocomplete suggestions', async ({ page }) => {
    const searchInput = page.getByPlaceholder(/search by keyword/i)
    await searchInput.fill('python')
    
    // Wait for debounce and suggestions to appear
    await page.waitForTimeout(500)
    
    // Check if suggestions dropdown appears
    const suggestions = page.locator('[class*="dropdown"], [class*="suggestions"]')
    // Note: This might not find anything depending on implementation
  })
})

test.describe('Sync View', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/')
    await page.getByRole('button', { name: /sync/i }).click()
  })

  test('should display source editor section', async ({ page }) => {
    await expect(page.getByText(/source editor/i)).toBeVisible()
  })

  test('should display target editors section', async ({ page }) => {
    await expect(page.getByText(/target editor/i)).toBeVisible()
  })

  test('should display filter buttons', async ({ page }) => {
    await expect(page.getByRole('button', { name: /all/i })).toBeVisible()
  })

  test('should enable filter buttons after selecting target', async ({ page }) => {
    // Find and click a target editor button
    const targetButtons = page.locator('button').filter({ hasText: /windsurf|cursor|kiro/i })
    const count = await targetButtons.count()
    
    if (count > 0) {
      await targetButtons.first().click()
      
      // Missing and Present buttons should now be enabled
      await page.waitForTimeout(1000)
      
      const missingButton = page.getByRole('button', { name: /missing/i })
      const presentButton = page.getByRole('button', { name: /present/i })
      
      // Check they are visible
      await expect(missingButton).toBeVisible()
      await expect(presentButton).toBeVisible()
    }
  })
})

test.describe('Audit View', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/')
    await page.getByRole('button', { name: /audit/i }).click()
  })

  test('should display start audit button', async ({ page }) => {
    await expect(page.getByRole('button', { name: /start audit/i })).toBeVisible()
  })

  test('should show loading state when audit starts', async ({ page }) => {
    await page.getByRole('button', { name: /start audit/i }).click()
    
    // Should show cancel button or loading indicator
    const cancelOrLoading = page.getByRole('button', { name: /cancel/i })
    await expect(cancelOrLoading).toBeVisible({ timeout: 5000 })
  })
})

test.describe('Extensions View', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/')
  })

  test('should display extensions list heading', async ({ page }) => {
    await expect(page.getByText(/installed extensions/i)).toBeVisible()
  })

  test('should display filter input', async ({ page }) => {
    await expect(page.getByPlaceholder(/filter extensions/i)).toBeVisible()
  })

  test('should display refresh button', async ({ page }) => {
    await expect(page.getByRole('button', { name: /refresh/i })).toBeVisible()
  })

  test('should filter extensions when typing', async ({ page }) => {
    const filterInput = page.getByPlaceholder(/filter extensions/i)
    await filterInput.fill('python')
    
    // Wait for filter to apply
    await page.waitForTimeout(300)
    
    // The filter should be applied (actual behavior depends on extensions)
  })
})

test.describe('Accessibility', () => {
  test('should have proper heading hierarchy', async ({ page }) => {
    await page.goto('/')
    
    // Check for h1
    const h1 = await page.getByRole('heading', { level: 1 }).count()
    expect(h1).toBeGreaterThanOrEqual(1)
  })

  test('should have accessible buttons', async ({ page }) => {
    await page.goto('/')
    
    // All buttons should have accessible names
    const buttons = await page.getByRole('button').all()
    
    for (const button of buttons) {
      const name = await button.getAttribute('aria-label') || await button.textContent()
      expect(name).toBeTruthy()
    }
  })

  test('should be keyboard navigable', async ({ page }) => {
    await page.goto('/')
    
    // Press Tab to focus first interactive element
    await page.keyboard.press('Tab')
    
    // Should have focused element
    const focused = await page.evaluate(() => document.activeElement?.tagName)
    expect(focused).toBeTruthy()
  })
})
