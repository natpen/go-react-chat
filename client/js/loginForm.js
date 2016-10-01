import React from 'react';

var LoginForm = React.createClass({
  getInitialState: function() {
    return {
      text: ''
    };
  },
  handleTextChange: function(e) {
    this.setState({
      text: e.target.value
    });
  },
  handleSubmit: function(e) {
    e.preventDefault();
    var text = this.state.text.trim();
    if (!text) {
      return;
    }
    this.props.onNameSubmit(text);
    this.setState({
      text: ''
    });
  },
  render: function() {
    return (
      <form className="loginForm" onSubmit={this.handleSubmit}>
        <input
      type="text"
      placeholder="Choose a name..."
      value={this.state.text}
      onChange={this.handleTextChange}
      />
        <input type="submit" value="Log in!" />
        </form>
      );
  }
});

export default LoginForm;
