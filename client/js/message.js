import React from 'react';
import moment from 'moment';
import Remarkable from 'remarkable';

var Message = React.createClass({
  composeMessageLine: function() {
    var message = this.props.message;
    if (message.name.length === 0) {
      return '*' + message.text + '*';
    }

    return '**[**' + moment(message.timestamp).format('HH:mm') + '**]** **' + message.name + '**: ' + message.text;
  },
  rawMarkup: function() {
    var md = new Remarkable();
    var rawMarkup = md.render(this.composeMessageLine());
    return {
      __html: rawMarkup
    };
  },
  render: function() {
    // If any non-displayable type messages snuck through, abort the render
    if (this.props.message.type !== 'user-message' && this.props.message.type !== 'system-message') {
      return null;
    }

    var className = 'message';
    className += ' ' + this.props.message.type;
    if (this.props.lastActive && moment(this.props.message.timestamp).isAfter(this.props.lastActive)) {
      className += ' new';
    }
    return (
      <div className={className}>
        <span dangerouslySetInnerHTML={this.rawMarkup()} />
        </div>
      );
  }
});

export default Message;
