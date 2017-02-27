import * as actionTypes from '../constants/actionTypes';

import {
  getTeamJobs,
  getTeam,
} from './teams';

function initialFetches(companyUuid, teamUuid) {
  return (dispatch) => {
    Promise.all([
      dispatch(getTeam(companyUuid, teamUuid)),
      dispatch(getTeamJobs(companyUuid, teamUuid)),
    ]);
  };
}

export function initializeSettings(
  companyUuid,
  teamUuid,
) {
  return (dispatch) => {
    // use promise to guarantee that current team is available in state
    dispatch(initialFetches(companyUuid, teamUuid));
  };
}

export function setColorPicker(colorPicker) {
  return {
    type: actionTypes.SET_COLOR_PICKER,
    colorPicker,
  };
}

export function setFilters(filters) {
  return {
    type: actionTypes.SET_SETTINGS_FILTERS,
    filters,
  };
}

export function setNewTeamJob(data) {
  return {
    type: actionTypes.SET_SETTINGS_NEW_TEAM_JOB,
    data,
  };
}
