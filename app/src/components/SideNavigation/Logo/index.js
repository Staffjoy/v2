import React, { PropTypes } from 'react';
import { Link } from 'react-router';
import {
  COMPANY_BASE,
  getRoute,
} from 'constants/paths';

require('./navigation-logo.scss');
const imgUrl = require(
  '../../../../../frontend_resources/images/staffjoy.png'
);

function NavigationLogo({ companyUuid }) {
  const route = getRoute(COMPANY_BASE, { companyUuid });
  return (
    <Link to={route} id="navigation-logo">
      <img role="presentation" alt="Staffjoy logo" src={imgUrl} />
    </Link>
  );
}

NavigationLogo.propTypes = {
  companyUuid: PropTypes.string.isRequired,
};

export default NavigationLogo;
