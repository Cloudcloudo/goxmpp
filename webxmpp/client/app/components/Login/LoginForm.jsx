import React from 'react';
import PropTypes from 'prop-types';
import TextField from 'material-ui/TextField'
import RaisedButton from 'material-ui/RaisedButton';

import history from '../../bundles/history'
import {login} from './rest';
import './LoginForm.scss';

export default class LoginForm extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      jid: '',
      password: '',
      submitDisabled: true,
    };
    this.handleJidChange = this.handleJidChange.bind(this);
    this.handlePassChange = this.handlePassChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  render() {
    return (
      <form className="login-form" onSubmit={this.handleSubmit}>
        <div className="jid-wrapper">
          <TextField
            className="jid"
            name="jid"
            floatingLabelText="Jabber ID *"
            value={this.state.jid}
            errorText={this.state.jidError}
            onChange={this.handleJidChange}
          />
          <span className="domain">@localhost</span>
        </div>

        <div className="jid-wrapper">
          <TextField
            className="password"
            floatingLabelText="Password *"
            type="password"
            value={this.state.password}
            errorText={this.state.passwordError}
            onChange={this.handlePassChange}
          />
        </div>

        <div className="jid-wrapper">
          <RaisedButton
            className="login-button"
            primary
            type="submit"
            label="Log in"
            disabled={this.state.submitDisabled}
          />
        </div>
      </form>
    );
  }
  handleJidChange(event) {
    this.setState({ jid: event.target.value }, this.validate());
  }

  handlePassChange(event) {
    this.setState({ password: event.target.value }, this.validate());
  }

  handleSubmit(event) {
    let userData = {
      jid: this.state.jid,
      password: this.state.password
    }

    login(userData).then((loginResp) => {
      if (loginResp) {
        sessionStorage.setItem('user', JSON.stringify(loginResp));

        history.push('/chat');
      } else {
        this.setState({ jidError: 'Invalid user or password' }, this.validate());
      }
    });
    event.preventDefault();
  }

  validate() {
    if (this.state.jidError != '' || this.state.jid.length >= 3 || this.state.password.length >= 6 ) {
      this.setState({ submitDisabled: false });
    } else {
      this.setState({ submitDisabled: true });
    }
  }
}
