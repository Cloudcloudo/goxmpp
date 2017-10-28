import React from 'react';
import PropTypes from 'prop-types';
import {Link} from 'react-router-dom';

import {Burger} from '../Burger/Burger';
import {NavList} from './components/NavList/NavList';
import LogoImg from '../../assets/images/gomxpp.png';

import './Nav.scss';

export default class Nav extends React.Component {

  constructor(props){
    super(props);
    this.state = {
      active: false,
    }

    this.toggleMobileNavigation = this.toggleMobileNavigation.bind(this);
  }

  render(){
    return(
      <nav className="main-page-menu">
        <div className="logo-wrapper">
          <Link to='/' className="nav-brand">
            <img src={LogoImg} alt="temporary project logo" />
          </Link>
          <Burger active={this.state.active} onClick={this.toggleMobileNavigation} />
        </div>
        <NavList links={this.props.links} active={this.state.active} />
      </nav>

    )
  }

  toggleMobileNavigation(){
    this.setState({ active: !this.state.active });
  }

}

Nav.propTypes = {
  links: PropTypes.arrayOf(PropTypes.object)
}

Nav.defaultProps = {
  links: [{}]
}
