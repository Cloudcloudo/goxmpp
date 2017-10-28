import React from 'react';

import LoginForm from '../../components/Login/LoginForm';
import {TextBanner} from '../../components/TextBanner/TextBanner';

export const LoginPage = () => {
  return (
    <div className="page">
      <TextBanner title="Logging Page" />
      <LoginForm />
    </div>
  );
}
