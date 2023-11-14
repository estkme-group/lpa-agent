import { ClipboardEvent, FC, useEffect, useState } from 'react'
import Form from 'react-bootstrap/Form'
import { DownloadRequest } from './types'
import { parseQRCode } from './utils'

interface ManuallyProps {
  onChange(request: DownloadRequest): void
}

export const Manually: FC<ManuallyProps> = ({ onChange }) => {
  const [smdp, setSMDP] = useState('')
  const [matchingId, setMatchingId] = useState('')
  const [imei, setIMEI] = useState('')
  const [confirmCode, setConfirmCode] = useState('')
  const onPaste = (event: ClipboardEvent) => {
    const matches = parseQRCode(event.clipboardData.getData('text'))
    if (matches) {
      event.preventDefault()
      const [smdp, matchingId] = matches
      setSMDP(smdp)
      setMatchingId(matchingId)
    }
  }
  useEffect(
    () => onChange({ smdp, matchingId, imei, confirmCode }),
    [smdp, matchingId, imei, confirmCode],
  )
  return (
    <Form>
      <Form.Group className='mb-3'>
        <Form.Label>SM-DP</Form.Label>
        <Form.Control
          type='text'
          value={smdp}
          onPaste={onPaste}
          onChange={(event) => setSMDP(event.target.value)}
        />
      </Form.Group>
      <Form.Group className='mb-3'>
        <Form.Label>Matching ID</Form.Label>
        <Form.Control
          type='text'
          value={matchingId}
          onPaste={onPaste}
          onChange={(event) => setMatchingId(event.target.value)}
        />
      </Form.Group>
      <Form.Group className='mb-3'>
        <Form.Label>IMEI</Form.Label>
        <Form.Control
          type='text'
          value={imei}
          onPaste={onPaste}
          onChange={(event) => setIMEI(event.target.value)}
        />
      </Form.Group>
      <Form.Group className='mb-3'>
        <Form.Label>Confirmation Code</Form.Label>
        <Form.Control
          type='text'
          value={confirmCode}
          onPaste={onPaste}
          onChange={(event) => setConfirmCode(event.target.value)}
        />
      </Form.Group>
    </Form>
  )
}
