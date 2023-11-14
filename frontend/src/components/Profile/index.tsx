import { css } from '@emotion/react'
import classNames from 'classnames'
import { FC, useEffect } from 'react'
import Button from 'react-bootstrap/Button'
import ButtonGroup from 'react-bootstrap/ButtonGroup'
import Table from 'react-bootstrap/Table'
import { useSelector } from 'react-redux'
import { useAppDispatch } from '../../hooks'
import { fetchProfileList, selectProfileList } from '../../reducers/euicc'
import { Profile } from '../../reducers/euicc-types'
import { FormattedICCID } from '../Formatted'
import { useOverlay } from '../Overlay'
import { DeleteProfile } from './Delete'
import { EditProfile } from './Edit'

export const ProfileList: FC = () => {
  const dispatch = useAppDispatch()
  const profiles = useSelector(selectProfileList)
  useEffect(() => {
    dispatch(fetchProfileList())
  }, [dispatch])
  return (
    <Table responsive bordered hover>
      <thead>
        <tr>
          <td>#</td>
          <th css={css`& { width: 28ex }`}>ICCID</th>
          <th css={css`& { width: 100% }`}>Name</th>
          <th>Provider</th>
          <th>State</th>
          <th>Operations</th>
        </tr>
      </thead>
      <tbody>
        {profiles.map((profile, index) => <ProfileRow key={profile.iccid} profile={profile} index={index + 1} />)}
      </tbody>
    </Table>
  )
}

const ProfileRow: FC<{ profile: Profile; index: number }> = ({ profile, index }) => {
  const overlay = useOverlay()
  const onEdit = () => {
    overlay.open(
      <EditProfile
        profile={profile}
        onClose={overlay.close}
      />,
    )
  }
  const onDelete = () => {
    overlay.open(
      <DeleteProfile
        profile={profile}
        onClose={overlay.close}
      />,
    )
  }
  const enabled = profile.profileState === 1
  const className = classNames(
    'align-middle',
    enabled && 'table-active',
  )
  return (
    <tr className={className}>
      <td className='text-nowrap font-monospace'>{index}</td>
      <td className='text-nowrap font-monospace'>
        <FormattedICCID iccid={profile.iccid} />
      </td>
      <td className='text-nowrap'>
        {profile.profileNickname ?? profile.profileName}
      </td>
      <td className='text-nowrap'>
        {profile.serviceProviderName}
      </td>
      <td>{['Disabled', 'Enabled'][profile.profileState]}</td>
      <td className='text-nowrap'>
        <ButtonGroup size='sm'>
          <Button onClick={onEdit}>
            Edit
          </Button>
          <Button onClick={onDelete} variant='danger'>
            Delete
          </Button>
        </ButtonGroup>
      </td>
    </tr>
  )
}
