import React from 'react';
import { storiesOf, action, linkTo } from '@kadira/storybook';

import DragDropContextMock from './mock/DragDropContextMock';

import EmptyCell from '../components/Scheduling/ShiftWeekTable/Section/Row/EmptyCell';

storiesOf('EmptyCell')
  .add('default', () => {
    return (
      <div>
        Hover over the area:
        <DragDropContextMock>
          <EmptyCell
            isOver={false}
            canDrop={false}
            connectDropTarget={() => {}}
            columnId={'42'}
            sectionUuid={'42'}
            timezone={'PST'}
            toggleSchedulingModal={() => {}}
            startDate={'2016-01-01'}
            tableSize={1000}
            viewBy={'yolo'}
            employees={{}}
            jobs={{}}
            modalFormData={{}}
            updateSchedulingModalFormData={() => {}}
            clearSchedulingModalFormData={() => {}}
            createTeamShift={() => {}}
          />
        </DragDropContextMock>
      </div>
    );
  });
