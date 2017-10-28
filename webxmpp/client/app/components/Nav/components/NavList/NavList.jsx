import React from 'react';
import PropTypes from 'prop-types';
import cn from 'classnames';

import {NavLink} from '../NavLink/NavLink';
import './NavList.scss';

export const NavList = (props) => {
  return (
    <div className={ cn( 'navigation-list-wrapper', {'is-active' : props.active}) } id="main-navigation">
      <ul className="navigation-list">
        { props.links.map( (link, index) => ( <NavLink {...link} key={index} active={false} /> ) ) }
      </ul>
    </div>
  )
}

NavList.propTypes = {
  active: PropTypes.bool.isRequired,
  links: PropTypes.arrayOf(PropTypes.object).isRequired,
}
