import React from 'react';
import PropTypes from 'prop-types';
import {Link} from 'react-router-dom';
import cn from 'classnames';

import './NavLink.scss';

export const NavLink = (props) => {
  return (
    <li className={cn('navigation-list-element', { 'is-active': props.active })}>
      <Link to={props.url} title={props.title}>
        <span>{props.name}</span>
      </Link>
    </li>
  )
}

NavLink.propTypes = {
  url: PropTypes.string.isRequired,
  name: PropTypes.string.isRequired,
  active: PropTypes.bool.isRequired,
  title: PropTypes.string
}

NavLink.defaultProps = {
  title: ''
};
