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
  Exported functions:
  * getTeamJobs
*/

// schemas!
const teamJobsSchema = new Schema('jobs', { idAttribute: 'uuid' });
const arrayOfTeamJobs = arrayOf(teamJobsSchema);

// team jobs

function requestTeamJobs(teamUuid) {
  return {
    type: actionTypes.REQUEST_TEAM_JOBS,
    teamUuid,
  };
}

function receiveTeamJobs(teamUuid, data) {
  return {
    type: actionTypes.RECEIVE_TEAM_JOBS,
    teamUuid,
    ...data,
  };
}

function fetchTeamJobs(companyUuid, teamUuid) {
  return (dispatch) => {
    // dispatch action to start the fetch
    dispatch(requestTeamJobs(teamUuid));
    const jobPath = `/v1/companies/${companyUuid}/teams/${teamUuid}/jobs`;

    return fetch(
      routeToMicroservice('company', jobPath),
      { credentials: 'include' })
      .then(checkStatus)
      .then(parseJSON)
      .then((data) => {
        const normalized = normalize(_.get(data, 'jobs', []), arrayOfTeamJobs);

        return dispatch(receiveTeamJobs(teamUuid, {
          data: normalized.entities.jobs,
          order: normalized.result,
          lastUpdate: Date.now(),
        }));
      });
  };
}

function shouldFetchTeamJobs(state, teamUuid) {
  const jobsData = state.teams.jobs;
  const teamJobs = _.get(jobsData, teamUuid, {});

  // no team employees have ever been fetched
  if (_.isEmpty(jobsData)) {
    return true;

  // the needed teamUuid is empty
  } else if (_.isEmpty(teamJobs)) {
    return true;

  // teamJobs is at least partially populated with a trusted object at this point
  // the order of these is related to how the 1st fetch might play out

  // this data set is currently being fetched
  } else if (teamJobs.isFetching) {
    return false;

  // this data set is not complete
  } else if (!teamJobs.completeSet) {
    return true;

  // this data set is stale
  } else if (!teamJobs.lastUpdate ||
    (timestampExpired(teamJobs.lastUpdate, 'TEAM_JOBS'))
  ) {
    return true;
  }

  // check if invalidated
  return teamJobs.didInvalidate;
}

export function getTeamJobs(companyUuid, teamUuid) {
  return (dispatch, getState) => {
    if (shouldFetchTeamJobs(getState(), teamUuid)) {
      return dispatch(fetchTeamJobs(companyUuid, teamUuid));
    }
    return emptyPromise();
  };
}

export function setTeamJob(teamUuid, jobUuid, data) {
  return {
    type: actionTypes.SET_TEAM_JOB,
    teamUuid,
    jobUuid,
    data,
  };
}

function updatingTeamJob(teamUuid, jobUuid, data) {
  return {
    type: actionTypes.UPDATING_TEAM_JOB,
    teamUuid,
    jobUuid,
    data,
  };
}

function updatedTeamJob(teamUuid, jobUuid, data) {
  return {
    type: actionTypes.UPDATED_TEAM_JOB,
    teamUuid,
    jobUuid,
    data,
  };
}

function updatingTeamJobField(jobUuid) {
  return {
    type: actionTypes.UPDATING_TEAM_JOB_FIELD,
    jobUuid,
  };
}

function updatedTeamJobField(jobUuid) {
  return {
    type: actionTypes.UPDATED_TEAM_JOB_FIELD,
    jobUuid,
  };
}

function hideTeamJobFieldSuccess(jobUuid) {
  return {
    type: actionTypes.HIDE_TEAM_JOB_FIELD_SUCCESS,
    jobUuid,
  };
}

export function updateTeamJob(
  companyUuid,
  teamUuid,
  jobUuid,
  newData,
  callback
) {
  return (dispatch, getState) => {
    const jobs = _.get(getState().teams.jobs, teamUuid, {});
    const job = _.get(jobs.data, jobUuid, {});
    const updateData = _.extend({}, job, newData);
    dispatch(updatingTeamJob(teamUuid, jobUuid, newData));

    const jobPath =
      `/v1/companies/${companyUuid}/teams/${teamUuid}/jobs/${jobUuid}`;

    return fetch(
      routeToMicroservice('company', jobPath),
      {
        credentials: 'include',
        method: 'PUT',
        body: JSON.stringify(updateData),
      })
      .then(checkStatus)
      .then(parseJSON)
      .then((data) => {
        if (callback) {
          callback.call(null, data, null);
        }

        dispatch(updatedTeamJob(teamUuid, jobUuid, data));
      });
  };
}

export function updateTeamJobField(companyUuid, teamUuid, jobUuid, newData) {
  return (dispatch) => {
    dispatch(updatingTeamJobField(jobUuid));

    return dispatch(
      updateTeamJob(
        companyUuid,
        teamUuid,
        jobUuid,
        newData,
        (response, error) => {
          if (!error) {
            dispatch(updatedTeamJobField(jobUuid));
            setTimeout(() => {
              dispatch(hideTeamJobFieldSuccess(jobUuid));
            }, 1000);
          }
        }
      )
    );
  };
}

function creatingTeamJob(teamUuid) {
  return {
    type: actionTypes.CREATING_TEAM_JOB,
    teamUuid,
  };
}

function createdTeamJob(teamUuid, jobUuid, data) {
  return {
    type: actionTypes.CREATED_TEAM_JOB,
    teamUuid,
    jobUuid,
    data,
  };
}

export function createTeamJob(companyUuid, teamUuid, jobPayload) {
  return (dispatch) => {
    dispatch(creatingTeamJob());
    const jobPath = `/v1/companies/${companyUuid}/teams/${teamUuid}/jobs`;

    return fetch(
      routeToMicroservice('company', jobPath),
      {
        credentials: 'include',
        method: 'POST',
        body: JSON.stringify(jobPayload),
      })
      .then(checkStatus)
      .then(parseJSON)
      .then((data) => {
        dispatch(createdTeamJob(teamUuid, data.uuid, data));

        setTimeout(() => {
          dispatch(hideTeamJobFieldSuccess(data.uuid));
        }, 1000);
      });
  };
}
