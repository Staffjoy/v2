import _ from 'lodash';
import React, { PropTypes } from 'react';
import { connect } from 'react-redux';
import { hashHistory } from 'react-router';
import * as actions from 'actions';
import LoadingScreen from 'components/LoadingScreen';
import SearchField from 'components/SearchField';
import { COMPANY_EMPLOYEE, getRoute } from 'constants/paths';
import CreateEmployeeModal from './CreateEmployeeModal';
import Table from './Table';
import * as rowTypes from './Table/Row/rowTypes';

require('./employees.scss');

class Employees extends React.Component {
  componentDidMount() {
    const { dispatch } = this.props;

    // get the employees for the whole company
    dispatch(actions.initializeEmployees(this.props.companyUuid));
  }

  render() {
    const {
      children,
      companyUuid,
      employees,
      isFetching,
      updateSearchFilter,
      teams,
      tableRowClicked,
    } = this.props;

    const columns = [
      {
        columnId: 'employees',
        colWidth: 4,
        displayName: 'Employees',
        component: rowTypes.PHOTO_NAME,
        propDataFields: {
          name: 'name',
          photoUrl: 'photo_url',
        },
      },
      {
        columnId: 'contact',
        colWidth: 3,
        displayName: 'Contact',
        component: rowTypes.CONTACT_INFO,
        propDataFields: {
          email: 'email',
          phoneNumber: 'phonenumber',
        },
      },
      {
        columnId: 'team',
        colWidth: 3,
        displayName: 'Team',
        component: rowTypes.INFO_LIST,
        propDataFields: {
          name: 'name',
          photoUrl: 'photo_url',
        },
      },
      {
        columnId: 'status',
        colWidth: 2,
        displayName: 'Status',
        component: rowTypes.BOOLEAN_LABEL,
        propDataFields: {
          booleanField: 'confirmed_and_active',
        },
        callback(fieldValue) { return (fieldValue) ? 'Active' : 'Inactive'; },
      },
    ];

    if (isFetching) {
      return (
        <LoadingScreen />
      );
    }

    // params have initialized if it's gotten this far

    return (
      <div className="employees">
        <div className="employees-container">
          <div className="employees-control-panel">
            <SearchField width={200} onChange={updateSearchFilter} />
            <div className="employees-control-panel-buttons">
              <CreateEmployeeModal companyUuid={companyUuid} teams={teams} />
            </div>
          </div>
          <div className="scrolling-panel">
            <Table
              columns={columns}
              rows={employees}
              onRowClick={tableRowClicked}
              uuidKeyName="user_uuid"
            />
          </div>
        </div>
        <div className="employees-sidebar">
          {children}
        </div>
      </div>
    );
  }
}

Employees.propTypes = {
  dispatch: PropTypes.func.isRequired,
  isFetching: PropTypes.bool.isRequired,
  companyUuid: PropTypes.string.isRequired,
  employees: PropTypes.array.isRequired,
  // filters: PropTypes.object.isRequired,
  updateSearchFilter: PropTypes.func.isRequired,
  teams: PropTypes.array.isRequired,
  children: PropTypes.element,
  tableRowClicked: PropTypes.func.isRequired,
};

function mapStateToProps(state, ownProps) {
  // apply filters to our list of employees
  const employees = [];
  const searchQuery = _.get(state.employees.filters, 'searchQuery', '');

  _.each(state.employees.data, (employee) => {
    if (employee.name.toLowerCase().includes(searchQuery) ||
        employee.email.includes(searchQuery)) {
      employees.push(employee);
    }
  });

  return {
    companyUuid: ownProps.routeParams.companyUuid,
    isFetching: !state.employees.lastUpdate || state.employees.isFetching,
    // filters: state.employees.filters,
    employees,
    teams: _.values(state.teams.data),
  };
}

const mapDispatchToProps = (dispatch, ownProps) => ({
  tableRowClicked: (event, employeeUuid) => {
    hashHistory.push(
      getRoute(COMPANY_EMPLOYEE, {
        companyUuid: ownProps.routeParams.companyUuid,
        employeeUuid,
      })
    );
  },
  updateSearchFilter: (event) => {
    dispatch(actions.updateEmployeesSearchFilter(event.target.value));
  },
  dispatch,
});

export default connect(mapStateToProps, mapDispatchToProps)(Employees);
