import { FC, useId, useState } from 'react'
import Button from 'react-bootstrap/Button'
import Form from 'react-bootstrap/Form'
import Modal from 'react-bootstrap/Modal'
import { useAppDispatch } from '../../hooks'
import { deleteProfile } from '../../reducers/euicc'
import { Profile } from '../../reducers/euicc-types'
import { ConfirmControl } from '../ConfirmControl'
import { FormattedICCID } from '../Formatted'

interface DeleteProfileProps {
  readonly profile: Profile
  onClose(): void
}

export const DeleteProfile: FC<DeleteProfileProps> = ({ profile, onClose }) => {
  const dispatch = useAppDispatch()
  const [deleting, setDeleting] = useState(false)
  const [confirmed, setConfirmed] = useState(false)
  const onConfirm = async () => {
    if (profile.profileState === 1) return
    if (!confirmed) return
    setDeleting(true)
    await dispatch(deleteProfile(profile))
    onClose()
  }
  const enabled = profile.profileState === 1
  return (
    <Modal show onHide={onClose}>
      <Modal.Header closeButton>
        <Modal.Title>
          Delete <q>{profile.profileName}</q> Profile
        </Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <p className='text-danger'>
          This operation <b>CANNOT</b> be undone
        </p>
        <p>
          ICCID:{' '}
          <span className='font-monospace'>
            <FormattedICCID iccid={profile.iccid} />
          </span>
        </p>
        <p>
          Delete the profile <q className='font-monospace'>{profile.profileName}</q>
          <span hidden={!enabled}>
            , need <b>DISABLE</b>
          </span>.
        </p>
        <ConfirmControl
          expected={profile.profileName}
          hidden={profile.profileState === 1}
          onChange={setConfirmed}
        />
      </Modal.Body>
      <Modal.Footer>
        <Button
          variant='danger'
          onClick={onConfirm}
          hidden={profile.profileState === 1}
          disabled={!confirmed || deleting}
        >
          {deleting ? 'Deleting' : 'Delete'}
        </Button>
        <Button variant='secondary' onClick={onClose}>
          Close
        </Button>
      </Modal.Footer>
    </Modal>
  )
}
