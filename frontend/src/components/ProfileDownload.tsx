import { FC, useState } from 'react'
import Button from 'react-bootstrap/Button'
import Form from 'react-bootstrap/Form'
import Modal from 'react-bootstrap/Modal'
import Tab from 'react-bootstrap/Tab'
import Tabs from 'react-bootstrap/Tabs'
import { useZxing } from 'react-zxing'
import { useAppDispatch } from '../hooks'
import { downloadProfile } from '../reducers/euicc'

const KEY_MANUALLY = 'manually'
const KEY_SCAN = 'scan'

interface Props {
  onClose(): void
}

export const DownloadProfile: FC<Props> = ({ onClose }) => {
  const dispatch = useAppDispatch()
  const [downloading, setDownloading] = useState(false)
  const [activeKey, setActiveKey] = useState(KEY_MANUALLY)
  const [smdp, setSMDP] = useState('')
  const [matchingId, setMatchingId] = useState('')
  const [imei, setIMEI] = useState('')
  const [code, setCode] = useState('')
  const zxing = useZxing({
    paused: activeKey != KEY_SCAN,
    onDecodeResult(result) {
      console.log(result.getText())
      const text = result.getText()
      if (!text.startsWith('LPA:')) return
      const segments = text.split('$')
      if (segments.length != 3) return
      setSMDP(segments[1])
      setMatchingId(segments[2])
      setActiveKey(KEY_MANUALLY)
    },
  })
  const onTabChange = (key: string | null) => {
    setActiveKey(key ?? KEY_MANUALLY)
  }
  const onDownload = async () => {
    setDownloading(true)
    await dispatch(downloadProfile({ smdp, matchingId, imei, confirmCode: code }))
    onClose()
  }
  return (
    <Modal show onHide={onClose}>
      <Modal.Header closeButton>
        <Modal.Title>Download Profile</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <Tabs activeKey={activeKey} onSelect={onTabChange}>
          <Tab eventKey={KEY_MANUALLY} title='Manually'>
            <Form>
              <Form.Group className='mb-3'>
                <Form.Label>SM-DP</Form.Label>
                <Form.Control type='text' value={smdp} onChange={(event) => setSMDP(event.target.value)} />
              </Form.Group>
              <Form.Group className='mb-3'>
                <Form.Label>Matching ID</Form.Label>
                <Form.Control
                  type='text'
                  value={matchingId}
                  onChange={(event) => setMatchingId(event.target.value)}
                />
              </Form.Group>
              <Form.Group className='mb-3'>
                <Form.Label>IMEI</Form.Label>
                <Form.Control
                  type='text'
                  value={imei}
                  onChange={(event) => setIMEI(event.target.value)}
                />
              </Form.Group>
              <Form.Group className='mb-3'>
                <Form.Label>Confirmation Code</Form.Label>
                <Form.Control
                  type='text'
                  value={code}
                  onChange={(event) => setCode(event.target.value)}
                />
              </Form.Group>
            </Form>
          </Tab>
          <Tab eventKey={KEY_SCAN} title='Scan'>
            <video ref={zxing.ref} width='100%' />
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
