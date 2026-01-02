import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import App from './App'

describe('App Component', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('Initial Render', () => {
    it('renders the app header', async () => {
      render(<App />)
      expect(screen.getByText('Vsynx')).toBeInTheDocument()
      expect(screen.getByText('Extension Manager')).toBeInTheDocument()
    })

    it('renders navigation buttons', async () => {
      render(<App />)
      // Use title attribute which is more specific
      expect(screen.getByTitle('Extensions')).toBeInTheDocument()
      expect(screen.getByTitle('Marketplace')).toBeInTheDocument()
      expect(screen.getByTitle('Audit')).toBeInTheDocument()
      expect(screen.getByTitle('Sync')).toBeInTheDocument()
    })

    it('starts with Extensions view active', async () => {
      render(<App />)
      // The extensions view shows a filter input
      await waitFor(() => {
        expect(screen.getByPlaceholderText(/filter/i)).toBeInTheDocument()
      })
    })
  })

  describe('Navigation', () => {
    it('switches to Marketplace view when clicked', async () => {
      const user = userEvent.setup()
      render(<App />)
      
      await user.click(screen.getByTitle('Marketplace'))
      
      await waitFor(() => {
        expect(screen.getByPlaceholderText(/search by keyword/i)).toBeInTheDocument()
      })
    })

    it('switches to Audit view when clicked', async () => {
      const user = userEvent.setup()
      render(<App />)
      
      await user.click(screen.getByTitle('Audit'))
      
      await waitFor(() => {
        // Multiple elements contain "Security Audit", use getAllByText
        const auditElements = screen.getAllByText(/Security Audit/i)
        expect(auditElements.length).toBeGreaterThan(0)
      })
    })

    it('switches to Sync view when clicked', async () => {
      const user = userEvent.setup()
      render(<App />)
      
      await user.click(screen.getByTitle('Sync'))
      
      await waitFor(() => {
        // Multiple elements contain "source editor", use getAllByText
        const elements = screen.getAllByText(/source editor/i)
        expect(elements.length).toBeGreaterThan(0)
      })
    })
  })

  describe('Marketplace Search', () => {
    it('renders search input', async () => {
      const user = userEvent.setup()
      render(<App />)
      
      await user.click(screen.getByTitle('Marketplace'))
      
      await waitFor(() => {
        expect(screen.getByPlaceholderText(/search by keyword/i)).toBeInTheDocument()
      })
    })

    it('renders quick search buttons', async () => {
      const user = userEvent.setup()
      render(<App />)
      
      await user.click(screen.getByTitle('Marketplace'))
      
      await waitFor(() => {
        expect(screen.getByRole('button', { name: 'python' })).toBeInTheDocument()
      })
    })
  })

  describe('Sync View', () => {
    it('renders source and target editor sections', async () => {
      const user = userEvent.setup()
      render(<App />)
      
      await user.click(screen.getByTitle('Sync'))
      
      await waitFor(() => {
        // Multiple elements contain "source editor", use getAllByText
        const elements = screen.getAllByText(/source editor/i)
        expect(elements.length).toBeGreaterThan(0)
      })
    })
  })

  describe('Audit View', () => {
    it('renders start audit button when no report exists', async () => {
      const user = userEvent.setup()
      render(<App />)
      
      await user.click(screen.getByTitle('Audit'))
      
      await waitFor(() => {
        expect(screen.getByRole('button', { name: /start audit/i })).toBeInTheDocument()
      })
    })
  })
})

describe('Extension Filtering', () => {
  it('filter input is rendered in extensions view', async () => {
    render(<App />)
    
    await waitFor(() => {
      expect(screen.getByPlaceholderText(/filter/i)).toBeInTheDocument()
    })
  })
})

describe('Editor Selection', () => {
  it('renders editor dropdown', async () => {
    render(<App />)
    
    await waitFor(() => {
      expect(screen.getByRole('combobox')).toBeInTheDocument()
    })
  })
})
