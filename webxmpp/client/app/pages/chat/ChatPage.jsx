import React from 'react';

import history from '../../bundles/history'
import {TextBanner} from '../../components/TextBanner/TextBanner';
import {ColorChanger} from './components/color_changer/ColorChanger';
import {MessageWindow} from './components/message_window/MessageWindow';

export class ChatPage extends React.Component {
  constructor(props) {
    super(props);
    if (!sessionStorage.user) {
      history.push('/login');
      this.state = {
        userData: {jid: ''},
        theme: 'default'
      }
    } else {
      this.state = {
        userData: JSON.parse(sessionStorage.user),
        theme: 'default'
      };
    }
    this.handleColorChange = this.handleColorChange.bind(this);
  }

  render(){
    return (
        <div>
          <MessageWindow userData={this.state.userData} theme={this.state.theme} />
          <ColorChanger onColorChange={this.handleColorChange} />
        </div>
    );
  }

  handleColorChange(color) {
    console.log(color);
    this.setState({theme: color});
  }
}
