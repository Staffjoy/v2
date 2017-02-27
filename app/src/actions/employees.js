import _ from 'lodash';
import 'whatwg-fetch';
import { normalize, Schema, arrayOf } from 'normalizr';
import { getTeams, createTeamEmployee } from './teams';
import { getAssociations } from './associations';
import * as actionTypes from '../constants/actionTypes';
import * as fieldUpdateStatus from '../constants/fieldUpdateStatus';
import { routeToMicroservice } from '../constants/paths';
import {
  emptyPromise,
  timestampExpired,
  checkStatus,
  parseJSON,
} from '../utility';

/*
  Exported functions:
  * getEmployees
  * getEmployee
  * editEmployee
  * initializeEmployees
  * initializeEmployeeSidePanel
  * updateEmployeesSearchFilter
*/

// schemas!
const employeesSchema = new Schema('employees', { idAttribute: 'user_uuid' });
const arrayOfEmployees = arrayOf(employeesSchema);

// employees

function requestEmployees() {
  return {
    type: actionTypes.REQUEST_EMPLOYEES,
  };
}

function receiveEmployees(data) {
  return {
    type: actionTypes.RECEIVE_EMPLOYEES,
    ...data,
  };
}

function requestEmployee() {
  return {
    type: actionTypes.REQUEST_EMPLOYEE,
  };
}

function receiveEmployee(data) {
  return {
    type: actionTypes.RECEIVE_EMPLOYEE,
    ...data,
  };
}

// state will update once an employeeUuid is available
function creatingEmployee() {
  return {
    type: actionTypes.CREATING_EMPLOYEE,
  };
}

function createdEmployee(data) {
  return {
    type: actionTypes.CREATED_EMPLOYEE,
    ...data,
  };
}

function updatingEmployee(data) {
  return {
    type: actionTypes.UPDATING_EMPLOYEE,
    ...data,
  };
}

function updatingEmployeeField(data) {
  return {
    type: actionTypes.UPDATING_EMPLOYEE_FIELD,
    ...data,
  };
}

function updatedEmployee(data) {
  return {
    type: actionTypes.UPDATED_EMPLOYEE,
    ...data,
  };
}

function fetchEmployees(companyUuid) {
  return (dispatch) => {
    // dispatch action to start the fetch
    dispatch(requestEmployees());
    const directoryPath = `/v1/companies/${companyUuid}/directory`;

    return fetch(
      routeToMicroservice('company', directoryPath),
      { credentials: 'include' })
      .then(checkStatus)
      .then(parseJSON)
      .then((data) => {
        const normalized = normalize(
          _.get(data, 'accounts', []),
          arrayOfEmployees
        );

        return dispatch(receiveEmployees({
          data: normalized.entities.employees,
          order: normalized.result,
          lastUpdate: Date.now(),
        }));
      });
  };
}

function fetchEmployee(companyUuid, employeeUuid) {
  return (dispatch) => {
    // dispatch action to start the fetch
    dispatch(requestEmployee());
    const directoryPath =
      `/v1/companies/${companyUuid}/directory/${employeeUuid}`;

    return fetch(
      routeToMicroservice('company', directoryPath),
      { credentials: 'include' })
      .then(checkStatus)
      .then(parseJSON)
      .then((data) => {
        const normalized = normalize(data, employeesSchema);

        return dispatch(receiveEmployee({
          data: normalized.entities.employees,
        }));
      });
  };
}

function shouldFetchEmployees(state) {
  const employeesState = state.employees;
  const employeesData = employeesState.data;

  // it has never been fetched
  if (_.isEmpty(employeesData)) {
    return true;

  // it's currently being fetched
  } else if (employeesState.isFetching) {
    return false;

  // it's been in the UI for more than the allowed threshold
  } else if (!employeesState.lastUpdate ||
    (timestampExpired(employeesState.lastUpdate, 'EMPLOYEES'))
  ) {
    return true;

  // make sure we have a complete collection too
  } else if (!employeesState.completeSet) {
    return true;
  }

  // otherwise, fetch if it's been invalidated
  return employeesState.didInvalidate;
}

function shouldFetchEmployee(state, employeeUuid) {
  const employeesState = state.employees;
  const employeesData = employeesState.data;

  // no employee has ever been fetched
  if (_.isEmpty(employeesData)) {
    return true;

  // the needed employeeUuid is not available
  } else if (!_.has(employeesData, employeeUuid)) {
    return true;

  // the collection has been in the UI for more than the allowed threshold
  } else if (_.has(employeesData, employeeUuid) &&
    (timestampExpired(employeesData.lastUpdate, 'EMPLOYEES'))
  ) {
    return true;
  }

  // otherwise, fetch if it's been invalidated
  return employeesState.didInvalidate;
}

// employee filters

function setFilters(filters) {
  return {
    type: actionTypes.SET_EMPLOYEE_FILTERS,
    data: filters,
  };
}

function initializeFilters() {
  return dispatch =>
    dispatch(setFilters({
      searchQuery: '',
      limitTeam: {},
      status: {},
    }));
}

export function updateEmployeesSearchFilter(query) {
  return setFilters({ searchQuery: query });
}

export function getEmployees(companyUuid) {
  return (dispatch, getState) => {
    if (shouldFetchEmployees(getState())) {
      return dispatch(fetchEmployees(companyUuid));
    }
    return emptyPromise();
  };
}

export function getEmployee(companyUuid, employeeUuid) {
  return (dispatch, getState) => {
    if (shouldFetchEmployee(getState())) {
      return dispatch(fetchEmployee(companyUuid, employeeUuid));
    }
    return emptyPromise();
  };
}

export function createEmployee(companyUuid, employeeData, callback) {
  return (dispatch) => {
    dispatch(creatingEmployee());
    const directoryPath = `/v1/companies/${companyUuid}/directory`;

    return fetch(
      routeToMicroservice('company', directoryPath),
      {
        credentials: 'include',
        method: 'POST',
        body: JSON.stringify(employeeData),
      })
      .then(checkStatus)
      .then(parseJSON)
      .then((data) => {
        if (callback) {
          callback(data);
        }

        const normalized = normalize(data, employeesSchema);
        return dispatch(createdEmployee({
          data: normalized.entities.employees,
        }));
      });
  };
}

export function updateEmployee(companyUuid, employeeUuid, newData, callback) {
  return (dispatch, getState) => {
    const employeeData = _.get(getState().employees.data, employeeUuid, {});
    const updateData = _.extend({}, employeeData, newData);
    dispatch(updatingEmployee({ data: { [employeeUuid]: updateData } }));

    const directoryPath =
      `/v1/companies/${companyUuid}/directory/${employeeUuid}`;

    return fetch(
      routeToMicroservice('company', directoryPath),
      {
        credentials: 'include',
        method: 'PUT',
        body: JSON.stringify(updateData),
      })
      .then(checkStatus)
      .then(parseJSON)
      .then((data) => {
        const normalized = normalize(data, employeesSchema);

        if (callback) {
          callback.call(null, data, null);
        }

        return dispatch(updatedEmployee({
          data: normalized.entities.employees,
        }));
      });
  };
}

export function createEmployeeFromForm(companyUuid) {
  return (dispatch, getState) => {
    const { values } = getState().form['create-employee'];
    const teams = _.pickBy(values.teams);

    // a name and a team must be selected
    if (_.has(values, 'full_name') && !_.isEmpty(teams)) {
      const payload = {
        name: values.full_name,
      };

      // will add email and phone number if they exist
      if (_.has(values, 'email')) {
        payload.email = values.email;
      }

      if (_.has(values, 'phonenumber')) {
        payload.phonenumber = values.phonenumber;
      }

      dispatch(createEmployee(companyUuid, payload, (response) => {
        _.forEach(teams, (value, keyUuid) => {
          dispatch(
            createTeamEmployee(companyUuid, keyUuid, response.user_uuid)
          );
        });
      }));
    }
  };
}

export function updateEmployeeField(companyUuid, employeeUuid, fieldName) {
  return (dispatch, getState) => {
    const { values } = getState().form['employee-side-panel'];
    const value = values[fieldName];
    dispatch(updatingEmployeeField({
      data: { [employeeUuid]: { [fieldName]: fieldUpdateStatus.UPDATING } },
    }));

    return dispatch(
      updateEmployee(
        companyUuid,
        employeeUuid,
        { [fieldName]: value },
        (response, error) => {
          if (!error) {
            dispatch(updatingEmployeeField({
              data: {
                [employeeUuid]: { [fieldName]: fieldUpdateStatus.SUCCESS },
              },
            }));
          }
        },
      )
    );
  };
}

function initialTableFetches(companyUuid) {
  return dispatch => Promise.all([
    dispatch(getTeams(companyUuid)),
    dispatch(getEmployees(companyUuid)),
    dispatch(getAssociations(companyUuid)),
  ]);
}

export function initializeEmployees(companyUuid) {
  return (dispatch) => {
    dispatch(initializeFilters());
    return dispatch(initialTableFetches(companyUuid));
  };
}

function initialEmployeeSidePanelFetches(companyUuid, employeeUuid) {
  return dispatch => Promise.all([
    dispatch(getTeams(companyUuid)),
    dispatch(getEmployee(companyUuid, employeeUuid)),
  ]);
}

export function initializeEmployeeSidePanel(companyUuid, employeeUuid) {
  return (dispatch) => {
    dispatch(initialEmployeeSidePanelFetches(companyUuid, employeeUuid));
  };
}
