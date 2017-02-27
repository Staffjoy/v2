import _ from 'lodash';
import React, { PropTypes } from 'react';
import { Link } from 'react-router';
import { teamNavLinks } from 'constants/sideNavigation';
import * as paths from 'constants/paths';
import { NO_TRANSPARENCY } from 'constants/config';
import { hexToRGBAString } from 'utility';

require('./side-navigation-team-section.scss');

function SideNavigationTeamSection({
  companyUuid,
  teamUuid,
  name,
  color,
  currentPath,
}) {
  const titleStyle = {
    color: hexToRGBAString(color, NO_TRANSPARENCY),
  };
  return (

    <div className="team-group">
      <div className="team-title" style={titleStyle}>
        {name}
      </div>
      {
        _.map(teamNavLinks, (link) => {
          const route = paths.getRoute(link.pathName, {
            companyUuid,
            teamUuid,
          });

          const className = (currentPath === route) ?
            'team-nav-link active' : 'team-nav-link';

          return (
            <Link
              key={route}
              to={route}
              className={className}
            >
              {link.displayName}
            </Link>
          );
        })
      }
    </div>
  );
}

SideNavigationTeamSection.propTypes = {
  companyUuid: PropTypes.string.isRequired,
  teamUuid: PropTypes.string.isRequired,
  name: PropTypes.string.isRequired,
  color: PropTypes.string.isRequired,
  currentPath: PropTypes.string.isRequired,
};

export default SideNavigationTeamSection;
