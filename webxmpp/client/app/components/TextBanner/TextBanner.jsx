import React from 'react';
import PropTypes from 'prop-types';

import './TextBanner.scss';

export const TextBanner = (props) => {
  return (
    <section className='text-banner'>
      <h1 className='banner-header'>{props.title}</h1>
    </section>
  );
}

TextBanner.propTypes = {
  title: PropTypes.string.isRequired
}
