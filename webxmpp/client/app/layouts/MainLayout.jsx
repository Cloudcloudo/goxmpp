import React from 'react';
import PropTypes from 'prop-types';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';

import muiTheme from './MainMuiTheme';
import Nav from '../components/Nav/Nav';

import './MainLayout.scss';

const example_link = {
  name: 'Login',
  url: '/login',
  title: '',
  action: ''
}

export const MainLayout = (props) => {
  return (
    <MuiThemeProvider muiTheme={muiTheme}>
      <section className="main-page-wrapper">
        <Nav links={[example_link]} />
        { props.children }
      </section>
    </MuiThemeProvider>
  );
}

MainLayout.propTypes = {
  children: PropTypes.element.isRequired
}
