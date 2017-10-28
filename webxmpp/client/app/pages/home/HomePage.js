import React from 'react';

import {TextBanner} from '../../components/TextBanner/TextBanner';
import RegisterForm from '../../components/Register/RegisterForm';

export const HomePage = () => {
  return (
    <div className="page">
      <TextBanner title="Register" />
      <RegisterForm />
    </div>
  );
}
