import _ from 'lodash';
import React, { PropTypes } from 'react';
import classNames from 'classnames';
import ShiftWeekTableCard from './Card';
import ShiftWeekTableEmptyCell from './EmptyCell';

require('./shift-week-table-row.scss');

function ShiftWeekTableRow({
  clearSchedulingModalFormData,
  createTeamShift,
  deleteTeamShift,
  droppedSchedulingCard,
  editTeamShift,
  employees,
  jobs,
  modalFormData,
  modalOpen,
  rowNumber,
  sectionUuid,
  shiftColumns,
  startDate,
  tableSize,
  timezone,
  toggleSchedulingModal,
  updateSchedulingModalFormData,
  viewBy,
  onCardZAxisChange,
  companyUuid,
  teamUuid,
}) {
  /*
    - A TableRow can display cards, with a max of 1 per column
    - It is dispatched by TableSection
  */

  const columnClasses = classNames({
    'table-column': true,
    [`col-${tableSize}`]: true,
  });

  return (
    <div className="shift-week-table-row table-row">
      {
        _.map(shiftColumns, (column) => {
          const sectionKey = (`${sectionUuid}-row-${rowNumber}-column-
            ${column.columnId}`);
          const cellKey = (`${sectionUuid}-row-${rowNumber}-column-
            ${column.columnId}-cell`);
          return (
            <div
              key={sectionKey}
              className={columnClasses}
            >
              {(() => {
                if (!_.isEmpty(column.shift)) {
                  return (
                    <ShiftWeekTableCard
                      key={cellKey}
                      timezone={timezone}
                      shiftStart={column.shift.start}
                      shiftStop={column.shift.stop}
                      published={column.shift.published}
                      employees={employees}
                      jobs={jobs}
                      viewBy={viewBy}
                      shiftUuid={column.shift.uuid}
                      jobUuid={column.shift.job_uuid}
                      userUuid={column.shift.user_uuid}
                      droppedSchedulingCard={droppedSchedulingCard}
                      columnId={column.columnId}
                      sectionUuid={sectionUuid}
                      deleteTeamShift={deleteTeamShift}
                      toggleSchedulingModal={toggleSchedulingModal}
                      modalOpen={modalOpen}
                      modalFormData={modalFormData}
                      editTeamShift={editTeamShift}
                      updateSchedulingModalFormData={
                        updateSchedulingModalFormData
                      }
                      clearSchedulingModalFormData={
                        clearSchedulingModalFormData
                      }
                      onZAxisChange={onCardZAxisChange}
                      companyUuid={companyUuid}
                      teamUuid={teamUuid}
                    />
                  );
                }
                return (
                  <ShiftWeekTableEmptyCell
                    key={cellKey}
                    toggleSchedulingModal={toggleSchedulingModal}
                    droppedSchedulingCard={droppedSchedulingCard}
                    columnId={column.columnId}
                    sectionUuid={sectionUuid}
                    timezone={timezone}
                    startDate={startDate}
                    tableSize={tableSize}
                    viewBy={viewBy}
                    employees={employees}
                    jobs={jobs}
                    modalFormData={modalFormData}
                    createTeamShift={createTeamShift}
                    updateSchedulingModalFormData={
                      updateSchedulingModalFormData
                    }
                    clearSchedulingModalFormData={
                      clearSchedulingModalFormData
                    }
                  />
                );
              })()}
            </div>
          );
        }
        )
      }
    </div>
  );
}

ShiftWeekTableRow.propTypes = {
  shiftColumns: PropTypes.array.isRequired,
  rowNumber: PropTypes.number.isRequired,
  tableSize: PropTypes.number.isRequired,
  timezone: PropTypes.string.isRequired,
  viewBy: PropTypes.string.isRequired,
  employees: PropTypes.object.isRequired,
  jobs: PropTypes.object.isRequired,
  sectionUuid: PropTypes.string.isRequired,
  droppedSchedulingCard: PropTypes.func.isRequired,
  deleteTeamShift: PropTypes.func.isRequired,
  toggleSchedulingModal: PropTypes.func.isRequired,
  editTeamShift: PropTypes.func.isRequired,
  createTeamShift: PropTypes.func.isRequired,
  modalOpen: PropTypes.bool.isRequired,
  modalFormData: PropTypes.object.isRequired,
  startDate: PropTypes.string.isRequired,
  updateSchedulingModalFormData: PropTypes.func.isRequired,
  clearSchedulingModalFormData: PropTypes.func.isRequired,
  onCardZAxisChange: PropTypes.func.isRequired,
  companyUuid: PropTypes.string,
  teamUuid: PropTypes.string,
};

export default ShiftWeekTableRow;
