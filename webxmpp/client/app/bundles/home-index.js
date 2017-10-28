import React from 'react';
import { render } from 'react-dom';
import { Router, Route, Switch } from 'react-router-dom';

import {HomePage} from '../pages/home/HomePage';
import {LoginPage} from '../pages/login/LoginPage';
import {ChatPage} from '../pages/chat/ChatPage';
import {NotFoundPage} from '../pages/not_found/NotFoundPage';
import {MainLayout} from '../layouts/MainLayout';

import history from './history';

const container = document.getElementById('app');

render ((
  <Router history={history}>
    <MainLayout>
      <Switch>
        <Route path='/' exact component={HomePage} />
        <Route path='/login' exact component={LoginPage} />
        <Route path='/chat' exact component={ChatPage} />
        <Route component={NotFoundPage} />
      </Switch>
    </MainLayout>
  </Router>
), container);
