import React from 'react';

import './MessageWindow.scss'

export class MessageWindow extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      ws: null,
      newMessage: '',
      messages: []
    }
    this.scrollToBottom = this.scrollToBottom.bind(this);
    this.pushMessage = this.pushMessage.bind(this);
    this.sendMessage = this.sendMessage.bind(this);
    this.handleReadMessage = this.handleReadMessage.bind(this);
    this.renderMsg = this.renderMsg.bind(this);
    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  componentDidMount() {
    let ws = new WebSocket('ws://localhost:8080/stream?token=' + this.props.userData.token);
    ws.addEventListener('message', this.handleReadMessage);
    let status = {type:"presence",
      payload:{
        from: this.props.userData.jid
      }
    }
    let away = {type:"presence",
      payload:{
        type: 'unavailable',
        from: this.props.userData.jid
      }
    }
    ws.onopen = () => ws.send(JSON.stringify(status));
    ws.onclose = () => ws.send(JSON.stringify(away));
    this.setState({ws: ws});
    this.scrollToBottom();
  }

  render() {
    console.log();
    return (
      <div className='message-window'>
        {this.renderMsg()}
        <form onSubmit={this.handleSubmit}>
          <input className="msg-box" value={this.state.newMessage} onChange={this.handleChange} />
        </form>
      </div>
    );
  }

  renderMsg() {
    let messages = '';
    if (this.state.messages) {
      messages = this.state.messages.map((message, index) => {
        return (
          <li key={message.id}>
            <span className={(message.from) ? this.props.theme : 'self'}>
              {message.body}
            </span>
          </li>
        );
      });
    }
    return (
      <ul className="messages" ref={(el) => { this.messagesEnd = el; }}>
        {messages}
      </ul>
    );
  }

  handleSubmit(event) {
    event.preventDefault();

    if (this.state.newMessage === ''){
      return;
    }

    let msg = {
      id: Math.random().toString(36).substring(7),
      type: 'chat',
      to: 'alicja@localhost/',
      body: this.state.newMessage
    }
    // connect and push messages
    // TODO push messages
    this.sendMessage('message', msg);

    this.setState(prevState => ({
      messages: [...this.state.messages, msg],
      newMessage: ''
    }), () => this.scrollToBottom());
  }

  handleChange(event) {
    this.setState({ newMessage: event.target.value });
  }

  sendMessage(type, msg) {
    console.log({type: type, payload: msg});

    this.state.ws.send(JSON.stringify({type: type, payload: msg}));
  }

  pushMessage(msg) {
    this.setState(prevState => ({
      messages: [...this.state.messages, msg]
    }), () => this.scrollToBottom());
  }

  handleReadMessage(event) {
    let msg = JSON.parse(event.data);
    if (msg.type == 'message') {

      if (msg.payload.body) {
        this.pushMessage(msg.payload);
      }
    }
    console.log(event);
  }

  scrollToBottom() {
    this.messagesEnd.scrollTop = this.messagesEnd.scrollHeight;
  }
}
