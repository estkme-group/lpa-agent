import { FC, useState } from 'react'
import Button from 'react-bootstrap/Button'
import Form from 'react-bootstrap/Form'
import Modal from 'react-bootstrap/Modal'
import { useAppDispatch } from '../../hooks'
import { putProfile } from '../../reducers/euicc'
import type { Profile } from '../../reducers/euicc-types'
import { FormattedICCID } from '../Formatted'

interface Props {
  readonly profile: Profile
  onClose(): void
}

export const EditProfile: FC<Props> = ({ onClose, profile }) => {
  const dispatch = useAppDispatch()
  const [name, setName] = useState(profile.profileNickname ?? '')
  const [state, setState] = useState(profile.profileState)
  const [saving, setSaving] = useState(false)
  const onChange = async () => {
    setSaving(true)
    await dispatch(putProfile({
      ...profile,
      profileNickname: name,
      profileState: state,
    }))
    onClose()
  }
  return (
    <Modal show onHide={onClose}>
      <Modal.Header closeButton>
        <Modal.Title>
          Edit <q>{profile.profileNickname ?? profile.profileName}</q> Profile
        </Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <Form>
          <Form.Group className='mb-3'>
            <Form.Label>ICCID</Form.Label>
            <Form.Control as='div' type='text' plaintext>
              <FormattedICCID iccid={profile.iccid} />
            </Form.Control>
          </Form.Group>
          <Form.Group className='mb-3'>
            <Form.Label>Profile Name</Form.Label>
            <Form.Control type='text' value={name} onChange={(event) => setName(event.target.value)} />
          </Form.Group>
          <Form.Group className='mb-3'>
            <Form.Label>Profile State</Form.Label>
            <Form.Select value={state} onChange={(event) => setState(Number.parseInt(event.target.value, 10))}>
              <option value={1}>Enabled</option>
              <option value={0}>Disabled</option>
            </Form.Select>
          </Form.Group>
        </Form>
      </Modal.Body>
      <Modal.Footer>
        <Button variant='primary' onClick={onChange} disabled={saving}>
          {saving ? 'Saving' : 'Save Changes'}
        </Button>
        <Button variant='secondary' onClick={onClose}>
          Close
        </Button>
      </Modal.Footer>
    </Modal>
  )
}
