import _ from 'lodash';
import React, { PropTypes } from 'react';
import { Menu, MenuItem } from 'react-mdl';
import {
  COMPANY_BASE,
  getRoute,
  routeToMicroservice,
} from 'constants/paths';

require('./side-navigation-user-context.scss');

function SideNavigationUserContext({
  companyUuid,
  companyName,
  companyPermissions,
  userName,
  userPhotoUrl,
}) {
  return (
    <div id="user-context-menu" className="user-context">
      <img
        className="profile-icon"
        role="presentation"
        src={userPhotoUrl}
      />
      <div className="user-menu-tag">
        <div className="user-name">{userName}</div>
        <div className="company-name">{companyName}</div>
      </div>
      <Menu target="user-context-menu" valign="top" align="right">
        {
          _.map(companyPermissions, (company) => {
            const companyPath = getRoute(
              COMPANY_BASE, { companyUuid: company.uuid }
            );
            const route = `/?uuid=${company.uuid}/#${companyPath}`;
            const className = (company.uuid === companyUuid) ? 'active' : '';
            const menuKey = `menu-${company.uuid}`;
            const linkKey = `link-${company.uuid}`;
            return (
              <a key={linkKey} href={route}>
                <MenuItem className={className} key={menuKey}>
                  {company.name}
                </MenuItem>
              </a>
            );
          })
        }

        <a href={routeToMicroservice('myaccount')}>
          <MenuItem className="separation">My Account</MenuItem>
        </a>
        <a href={routeToMicroservice('www', '/logout/')}>
          <MenuItem className="separation">Log Out</MenuItem>
        </a>
      </Menu>
    </div>
  );
}

SideNavigationUserContext.propTypes = {
  userName: PropTypes.string.isRequired,
  userPhotoUrl: PropTypes.string.isRequired,
  companyName: PropTypes.string.isRequired,
  companyPermissions: PropTypes.array.isRequired,
  companyUuid: PropTypes.string.isRequired,
};

export default SideNavigationUserContext;
