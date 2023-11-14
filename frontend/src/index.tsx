import ready from 'domready'
import { createRoot } from 'react-dom/client'
import { Provider as ReduxProvider } from 'react-redux'
import { OverlayProvider } from './components/Overlay'
import { Entry } from './Entry'
import { store } from './store'

import 'bootstrap/dist/css/bootstrap.min.css'
import { purge } from './reducers/euicc'

ready(() => {
  const container = document.createElement('main')
  document.body = document.createElement('body')
  document.body.append(container)
  const root = createRoot(container)
  root.render(
    <ReduxProvider store={store}>
      <OverlayProvider>
        <Entry />
      </OverlayProvider>
    </ReduxProvider>,
  )
})

Reflect.set(globalThis, 'EUICC', {
  purge: () => store.dispatch(purge()),
})
