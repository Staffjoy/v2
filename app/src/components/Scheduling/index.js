import _ from 'lodash';
import React, { PropTypes } from 'react';
import { connect } from 'react-redux';
import { DragDropContext as dragDropContext } from 'react-dnd';
import HTML5Backend from 'react-dnd-html5-backend';
import $ from 'npm-zepto';
import * as actions from 'actions';
import LoadingScreen from 'components/LoadingScreen';
import StaffjoyButton from 'components/StaffjoyButton';
import SearchField from 'components/SearchField';
import ShiftWeekTable from './ShiftWeekTable';
import SchedulingDateController from './DateController';
import SchedulingViewByController from './ViewByController';
import CreateShiftModal from './CreateShiftModal';

require('./scheduling.scss');

class Scheduling extends React.Component {

  componentDidMount() {
    const { dispatch, companyUuid, teamUuid, routeQuery } = this.props;

    // initialize and get shifts
    // this will get the team, initialize params, initialize filters
    dispatch(actions.initializeScheduling(companyUuid, teamUuid, routeQuery));
  }

  componentWillReceiveProps(nextProps) {
    const { dispatch, companyUuid, teamUuid, routeQuery } = this.props;
    const newTeamUuid = nextProps.teamUuid;

    // if changing between different scheduling views, the same component
    // instance is used, so this is checks to get data
    if (newTeamUuid !== teamUuid) {
      dispatch(
        actions.initializeScheduling(companyUuid, newTeamUuid, routeQuery)
      );
    }
  }

  render() {
    const { isFetching, updateSearchFilter, params, filters, employees, jobs,
      shifts, timezone, stepDateRange, changeViewBy, droppedSchedulingCard,
      deleteTeamShift, toggleSchedulingModal, modalOpen, editTeamShift,
      updateSchedulingModalFormData, createTeamShift, modalFormData,
      clearSchedulingModalFormData, publishTeamShifts, isSaving, companyUuid,
      teamUuid } = this.props;
    const tableSize = 7;
    const viewBy = filters.viewBy;
    const startDate = params.startDate;

    const allShiftsPublished = !_.some(shifts, x => !x.published);
    let publishAction = 'Publish Week';
    let publishButtonStyle = 'primary';

    if (allShiftsPublished && shifts.length > 0) {
      publishAction = 'Unpublish Week';
      publishButtonStyle = 'outline-error';
    }

    if (isFetching) {
      return (
        <LoadingScreen />
      );
    }

    // TODO - add publish button into top controls

    return (
      <div className="scheduling-container">
        <ul className="scheduling-controls">
          <li className="control-unit">
            <SchedulingDateController
              queryStart={params.range.start}
              queryStop={params.range.stop}
              timezone={timezone}
              stepDateRange={stepDateRange}
              disabled={isSaving}
            />
          </li>
          <li className="control-unit">
            <SchedulingViewByController
              onClick={changeViewBy}
              viewBy={viewBy}
              disabled={isSaving}
            />
          </li>
          <li className="control-unit control-unit-hidden-on-collapse">
            <SearchField
              width={200}
              onChange={updateSearchFilter}
              darkBackground
              disabled={isSaving}
            />
          </li>
          <li className="publish-week-btn control-unit-hidden-on-collapse">
            <StaffjoyButton
              buttonType={publishButtonStyle}
              onClick={publishTeamShifts}
              disabled={isSaving}
            >
              {publishAction}
            </StaffjoyButton>
          </li>
          <li className="create-shift-btn control-unit-hidden-on-collapse">
            <CreateShiftModal
              tableSize={tableSize}
              startDate={startDate}
              timezone={timezone}
              modalCallbackToggle={toggleSchedulingModal}
              containerComponent="button"
              containerProps={{
                buttonType: 'neutral',
                disabled: isSaving,
              }}
              viewBy={viewBy}
              employees={employees}
              jobs={jobs}
              onSave={createTeamShift}
              modalFormData={modalFormData}
              updateSchedulingModalFormData={updateSchedulingModalFormData}
              clearSchedulingModalFormData={clearSchedulingModalFormData}
            />
          </li>
        </ul>
        {(() =>
          // TODO when we have more views, determine which view type to use
          // if (props.params.viewType === 'week') {
          <ShiftWeekTable
            droppedSchedulingCard={droppedSchedulingCard}
            startDate={startDate}
            tableSize={tableSize}
            timezone={timezone}
            employees={employees}
            jobs={jobs}
            shifts={shifts}
            filters={filters}
            viewBy={viewBy}
            deleteTeamShift={deleteTeamShift}
            toggleSchedulingModal={toggleSchedulingModal}
            modalOpen={modalOpen}
            modalFormData={modalFormData}
            editTeamShift={editTeamShift}
            createTeamShift={createTeamShift}
            updateSchedulingModalFormData={updateSchedulingModalFormData}
            clearSchedulingModalFormData={clearSchedulingModalFormData}
            onCardZAxisChange={this.props.handleCardZAxisChange}
            isSaving={isSaving}
            companyUuid={companyUuid}
            teamUuid={teamUuid}
          />
        )()}
      </div>
    );
  }
}

Scheduling.propTypes = {
  dispatch: PropTypes.func.isRequired,
  isFetching: PropTypes.bool.isRequired,
  isSaving: PropTypes.bool.isRequired,
  routeQuery: PropTypes.object.isRequired,
  companyUuid: PropTypes.string.isRequired,
  teamUuid: PropTypes.string.isRequired,
  params: PropTypes.object.isRequired,
  filters: PropTypes.object.isRequired,
  employees: PropTypes.object.isRequired,
  jobs: PropTypes.object.isRequired,
  shifts: PropTypes.arrayOf(PropTypes.object).isRequired,
  timezone: PropTypes.string.isRequired,
  updateSearchFilter: PropTypes.func.isRequired,
  stepDateRange: PropTypes.func.isRequired,
  changeViewBy: PropTypes.func.isRequired,
  droppedSchedulingCard: PropTypes.func.isRequired,
  deleteTeamShift: PropTypes.func.isRequired,
  toggleSchedulingModal: PropTypes.func.isRequired,
  editTeamShift: PropTypes.func.isRequired,
  createTeamShift: PropTypes.func.isRequired,
  modalOpen: PropTypes.bool.isRequired,
  modalFormData: PropTypes.object.isRequired,
  updateSchedulingModalFormData: PropTypes.func.isRequired,
  clearSchedulingModalFormData: PropTypes.func.isRequired,
  publishTeamShifts: PropTypes.func.isRequired,
  handleCardZAxisChange: PropTypes.func.isRequired,
};

function mapStateToProps(state, ownProps) {
  const teamUuid = ownProps.routeParams.teamUuid;

  // consts for team data
  const teamData = _.get(state.teams.data, teamUuid, {});
  const isTeamFetching = state.teams.isFetching;
  const timezone = _.get(teamData, 'timezone', 'UTC');

  // consts for shift data
  const shiftState = _.get(state.teams.shifts, teamUuid, {});
  const shifts = _.values(_.get(shiftState, 'data', {}));
  const isShiftSaving = _.get(shiftState, 'isSaving', false);

  // consts for job data
  const jobState = _.get(state.teams.jobs, teamUuid, {});
  const isJobFetching = _.get(jobState, 'isFetching', true);
  const jobs = _.get(jobState, 'data', {});

  // consts for employee data
  const employeeState = _.get(state.teams.employees, teamUuid, {});
  const isEmployeeFetching = _.get(employeeState, 'isFetching', true);
  const employees = _.get(employeeState, 'data', {});

  // scheduling
  const schedulingState = state.scheduling;
  const schedulingParams = schedulingState.params;
  const schedulingFilters = schedulingState.filters;

  const isSchedulingParamsFetching = schedulingParams.isFetching;
  const isSchedulingFiltersFetching = schedulingFilters.isFetching;

  const isFetching =
    isTeamFetching ||
    isJobFetching ||
    isEmployeeFetching ||
    isSchedulingFiltersFetching ||
    isSchedulingParamsFetching ||
    _.isEmpty(schedulingParams);

  const isSaving =
    isShiftSaving;

  return {
    companyUuid: ownProps.routeParams.companyUuid,
    routeQuery: ownProps.location.query,
    isFetching,
    isSaving,
    params: schedulingParams,
    filters: schedulingFilters,
    modalOpen: schedulingState.modal.modalOpen,
    modalFormData: schedulingState.modal.formData,
    teamUuid,
    timezone,
    jobs,
    employees,
    shifts,
  };
}

const mapDispatchToProps = (dispatch, ownProps) => ({
  updateSearchFilter: (event) => {
    dispatch(actions.updateSchedulingSearchFilter(event.target.value));
  },
  changeViewBy: (event) => {
    const newView = $(event.target).data('id');
    const { teamUuid } = ownProps.routeParams;

    dispatch(actions.changeViewBy(newView, teamUuid));
  },
  stepDateRange: (event) => {
    const { companyUuid, teamUuid } = ownProps.routeParams;
    const direction = $(event.target)
      .closest('.square-button')
      .data('direction');

    dispatch(actions.stepDateRange(companyUuid, teamUuid, direction));
  },
  droppedSchedulingCard: (shiftUuid, oldColumnId, sectionUuid, newColumnId) => {
    const { companyUuid, teamUuid } = ownProps.routeParams;

    dispatch(actions.droppedSchedulingCard(
      companyUuid,
      teamUuid,
      shiftUuid,
      oldColumnId,
      sectionUuid,
      newColumnId
    ));
  },
  editTeamShift: (shiftUuid, timezone) => {
    const { companyUuid, teamUuid } = ownProps.routeParams;

    dispatch(actions.editTeamShiftFromModal(
      companyUuid,
      teamUuid,
      shiftUuid,
      timezone
    ));
  },
  createTeamShift: (timezone) => {
    const { companyUuid, teamUuid } = ownProps.routeParams;
    dispatch(
      actions.createTeamShiftsFromModal(companyUuid, teamUuid, timezone)
    );
  },
  deleteTeamShift: (shiftUuid) => {
    const { companyUuid, teamUuid } = ownProps.routeParams;

    dispatch(actions.deleteTeamShift(companyUuid, teamUuid, shiftUuid));
  },
  publishTeamShifts: () => {
    const { companyUuid, teamUuid } = ownProps.routeParams;

    dispatch(actions.publishTeamShifts(companyUuid, teamUuid));
  },
  toggleSchedulingModal: (value) => {
    dispatch(actions.toggleSchedulingModal(value));
  },
  updateSchedulingModalFormData: (data) => {
    dispatch(actions.updateSchedulingModalFormData(data));
  },
  clearSchedulingModalFormData: (data) => {
    dispatch(actions.clearSchedulingModalFormData(data));
  },
  handleCardZAxisChange: ({ key, shiftUuid, value }) => {
    const { companyUuid, teamUuid } = ownProps.routeParams;
    const newData = { [key]: value };

    dispatch(actions.updateTeamShift(
      companyUuid,
      teamUuid,
      shiftUuid,
      newData
    ));
  },
  dispatch,
});

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(dragDropContext(HTML5Backend)(Scheduling));
