import { css } from '@emotion/react'
import type { FC } from 'react'
import Container from 'react-bootstrap/Container'
import Row from 'react-bootstrap/Row'
import { Header } from './components/Header'
import { ProfileList } from './components/Profile'

export const Entry: FC = () => (
  <>
    <Header />
    <Container fluid='md'>
      <Row css={css`margin-top: 1em`}>
        <ProfileList />
      </Row>
    </Container>
  </>
)
