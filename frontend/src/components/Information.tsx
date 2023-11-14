import { FC, useEffect } from 'react'
import Form from 'react-bootstrap/Form'
import Offcanvas from 'react-bootstrap/Offcanvas'
import { useSelector } from 'react-redux'
import { useAppDispatch } from '../hooks'
import { fetchInformation, selectInformation } from '../reducers/euicc'
import { FormattedEID } from './Formatted'

interface InformationProps {
  onClose(): void
}

export const Information: FC<InformationProps> = ({ onClose }) => {
  const dispatch = useAppDispatch()
  const information = useSelector(selectInformation)
  useEffect(() => {
    dispatch(fetchInformation())
  }, [dispatch])
  return (
    <Offcanvas show onHide={onClose}>
      <Offcanvas.Header closeButton>
        <Offcanvas.Title>Card Information</Offcanvas.Title>
      </Offcanvas.Header>
      <Offcanvas.Body>
        <Form>
          <Form.Group className='mb-3'>
            <Form.Label>EID</Form.Label>
            <Form.Control as='div' plaintext className='font-monospace'>
              <FormattedEID eid={information.eid} />
            </Form.Control>
          </Form.Group>
          <Form.Group className='mb-3'>
            <Form.Label>Default SM-DS</Form.Label>
            <Form.Control as='div' plaintext className='font-monospace'>
              {information.default_smds}
            </Form.Control>
          </Form.Group>
          <Form.Group className='mb-3' hidden={information.default_smdp === undefined}>
            <Form.Label>Default SM-DP</Form.Label>
            <Form.Control as='div' plaintext className='font-monospace'>
              {information.default_smdp}
            </Form.Control>
          </Form.Group>
        </Form>
      </Offcanvas.Body>
    </Offcanvas>
  )
}
