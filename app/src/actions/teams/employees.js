import _ from 'lodash';
import 'whatwg-fetch';
import { normalize, Schema, arrayOf } from 'normalizr';
import { invalidateAssociations } from '../associations';
import * as actionTypes from '../../constants/actionTypes';
import { routeToMicroservice } from '../../constants/paths';
import {
  emptyPromise,
  timestampExpired,
  checkStatus,
  parseJSON,
} from '../../utility';

/*
  Exported functions:
  * getTeamEmployees
  * createTeamEmployee
*/

// schemas!
const teamEmployeesSchema = new Schema(
  'employees',
  { idAttribute: 'user_uuid' }
);
const arrayOfTeamEmployees = arrayOf(teamEmployeesSchema);

// team employees

function requestTeamEmployees(teamUuid) {
  return {
    type: actionTypes.REQUEST_TEAM_EMPLOYEES,
    teamUuid,
  };
}

function receiveTeamEmployees(teamUuid, data) {
  return {
    type: actionTypes.RECEIVE_TEAM_EMPLOYEES,
    teamUuid,
    ...data,
  };
}

function creatingTeamEmployee(teamUuid) {
  return {
    type: actionTypes.CREATING_TEAM_EMPLOYEE,
    teamUuid,
  };
}

function createdTeamEmployee(teamUuid, userUuid, data) {
  return {
    type: actionTypes.CREATED_TEAM_EMPLOYEE,
    teamUuid,
    userUuid,
    ...data,
  };
}

function fetchTeamEmployees(companyUuid, teamUuid) {
  return (dispatch) => {
    // dispatch action to start the fetch
    dispatch(requestTeamEmployees(teamUuid));
    const teamEmployeePath =
      `/v1/companies/${companyUuid}/teams/${teamUuid}/workers`;

    return fetch(
      routeToMicroservice('company', teamEmployeePath),
      { credentials: 'include' })
      .then(checkStatus)
      .then(parseJSON)
      .then((data) => {
        const normalized = normalize(
          _.get(data, 'workers', []),
          arrayOfTeamEmployees
        );

        return dispatch(receiveTeamEmployees(teamUuid, {
          data: normalized.entities.employees,
          order: normalized.result,
          lastUpdate: Date.now(),
        }));
      });
  };
}

function shouldFetchTeamEmployees(state, teamUuid) {
  const employeesData = state.teams.employees;
  const teamEmployees = _.get(employeesData, teamUuid, {});

  // no team employees have ever been fetched
  if (_.isEmpty(employeesData)) {
    return true;

  // the needed teamUuid is empty
  } else if (_.isEmpty(teamEmployees)) {
    return true;

  // teamEmployees is at least partially populated with a trusted object at this point
  // the order of these is related to how the 1st fetch might play out

  // this data set is currently being fetched
  } else if (teamEmployees.isFetching) {
    return false;

  // this data set is not complete
  } else if (!teamEmployees.completeSet) {
    return true;

  // this data set is stale
  } else if (!teamEmployees.lastUpdate ||
    (timestampExpired(teamEmployees.lastUpdate, 'TEAM_EMPLOYEES'))
  ) {
    return true;
  }

  // check if invalidated
  return teamEmployees.didInvalidate;
}

export function getTeamEmployees(companyUuid, teamUuid) {
  return (dispatch, getState) => {
    if (shouldFetchTeamEmployees(getState(), teamUuid)) {
      return dispatch(fetchTeamEmployees(companyUuid, teamUuid));
    }
    return emptyPromise();
  };
}

export function createTeamEmployee(companyUuid, teamUuid, userUuid) {
  return (dispatch) => {
    dispatch(creatingTeamEmployee(teamUuid));
    const workerPath = `/v1/companies/${companyUuid}/teams/${teamUuid}/workers`;

    return fetch(
      routeToMicroservice('company', workerPath), {
        credentials: 'include',
        method: 'POST',
        body: JSON.stringify({ user_uuid: userUuid }),
      })
      .then(checkStatus)
      .then(parseJSON)
      .then((data) => {
        dispatch(invalidateAssociations());
        return dispatch(createdTeamEmployee(teamUuid, data.uuid, {
          data,
        }));
      });
  };
}
