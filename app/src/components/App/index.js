import _ from 'lodash';
import React, { PropTypes } from 'react';
import { connect } from 'react-redux';
import { Layout, Content } from 'react-mdl';
import NavigationSide from 'components/SideNavigation';
import Intercom from 'components/Intercom';
import * as actions from 'actions';

require('./app.scss');

class App extends React.Component {
  componentDidMount() {
    const { dispatch, companyUuid } = this.props;

    // query whoami endpoint if needed
    dispatch(actions.getWhoAmI());

    // get user data too
    dispatch(actions.getUser());

    // get company info because we are now at the company level
    dispatch(actions.getCompany(companyUuid));

    // get team data because it's needed for side nav paths
    dispatch(actions.getTeams(companyUuid));

    // get intercom settings
    dispatch(actions.fetchIntercomSettings());
  }

  render() {
    const { children, companyUuid, intercomSettings } = this.props;

    return (
      <Layout fixedDrawer>
        <NavigationSide companyUuid={companyUuid} />
        <Content>
          {children}
        </Content>
        {!_.isEmpty(intercomSettings)
        &&
        <Intercom
          {...intercomSettings}
          appID={intercomSettings.app_id}
        />}
      </Layout>
    );
  }
}

App.propTypes = {
  dispatch: PropTypes.func.isRequired,
  children: PropTypes.element,
  companyUuid: PropTypes.string.isRequired,
  intercomSettings: PropTypes.object.isRequired,
};

function mapStateToProps(state, ownProps) {
  return {
    companyUuid: ownProps.routeParams.companyUuid,
    intercomSettings: state.whoami.intercomSettings,
  };
}

export default connect(mapStateToProps)(App);
