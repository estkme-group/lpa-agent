import { ReactNode } from 'react'

export interface DialogProps {
  title: ReactNode
  children: ReactNode
  onConfirm(): void
  onClose(): void
}
