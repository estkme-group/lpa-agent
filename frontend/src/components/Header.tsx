import { FC } from 'react'
import Container from 'react-bootstrap/Container'
import Nav from 'react-bootstrap/Nav'
import Navbar from 'react-bootstrap/Navbar'
import { Information } from './Information'
import { useOverlay } from './Overlay'
import { DownloadProfile } from './Profile/Download'

export const Header: FC = () => {
  const overlay = useOverlay()
  const onDownload = () => {
    overlay.open(<DownloadProfile onClose={overlay.close} />)
  }
  const onCardInfo = () => {
    overlay.open(<Information onClose={overlay.close} />)
  }
  return (
    <Navbar className='bg-body-tertiary'>
      <Container fluid='md'>
        <Navbar.Brand>eSIM LPA Agent</Navbar.Brand>
        <Navbar.Toggle />
        <Navbar.Collapse>
          <Nav className='me-auto'>
            <Nav.Link onClick={onDownload}>Download</Nav.Link>
            <Nav.Link onClick={onCardInfo}>Card Info</Nav.Link>
          </Nav>
        </Navbar.Collapse>
      </Container>
    </Navbar>
  )
}
