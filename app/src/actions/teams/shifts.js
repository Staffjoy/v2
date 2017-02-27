import _ from 'lodash';
import 'whatwg-fetch';
import { normalize, Schema, arrayOf } from 'normalizr';
import * as actionTypes from '../../constants/actionTypes';
import { routeToMicroservice } from '../../constants/paths';
import {
  emptyPromise,
  timestampExpired,
  checkStatus,
  parseJSON,
} from '../../utility';

/*
  Exported Actions:
  * getShifts
  * createTeamShift
  * updateTeamShift
  * deleteTeamShift
*/

// schema!
const shiftSchema = new Schema('shifts', { idAttribute: 'uuid' });
const arrayOfShifts = arrayOf(shiftSchema);

// shifts
function requestTeamShifts(teamUuid, params) {
  return {
    type: actionTypes.REQUEST_TEAM_SHIFTS,
    teamUuid,
    params,
  };
}

function receiveTeamShifts(teamUuid, data) {
  return {
    type: actionTypes.RECEIVE_TEAM_SHIFTS,
    teamUuid,
    ...data,
  };
}

// state will update once a shiftUuid is available
function creatingTeamShift(teamUuid) {
  return {
    type: actionTypes.CREATING_TEAM_SHIFT,
    teamUuid,
  };
}

function createdTeamShift(teamUuid, shiftUuid, data) {
  return {
    type: actionTypes.CREATED_TEAM_SHIFT,
    teamUuid,
    shiftUuid,
    ...data,
  };
}

// state will update with the response
function bulkUpdatingTeamShifts(teamUuid, params) {
  return {
    type: actionTypes.BULK_UPDATING_TEAM_SHIFTS,
    teamUuid,
    params,
  };
}

function bulkUpdatedTeamShifts(teamUuid, data) {
  return {
    type: actionTypes.BULK_UPDATED_TEAM_SHIFTS,
    teamUuid,
    ...data,
  };
}

// state will update before the request is made
function updatingTeamShift(teamUuid, shiftUuid, data) {
  return {
    type: actionTypes.UPDATING_TEAM_SHIFT,
    teamUuid,
    shiftUuid,
    ...data,
  };
}

function updatedTeamShift(teamUuid, shiftUuid, data) {
  return {
    type: actionTypes.UPDATED_TEAM_SHIFT,
    teamUuid,
    shiftUuid,
    ...data,
  };
}

function deletingTeamShift(teamUuid, shiftUuid) {
  return {
    type: actionTypes.DELETING_TEAM_SHIFT,
    teamUuid,
    shiftUuid,
  };
}

function deletedTeamShift(teamUuid, shiftUuid) {
  return {
    type: actionTypes.DELETED_TEAM_SHIFT,
    teamUuid,
    shiftUuid,
  };
}

function fetchTeamShifts(companyUuid, teamUuid, params) {
  return (dispatch) => {
    // dispatch action to start the fetch
    dispatch(requestTeamShifts(teamUuid, params));
    const shiftPath = `/v1/companies/${companyUuid}/teams/${teamUuid}/shifts`;

    return fetch(
      routeToMicroservice('company', shiftPath, params),
      { credentials: 'include' })
      .then(checkStatus)
      .then(parseJSON)
      .then((data) => {
        const normalized = normalize(_.get(data, 'shifts', []), arrayOfShifts);

        return dispatch(receiveTeamShifts(teamUuid, {
          data: normalized.entities.shifts,
          order: normalized.result,
          lastUpdate: Date.now(),
        }));
      });
  };
}

function shouldFetchTeamShifts(state, teamUuid, params) {
  const shiftsData = state.teams.shifts;
  const teamShifts = _.get(shiftsData, teamUuid, {});

  // it has never been fetched before
  if (_.isEmpty(shiftsData)) {
    return true;

  // the needed teamUuid is empty
  } else if (_.isEmpty(teamShifts)) {
    return true;

  // teamShifts is at least partially populated with a trusted object at this point
  // the order of these is related to how the 1st fetch might play out

  // the params must be the same as last time
  } else if (!_.isEqual(shiftsData.params, params)) {
    return true;

  // this data set is currently being fetched
  } else if (teamShifts.isFetching) {
    return false;

  // this data set is not complete
  } else if (!teamShifts.completeSet) {
    return true;

  // this data set is stale
  } else if (!teamShifts.lastUpdate ||
    (timestampExpired(teamShifts.lastUpdate, 'TEAM_SHIFTS'))
  ) {
    return true;
  }

  // check if invalidated
  return teamShifts.didInvalidate;
}

export function getTeamShifts(companyUuid, teamUuid, params) {
  return (dispatch, getState) => {
    if (shouldFetchTeamShifts(getState(), teamUuid, params)) {
      return dispatch(fetchTeamShifts(companyUuid, teamUuid, params));
    }
    return emptyPromise();
  };
}

export function createTeamShift(companyUuid, teamUuid, shiftPayload) {
  return (dispatch) => {
    dispatch(creatingTeamShift(teamUuid));
    const shiftPath = `/v1/companies/${companyUuid}/teams/${teamUuid}/shifts`;

    return fetch(
      routeToMicroservice('company', shiftPath),
      {
        credentials: 'include',
        method: 'POST',
        body: JSON.stringify(shiftPayload),
      })
      .then(checkStatus)
      .then(parseJSON)
      .then(data =>
        dispatch(createdTeamShift(teamUuid, data.uuid, {
          data,
        }))
      );
  };
}

export function updateTeamShift(companyUuid, teamUuid, shiftUuid, newData) {
  return (dispatch, getState) => {
    const shifts = _.get(getState().teams.shifts, teamUuid, {});
    const shift = _.get(shifts.data, shiftUuid, {});
    const updateData = _.extend({}, shift, newData);
    dispatch(updatingTeamShift(teamUuid, shiftUuid, { data: updateData }));

    const shiftPath =
      `/v1/companies/${companyUuid}/teams/${teamUuid}/shifts/${shiftUuid}`;

    return fetch(
      routeToMicroservice('company', shiftPath),
      {
        credentials: 'include',
        method: 'PUT',
        body: JSON.stringify(updateData),
      })
      .then(checkStatus)
      .then(parseJSON)
      .then(data =>
        dispatch(updatedTeamShift(teamUuid, shiftUuid, {
          data,
        }))
      );
  };
}

export function bulkUpdateTeamShifts(companyUuid, teamUuid, putBody) {
  return (dispatch) => {
    dispatch(bulkUpdatingTeamShifts(teamUuid, putBody));

    const shiftPath = `/v1/companies/${companyUuid}/teams/${teamUuid}/shifts`;

    return fetch(
      routeToMicroservice('company', shiftPath),
      {
        credentials: 'include',
        method: 'PUT',
        body: JSON.stringify(putBody),
      })
      .then(checkStatus)
      .then(parseJSON)
      .then((data) => {
        const normalized = normalize(_.get(data, 'shifts', []), arrayOfShifts);

        return dispatch(bulkUpdatedTeamShifts(teamUuid, {
          data: normalized.entities.shifts,
        }));
      });
  };
}

export function deleteTeamShift(companyUuid, teamUuid, shiftUuid) {
  return (dispatch) => {
    dispatch(deletingTeamShift(teamUuid, shiftUuid));

    const shiftPath =
      `/v1/companies/${companyUuid}/teams/${teamUuid}/shifts/${shiftUuid}`;

    return fetch(
      routeToMicroservice('company', shiftPath),
      {
        credentials: 'include',
        method: 'DELETE',
      })
      .then(checkStatus)
      .then(parseJSON)
      .then(() =>
        dispatch(deletedTeamShift(teamUuid, shiftUuid))
      );
  };
}
