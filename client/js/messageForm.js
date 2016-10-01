import React from 'react';

var MessageForm = React.createClass({
  getInitialState: function() {
    return {
      text: ''
    };
  },
  handleTextChange: function(e) {
    this.setState({
      text: e.target.value
    });
    this.props.onMessageFormActivity();
  },
  handleClick: function(e) {
    this.props.onMessageFormActivity();
  },
  handleSubmit: function(e) {
    e.preventDefault();
    var text = this.state.text.trim();
    if (!text) {
      return;
    }
    this.props.onMessageSubmit(text);
    this.setState({
      text: ''
    });
  },
  render: function() {
    return (
      <form className="messageForm" onClick={this.handleClick} onSubmit={this.handleSubmit}>
        <input
      type="text"
      placeholder="Say something..."
      value={this.state.text}
      onChange={this.handleTextChange}
      />
        <input type="submit" value="Post" />
        </form>
      );
  }
});

export default MessageForm;
