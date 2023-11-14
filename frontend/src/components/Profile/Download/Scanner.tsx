import { css } from '@emotion/react'
import { BarcodeFormat } from '@zxing/library'
import { FC } from 'react'
import { DecodeHintType, useZxing } from 'react-zxing'
import { DownloadRequest } from './types'
import { parseQRCode } from './utils'

interface ScannerProps {
  parsed: boolean
  onDetect(request: DownloadRequest): void
}

export const Scanner: FC<ScannerProps> = (props) => {
  const zxing = useZxing({
    paused: props.parsed,
    hints: new Map([
      [DecodeHintType.POSSIBLE_FORMATS, [BarcodeFormat.QR_CODE]],
    ]),
    onDecodeResult: (result) => {
      const matches = parseQRCode(result.getText())
      if (!matches) return
      const [smdp, matchingId] = matches
      props.onDetect({ smdp, matchingId })
    },
  })
  return <video ref={zxing.ref} css={css`width: 100%`} />
}
