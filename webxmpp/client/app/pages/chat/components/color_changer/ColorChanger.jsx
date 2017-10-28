import React from 'react';

import './ColorChanger.scss'

export class ColorChanger extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      className: 'default'
    }
    this.handleClick = this.handleClick.bind(this);
  }

  render(){
    return (
        <ul className='color-changer'>
          <li data-color="red" className='red' onClick={this.handleClick}></li>
          <li data-color="blue" className='blue' onClick={this.handleClick}></li>
          <li data-color="purple" className='purple' onClick={this.handleClick}></li>
          <li data-color="green" className='green' onClick={this.handleClick}></li>
        </ul>
    );
  }

  handleClick(event) {
    this.props.onColorChange(event.currentTarget.dataset.color);
  }
}
