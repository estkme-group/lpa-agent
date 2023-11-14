import { FC, useState } from 'react'
import Button from 'react-bootstrap/Button'
import Modal from 'react-bootstrap/Modal'
import Tab from 'react-bootstrap/Tab'
import Tabs from 'react-bootstrap/Tabs'
import { useAppDispatch } from '../../../hooks'
import { downloadProfile } from '../../../reducers/euicc'
import { Manually } from './Manually'
import { Scanner } from './Scanner'
import { DownloadRequest } from './types'

const KEY_MANUALLY = 'manually'
const KEY_SCAN = 'scan'

interface Props {
  onClose(): void
}

export const DownloadProfile: FC<Props> = ({ onClose }) => {
  const dispatch = useAppDispatch()
  const [downloading, setDownloading] = useState(false)
  const [activeKey, setActiveKey] = useState(KEY_MANUALLY)
  const [request, setRequest] = useState<DownloadRequest>({})
  const onQRCodeDetect = (request: DownloadRequest) => {
    setRequest(request)
    setActiveKey(KEY_MANUALLY)
  }
  const onDownload = async () => {
    setDownloading(true)
    await dispatch(downloadProfile(request))
    onClose()
  }
  return (
    <Modal show onHide={onClose}>
      <Modal.Header closeButton>
        <Modal.Title>Download Profile</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <Tabs activeKey={activeKey} onSelect={(key) => setActiveKey(key ?? KEY_MANUALLY)}>
          <Tab eventKey={KEY_MANUALLY} title='Manually'>
            <Manually onChange={setRequest} />
          </Tab>
          <Tab eventKey={KEY_SCAN} title='Scan'>
            <Scanner onDetect={onQRCodeDetect} parsed={activeKey != KEY_SCAN} />
          </Tab>
        </Tabs>
      </Modal.Body>
      <Modal.Footer>
        <Button variant='primary' onClick={onDownload} disabled={downloading}>
          {downloading ? 'Downloading' : 'Download'}
        </Button>
        <Button variant='secondary' onClick={onClose}>
          Close
        </Button>
      </Modal.Footer>
    </Modal>
  )
}
