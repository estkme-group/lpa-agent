import { cloneElement, createContext, FC, PropsWithChildren, ReactElement, useContext, useState } from 'react'

interface OverlayController {
  open(element: ReactElement): void
  close(): void
}

const OverlayContext = createContext<OverlayController | null>(null)

export const OverlayProvider: FC<PropsWithChildren> = (props) => {
  const [element, setElement] = useState<ReactElement | null>(null)
  const controller: OverlayController = {
    open(element) {
      setElement(element)
    },
    close() {
      setElement(null)
    },
  }
  return (
    <OverlayContext.Provider value={controller}>
      {props.children}
      {element}
    </OverlayContext.Provider>
  )
}

export const useOverlay = (): OverlayController => {
  const context = useContext(OverlayContext)
  if (context == null) {
    throw new Error('overlay uninitialized')
  }
  return context
}
