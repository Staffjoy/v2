import _ from 'lodash';
import $ from 'npm-zepto';
import React, { PropTypes } from 'react';
import ModalListSelectableItem from './SelectableItem';

require('./selectable-modal-list.scss');

class SelectableModalList extends React.Component {

  constructor(props) {
    super(props);
    this.selectElement = this.selectElement.bind(this);
    this.state = {
      selections: {},
    };
  }

  componentWillMount() {
    const { records, selectedUuid, formField, formCallback,
    uuidKey } = this.props;
    const selections = {};
    _.forEach(records, (record) => {
      selections[record[uuidKey]] = false;
    });

    if (_.has(selections, selectedUuid)) {
      selections[selectedUuid] = true;
    }

    this.setState({ selections });
    formCallback({ [formField]: selections });
  }

  selectElement(event) {
    const { formField, formCallback } = this.props;
    const newUuid = $(event.target)
                      .closest('.modal-list-selectable-item')
                      .data('uuid');
    const selections = _.extend({}, this.state.selections);
    selections[newUuid] = !selections[newUuid];
    this.setState({ selections });
    formCallback({ [formField]: selections });
  }

  render() {
    const { error, records, displayByProperty, uuidKey } = this.props;
    const { selections } = this.state;

    let errorMessage;
    if (error) {
      errorMessage = (
        <div className="error-message">
          {error}
        </div>
      );
    }

    return (
      <div className="modal-selectable-list">
        {errorMessage}
        {
          _.map(records, (record) => {
            const selectorKey = `modal-list-selectable-item-${record[uuidKey]}`;

            return (
              <ModalListSelectableItem
                key={selectorKey}
                selected={selections[record[uuidKey]]}
                changeFunction={this.selectElement}
                uuid={record[uuidKey]}
                name={record[displayByProperty]}
              />
            );
          })
        }
      </div>
    );
  }
}

SelectableModalList.propTypes = {
  error: PropTypes.string,
  records: PropTypes.arrayOf(React.PropTypes.object),
  displayByProperty: PropTypes.string.isRequired,
  selectedUuid: PropTypes.string,
  formCallback: PropTypes.func.isRequired,
  formField: PropTypes.string.isRequired,
  uuidKey: PropTypes.string.isRequired,
};

export default SelectableModalList;
