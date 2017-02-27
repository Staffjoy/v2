import _ from 'lodash';
import React, { PropTypes } from 'react';
import { Link } from 'react-router';
import { DragSource as dragSource, DropTarget as dropTarget } from 'react-dnd';
import { ScaleModal } from 'boron';
import moment from 'moment';
import 'moment-timezone';
import classNames from 'classnames';
import StaffjoyButton from 'components/StaffjoyButton';
import TimeSelector from 'components/TimeSelector';
import { ModalLayoutSingleColumn } from 'components/ModalLayout';
import {
  MOMENT_SHIFT_CARD_TIMES,
  UNASSIGNED_SHIFT_NAME,
  NO_TRANSPARENCY,
} from 'constants/config';
import {
  formattedDifferenceFromMoment,
  hexToRGBAString,
} from 'utility';
import * as paths from 'constants/paths';
import SchedulingTablePhotoName from './PhotoName';

require('./shift-week-table-card.scss');

const unassignedShiftPhoto = require(
  '../../../../../../../../frontend_resources/images/unassigned_shift_icon.png'
);

class ShiftWeekTableCard extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      selectedZItem: null,
      zAxisOpened: false,
    };
    this.handleZAxisChange = this.handleZAxisChange.bind(this);
    this.openZAxisPicker = this.openZAxisPicker.bind(this);
    this.showEditShiftModal = this.showEditShiftModal.bind(this);
    this.deleteShiftButton = this.deleteShiftButton.bind(this);
    this.saveChangesButton = this.saveChangesButton.bind(this);
    this.onModalClose = this.onModalClose.bind(this);
    this.closeModal = this.closeModal.bind(this);
  }


  onModalClose() {
    this.props.toggleSchedulingModal(false);
    this.props.clearSchedulingModalFormData();
  }

  closeModal() {
    this.modal.hide();
  }

  showEditShiftModal() {
    this.modal.show();
    this.props.toggleSchedulingModal(true);
  }

  deleteShiftButton() {
    const { deleteTeamShift, shiftUuid } = this.props;
    this.closeModal();
    deleteTeamShift(shiftUuid);
  }

  saveChangesButton() {
    const { editTeamShift, shiftUuid, timezone } = this.props;

    editTeamShift(shiftUuid, timezone);
    this.closeModal();
  }

  openZAxisPicker() {
    // Do NOT expand the
    if (!this.state.zAxisOpened) {
      this.setState({ zAxisOpened: true });
    }
  }

  handleZAxisChange(data) {
    this.props.onZAxisChange(data);
  }

  get pickerMenuItems() {
    const {
      viewBy,
      jobs,
      employees,
      companyUuid,
      teamUuid,
    } = this.props;

    const entities = viewBy === 'employee' ?
      _.pickBy(jobs, job => !job.archived)
      : employees;

    if (
      viewBy === 'employee'
      &&
      _.isEmpty(entities)
    ) {
      const route = paths.getRoute(
        paths.TEAM_SETTINGS,
        {
          companyUuid,
          teamUuid,
        }
      );

      return (
        <li>
          <Link
            className="z-axis-picker-option add-job"
            to={route}
          >
            Add New Job
          </Link>
        </li>
      );
    }

    return _.map(
      entities,
      (entity) => {
        const entityUuid = viewBy === 'employee'
                                          ? entity.uuid
                                          : entity.user_uuid;

        const textStyle = {};
        if (viewBy === 'employee') {
          textStyle.color = hexToRGBAString(
            entity.color, NO_TRANSPARENCY
          );
        }

        return (
          <li key={entityUuid}>
            <span
              className="z-axis-picker-option"
              data-name={entity.name}
              data-uuid={entityUuid}
              style={textStyle}
            >
              {entity.name}
            </span>
          </li>
        );
      }
    );
  }

  render() {
    const { employees, jobs, shiftStart, shiftStop, timezone, published,
      viewBy, jobUuid, userUuid, connectDragSource, connectDropTarget,
      isDragging, columnId, updateSchedulingModalFormData, modalFormData,
      isOver } = this.props;
    const { zAxisOpened } = this.state;
    const startMoment = moment.utc(shiftStart).tz(timezone);
    const startDisplay = startMoment.format(MOMENT_SHIFT_CARD_TIMES);
    const stopMoment = moment.utc(shiftStop).tz(timezone);
    const stopDisplay = stopMoment.format(MOMENT_SHIFT_CARD_TIMES);
    const formattedDuration = formattedDifferenceFromMoment(
                                startMoment,
                                stopMoment
                              );

    // determine whether the save button on the modal should be enabled
    const modalStartMoment = moment(
      _.get(modalFormData, 'startFieldText', ''),
      MOMENT_SHIFT_CARD_TIMES
    );
    const modalStopMoment = moment(
      _.get(modalFormData, 'stopFieldText', ''),
      MOMENT_SHIFT_CARD_TIMES
    );

    const disabledSave = !modalStartMoment.isValid() ||
                         !modalStopMoment.isValid() ||
                         modalStartMoment.isSameOrAfter(modalStopMoment);

    let zAxisElement;
    const coloredProperty = (published) ? 'backgroundColor' : 'color';

    if (viewBy === 'employee') {
      if (jobUuid) {
        zAxisElement =
          (<div
            className="job-label"
            style={{
              [coloredProperty]: hexToRGBAString(
                jobs[jobUuid].color, NO_TRANSPARENCY
              ),
            }}
          >
            <span>{jobs[jobUuid].name}</span>
          </div>);
      } else if (this.state.selectedZItem) {
        zAxisElement =
          (<div className="job-label">
            <span>{this.state.selectedZItem.name}</span>
          </div>);
      } else {
        zAxisElement =
          (<div
            className="job-label job-label-none"
            title="No job assigned to this shift."
          >
            <span>- -</span>
          </div>);
      }
    } else if (viewBy === 'job') {
      const photoUrl = (userUuid !== '') ?
        employees[userUuid].photo_url : unassignedShiftPhoto;
      const userName = (userUuid !== '') ?
        employees[userUuid].name : UNASSIGNED_SHIFT_NAME;
      zAxisElement = (
        <SchedulingTablePhotoName
          name={userName}
          photoUrl={photoUrl}
        />
      );
    }

    const classes = classNames({
      'shift-week-table-card': true,
      'card-is-dragging': isDragging,
      published,
      isOver,
    });

    const zAxisClasses = classNames({
      'shift-z-axis': true,
      'z-axis-expanded': zAxisOpened,
      published,
    });

    return connectDropTarget(connectDragSource(
      <div className={classes}>
        <ScaleModal
          ref={(modal) => { this.modal = modal; }}
          modalStyle={{ width: '420px' }}
          contentStyle={{ borderRadius: '3px' }}
          onHide={this.onModalClose}
        >
          <ModalLayoutSingleColumn
            buttons={[
              <StaffjoyButton
                buttonType="outline"
                key="delete-button"
                size="small"
                onClick={this.deleteShiftButton}
              >
                Delete
              </StaffjoyButton>,
              <StaffjoyButton
                buttonType="primary"
                size="small"
                key="save-button"
                onClick={this.saveChangesButton}
                disabled={disabledSave}
              >
                Save
              </StaffjoyButton>,
            ]}
          >
            <TimeSelector
              timezone={timezone}
              start={shiftStart}
              stop={shiftStop}
              date={columnId}
              formCallback={updateSchedulingModalFormData}
            />
          </ModalLayoutSingleColumn>
        </ScaleModal>
        <div className="shift-details" onClick={this.showEditShiftModal}>
          <span className="duration">{formattedDuration}</span>
          <div>
            <div className="card-label">Start:</div>
            <div className="card-time">{startDisplay}</div>
          </div>
          <div>
            <div className="card-label">End:</div>
            <div className="card-time">{stopDisplay}</div>
          </div>
        </div>
        <div className={zAxisClasses} onClick={this.openZAxisPicker}>
          {
            zAxisOpened
            ? <div
              className="z-axis-picker-wrapper"
              key="z-axis-picker"
              ref={(element) => {
                const el = element;
                if (el && el.parentElement) {
                  const scrollWidth = el.parentElement.scrollWidth;
                  el.style.width = `${scrollWidth}px`;
                }
              }}
            >
              <div
                className="z-axis-picker"
                ref={(element) => { this.pickerElement = element; }}
                onClick={(event) => {
                  const element = event.target;
                  if (element.className === 'z-axis-picker-option') {
                    this.setState({
                      selectedZItem: {
                        name: element.dataset.name,
                        uuid: element.dataset.uuid,
                      },
                      zAxisOpened: false,
                    });
                    const itemUuid =
                        this.props.viewBy === 'employee'
                        ? 'job_uuid'
                        : 'user_uuid';

                    this.handleZAxisChange({
                      key: itemUuid,
                      value: element.dataset.uuid,
                      shiftUuid: this.props.shiftUuid,
                    });
                  }
                }}
                onMouseLeave={() => { this.setState({ zAxisOpened: false }); }}
              >
                <ul className="z-axis-picker-menu">
                  {this.pickerMenuItems}
                </ul>
              </div>
            </div>
            : zAxisElement
          }
        </div>
      </div>
    ));
  }
}

ShiftWeekTableCard.propTypes = {
  columnId: PropTypes.string.isRequired,
  timezone: PropTypes.string.isRequired,
  shiftStart: PropTypes.string.isRequired,
  shiftStop: PropTypes.string.isRequired,
  shiftUuid: PropTypes.string.isRequired,
  jobUuid: PropTypes.string,
  userUuid: PropTypes.string,
  viewBy: PropTypes.string.isRequired,
  employees: PropTypes.object.isRequired,
  jobs: PropTypes.object.isRequired,
  isDragging: PropTypes.bool.isRequired,
  connectDragSource: PropTypes.func.isRequired,
  isOver: PropTypes.bool.isRequired,
  connectDropTarget: PropTypes.func.isRequired,
  deleteTeamShift: PropTypes.func.isRequired,
  toggleSchedulingModal: PropTypes.func.isRequired,
  modalFormData: PropTypes.object.isRequired,
  updateSchedulingModalFormData: PropTypes.func.isRequired,
  clearSchedulingModalFormData: PropTypes.func.isRequired,
  editTeamShift: PropTypes.func.isRequired,
  published: PropTypes.bool.isRequired,
  onZAxisChange: PropTypes.func.isRequired,
  companyUuid: PropTypes.string,
  teamUuid: PropTypes.string,
};

/*
  There are some props needed for the drag and drop wrappers, but they aren't
  used from inside the actual component. We will just list them here for
  reference.

  * sectionUuid: PropTypes.string.isRequired
  * droppedSchedulingCard: PropTypes.func.isRequired
  * modalOpen: PropTypes.bool.isRequired
*/

// react dnd - dragging
const cardDragSpec = {
  canDrag(props) {
    return !props.modalOpen;
  },
  beginDrag(props) {
    return {
      shiftUuid: props.shiftUuid,
      oldColumnId: props.columnId,
    };
  },
};

function collectDrag(connect, monitor) {
  return {
    connectDragSource: connect.dragSource(),
    isDragging: monitor.isDragging(),
  };
}

const cardDropSpec = {
  drop(props, monitor) {
    const { columnId, sectionUuid, droppedSchedulingCard } = props;
    const { shiftUuid, oldColumnId } = monitor.getItem();

    // trigger our drag action/prop
    droppedSchedulingCard(shiftUuid, oldColumnId, sectionUuid, columnId);
  },
};

function collectDrop(connect, monitor) {
  return {
    connectDropTarget: connect.dropTarget(),
    isOver: monitor.isOver(),
    canDrop: monitor.canDrop(),
  };
}

export default _.flow(
  dragSource('card', cardDragSpec, collectDrag),
  dropTarget('card', cardDropSpec, collectDrop)
)(ShiftWeekTableCard);
