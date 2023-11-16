import { FC, useId } from 'react'
import Form from 'react-bootstrap/Form'

interface ConfirmControlProps {
  hidden?: boolean
  expected: string
  onChange(matched: boolean): void
}

export const ConfirmControl: FC<ConfirmControlProps> = (props) => {
  const id = useId()
  if (props.hidden) return
  return (
    <>
      <Form.Label htmlFor={id}>
        Type <q className='font-monospace'>{props.expected}</q> to confirm.
      </Form.Label>
      <Form.Control
        id={id}
        type='text'
        className='font-monospace'
        onChange={(event) => props.onChange(event.target.value === props.expected)}
      />
    </>
  )
}
