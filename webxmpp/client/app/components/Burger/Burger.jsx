import React from 'react';
import PropTypes from 'prop-types';

import cn from 'classnames';
import './Burger.scss'


export const Burger = (props) => {
  return (
    <button className={ cn('c-hamburger', 'c-hamburger--htx', { 'is-active': props.active }) } onClick={props.onClick}>
      <span>toggle menu</span>
    </button>
  )
}

Burger.propTypes = {
  active: PropTypes.bool,
  onClick: PropTypes.func.isRequired
};

Burger.defaultProps = {
  active: false
};
