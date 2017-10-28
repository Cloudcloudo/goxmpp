import React from 'react';
import PropTypes from 'prop-types';
import TextField from 'material-ui/TextField'
import RaisedButton from 'material-ui/RaisedButton';
import _ from 'lodash';

import history from '../../bundles/history';
import {jidTaken, createUser} from './rest';
import './RegisterForm.scss';

export default class RegisterForm extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      jid: '',
      password: '',
      confirmPassword: '',
      jidError: '',
      passwordError: '',
      confirmPasswordError: '',
      submitDisabled: true,
    };
    this.checkJidtaken = _.debounce(this.checkJidtaken.bind(this), 1000);

    this.validate = this.validate.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.handleJidChange = this.handleJidChange.bind(this);
    this.handlePassChange = this.handlePassChange.bind(this);
    this.handleRepeatPassChange = this.handleRepeatPassChange.bind(this);
  }

  render() {
    return (
      <form className="register-form" onSubmit={this.handleSubmit}>
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
          <TextField
            className="password"
            floatingLabelText="Repeat password *"
            type="password"
            errorText={this.state.confirmPasswordError}
            value={this.state.confirmPassword}
            onChange={this.handleRepeatPassChange}
          />
        </div>
        <span className="required-message">* Field required</span>

        <div className="jid-wrapper">
          <RaisedButton
            className="register-button"
            primary
            type="submit"
            label="Register me"
            disabled={this.state.submitDisabled}
          />
        </div>
      </form>
    );
  }
  checkJidtaken(jid) {
    jidTaken(jid).then((taken) => {
      if (taken) {
        this.setState({ jidError: 'JID taken' }, this.validate());
      } else {
        this.setState({ jidError: '' }, this.validate());
      }
    });
  }

  handleJidChange(event) {
    this.setState({ jid: event.target.value }, this.checkJidtaken(event.target.value));
  }

  handlePassChange(event) {
    this.setState({ password: event.target.value }, this.validate());
  }
  handleRepeatPassChange(event) {
    let confirmPasswordError;
    if (this.state.password != event.target.value) {
      confirmPasswordError = 'Passwords do not match';
    } else {
      confirmPasswordError = '';
    }
    this.setState({ confirmPassword: event.target.value, confirmPasswordError: confirmPasswordError }, this.validate());
  }

  handleSubmit(event) {
    let userData = {
      jid: this.state.jid,
      password: this.state.password
    }
    createUser(userData).then((created) => {
      if (created) {
        history.push('/login');
      } else {
        this.setState({ jidError: 'Failed to create user' }, this.validate());
      }
    });
    event.preventDefault();
  }

  validate() {
    if ((this.state.jidError == '' && this.state.jid.length >= 3)
        || ( this.state.confirmPasswordError == '' && this.state.confirmPassword.length >= 6 )
        || ( this.state.passwordError == '' && this.state.password.length >= 6 )) {

      this.setState({ submitDisabled: false });
    } else {
      this.setState({ submitDisabled: true });
    }
  }
}
