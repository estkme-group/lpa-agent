import { css } from '@emotion/react'
import type { FC } from 'react'

interface FormattedICCIDProps {
  readonly iccid: string
}

export const FormattedICCID: FC<FormattedICCIDProps> = ({ iccid }) => {
  if (!iccid) return 'N/A'
  const segments = [
    iccid.slice(0, 2),
    iccid.slice(2, 5),
    iccid.slice(5, 9),
    iccid.slice(9, 18),
    iccid.slice(18),
  ]
  return (
    <span css={css`span:not(:first-of-type) { margin-left: 1ex }`}>
      {segments.map((segment, index) => <span key={index}>{segment}</span>)}
    </span>
  )
}

interface FormattedEIDProps {
  readonly eid: string
}

export const FormattedEID: FC<FormattedEIDProps> = ({ eid }) => {
  if (!eid) return 'N/A'
  const segments = [
    eid.slice(0, 2),
    eid.slice(2, 5),
    eid.slice(5, 13),
    eid.slice(13, 16),
    eid.slice(16, 18),
    eid.slice(18, 30),
    eid.slice(30),
  ]
  return (
    <span css={css`span:not(:first-of-type) { margin-left: 1ex }`}>
      {segments.map((segment, index) => <span key={index}>{segment}</span>)}
    </span>
  )
}
